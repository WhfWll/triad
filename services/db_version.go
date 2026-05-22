package services

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"smart/tools/enums"
)

var semverRe = regexp.MustCompile(`(\d+\.\d+\.\d+(?:[-\w.]*)?)`)

// GetServerVersion 探测目标数据库版本字符串（用于 CVE 匹配与概览展示）
func (m *DBConnManager) GetServerVersion(ctx context.Context, conn *DBConnection) string {
	if conn == nil || conn.Config == nil {
		return ""
	}
	switch conn.Config.DBType {
	case enums.DBSupportTypeMySQL:
		return normalizeDBVersion(firstQueryCell(ctx, m, conn, "SELECT VERSION() AS version"))
	case enums.DBSupportTypePostgreSQL:
		return normalizeDBVersion(firstQueryCell(ctx, m, conn, "SHOW server_version"))
	case enums.DBSupportTypeMongoDB:
		return normalizeDBVersion(mongoVersionFromHTTP(ctx, m, conn))
	case enums.DBSupportTypeRedis:
		return normalizeDBVersion(redisVersionFromINFO(ctx, m, conn))
	case enums.DBSupportTypeCouchDB:
		return normalizeDBVersion(couchVersionFromHTTP(ctx, m, conn))
	default:
		return ""
	}
}

func firstQueryCell(ctx context.Context, m *DBConnManager, conn *DBConnection, query string) string {
	rows, err := m.ExecuteQuery(ctx, conn, query)
	if err != nil || len(rows) == 0 {
		return ""
	}
	for _, v := range rows[0] {
		if strings.TrimSpace(v) != "" && !strings.EqualFold(v, "NULL") {
			return v
		}
	}
	return ""
}

func redisVersionFromINFO(ctx context.Context, m *DBConnManager, conn *DBConnection) string {
	rows, err := m.ExecuteQuery(ctx, conn, "INFO server")
	if err != nil || len(rows) == 0 {
		return ""
	}
	text := rows[0]["result"]
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "redis_version:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "redis_version:"))
		}
	}
	return text
}

func mongoVersionFromHTTP(ctx context.Context, m *DBConnManager, conn *DBConnection) string {
	rows, err := m.ExecuteQuery(ctx, conn, "/buildInfo")
	if err != nil || len(rows) == 0 {
		rows, err = m.ExecuteQuery(ctx, conn, "/")
		if err != nil || len(rows) == 0 {
			return ""
		}
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(rows[0]["response"]), &payload); err != nil {
		return ""
	}
	for _, key := range []string{"version", "versionArray"} {
		if v, ok := payload[key]; ok {
			switch t := v.(type) {
			case string:
				return t
			case []interface{}:
				parts := make([]string, 0, len(t))
				for _, p := range t {
					parts = append(parts, fmtAny(p))
				}
				return strings.Join(parts, ".")
			}
		}
	}
	return ""
}

func couchVersionFromHTTP(ctx context.Context, m *DBConnManager, conn *DBConnection) string {
	rows, err := m.ExecuteQuery(ctx, conn, "/")
	if err != nil || len(rows) == 0 {
		return ""
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(rows[0]["response"]), &payload); err != nil {
		return ""
	}
	if v, ok := payload["version"].(string); ok {
		return v
	}
	return ""
}

func fmtAny(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}

func normalizeDBVersion(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if m := semverRe.FindStringSubmatch(raw); len(m) > 1 {
		v := m[1]
		// MySQL: 8.0.32-log -> 8.0.32
		if idx := strings.IndexAny(v, "-_"); idx > 0 {
			v = v[:idx]
		}
		return v
	}
	return raw
}
