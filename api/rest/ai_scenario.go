package rest

import (
	"context"
	"smart/api/typespec"
	"smart/application"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
)

// AiScenarioList AI应用场景列表
func AiScenarioList(c *gin.Context) {
	var (
		req typespec.AiScenarioListReq
		app application.AiScenario
	)
	ctx := server.NewContext(context.Background(), c)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("AiScenarioList parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	list, err := app.AiScenarioList(ctx)
	if err != nil {
		log.Error("AiScenarioList error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	res := typespec.AiScenarioListRes{
		List: list,
	}

	server.RespSuccess(c, res)
}

// AiScenarioConfig AI应用场景配置
func AiScenarioConfig(c *gin.Context) {
	var (
		req typespec.AiScenarioConfigReq
		app application.AiScenario
	)
	ctx := server.NewContext(context.Background(), c)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("AiScenarioConfig parameter error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	err := app.AiScenarioConfig(ctx, req.ID, req.Description, req.IsEnabled, req.LlmModelID, req.Prompt)
	if err != nil {
		log.Error("AiScenarioConfig error, the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	res := typespec.AiScenarioConfigRes{}
	server.RespSuccess(c, res)
}