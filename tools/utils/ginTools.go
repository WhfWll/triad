package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// GetClientIp 获取用户请求的ip地址
func GetClientIp(c *gin.Context) string {
	ip := c.Request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = c.Request.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	}
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}
