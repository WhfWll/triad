// Package enums
// @Author bcy2007  2025/8/16 15:36
package enums

import "strings"

type conn struct{}

var ConnEnum conn

const (
	ConnMethodSSH = iota + 1
)

func (c *conn) GetConnMethodEnums() interface{} {
	return []struct {
		Label string `json:"label"`
		Value int    `json:"value"`
	}{
		{Label: "ssh", Value: ConnMethodSSH},
	}
}

func (c *conn) AllConnMethodEnum() map[int]string {
	res := map[int]string{
		ConnMethodSSH: "ssh",
	}
	return res
}

func (c *conn) GetConnMethodEnum(connType int) string {
	enum := c.AllConnMethodEnum()
	if res, ok := enum[connType]; ok {
		return res
	}
	return ""
}

const (
	ConnStatusOK = iota + 1
	ConnStatusTimeout
	ConnStatusHandShakeFail
	ConnStatusRefuse
	ConnStatusOther
)

var targetStatusNameMap = map[int]string{
	ConnStatusOK:            "连接成功",
	ConnStatusTimeout:       "连接超时",
	ConnStatusHandShakeFail: "握手失败",
	ConnStatusRefuse:        "连接拒绝",
	ConnStatusOther:         "其他连接错误",
}

func (c *conn) GetConnStatusName(status int) string {
	if desc, ok := targetStatusNameMap[status]; ok {
		return desc
	}
	return "其他连接错误"
}

func (c *conn) CheckConnErrorInfo(err error) int {
	if strings.Contains(err.Error(), "i/o timeout") {
		return ConnStatusTimeout
	} else if strings.Contains(err.Error(), "handshake failed") {
		return ConnStatusHandShakeFail
	} else if strings.Contains(err.Error(), "connection refused") {
		return ConnStatusRefuse
	}
	return ConnStatusOther
}
