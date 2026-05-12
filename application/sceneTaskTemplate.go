package application

import (
	"context"
	"errors"
	"fmt"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/data"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
)

type SceneTaskTemplate struct {
}

// 枚举
func (a *SceneTaskTemplate) SceneEnums(res *typespec.SceneEnumsRes) (err error) {
	var service services.SceneTaskTemplate
	// 主机探活
	res.AliveProbe.AliveProbeType = service.AliveProbeTypeEnum()
	res.AliveProbe.AliveProbePortRange = service.AliveProbePortRangeEnum()

	// 端口扫描
	res.PortScan.TcpScanType = service.PortScanTcpScanTypeEnum()
	res.PortScan.PortRange = service.PortScanPortRangeEnum()
	res.PortScan.PortRangeValue = service.PortScanPortRangeValueEnum()
	res.PortScan.Timeout = service.PortScanTimeoutEnum()
	res.PortScan.Concurrent = service.PortScanConcurrentEnum()

	// 爬虫
	res.Crawler.MaxDeep = service.CrawlerDeepEnum()
	res.Crawler.MaxUrl = service.CrawlerMaxUrlEnum()
	res.Crawler.ScanRange = service.CrawlerScanRangeEnum()
	res.Crawler.Timeout = service.CrawlerSingleTimeout()
	res.Crawler.FullTimeout = service.CrawlerFullTimeout()
	res.Crawler.ScanRepeat = service.CrawlerRepeatEnum()
	res.Crawler.CrawlerSpeed = service.CrawlerSpeedEnum()

	// web路径爆破
	res.WebPathScan.Speed = service.WebPathScanSpeedEnum()
	res.WebPathScan.Times = service.WebPathScanTimeEnum()
	//res.WebPathScan.ScanDict, err = service.WebPathScanDictEnum()
	if err != nil {
		return err
	}

	// 子域名收集
	res.SubdomainDictCollect, err = service.SubdomainDictCollectEnum()

	// web弱口令
	res.WeakPass.Services, err = service.WeakPassServiceEnum() // 服务
	if err != nil {
		return err
	}
	res.WeakPass.CommonUserDict, res.WeakPass.CommonPassDict, err = service.WeakPassCommonDictEnum()
	if err != nil {
		return err
	}
	res.WeakPass.Type = service.WeakPassDictType()
	res.WeakPass.OnlyUseAdd = false
	res.WeakPass.GuessNum = service.WeakPassGuessNumber()
	res.WeakPass.GuessTimeout = service.WeakPassGuessTime()
	res.WeakPass.GuessRate = service.WeakPassGuessRate()
	//横向移动

	return nil
}

// List 任务场景列表
func (a *SceneTaskTemplate) List(ctx context.Context, req *typespec.SceneTaskTemplateListReq, resp *typespec.SceneTaskTemplateListRes) error {
	var srvTemplate services.SceneTaskTemplate
	taskTemplateList, total, err := srvTemplate.Page(ctx, req.Page, req.Size, req.Search, []int{enums.TaskTemplateSourceUser, enums.TaskTemplateSourceSystem}, enums.TaskTemplateStatusSuccess)
	resp.Total = total
	if err != nil {
		return err
	}

	// 获取用户数据
	var userMaps map[int]mysqls.User
	userIds := make([]int, 0)
	for _, v := range taskTemplateList {
		if v.UserID > 0 {
			userIds = append(userIds, v.UserID)
		}
	}
	if len(userIds) > 0 {
		var srvUser services.User
		userMaps, err = srvUser.AllForIds(ctx, userIds)
		if err != nil {
			return err
		}
	}

	// 组合返回数据
	for _, v := range taskTemplateList {
		// 用户名获取
		username := ""
		if user, ok := userMaps[v.UserID]; ok {
			username = user.Username
		}

		temp := typespec.SceneTaskTemplateListResItem{
			Id:            v.ID,
			TemplateName:  v.TemplateName,
			Describe:      v.Describe,
			IsDefault:     int(v.IsDefault),
			IsDefaultName: enums.TaskTemplate.IsDefaultStr(v.IsDefault),
			CreateTime:    v.CreateTime.Format(utils.DateTime),
			UserId:        v.UserID,
			UserName:      username,
		}

		resp.List = append(resp.List, temp)
	}

	return nil
}

// 创建 或 更新
func (a *SceneTaskTemplate) CreateOrUpdate(ctx context.Context, req *typespec.SceneTaskTemplateCreateReq, resp *typespec.SceneTaskTemplateCreateRes) error {
	var srvTemplate services.SceneTaskTemplate
	template, err := srvTemplate.GetByName(ctx, req.TemplateName)
	if err != nil {
		return err
	}
	configStruct, err := data.TaskCheckTaskConfig.VerifyConfig(req.Config)
	if err != nil {
		return err
	}
	if req.TaskTemplateId != 0 { // 更新
		if template.ID != 0 && template.Estate == enums.TaskTemplateStatusSuccess && template.ID != req.TaskTemplateId {
			return errors.New("任务场景名字已经存在，请更换场景名字")
		}
		resp.TaskTemplateId, err = srvTemplate.Update(ctx, req.TaskTemplateId, configStruct, req.TemplateName, req.Describe, req.UserId)
		if err != nil {
			return err
		}
	} else { // 创建
		if template.ID > 0 && template.Estate == enums.TaskTemplateStatusSuccess {
			return errors.New("任务场景名字已经存在，请更换场景名字")
		}
		resp.TaskTemplateId, err = srvTemplate.Create(ctx, configStruct, req.TemplateName, req.Describe, req.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}

// Detail 任务场景详情
func (a *SceneTaskTemplate) Detail(ctx context.Context, req *typespec.SceneTaskTemplateDetailReq, resp *typespec.SceneTaskTemplateDetailRes) error {
	var srvTemplate services.SceneTaskTemplate
	templateData, configData, err := srvTemplate.Detail(ctx, req.TaskTemplateId, enums.TaskTemplateStatusSuccess)
	if err != nil {
		return err
	}
	// 场景数据
	resp.Id = templateData.ID
	resp.TemplateName = templateData.TemplateName
	resp.Describe = templateData.Describe
	resp.IsDefault = templateData.IsDefault
	resp.UserId = templateData.UserID
	resp.CreateTime = templateData.CreateTime.Format(enums.TimeLayout)
	resp.Config = configData
	return nil
}

// Copy 拷贝
func (a *SceneTaskTemplate) Copy(ctx context.Context, req *typespec.SceneTaskTemplateCopyReq, resp *typespec.SceneTaskTemplateCopyRes) error {
	var srvTemplate services.SceneTaskTemplate
	templateData, configData, err := srvTemplate.Detail(ctx, req.TaskTemplateId, enums.TaskTemplateStatusSuccess)
	if err != nil {
		return err
	}
	if templateData.ID == 0 {
		return errors.New("任务场景模版不存在，复制失败")
	}
	newTemplateName := fmt.Sprintf("copy_%s", templateData.TemplateName)

	resp.TaskTemplateId, err = srvTemplate.Create(ctx, configData, newTemplateName, templateData.Describe, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

// SetDefault 设置默认
func (a *SceneTaskTemplate) SetDefault(ctx context.Context, req *typespec.SceneTaskTemplateSetDefaultReq) error {
	var srv services.SceneTaskTemplate
	taskTemplate, err := srv.GetByTemplateId(ctx, req.TaskTemplateId, enums.TaskTemplateStatusSuccess)
	if err != nil {
		return err
	}
	if taskTemplate.ID == 0 {
		return errors.New("任务场景不存在")
	}

	return srv.SetDefault(ctx, req.TaskTemplateId)
}

func (a *SceneTaskTemplate) Del(ctx context.Context, req *typespec.SceneTaskTemplateDelReq) error {
	templateIds := make([]int, 0)
	paramTemplateIds := strings.Split(req.TaskTemplateIds, ",")
	for _, item := range paramTemplateIds {
		templateId, _ := strconv.Atoi(item)
		templateIds = append(templateIds, templateId)
	}
	var srv services.SceneTaskTemplate
	return srv.Del(ctx, templateIds)
}

// 任务中的模板下拉菜单
func (a *SceneTaskTemplate) ToTaskTemplateOptions(ctx context.Context, res *typespec.SceneTaskTemplateToCheckTaskOptionsRes) error {
	var srv services.SceneTaskTemplate
	templates, _, err := srv.All(ctx, enums.TaskTemplateStatusSuccess)
	if err != nil {
		return err
	}
	for _, item := range templates {
		res.TaskTemplate = append(res.TaskTemplate, typespec.TemplateParamsNode{
			Id:        item.ID,
			Name:      item.TemplateName,
			IsDefault: item.IsDefault,
			Describe:  item.Describe,
		})
	}
	return nil
}

// 任务中的模板下拉菜单
func (a *SceneTaskTemplate) Graph(ctx context.Context, req *typespec.GraphReq, res *typespec.GraphRes) error {
	var srv services.SceneTaskTemplate
	vulIds := srv.GetTemplateVulIds(ctx, req.TaskTemplateId)
	var graphSrv services.AttackGraph
	nodes, links, err := graphSrv.BuildGraphByVulIds(ctx, vulIds)
	if err != nil {
		return err
	}
	res.Nodes = nodes
	res.Links = links
	return nil
}
