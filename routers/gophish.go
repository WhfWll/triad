package routers

import (
	"smart/api/rest/gofish"

	"github.com/gin-gonic/gin"
)

func GoPhishRouter(r *gin.RouterGroup) {
	// 模块前缀
	group := r.Group("/gophish")
	api := gofish.GoPhishAPI{}

	group.GET("/group/all", api.GetGroupAll)
	group.GET("/group/detail", api.GetGroupDetail)
	group.POST("/group/create", api.GroupCreate)
	group.POST("/group/update", api.GroupUpdate)
	group.POST("/group/delete", api.GroupDelete)

	// templates
	group.GET("/template/all", api.TemplateAll)
	group.GET("/template/detail", api.TemplateDetail)
	group.POST("/template/create", api.TemplateCreate)
	group.POST("/template/update", api.TemplateUpdate)
	group.POST("/template/delete", api.TemplateDelete)
	group.POST("/template/import_email", api.TemplateImportEmail)

	// landing pages
	group.GET("/page/all", api.PageAll)
	group.GET("/page/detail", api.PageDetail)
	group.POST("/page/create", api.PageCreate)
	group.POST("/page/update", api.PageUpdate)
	group.POST("/page/delete", api.PageDelete)
	group.POST("/page/import_site", api.PageImportSite)

	// profiles (smtp)
	group.GET("/profile/all", api.ProfileAll)
	group.GET("/profile/detail", api.ProfileDetail)
	group.POST("/profile/create", api.ProfileCreate)
	group.POST("/profile/update", api.ProfileUpdate)
	group.POST("/profile/delete", api.ProfileDelete)
	group.POST("/profile/send_test_email", api.ProfileSendTestEmail)

	// campaigns
	group.GET("/campaign/all", api.CampaignAll)
	group.GET("/campaign/detail", api.CampaignDetail)
	group.POST("/campaign/create", api.CampaignCreate)
	group.POST("/campaign/update", api.CampaignUpdate)
	group.POST("/campaign/delete", api.CampaignDelete)
	group.POST("/campaign/launch", api.CampaignLaunch)
	group.POST("/campaign/complete", api.CampaignComplete)
	group.GET("/campaign/result", api.CampaignResult)
}
