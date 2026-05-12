package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
)

// UserGroupList 用户组列表
func UserGroupList(c *gin.Context) {
	var (
		req typespec.UserGroupListReq
		res typespec.UserGroupListRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserGroupList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroup
	if err := app.GetUserGroupList(ctx, &req, &res); err != nil {
		log.Error("UserGroupList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// UserGroupCreate 新建用户组
func UserGroupCreate(c *gin.Context) {
	var (
		req typespec.UserGroupCreateReq
		res typespec.UserGroupCreateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserGroupCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("UpdateUserExp parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroup
	if err := app.UserGroupCreate(ctx, uid.(int), &req, &res); err != nil {
		log.Error("UserGroupCreate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "操作成功")
}

// UserGroupUpdate 用户组 - 编辑
func UserGroupUpdate(c *gin.Context) {
	var (
		req typespec.UserGroupUpdateReq
		res typespec.UserGroupUpdateRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserGroupUpdate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("UserGroupUpdate parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroup
	if err := app.UserGroupEdit(ctx, uid.(int), &req, &res); err != nil {
		log.Error("UserGroupUpdate parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "操作成功")
}

// GroupSelect 用户组-层级选择
func GroupSelect(c *gin.Context) {
	var res typespec.UserGroupSelectRes
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroup
	if err := app.UserGroupSelect(ctx, &res); err != nil {
		log.Error("GroupSelect parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, res)
}

// UserGroupUpdateStatus 用户组 - 修改状态
func UserGroupUpdateStatus(c *gin.Context) {
	var (
		req typespec.UserGroupUpdateStatusReq
		res typespec.UserGroupUpdateStatusRes
	)
	if err := c.ShouldBind(&req); err != nil {
		log.Error("UserGroupUpdateStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}
	uid, ok := c.Get("uid")
	if !ok {
		log.Error("UserGroupUpdateStatus parameter error,uid is false：")
		server.RespFail(c, 4000, "参数错误,错误原因：token get uid is false")
		return
	}
	ctx := server.NewContext(context.Background(), c)
	var app application.UserGroup
	if err := app.UserGroupUpdateStatus(ctx, uid.(int), &req, &res); err != nil {
		log.Error("UserGroupUpdateStatus parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}
	server.RespSuccess(c, "操作成功")
}
