package typespec

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// FlexInt 兼容 JSON 数字或字符串（前端 el-input type=number 常为字符串）
type FlexInt int

func (f *FlexInt) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		*f = 0
		return nil
	}
	var n int
	if err := json.Unmarshal(b, &n); err == nil {
		*f = FlexInt(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		s = strings.TrimSpace(s)
		if s == "" {
			*f = 0
			return nil
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid integer: %s", s)
		}
		*f = FlexInt(v)
		return nil
	}
	return fmt.Errorf("invalid integer json: %s", string(b))
}

func (f FlexInt) Int() int {
	return int(f)
}
