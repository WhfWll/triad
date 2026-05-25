package services

import (
	"encoding/json"
	"strconv"
	"strings"
)

type httpProbePayload struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Host       string            `json:"host,omitempty"`
	Port       int               `json:"port,omitempty"`
	Error      string            `json:"error,omitempty"`
}

type mongoProbePayload struct {
	Command string                 `json:"command,omitempty"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Host    string                 `json:"host,omitempty"`
	Port    int                    `json:"port,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

func unwrapProbeResponse(actual string) string {
	actual = strings.TrimSpace(actual)
	if actual == "" {
		return ""
	}
	if strings.HasPrefix(actual, "response=") {
		return strings.TrimSpace(strings.TrimPrefix(actual, "response="))
	}
	if idx := strings.Index(actual, "response="); idx >= 0 {
		return strings.TrimSpace(actual[idx+len("response="):])
	}
	return actual
}

func parseHTTPProbePayload(actual string) (httpProbePayload, bool) {
	var payload httpProbePayload
	raw := unwrapProbeResponse(actual)
	if raw == "" {
		return payload, false
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return payload, false
	}
	return payload, true
}

func parseMongoProbePayload(actual string) (mongoProbePayload, bool) {
	var payload mongoProbePayload
	raw := unwrapProbeResponse(actual)
	if raw == "" {
		return payload, false
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return payload, false
	}
	return payload, true
}

func jsonMap(raw string) map[string]interface{} {
	out := map[string]interface{}{}
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

func nestedValue(m map[string]interface{}, path ...string) interface{} {
	var cur interface{} = m
	for _, key := range path {
		obj, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur, ok = obj[key]
		if !ok {
			return nil
		}
	}
	return cur
}

func nestedString(m map[string]interface{}, path ...string) string {
	v := nestedValue(m, path...)
	switch t := v.(type) {
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}

func nestedBool(m map[string]interface{}, path ...string) (bool, bool) {
	v := nestedValue(m, path...)
	b, ok := v.(bool)
	return b, ok
}

func nestedNumber(m map[string]interface{}, path ...string) (float64, bool) {
	v := nestedValue(m, path...)
	n, ok := v.(float64)
	return n, ok
}

func isUnauthorizedText(s string) bool {
	text := strings.ToLower(strings.TrimSpace(s))
	return strings.Contains(text, "unauthorized") ||
		strings.Contains(text, "not authorized") ||
		strings.Contains(text, "requires authentication") ||
		strings.Contains(text, "authentication failed") ||
		strings.Contains(text, "401 unauthorized") ||
		strings.Contains(text, "403 forbidden")
}

func isLoopbackHost(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	return h == "127.0.0.1" || h == "::1" || h == "localhost"
}
