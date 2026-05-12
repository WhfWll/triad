package services

import (
	"context"
	"encoding/json"
	"errors"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/tools/enums"
	"time"
)

// 场景管理 - 任务场景

type SceneTaskTemplate struct {
}

//// GetSceneEnum 获取场景枚举值
//func (a *SceneTaskTemplate) GetSceneEnum(ctx context.Context) (string, error) {
//	var mapSet mysqls.MapSet
//	mapValue, err := mapSet.GetsByObjKey(ctx, enums.ConfigSceneMapSetKey)
//	if err != nil {
//		return "", err
//	}
//	return mapValue.ObjValue, nil
//}

// 枚举信息 - 端口扫描 - TCP扫描方式
func (a *SceneTaskTemplate) PortScanTcpScanTypeEnum() []typespec.GlobalOptionsItemRes {
	var portScan enums.PortScanConfig
	return toolsSort(portScan.AllTcpScanTypeEnum())
}

// 枚举信息 - 端口扫描 - 端口范围
func (a *SceneTaskTemplate) PortScanPortRangeEnum() []typespec.GlobalOptionsItemRes {
	var portScan enums.PortScanConfig
	return toolsSort(portScan.AllPortScanTypeEnum())
}

// 枚举信息 - 端口扫描 - 端口值
func (a *SceneTaskTemplate) PortScanPortRangeValueEnum() []typespec.GlobalOptionsItemRes {
	var portScan enums.PortScanConfig
	return toolsSort(portScan.AllPortScanTypeValue())
}

// 枚举信息 - 端口扫描 - 超时时间
func (a *SceneTaskTemplate) PortScanTimeoutEnum() []typespec.GlobalOptionsItemRes {
	var portScan enums.PortScanConfig
	return toolsSort(portScan.AllPortScanTimeoutEnum())
}

// 枚举信息 - 端口扫描 - 并发数
func (a *SceneTaskTemplate) PortScanConcurrentEnum() []typespec.GlobalOptionsItemRes {
	var portScan enums.PortScanConfig
	return toolsSort(portScan.AllPortScanConcurrentEnum())
}

// 枚举信息 - 爬虫 - 爬取范围
func (a *SceneTaskTemplate) CrawlerScanRangeEnum() []typespec.GlobalOptionsItemRes {
	var crawler enums.WebCrawlerConfig
	return toolsSort(crawler.AllCrawlerScanRange())
}

// 枚举信息 - 爬虫 - 爬取深度
func (a *SceneTaskTemplate) CrawlerDeepEnum() []typespec.GlobalOptionsItemRes {
	var crawler enums.WebCrawlerConfig
	return toolsSort(crawler.AllCrawlerDeep())
}

// 枚举信息 - 爬虫 - 连接总数
func (a *SceneTaskTemplate) CrawlerMaxUrlEnum() []typespec.GlobalOptionsItemRes {
	var crawler enums.WebCrawlerConfig
	return toolsSort(crawler.AllCrawlerMaxUrl())
}

// 枚举信息 - 爬虫 - url去重
func (a *SceneTaskTemplate) CrawlerRepeatEnum() []typespec.GlobalOptionsItemRes {
	var crawler enums.WebCrawlerConfig
	return toolsSort(crawler.AllCrawlerRepeat())
}

// 枚举信息 - 爬虫 - 爬取速率
func (a *SceneTaskTemplate) CrawlerSpeedEnum() []typespec.GlobalOptionsItemRes {
	var crawler enums.WebCrawlerConfig
	return toolsSort(crawler.AllCrawlerSpeed())
}

// 枚举信息 - 爬虫 - 单链接超时
func (a *SceneTaskTemplate) CrawlerSingleTimeout() []typespec.GlobalOptionsItemRes {
	var crawler enums.WebCrawlerConfig
	return toolsSort(crawler.AllSingleTimeout())
}

// 枚举信息 - 爬虫 - 全局超时时间
func (a *SceneTaskTemplate) CrawlerFullTimeout() []typespec.GlobalOptionsItemRes {
	var crawler enums.WebCrawlerConfig
	return toolsSort(crawler.AllFullTimeoutEnum())
}

// 枚举信息 - 爬虫 - 敏感词
func (a *SceneTaskTemplate) CrawlerSensitive() string {
	return enums.Sensitive
}

// 枚举信息 - 爬虫 - 黑名单
func (a *SceneTaskTemplate) CrawlerWhiteList() string {
	return enums.WhiteList
}

// 枚举信息 - 爬虫 - 白名单
func (a *SceneTaskTemplate) CrawlerBlackList() string {
	return enums.BlackList
}

// 枚举信息 - Web路径爆破 - 猜测速率
func (a *SceneTaskTemplate) WebPathScanSpeedEnum() []typespec.GlobalOptionsItemRes {
	var webPathScan enums.WebPathScanConfig
	return toolsSort(webPathScan.AllWebPathScanSpeed())
}

// 枚举信息 - Web路径爆破 - 猜测时长
func (a *SceneTaskTemplate) WebPathScanTimeEnum() []typespec.GlobalOptionsItemRes {
	var webPathScan enums.WebPathScanConfig
	return toolsSort(webPathScan.AllWebPathScanTime())
}

// 枚举信息 - Web路径爆破 - 路径字典
func (a *SceneTaskTemplate) WebPathScanDictEnum() ([]typespec.GlobalOptionsItemRes, error) {
	var req httpclients.OpenDictionaryGetByTypeReq
	req.Type = append(req.Type, enums.DictionaryTypeWebPathScan)
	data, err := httpclients.GetDecisionDictionaryByType(req)
	if err != nil {
		return nil, err
	}
	recoverData := make([]typespec.GlobalOptionsItemRes, 0)
	for _, item := range data.Data.List {
		recoverData = append(recoverData, typespec.GlobalOptionsItemRes{
			Value: item.Id,
			Label: item.Name,
		})
	}
	return recoverData, nil
}

// 枚举信息 - 子域名收集
func (a *SceneTaskTemplate) SubdomainDictCollectEnum() ([]typespec.GlobalOptionsItemRes, error) {
	// OpenDictionaryGetByType 依据类型获取每个服务下默认的配置 types = 类型 and is_default=默认 group by services
	// 获取列表
	var dicModel mysqls.Dictionary
	ctx := context.Background()
	typeInt := []int{enums.DictionaryTypeSubdomainScan}
	dictList := dicModel.GetByTypeAndIsDefault(ctx, typeInt, 0)
	recoverData := make([]typespec.GlobalOptionsItemRes, 0)
	for _, item := range dictList {
		recoverData = append(recoverData, typespec.GlobalOptionsItemRes{
			Value: item.ID,
			Label: item.Name,
		})
	}
	return recoverData, nil
}

// web弱口令 - 服务
func (a *SceneTaskTemplate) WeakPassServiceEnum() ([]typespec.GlobalOptionsItemRes, error) {
	var req httpclients.OpenDictionaryGetServiceEnumByTypeReq
	req.Type = enums.DictionaryTypeUser
	tempDict := make([]typespec.GlobalOptionsItemRes, 0)
	dataList := a.GetServiceEnumByType(req.Type)
	for _, item := range dataList {
		item.IsDefault = true
		tempDict = append(tempDict, item)
	}
	return tempDict, nil
}

// GetServiceEnumByType 通过类型获取服务枚举
func (a *SceneTaskTemplate) GetServiceEnumByType(types int) []typespec.GlobalOptionsItemRes {
	var data map[int]string
	switch types {
	case enums.DictionaryTypeUser, enums.DictionaryTypePassword:
		data = enums.AllDictionaryServiceWeakPassEnum(false)
	default:
		return nil
	}
	return toolsSort(data)
}

// web弱口令 - 公共字典
func (a *SceneTaskTemplate) WeakPassCommonDictEnum() ([]typespec.GlobalOptionsItemRes, []typespec.GlobalOptionsItemRes, error) {
	var req httpclients.OpenDictionaryGetByTypeAndServiceReq
	req.Type = append(req.Type, enums.DictionaryTypeUser)
	req.Type = append(req.Type, enums.DictionaryTypePassword)
	req.Service = enums.DictionaryServiceWeakPassCommon

	ctx := context.Background()
	dictList := a.OpenDictionaryGetByTypeAndService(ctx, req.Type, req.Service)
	// 服务数据
	recoverUserData := make([]typespec.GlobalOptionsItemRes, 0)
	recoverPassData := make([]typespec.GlobalOptionsItemRes, 0)
	for _, item := range dictList {
		isDefault := false
		if item.IsDefault == enums.DictionaryDefaultYes {
			isDefault = true
		}
		switch item.Types {
		case enums.DictionaryTypeUser: // 公共 - 用户字典
			recoverUserData = append(recoverUserData, typespec.GlobalOptionsItemRes{
				Value:     item.ID,
				Label:     item.Name,
				IsDefault: isDefault,
			})
		case enums.DictionaryTypePassword: // 公共 - 用户字典
			recoverPassData = append(recoverPassData, typespec.GlobalOptionsItemRes{
				Value:     item.ID,
				Label:     item.Name,
				IsDefault: isDefault,
			})
		}
	}
	return recoverUserData, recoverPassData, nil
}

// OpenDictionaryGetByTypeAndService 依据类型与服务获取所有数据
func (a *SceneTaskTemplate) OpenDictionaryGetByTypeAndService(ctx context.Context, typeInt []int, service int) []mysqls.Dictionary {
	// 获取列表
	var dicModel mysqls.Dictionary
	return dicModel.GetByTypeAndService(ctx, typeInt, service)
}

// web弱口令 - 字典类型
func (a *SceneTaskTemplate) WeakPassDictType() []typespec.GlobalOptionsItemRes {
	var weakPass enums.WeakPassConfig
	return toolsSort(weakPass.AllWeakPassDictType())
}

// web弱口令 - 猜测次数
func (a *SceneTaskTemplate) WeakPassGuessNumber() []typespec.GlobalOptionsItemRes {
	var weakPass enums.WeakPassConfig
	return toolsSort(weakPass.AllWeakPassGuessNumber())
}

// web弱口令 - 猜测时间
func (a *SceneTaskTemplate) WeakPassGuessTime() []typespec.GlobalOptionsItemRes {
	var weakPass enums.WeakPassConfig
	return toolsSort(weakPass.AllWeakPassGuessTime())
}

// web弱口令 - 猜测速率
func (a *SceneTaskTemplate) WeakPassGuessRate() []typespec.GlobalOptionsItemRes {
	var weakPass enums.WeakPassConfig
	return toolsSort(weakPass.AllWeakPassGuessRate())
}

// 场景 - 所有
func (a *SceneTaskTemplate) All(ctx context.Context, estate string) ([]mysqls.TaskTemplate, int64, error) {
	var taskTemplateModel mysqls.TaskTemplate
	return taskTemplateModel.AllTaskTemplate(ctx, "*", estate)
}

// 场景 - 分页
func (a *SceneTaskTemplate) Page(ctx context.Context, page, size int, templateName string, source []int, estate string) ([]mysqls.TaskTemplate, int64, error) {
	var taskTemplateModel mysqls.TaskTemplate
	return taskTemplateModel.GetTaskTemplateList(ctx, page, size, templateName, source, estate)
}

// 场景 - 根据场景名获取单个
func (a *SceneTaskTemplate) GetByName(ctx context.Context, templateName string) (mysqls.TaskTemplate, error) {
	var taskTemplateModel mysqls.TaskTemplate
	return taskTemplateModel.GetByName(ctx, templateName)
}

// 场景 - id获取单个
func (a *SceneTaskTemplate) GetByTemplateId(ctx context.Context, templateId int, estate string) (mysqls.TaskTemplate, error) {
	var taskTemplateModel mysqls.TaskTemplate
	return taskTemplateModel.GetTaskTemplate(ctx, templateId, estate)
}

// 场景 - 创建
func (a *SceneTaskTemplate) Create(ctx context.Context, config enums.ConfigJson, templateName, describe string, userId int) (int, error) {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	//创建场景模版
	var templateModel mysqls.TaskTemplate
	templateModel.Estate = enums.TaskTemplateStatusSuccess
	templateModel.TemplateName = templateName
	templateModel.Describe = describe
	templateModel.IsDefault = enums.TaskTemplateIsDefaultN
	templateModel.UserID = userId
	templateModel.Source = enums.TaskTemplateSourceUser
	templateModel.CreateTime = time.Now()
	templateModel.UpdateTime = time.Now()
	err := templateModel.AddTaskTemplate(dCtx)
	if err != nil {
		return 0, err
	}

	//创建配置
	// struct转map
	configByte, err := json.Marshal(config)
	if err != nil {
		return 0, err
	}
	configMap := make(map[string]interface{})
	if err := json.Unmarshal(configByte, &configMap); err != nil {
		return 0, err
	}
	configDataList := make([]mysqls.TaskConfiguration, 0)
	for k, v := range configMap {
		if itemConfigByte, err := json.Marshal(v); err == nil {
			configDataList = append(configDataList, mysqls.TaskConfiguration{
				ObjId:      templateModel.ID,
				ConfigKey:  k,
				ConfigJson: string(itemConfigByte),
				UserId:     userId,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
		}
	}
	if len(configDataList) > 0 {
		var confModel mysqls.TaskConfiguration
		err = confModel.AddAll(dCtx, configDataList)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit().Error; err != nil { //提交事务
		return 0, err
	}

	return templateModel.ID, nil
}

// 场景 - 更新
func (a *SceneTaskTemplate) Update(ctx context.Context, templateId int, config enums.ConfigJson, templateName, describe string, userId int) (int, error) {
	var templateModel mysqls.TaskTemplate
	templateData, err := templateModel.GetTaskTemplate(ctx, templateId, enums.TaskTemplateStatusSuccess)
	if err != nil {
		return 0, err
	}
	if templateData.ID == 0 {
		return 0, errors.New("未知的场景")
	}

	// 更新配置
	// struct转map
	configByte, err := json.Marshal(config)
	if err != nil {
		return 0, err
	}
	configMap := make(map[string]interface{})
	if err := json.Unmarshal(configByte, &configMap); err != nil {
		return 0, err
	}

	// 获取所有配置
	var confModel mysqls.TaskConfiguration
	alreadyConfig := confModel.AllByObjId(ctx, templateId, "")

	// 开启事物
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	// 更新场景模版
	templateData.TemplateName = templateName
	templateData.Describe = describe
	templateData.UserID = userId
	err = templateData.UpdateTaskTemplate(dCtx)
	if err != nil {
		return 0, err
	}

	// 更新已存在的配置 并删除不存在的配置
	for _, item := range alreadyConfig {
		if data, ok := configMap[item.ConfigKey]; ok {
			// 数据存在，需要更新
			itemConfigByte, err := json.Marshal(data)
			if err != nil {
				return 0, err
			}

			item.ConfigJson = string(itemConfigByte)
			item.UpdateTime = time.Now()
			err = item.UpdateTaskConfiguration(dCtx)
			if err != nil {
				return 0, err
			}

			delete(configMap, item.ConfigKey)
		} else {
			// 数据不存在，直接删除
			err = item.DeleteTaskConfiguration(dCtx)
			if err != nil {
				return 0, err
			}
		}
	}
	// 添加不存在的配置
	configDataList := make([]mysqls.TaskConfiguration, 0)
	for k, v := range configMap {
		itemConfigByte, err := json.Marshal(v)
		if err != nil {
			return 0, err
		}
		configDataList = append(configDataList, mysqls.TaskConfiguration{
			ObjId:      templateData.ID,
			ConfigKey:  k,
			ConfigJson: string(itemConfigByte),
			UserId:     userId,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
	}
	if len(configDataList) > 0 {
		if err = confModel.AddAll(dCtx, configDataList); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit().Error; err != nil { //提交事务
		return 0, err
	}

	return templateId, nil
}

// 场景详情
func (a *SceneTaskTemplate) Detail(ctx context.Context, templateId int, estate string) (templateData mysqls.TaskTemplate, configStruct enums.ConfigJson, err error) {
	// 场景
	var templateModel mysqls.TaskTemplate
	templateData, err = templateModel.GetTaskTemplate(ctx, templateId, estate)
	if err != nil {
		return
	}
	if templateData.ID == 0 {
		err = errors.New("任务场景模版不存在")
		return
	}
	// 配置
	var configModel mysqls.TaskConfiguration
	configData := configModel.AllByObjId(ctx, templateData.ID, "")
	if err != nil {
		return
	}
	configMap := make(map[string]string)
	for _, item := range configData {
		configMap[item.ConfigKey] = item.ConfigJson
	}
	if err = configStruct.Decode(configMap); err != nil {
		return mysqls.TaskTemplate{}, enums.ConfigJson{}, err
	}
	return
}

// 场景 - 设置默认
func (a *SceneTaskTemplate) SetDefault(ctx context.Context, templateId int) error {
	var taskTemplateModel mysqls.TaskTemplate
	//开启事务，先取消默认场景，后设置默认场景
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	err := taskTemplateModel.CancelDefault(dCtx)
	if err != nil {
		return err
	}
	err = taskTemplateModel.SetDefault(dCtx, templateId)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务
		return err
	}
	return nil
}

// 场景 - 删除
func (a *SceneTaskTemplate) Del(ctx context.Context, templateId []int) error {
	var taskTemplateModel mysqls.TaskTemplate
	return taskTemplateModel.UpdateStatusByIds(ctx, templateId, enums.TaskTemplateStatusDel)
}

// 场景下配置到漏洞信息
func (a *SceneTaskTemplate) GetTemplateVulIds(ctx context.Context, templateId int) (vulIds []int) {
	var taskConfiguration mysqls.TaskConfiguration
	configs := taskConfiguration.AllByObjId(ctx, templateId, enums.ConfigJsonVulIdsKey)
	if len(configs) > 0 {
		json.Unmarshal([]byte(configs[0].ConfigJson), &vulIds)
	}
	return
}

// GetTaskTemplateById 根据id查询模板信息
func (a *SceneTaskTemplate) GetTaskTemplateById(ctx context.Context, templateId int) (mysqls.TaskTemplate, error) {
	var taskTemplateModels mysqls.TaskTemplate
	return taskTemplateModels.GetTaskTemplateById(ctx, templateId)
}

// GetTaskConfigById 根据id查询模板信息
func (a *SceneTaskTemplate) GetTaskConfigById(ctx context.Context, templateId int) []mysqls.TaskConfiguration {
	var taskConfigModels mysqls.TaskConfiguration
	return taskConfigModels.AllByObjId(ctx, templateId, "")
}

// GetTaskSceneCount 获取任务场景数量
func (a *SceneTaskTemplate) GetTaskSceneCount(ctx context.Context) int {
	var taskSceneModel mysqls.TaskTemplate
	return int(taskSceneModel.GetTaskSceneCount(ctx, enums.TaskTemplateStatusSuccess))
}

// 枚举信息 - 存活探测 - 存活探测枚举
func (a *SceneTaskTemplate) AliveProbeTypeEnum() []typespec.GlobalOptionsItemRes {
	var aliveProbe enums.AliveProbeConfig
	return toolsSort(aliveProbe.AllAliveProbeTypeEnum())
}

// 枚举信息 - 存活探测 - 端口范围枚举
func (a *SceneTaskTemplate) AliveProbePortRangeEnum() []typespec.GlobalOptionsItemRes {
	var aliveProbe enums.AliveProbeConfig
	return toolsSort(aliveProbe.AllAliveProbePortRangeEnum())
}
