package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// GroupUserPreselectionList 用户组用与预选列表
func GroupUserPreselectionList(c *gin.Context) {
	var (
		req typespec.UserGroupUserPreselectionReq
		res typespec.UserGroupUserPreselectionRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("GroupUserPreselectionList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroupUserList
	if err := app.GroupUserPreselectionList(ctx, &req, &res); err != nil {
		log.Error("UserGroupUpdateStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// GroupUserAlreadyList 组内被选择用户信息
func GroupUserAlreadyList(c *gin.Context) {
	var (
		req typespec.UserGroupUserAlreadyReq
		res typespec.UserGroupUserAlreadyRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("GroupUserAlreadyList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroupUserList
	if err := app.GroupUserAlreadyList(ctx, &req, &res); err != nil {
		log.Error("GroupUserAlreadyList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// GroupUserRelation 保存至用户组
func GroupUserRelation(c *gin.Context) {
	var (
		req typespec.UserGroupUserRelationReq
		res typespec.UserGroupUserRelationRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("GroupUserRelation parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("GroupUserRelation parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroupUserList
	if err := app.GroupUserRelation(ctx, uid.(int), &req, &res); err != nil {
		log.Error("GroupUserRelation parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "操作成功")
}
