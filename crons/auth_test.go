package crons

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"smart/tools/enums"
)

type mockAuth struct {
	infoErr        error
	info           map[string]string
	serial         string
	serialErr      error
	decryptMap     map[string]string
	decryptErr     error
	lastState      string
	lastInfo       map[string]string
	updateInfoErr  error
	updateStateErr error
}

func (m *mockAuth) GetAuthInfo(_ context.Context) (map[string]string, error) {
	return m.info, m.infoErr
}
func (m *mockAuth) GenerateSystemSerialNumber(_ context.Context) (string, error) {
	return m.serial, m.serialErr
}
func (m *mockAuth) RsaDecrypt(_ context.Context, authCode string) (map[string]string, error) {
	if authCode == "" {
		return nil, errors.New("empty")
	}
	return m.decryptMap, m.decryptErr
}
func (m *mockAuth) UpdateAuthState(_ context.Context, s string) error {
	m.lastState = s
	return m.updateStateErr
}
func (m *mockAuth) UpdateAuthInfo(_ context.Context, info map[string]string) error {
	m.lastInfo = info
	return m.updateInfoErr
}
func (m *mockAuth) RsaEncrypt(data, keyBytes []byte) []byte {
	return []byte("MOCKED_ENCRYPTED")
}
func (m *mockAuth) CleanAuthRecord(_ context.Context) error {
	return nil
}

func withMock(t *testing.T, m *mockAuth, fn func()) {
	old := authProvider
	authProvider = m
	defer func() { authProvider = old }()
	fn()
}

func Test_CheckSystemAuth_InitOnMissingInfo(t *testing.T) {
	m := &mockAuth{
		infoErr: errors.New("notfound"),
		serial:  "SERIAL_X",
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateFailed {
			t.Fatalf("state=%s", m.lastState)
		}
		if m.lastInfo == nil || m.lastInfo["productID"] != "SERIAL_X" || m.lastInfo["authCode"] != "" || m.lastInfo["authTime"] != "未授权" {
			t.Fatalf("info=%v", m.lastInfo)
		}
	})
}

func Test_CheckSystemAuth_ProductIDMismatch(t *testing.T) {
	m := &mockAuth{
		info:   map[string]string{"productID": "A", "authCode": "x", "authTime": time.Now().Format(enums.ResTimeDayLayout)},
		serial: "B",
		decryptMap: map[string]string{
			"productID": "A",
			"authDays":  "1",
		},
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateFailed {
			t.Fatalf("state=%s", m.lastState)
		}
		// Verify auto-init logic
		if m.lastInfo == nil {
			t.Fatal("expected auto-init info update")
		}
		if m.lastInfo["productID"] != "B" {
			t.Fatalf("expected productID reset to B, got %s", m.lastInfo["productID"])
		}
		if m.lastInfo["authTime"] != "未授权" {
			t.Fatal("expected authTime reset to 未授权")
		}
	})
}

func Test_CheckSystemAuth_EmptyAuthCode(t *testing.T) {
	m := &mockAuth{
		info:   map[string]string{"productID": "A", "authCode": "", "authTime": time.Now().Format(enums.ResTimeDayLayout)},
		serial: "A",
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateFailed {
			t.Fatalf("state=%s", m.lastState)
		}
	})
}

func Test_CheckSystemAuth_DecryptError(t *testing.T) {
	m := &mockAuth{
		info:       map[string]string{"productID": "A", "authCode": "x", "authTime": time.Now().Format(enums.ResTimeDayLayout)},
		serial:     "A",
		decryptErr: errors.New("bad"),
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateFailed {
			t.Fatalf("state=%s", m.lastState)
		}
	})
}

func Test_CheckSystemAuth_DaysParseError(t *testing.T) {
	m := &mockAuth{
		info:   map[string]string{"productID": "A", "authCode": "x", "authTime": time.Now().Format(enums.ResTimeDayLayout)},
		serial: "A",
		decryptMap: map[string]string{
			"productID": "A",
			"authDays":  "abc",
		},
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateFailed {
			t.Fatalf("state=%s", m.lastState)
		}
	})
}

func Test_CheckSystemAuth_TimeParseError(t *testing.T) {
	m := &mockAuth{
		info:   map[string]string{"productID": "A", "authCode": "x", "authTime": "2025/01/01"},
		serial: "A",
		decryptMap: map[string]string{
			"productID": "A",
			"authDays":  "1",
		},
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateFailed {
			t.Fatalf("state=%s", m.lastState)
		}
	})
}

func Test_CheckSystemAuth_Expired(t *testing.T) {
	authTime := time.Now().AddDate(0, 0, -5).Format(enums.ResTimeDayLayout)
	m := &mockAuth{
		info:   map[string]string{"productID": "A", "authCode": "x", "authTime": authTime},
		serial: "A",
		decryptMap: map[string]string{
			"productID": "A",
			"authDays":  "1",
		},
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateFailed {
			t.Fatalf("state=%s", m.lastState)
		}
	})
}

func Test_CheckSystemAuth_Success(t *testing.T) {
	authTime := time.Now().Format(enums.ResTimeDayLayout)
	m := &mockAuth{
		info:   map[string]string{"productID": "A", "authCode": "x", "authTime": authTime},
		serial: "A",
		decryptMap: map[string]string{
			"productID": "A",
			"authDays":  "3",
		},
	}
	withMock(t, m, func() {
		CheckSystemAuth()
		if m.lastState != enums.ProductAuthStateSuccess {
			t.Fatalf("state=%s", m.lastState)
		}
		if m.lastInfo == nil {
			t.Fatalf("info=nil")
		}
		days := m.lastInfo["authDays"]
		if !strings.HasSuffix(days, "天") {
			t.Fatalf("authDays=%s", days)
		}
		left := m.lastInfo["leftDays"]
		if _, err := strconv.Atoi(left); err != nil {
			t.Fatalf("leftDays=%s", left)
		}
	})
}
