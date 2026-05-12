package middleware

import (
	"smart/tools/errno"

	"github.com/gin-gonic/gin"
	"gitlabee.4dogs.cn/common/server"
)

// RoleAuth 角色权限校验中间件
func RoleAuth(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			// 如果上下文中没有角色信息
			server.RespFail(c, errno.LoginAuthErr, "无法获取用户角色信息，禁止访问")
			c.Abort()
			return
		}

		userRole, ok := role.(int)
		if !ok {
			server.RespFail(c, errno.LoginAuthErr, "用户角色信息异常")
			c.Abort()
			return
		}

		// 2. 检查角色是否在允许列表中
		allowed := false
		for _, r := range roles {
			if userRole == r {
				allowed = true
				break
			}
		}

		if !allowed {
			server.RespFail(c, errno.BusinessAuthErr, "您没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}
