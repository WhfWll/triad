package network

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const (
	CRLF = "\r\n"
)

var (
	TAIL = []byte{0, 0, 0xff, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff}
)

// parseRequestLine parses "GET /foo HTTP/1.1" into its three parts.
func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.LastIndex(line[s1+1:], " ")
	if s1 < 0 {
		return
	}

	var httpVersion = "HTTP/1.1"
	if s2 < 0 {
		return line[:s1], line[s1+1:], httpVersion, true
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

func TrimLeftHTTPPacket(raw []byte) []byte {
	return bytes.TrimLeft(raw, "\t \n\v\f\n\b\r")
}

func TrimSpaceHTTPPacket(raw []byte) []byte {
	return bytes.Trim(raw, "\t \n\v\f\n\b\r")
}

func SplitHTTPHeadersAndBodyFromPacketEx(raw []byte, mf func(method string, requestUri string, proto string) error, hook ...func(line string)) (string, []byte) {
	raw = TrimLeftHTTPPacket(raw)
	reader := bufio.NewReader(bytes.NewBuffer(raw))
	var err error
	firstLineBytes, err := BufioReadLine(reader)
	if err != nil {
		return "", nil
	}
	firstLineBytes = TrimSpaceHTTPPacket(firstLineBytes)

	var headers []string
	headers = append(headers, string(firstLineBytes))
	method, requestURI, proto, _ := parseRequestLine(string(firstLineBytes))
	if mf != nil {
		err := mf(method, requestURI, proto)
		if err != nil {
			return "", nil
		}
	}

	for {
		//lineBytes, _, err := reader.ReadLine()
		lineBytes, err := BufioReadLine(reader)
		if err != nil && err != io.EOF {
			break
		}
		if bytes.TrimSpace(lineBytes) == nil {
			break
		}

		for _, h := range hook {
			h(string(lineBytes))
		}

		headers = append(headers, string(lineBytes))
	}
	headersRaw := strings.Join(headers, CRLF) + CRLF + CRLF
	bodyRaw, _ := ioutil.ReadAll(reader)
	if bodyRaw == nil {
		return headersRaw, nil
	}

	if bytes.HasSuffix(bodyRaw, []byte(CRLF+CRLF)) {
		bodyRaw = bodyRaw[:len(bodyRaw)-4]
	}

	// 单独修复请求中的问题
	if !strings.HasPrefix(headersRaw, "HTTP/") {
		if bytes.HasSuffix(bodyRaw, []byte("\n\n")) {
			bodyRaw = bodyRaw[:len(bodyRaw)-2]
		}
	}

	return headersRaw, bodyRaw
}

func SplitHTTPHeadersAndBodyFromPacket(raw []byte, hook ...func(line string)) (string, []byte) {
	return SplitHTTPHeadersAndBodyFromPacketEx(raw, nil, hook...)
}

func BufioReadLine(reader *bufio.Reader) ([]byte, error) {
	if reader == nil {
		return nil, fmt.Errorf("empty reader(bufio)")
	}

	var buf bytes.Buffer
	for {
		tmp, isPrefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		buf.Write(tmp)
		if !isPrefix {
			return buf.Bytes(), nil
		}
	}
}
