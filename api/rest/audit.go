package rest

import (
	"context"
	"fmt"

	"smart/services"
	"smart/tools/enums"

	"github.com/gin-gonic/gin"
)

func addOperateAudit(c *gin.Context, ctx context.Context, content string) {
	username, _ := c.Get("username")
	ip := c.ClientIP()
	var svc services.LogAudit
	_ = svc.LogAuditAdd(ctx, enums.LogAuditTypeOperate, content, fmt.Sprintf("%v", username), ip)
}

func addOperateAuditf(c *gin.Context, ctx context.Context, format string, args ...interface{}) {
	addOperateAudit(c, ctx, fmt.Sprintf(format, args...))
}
