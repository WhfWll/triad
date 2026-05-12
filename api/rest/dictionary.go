package rest

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/server"
	"smart/api/typespec"
	"smart/application"
	"strings"
)

// ToolsDictionaryEnum 字典库 - 枚举
func ToolsDictionaryEnum(c *gin.Context) {
	var res typespec.DictionaryEnumRes
	var app application.ToolsDictionary
	if err := app.Enums(&res); err != nil {
		log.Error("ToolsDictionaryEnum parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// ToolsDictionaryList 字典库 - 列表页
func ToolsDictionaryList(c *gin.Context) {
	var (
		req typespec.DictionaryListReq
		res typespec.DictionaryListRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("ToolsDictionaryList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.ToolsDictionary
	if err := app.List(&req, &res); err != nil {
		log.Error("ToolsDictionaryList parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// ToolsDictionaryDetail 字典库 - 详情页
func ToolsDictionaryDetail(c *gin.Context) {
	var (
		req typespec.DictionaryDetailReq
		res typespec.DictionaryDetailRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("ToolsDictionaryDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.ToolsDictionary
	if err := app.Detail(c, &req, &res); err != nil {
		log.Error("ToolsDictionaryDetail parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// ToolsDictionaryDelete 字典库 - 删除
func ToolsDictionaryDelete(c *gin.Context) {
	var (
		req typespec.DictionaryDeleteReq
		res typespec.DictionaryDeleteRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("ToolsDictionaryDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.ToolsDictionary
	if err := app.Delete(c, &req); err != nil {
		log.Error("ToolsDictionaryDelete parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// ToolsDictionarySetDefault 字典库 - 设置默认字典
func ToolsDictionarySetDefault(c *gin.Context) {
	var (
		req typespec.DictionarySetDefaultReq
		res typespec.DictionarySetDefaultRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("ToolsDictionarySetDefault parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	var app application.ToolsDictionary
	if err := app.SetDefault(&req); err != nil {
		log.Error("ToolsDictionarySetDefault parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}

// ToolsDictionaryAddOrEdit 字典库 - 增加或编辑
func ToolsDictionaryAddOrEdit(c *gin.Context) {
	var (
		req typespec.DictionaryAddOrEditReq
		res typespec.DictionaryAddOrEditRes
	)

	if err := c.ShouldBind(&req); err != nil {
		log.Error("ToolsDictionaryAddOrEdit parameter error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, "参数错误,错误原因："+err.Error())
		return
	}

	//处理file
	fileHeader := req.File
	if fileHeader != nil {
		if fileHeader.Size > 1024*1024 {
			log.Error("ToolsDictionaryAddOrEdit request error,the fail reason is：" + "文件大小不能超过1M")
			server.RespFail(c, 4000, "文件大小不能超过1M")
			return
		}
		if !strings.HasSuffix(fileHeader.Filename, ".txt") {
			log.Error("ToolsDictionaryAddOrEdit request error,the fail reason is：" + "上传文件格式为txt文件")
			server.RespFail(c, 4000, "上传文件格式为txt文件")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			log.Error("AddOrEdit file open error,the fail reason is：" + err.Error())
			server.RespFail(c, 4000, err.Error())
			return
		}
		defer file.Close()

		dataByte := make([]byte, 0)
		_, err = file.Read(dataByte)
		if err != nil {
			log.Error("AddOrEdit file read error,the fail reason is：" + err.Error())
			server.RespFail(c, 4000, err.Error())
			return
		}
		data := string(dataByte)

		//content为空，file不为空时，把file中的内容赋值给content
		if req.Content == "" && data != "" {
			req.Content = data
		}
	}

	var app application.ToolsDictionary
	if err := app.AddOrEdit(c, &req); err != nil {
		log.Error("ToolsDictionaryAddOrEdit request error,the fail reason is：" + err.Error())
		server.RespFail(c, 4000, err.Error())
		return
	}

	server.RespSuccess(c, res)
}
