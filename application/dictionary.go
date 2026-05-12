package application

import (
	"context"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
	"strings"
)

type ToolsDictionary struct {
}

// Enums 字典库 - 枚举
func (d *ToolsDictionary) Enums(res *typespec.DictionaryEnumRes) error {
	// 调用服务层获取字典枚举数据
	var dictionaryService services.Dictionary
	enumData := dictionaryService.GetDictionaryEnums()

	// 组织返回的数据
	res.Types = enumData["types"]

	serviceData := enumData["service"].(map[string]interface{})
	res.Service.WeakPass = serviceData["weakPass"]
	res.Service.WebPathScan = serviceData["webPathScan"]
	res.Service.SubdomainScan = serviceData["subdomainScan"]
	res.Service.Wifi = serviceData["wifi"]

	return nil
}

// List 字典库 - 列表页
func (d *ToolsDictionary) List(req *typespec.DictionaryListReq, res *typespec.DictionaryListRes) error {
	// 调用服务层获取字典列表
	var dictionaryService services.Dictionary

	dictionaryList, total, err := dictionaryService.GetDictionaryList(context.Background(), req.Page, req.Size, req.Types, req.Search)
	if err != nil {
		return err
	}

	// 转换数据格式
	var list []map[string]interface{}
	for _, item := range dictionaryList {
		dictItem := map[string]interface{}{
			"id":            item.ID,
			"sources":       item.Sources,
			"service":       item.Service,
			"serviceName":   enums.GetDictionaryService(item.Service),
			"name":          item.Name,
			"types":         item.Types,
			"typesName":     enums.GetDictionaryType(item.Types),
			"isDefault":     item.IsDefault,
			"isDefaultName": enums.GetDictionaryDefault(item.IsDefault),
			//"content":    item.Content,
			"createTime": item.CreateTime.Format("2006-01-02 15:04:05"),
			"updateTime": item.UpdateTime.Format("2006-01-02 15:04:05"),
		}
		list = append(list, dictItem)
	}

	res.List = list
	res.Total = total

	return nil
}

// Detail 字典库 - 详情页
func (d *ToolsDictionary) Detail(ctx context.Context, req *typespec.DictionaryDetailReq, res *typespec.DictionaryDetailRes) error {
	var dictSrv services.Dictionary
	dictModel, err := dictSrv.OpenDictionaryDetail(ctx, req.DictId)
	if err != nil {
		return err
	}

	//组织返回的数据
	res.Id = dictModel.ID
	res.IsDefault = dictModel.IsDefault
	res.Types = dictModel.Types
	res.Service = dictModel.Service
	res.Content = dictModel.Content
	res.Sources = dictModel.Sources
	res.Name = dictModel.Name
	res.CreateTime = dictModel.CreateTime.Format(enums.TimeLayout)
	res.UpdateTime = dictModel.UpdateTime.Format(enums.TimeLayout)

	return nil
}

//
//// Delete 字典库 - 删除
//func (d *ToolsDictionary) Delete(req *typespec.DictionaryDeleteReq) error {
//	//组织请求参数
//	param := map[string]interface{}{
//		"dictIds": req.DictIds,
//	}
//
//	//请求查询
//	decisionRes, err := httpclients.OpenDictionaryDelete(param)
//	if err != nil {
//		return err
//	}
//	if decisionRes.Code != 200 {
//		return errors.New(decisionRes.Msg)
//	}
//
//	return nil
//}

// Delete 字典库 - 删除
func (d *ToolsDictionary) Delete(ctx context.Context, req *typespec.DictionaryDeleteReq) error {
	dictIdArray := strings.Split(req.DictIds, ",")
	var dictService services.Dictionary
	return dictService.OpenDictionaryMultiDelete(ctx, dictIdArray)
}

// SetDefault 字典库 - 设置默认字典
func (d *ToolsDictionary) SetDefault(req *typespec.DictionarySetDefaultReq) error {
	var dictService services.Dictionary
	return dictService.SetDefaultDictionary(context.Background(), req.DictId, req.Service, req.Types)
}

// AddOrEdit 字典库 - 新增或编辑
func (d *ToolsDictionary) AddOrEdit(ctx context.Context, req *typespec.DictionaryAddOrEditReq) error {
	var dictSrv services.Dictionary
	if req.Id == 0 {
		return dictSrv.OpenDictionaryAdd(ctx, req.Types, req.Service, req.Sources, req.Name, req.Content)
	} else {
		return dictSrv.OpenDictionaryEdit(ctx, req.Id, req.Types, req.Service, req.Sources, req.Name, req.Content)
	}
}
