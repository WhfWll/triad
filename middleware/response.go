package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/redis"
)

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

// 记录登陆信息
func RecordLogin(c *gin.Context) {
	if c.Request.RequestURI == "/smart/user/login" {
		// 重写响应体
		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer

		// 继续走
		c.Next()

		// 获取响应体
		responseData := struct {
			Code int         `json:"code"`
			Data interface{} `json:"data"`
			Msg  string      `json:"msg"`
		}{}
		if err := json.Unmarshal(writer.b.Bytes(), &responseData); err == nil {
			username := c.DefaultPostForm("username", "未知的用户")
			if responseData.Code == 200 {
				responseData.Msg = "成功"
			}
			var logAuditService services.LogAudit //增加审计日志-登录日志
			_ = logAuditService.LogAuditAdd(c, enums.LogAuditTypeLogin, username+"登录："+responseData.Msg, username, utils.GetClientIp(c))
		}
		return
	}

	log.Infof("[RecordLogin] non-login request passing through: %s %s", c.Request.Method, c.Request.URL.Path)
	c.Next()
}

// 记录用户每个操作日志 RequestCn
// 只记录 POST/PUT/DELETE 等写操作，GET 查询操作不记录（避免产生过多噪音）
func RecordUserLog(c *gin.Context) {
	c.Next()

	// 只记录写操作（创建、修改、删除），跳过 GET 查询
	if c.Request.Method == "GET" {
		return
	}

	// 这几个API不记录日志
	whiteApi := make(map[string]struct{})
	whiteApi["/smart/user/login"] = struct{}{}        // 登陆接口 由独立记录登陆信息的中间件记录
	whiteApi["/smart/user/logincaptcha"] = struct{}{} // 验证码接口

	urlPath := c.Request.URL.Path
	if _, ok := whiteApi[urlPath]; !ok {
		// 获取操作名称 为空的不记录
		actionName := utils.RequestCn.GetCnName(urlPath)
		// 获取操作用户ID
		uid, existsUid := c.Get("uid")
		if actionName != "" && existsUid {
			uidInt, ok := uid.(int)
			if !ok {
				return
			}
			clientIP := utils.GetClientIp(c)
			go func(ctx context.Context, actionName string, uid int, clientIP string) {
				// 为避免每次查询用户名，这是使用缓存
				redisClient, err := redis.NewClient()
				if err != nil {
					return
				}
				redisUsername := redisClient.Get(ctx, "username_"+strconv.Itoa(uid)).Val()
				if redisUsername == "" {
					var srv services.User
					userData, err := srv.GetUserForId(ctx, uid)
					if err != nil || userData.Username == "" {
						// todo 缓存穿透，暂不考虑
						return
					}
					redisUsername = userData.Username
					// 缓存1小时
					redisClient.Set(ctx, "username_"+strconv.Itoa(uid), redisUsername, 3600*time.Second)
				}

				var logAuditService services.LogAudit
				_ = logAuditService.LogAuditAdd(ctx, enums.LogAuditTypeOperate, redisUsername+" 请求["+actionName+"]", redisUsername, clientIP)
			}(context.Background(), actionName, uidInt, clientIP)
		}
	}
}
