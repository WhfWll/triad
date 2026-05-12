package services

import (
	"context"
	"errors"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

type Dictionary struct {
}

// GetDictionaryCount 获取字典总数
func (d *Dictionary) GetDictionaryCount(ctx context.Context) (int64, error) {
	var dictionaryModel mysqls.Dictionary
	return dictionaryModel.Count(ctx)
}

// GetDictionaryList 获取字典列表
func (d *Dictionary) GetDictionaryList(ctx context.Context, page, size, types int, search string) ([]mysqls.Dictionary, int64, error) {
	var dictionaryModel mysqls.Dictionary
	return dictionaryModel.DictionaryList(ctx, page, size, types, search)
}

// GetDictionaryDetail 获取字典详情
func (d *Dictionary) GetDictionaryDetail(ctx context.Context, dictId string) (mysqls.Dictionary, error) {
	var dictionaryModel mysqls.Dictionary
	return dictionaryModel.GetDictById(ctx, dictId)
}

// DeleteDictionaries 删除字典
func (d *Dictionary) DeleteDictionaries(ctx context.Context, ids []string) error {
	var dictionaryModel mysqls.Dictionary
	return dictionaryModel.MultiDeleteDictionary(ctx, ids)
}

// SetDefaultDictionary 设置默认字典
func (d *Dictionary) SetDefaultDictionary(ctx context.Context, dictId, service, types int) error {
	var dictionaryModel mysqls.Dictionary

	// 先取消同类型和同服务的默认字典
	if err := dictionaryModel.CancelDefaultDictionary(ctx, types, service); err != nil {
		return err
	}

	// 设置新的默认字典
	return dictionaryModel.SetDefaultDictionary(ctx, dictId)
}

// AddOrEditDictionary 新增或编辑字典
func (d *Dictionary) AddOrEditDictionary(ctx context.Context, id int, sources, types, service int, name, content string) error {
	var dictionaryModel mysqls.Dictionary

	if id == 0 {
		// 新增
		dictionaryModel.Sources = sources
		dictionaryModel.Types = types
		dictionaryModel.Service = service
		dictionaryModel.Name = name
		dictionaryModel.Content = content
		dictionaryModel.IsDefault = enums.DictionaryDefaultNo
		return dictionaryModel.AddDictionary(ctx)
	} else {
		// 编辑
		param := map[string]interface{}{
			"id":      id,
			"sources": sources,
			"types":   types,
			"service": service,
			"name":    name,
			"content": content,
		}
		return dictionaryModel.UpdateDictionary(ctx, param)
	}
}

// GetDictionaryEnums 获取字典枚举
func (d *Dictionary) GetDictionaryEnums() map[string]interface{} {
	return map[string]interface{}{
		"types": []map[string]interface{}{
			{"label": "用户字典", "value": 1},
			{"label": "密码字典", "value": 2},
			{"label": "web路径爆破字典", "value": 4},
			{"label": "子域名爆破字典", "value": 5},
		},
		"service": map[string]interface{}{
			"weakPass": []map[string]interface{}{
				{"label": "通用", "value": 1},
				{"label": "ssh", "value": 2},
				{"label": "ftp", "value": 3},
				{"label": "memcached", "value": 4},
				{"label": "mongodb", "value": 5},
				{"label": "mssql", "value": 6},
				{"label": "mysql", "value": 7},
				{"label": "postgres", "value": 8},
				{"label": "rdp", "value": 9},
				{"label": "redis", "value": 10},
				{"label": "smb", "value": 11},
				{"label": "telnet", "value": 12},
				{"label": "tomcat", "value": 13},
				{"label": "vnc", "value": 14},
				{"label": "snmpv2", "value": 15},
				{"label": "snmpv3_md5", "value": 16},
				{"label": "snmpv3_sha", "value": 17},
				{"label": "snmpv3_sha_224", "value": 18},
				{"label": "snmpv3_sha_256", "value": 19},
				{"label": "snmpv3_sha_384", "value": 20},
				{"label": "snmpv3_sha_512", "value": 21},
				{"label": "HTTP", "value": 22},
			},
			"wifi": nil,
			"webPathScan": []map[string]interface{}{
				{"label": "通用", "value": 101},
				{"label": "php", "value": 102},
				{"label": "asp", "value": 103},
				{"label": "aspx", "value": 104},
				{"label": "jsp", "value": 105},
			},
			"subdomainScan": []map[string]interface{}{
				{"label": "通用", "value": 201},
			},
		},
	}
}

// OpenDictionaryDetail 根据id获取字典的详情数据
func (d *Dictionary) OpenDictionaryDetail(ctx context.Context, dictId int) (mysqls.Dictionary, error) {
	var dictModel mysqls.Dictionary
	return dictModel.DictionaryRecord(ctx, dictId)
}

// OpenDictionaryAdd 添加一条字典数据
func (d *Dictionary) OpenDictionaryAdd(ctx context.Context, types, service, sources int, name, content string) error {
	var dictModel = mysqls.Dictionary{
		Name:       name,
		Types:      types,
		Service:    service,
		Sources:    sources,
		Content:    content,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		IsDefault:  enums.DictionaryDefaultYes,
	}
	//根据types、name和service查询字典，判断字典名称是否已存在
	dict, err := dictModel.DictionaryByTypesAndServiceAndName(ctx, types, service, name)
	if err != nil {
		return err
	}
	if dict.ID != 0 {
		return errors.New("字典名称已存在")
	}
	//根据types、is_default和service查询字典，判断默认字典是否已存在
	dict, err = dictModel.DictionaryByTypesAndServiceAndIsDefault(ctx, types, service)
	if err != nil {
		return err
	}
	if dict.ID != 0 {
		dictModel.IsDefault = enums.DictionaryDefaultNo
	}

	return dictModel.AddDictionary(ctx)
}

// OpenDictionaryEdit 根据id更新一条字典数据
func (d *Dictionary) OpenDictionaryEdit(ctx context.Context, dictId, types, service, sources int, name, content string) error {
	var dictModel mysqls.Dictionary
	var param = map[string]interface{}{
		"id":          dictId,
		"name":        name,
		"sources":     sources,
		"content":     content,
		"update_time": time.Now(),
	}
	//根据types和service查询字典，获取字典列表
	dict, err := dictModel.DictionaryByTypesAndServiceAndName(ctx, types, service, name)
	if err != nil {
		return err
	}
	if dict.ID != 0 && dict.ID != dictId {
		return errors.New("字典名称已存在")
	}

	return dictModel.UpdateDictionary(ctx, param)
}

// OpenDictionaryMultiDelete 批量删除字典
func (d *Dictionary) OpenDictionaryMultiDelete(ctx context.Context, dictIds any) error {
	var dictModel mysqls.Dictionary
	//判定是否有不可删除字典
	dictionary, err := dictModel.DictionaryByIdList(ctx, dictIds)
	if err != nil {
		return err
	}
	if dictionary.ID != 0 {
		return errors.New("系统字典或默认字典不能删除")
	}
	//删除字典
	return dictModel.MultiDeleteDictionary(ctx, dictIds)
}
