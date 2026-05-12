package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

/** gin 框架拓展方法
 */

// GlobalOptions 全局枚举接口
func GlobalOptions(c *gin.Context) {
	var (
		res typespec.GlobalOptionsRes
	)

	ctx := server.NewContext(context.Background(), c)

	common := application.Common{}
	if err := common.GlobalOptions(ctx, &res); err != nil {
		log.Error("GlobalOptions parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}
