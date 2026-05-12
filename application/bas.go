package application

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/api/typespec"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/network"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"
)

type Bas struct {
}

// 规则导入
func (b *Bas) BasRuleImport(ctx context.Context, excelFile *excelize.File) error {
	rows, err := excelFile.GetRows("chaos_maker_rules")
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return errors.New("请确认文件中是否有chaos_maker_rules的sheet，且此sheet中是否有数据")
	}

	// 先获取指定字段所在位置
	var hashIndex, createdAtIndex, updatedAtIndex, protocolIndex, nameIndex, nameZHIndex, classTypeIndex, suricataRawIndex, keywordsIndex, keywordsZHIndex, descriptionIndex, descriptionZHIndex, cveIndex, rawTrafficBeyondIpPacketBase64Index, rawTrafficBeyondHttpBase64Index int
	createdAtName := "created_at"
	updatedAtName := "updated_at"
	protocolName := "protocol"
	nameName := "name"
	nameZhName := "name_zh"
	classTypeName := "class_type"
	suricataRawName := "suricata_raw"
	keywordsName := "keywords"
	keywordsZHName := "keywords_zh"
	descriptionName := "description"
	descriptionZHName := "description_zh"
	rawTrafficBeyondIpPacketBase64Name := "raw_traffic_beyond_ip_packet_base64"
	rawTrafficBeyondHttpBase64Name := "raw_traffic_beyond_http_base64"
	cve := "cve"
	hashName := "hash"
	for k, item := range rows[0] {
		switch item {
		case createdAtName:
			createdAtIndex = k
		case updatedAtName:
			updatedAtIndex = k
		case protocolName:
			protocolIndex = k
		case nameName:
			nameIndex = k
		case nameZhName:
			nameZHIndex = k
		case classTypeName:
			classTypeIndex = k
		case suricataRawName:
			suricataRawIndex = k
		case keywordsName:
			keywordsIndex = k
		case keywordsZHName:
			keywordsZHIndex = k
		case descriptionName:
			descriptionIndex = k
		case descriptionZHName:
			descriptionZHIndex = k
		case cve:
			cveIndex = k
		case rawTrafficBeyondIpPacketBase64Name:
			rawTrafficBeyondIpPacketBase64Index = k
		case rawTrafficBeyondHttpBase64Name:
			rawTrafficBeyondHttpBase64Index = k
		case hashName:
			hashIndex = k
		}
	}

	if createdAtIndex == 0 || updatedAtIndex == 0 || protocolIndex == 0 || nameIndex == 0 || suricataRawIndex == 0 {
		return errors.New("请确认文件中是否有chaos_maker_rules的sheet中是否有" + createdAtName + "、" + updatedAtName + "、" + protocolName + "、" + nameName + "、" + suricataRawName + "字段")
	}
	// 组合数据
	processDatas := make([]services.BasImportRules, 0)
	// 将名称储存，校验名称是否存在，已存在的进行修改，不存在的进行添加
	names := make([]string, 0)

	for k, item := range rows {
		if k == 0 {
			continue
		}

		var tempImportRules services.BasImportRules

		// 某些元素可能不存在
		if len(item) > suricataRawIndex {
			tempImportRules.Content = item[suricataRawIndex]
		}
		if len(item) > nameIndex {
			tempImportRules.Name = item[nameIndex]
		}
		if len(item) > nameZHIndex {
			tempImportRules.NameZh = item[nameZHIndex]
		}
		if len(item) > classTypeIndex {
			tempImportRules.ClassType = item[classTypeIndex]
		}
		if len(item) > protocolIndex {
			tempImportRules.Protocal = item[protocolIndex]
		}
		if len(item) > keywordsIndex {
			tempImportRules.Keywords = item[keywordsIndex]
		}
		if len(item) > keywordsZHIndex {
			tempImportRules.KeywordsZH = item[keywordsZHIndex]
		}
		if len(item) > descriptionIndex {
			tempImportRules.Description = item[descriptionIndex]
		}
		if len(item) > descriptionZHIndex {
			tempImportRules.DescriptionZH = item[descriptionZHIndex]
		}
		if len(item) > cveIndex {
			tempImportRules.Cve = item[cveIndex]
		}
		if len(item) > rawTrafficBeyondIpPacketBase64Index {
			tempImportRules.RawTrafficBeyondIpPacketBase64Name = item[rawTrafficBeyondIpPacketBase64Index]
		}
		if len(item) > rawTrafficBeyondHttpBase64Index {
			tempImportRules.RawTrafficBeyondHttpBase64Name = item[rawTrafficBeyondHttpBase64Index]
		}
		if len(item) > hashIndex {
			tempImportRules.Hash = item[hashIndex]
		}
		tempImportRules.CreateTime = time.Now()
		tempImportRules.UpdateTime = time.Now()

		processDatas = append(processDatas, tempImportRules)
		names = append(names, item[nameIndex])
	}

	var srv services.Bas

	inserts := make([]services.BasImportRules, 0)
	Updates := make([]services.BasImportRules, 0)
	alreadyDatas := srv.BasGetRuleByNames(ctx, names)
	for _, item := range processDatas {
		if aData, ok := alreadyDatas[item.Hash]; ok {
			item.Id = aData.ID
			Updates = append(Updates, item)
		} else {
			inserts = append(inserts, item)
		}
	}

	// 新增
	srv.BasInsertRule(ctx, inserts)
	// 修改
	srv.BasUpdateRule(ctx, Updates)

	return nil
}

// 规则枚举
func (b *Bas) BasRuleEnum(ctx context.Context, res *typespec.BasRuleEnumRes) error {
	var basSrv services.Bas
	// 规则分类
	res.Class = basSrv.BasRuleClassEnum()
	// 影响级别
	res.RiskLevel = basSrv.BasRuleLevelRiskEnum()
	// 攻击阶段
	res.AttackStage = basSrv.BasRuleAttackStageEnum()

	res.Status = basSrv.BasRuleStatusEnum()

	// 获取漏洞服务枚举值
	decisionOptions, err := httpclients.GetDecisionOptions(ctx)
	if err != nil {
		return errors.New("decision服务通信错误：" + err.Error())
	}
	// 攻击类型取漏洞的漏洞类型 - 启兵给的需求
	res.AttackType = decisionOptions.Data.VulLibrariesType

	return nil
}

// 规则查询列表
func (b *Bas) BasRuleGet(ctx context.Context, req *typespec.BasRuleListReq, res *typespec.BasRuleListRes) error {
	var basSrv services.Bas
	var idArray = make([]string, 0)
	if len(req.Ids) > 0 {
		idArray = strings.Split(req.Ids, ",")
	}
	lists, total, _ := basSrv.BasRuleGet(ctx, req.Page, req.Size, req.AttackStage, req.AttackType, req.RiskLevel, enums.BasRuleEnum.GetClassEnum(req.ClassType), req.Search, idArray)
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total

	// 获取漏洞服务枚举值
	decisionOptions, err := httpclients.GetDecisionOptions(ctx)
	if err != nil {
		return errors.New("decision服务通信错误：" + err.Error())
	}
	// 攻击类型取漏洞的漏洞类型 - 启兵给的需求
	attackMode := decisionOptions.Data.VulLibrariesType
	attackModeMap := make(map[int]string)
	for _, item := range attackMode {
		attackModeMap[int(item.Value.(float64))] = item.Label.(string)
	}

	for _, item := range lists {
		res.List = append(res.List, typespec.BasRuleListItem{
			Id:                             item.ID,
			Content:                        item.Content,
			Name:                           item.Name,
			NameZh:                         item.NameZh,
			Hash:                           item.Hash,
			ClassType:                      item.ClassType,
			Protocol:                       item.Protocol,
			Keywords:                       item.Keywords,
			KeywordsZh:                     item.KeywordsZh,
			Description:                    item.Description,
			DescriptionZh:                  item.DescriptionZh,
			Cve:                            item.Cve,
			RawTrafficBeyondIpPacketBase64: item.RawTrafficBeyondIpPacketBase64,
			RawTrafficBeyondHttpBase64:     item.RawTrafficBeyondHttpBase64,
			AttackMode:                     item.AttackMode,
			AttackModeEnum:                 attackModeMap[item.AttackMode],
			AttackStage:                    item.AttackStage,
			AttackStageEnum:                enums.BasRuleEnum.GetBasRuleAttackStage(item.AttackStage),
			RiskLevel:                      item.RiskLevel,
			RiskLevelEnum:                  enums.BasRuleEnum.GetBasRuleRisk(item.RiskLevel),
			AffectTarget:                   item.AffectTarget,
			AffectScore:                    item.AffectScore,
			RelationAttackMethod:           item.RelationAttackMethod,
			RefUrl:                         item.RefUrl,
			FixSuggest:                     item.FixSuggest,
			CreateTime:                     item.CreateTime.Format(utils.DateTime),
			UpdateTime:                     item.UpdateTime.Format(utils.DateTime),
		})
	}
	return nil
}

// 规则详情
func (b *Bas) BasRuleInfo(ctx context.Context, req *typespec.BasRuleInfoReq, resp *typespec.BasRuleInfoResp) error {
	var basSrv services.Bas
	basRuleRes, err := basSrv.GetBasRuleById(ctx, req.BasRuleId)
	if err != nil {
		return err
	}
	// 获取漏洞服务枚举值
	decisionOptions, err := httpclients.GetDecisionOptions(ctx)
	if err != nil {
		return errors.New("decision服务通信错误：" + err.Error())
	}
	// 攻击类型取漏洞的漏洞类型 - 启兵给的需求
	attackMode := decisionOptions.Data.VulLibrariesType
	attackModeMap := make(map[int]string)
	for _, item := range attackMode {
		attackModeMap[int(item.Value.(float64))] = item.Label.(string)
	}
	resp.Id = basRuleRes.ID
	resp.Content = basRuleRes.Content
	resp.Name = basRuleRes.Name
	resp.NameZh = basRuleRes.NameZh
	if len(resp.NameZh) == 0 {
		resp.NameZh = basRuleRes.Name
	}
	resp.Hash = basRuleRes.Hash
	resp.ClassType = basRuleRes.ClassType
	resp.Protocol = basRuleRes.Protocol
	resp.Keywords = basRuleRes.Keywords
	resp.KeywordsZh = basRuleRes.KeywordsZh
	resp.Description = basRuleRes.Description
	resp.DescriptionZh = basRuleRes.DescriptionZh
	resp.Cve = basRuleRes.Cve
	resp.RawTrafficBeyondIpPacketBase64 = basRuleRes.RawTrafficBeyondIpPacketBase64
	resp.RawTrafficBeyondHttpBase64 = basRuleRes.RawTrafficBeyondHttpBase64
	resp.AttackMode = basRuleRes.AttackMode
	resp.AttackModeEnum = attackModeMap[basRuleRes.AttackMode]
	resp.AttackStage = basRuleRes.AttackStage
	resp.AttackStageEnum = enums.BasRuleEnum.GetBasRuleAttackStage(basRuleRes.AttackStage)
	resp.RiskLevel = basRuleRes.RiskLevel
	resp.RiskLevelEnum = enums.BasRuleEnum.GetBasRuleRisk(basRuleRes.RiskLevel)
	resp.AffectTarget = basRuleRes.AffectTarget
	resp.AffectScore = basRuleRes.AffectScore
	resp.RelationAttackMethod = basRuleRes.RelationAttackMethod
	resp.RefUrl = basRuleRes.RefUrl
	resp.FixSuggest = basRuleRes.FixSuggest
	resp.CreateTime = basRuleRes.CreateTime.Format(utils.DateTime)
	resp.UpdateTime = basRuleRes.UpdateTime.Format(utils.DateTime)
	return nil
}

// 规则编辑
func (b *Bas) BasRuleEdit(ctx context.Context, req *typespec.BasRuleEditReq, res *typespec.BasRuleEditRes) error {
	var basSrv services.Bas
	return basSrv.BasRuleEdit(ctx, req.Id, req.AffectTarget, req.AttackMode, req.AttackStage, req.RiskLevel, req.AffectScore, req.RelationAttackMethod, req.FixSuggest, req.RefUrl)
}

// 剧本集创建
func (b *Bas) BasTemplateCreate(ctx context.Context, req *typespec.BasTemplateCreateReq) error {
	ruleIdsJson, err := json.Marshal(req.RuleIds)
	if err != nil {
		return errors.New("ruleIds解析错误：" + err.Error())
	}
	var basSrv services.Bas
	basTemplate, err := basSrv.BasGetTemplateByName(ctx, req.Name)
	if err != nil {
		return err
	}

	if req.Id == 0 {
		// 添加
		if basTemplate.ID != 0 {
			return errors.New("该剧本名称已存在")
		}
		return basSrv.BasTemplateCreate(ctx, req.Name, req.Desc, string(ruleIdsJson))
	} else {
		// 更新
		if basTemplate.ID != 0 && basTemplate.ID != req.Id {
			return errors.New("该剧本名称已存在")
		}
		return basSrv.BasTemplateUpdate(ctx, req.Name, req.Desc, string(ruleIdsJson), req.Id)
	}
}

// 剧本集列表
func (b *Bas) BasTemplateGet(ctx context.Context, req *typespec.BasTemplateListReq, res *typespec.BasTemplateListRes) error {
	var basSrv services.Bas
	lists, total, _ := basSrv.BasTemplateGet(ctx, req.Page, req.Size, req.Search)
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total
	for _, item := range lists {
		res.List = append(res.List, typespec.BasTemplateListItem{
			Id:            item.ID,
			Name:          item.Name,
			Desc:          item.Desc,
			IsDefault:     item.IsDefault,
			IsDefaultEnum: enums.BasTemplateEnum.GetIsDefault(item.IsDefault),
			CreateTime:    item.CreateTime.Format(utils.DateTime),
		})
	}
	return nil
}

// 剧本集详情
func (b *Bas) BasGetTemplateById(ctx context.Context, req *typespec.BasTemplateGetReq, res *typespec.BasTemplateGetRes) error {
	var basSrv services.Bas
	template, err := basSrv.BasGetTemplateById(ctx, req.Id)
	if err != nil {
		return err
	}
	if template.ID == 0 {
		return errors.New("未知的剧本集")
	}
	ruleIds := make([]int, 0)
	if err := json.Unmarshal([]byte(template.RuleIds), &ruleIds); err != nil {
		return errors.New("解析规则ID错误")
	}
	res.Id = template.ID
	res.Name = template.Name
	res.Desc = template.Desc
	res.RuleIds = ruleIds
	res.CreateTime = time.Now().Format(utils.DateTime)
	res.UpdateTime = time.Now().Format(utils.DateTime)
	return nil
}

// 剧本集删除
func (b *Bas) BasDelTemplateById(ctx context.Context, req *typespec.BasTemplateDelReq) error {
	idsStr := strings.Split(req.Id, ",")
	if len(idsStr) == 0 {
		return nil
	}

	ids := make([]int, 0)
	for _, idS := range idsStr {
		id, err := strconv.Atoi(idS)
		if err != nil {
			return errors.New("请确认ID是否为number:" + err.Error())
		}
		ids = append(ids, id)
	}

	var basSrv services.Bas
	return basSrv.BasDelTemplateById(ctx, ids)
}

// 剧本集设置是否默认
func (b *Bas) BasTemplateSetDefault(ctx context.Context, req *typespec.BasTemplateSetDefaultReq) error {
	var basSrv services.Bas
	basTemplate, err := basSrv.BasGetTemplateById(ctx, req.Id)
	if err != nil {
		return err
	}
	if basTemplate.ID == 0 {
		return errors.New("未知的剧本集")
	}
	return basSrv.BasTemplateSetDefault(ctx, req.Id)
}

// Bas agent 是否在线
func (b *Bas) BasAgentIsOnline(ctx context.Context, req *typespec.BasAgentIsOnlineReq) error {
	if !network.TelnetIsOpen(req.Host, req.Port) {
		return errors.New(req.Host + ":" + req.Port + "未在线或填写错误")
	}
	return nil
}

// BAS任务创建
func (b *Bas) BasCreateTask(ctx context.Context, req *typespec.BasTaskCreateReq, uid int) error {
	var basSrv services.Bas
	//查询方案数据
	template, err := basSrv.BasGetTemplateById(ctx, req.BasTemplateId)
	if err != nil || template.ID == 0 {
		return errors.New("查询评估方案出错...")
	}
	//查询节点数据
	agentIps, err := basSrv.GetBasNodeByIds(ctx, req.BasNodeIds, enums.BasNodeOnlineStatusOnline, enums.BasNodeStatusEnable)
	if err != nil || (len(agentIps) != len(req.BasNodeIds)) {
		return errors.New("所选节点离线或查询节点数据失败...")
	}
	//创建数据
	err = basSrv.BasCreateTask(ctx, req.Name, req.BasTemplateId, template.RuleIds, agentIps, uid)
	if err != nil {
		return err
	}
	return nil
}

// bas心跳及检测结果接收
func (b *Bas) BasReceivResult(ctx context.Context, req *typespec.BasReceivResultReq) error {
	log.Info("BasReceivResult receive req:", req)
	var basSrv services.Bas
	if req.Type == "heartbeat" {
		var tmp BasReceivHeartBeat
		err := json.Unmarshal([]byte(req.Content), &tmp)
		if err != nil {
			return err
		}
		err = basSrv.BasReceivHeartBeat(ctx, tmp.Ip)
		if err != nil {
			return err
		}
	} else if req.Type == "result" {
		var tmp BasReceivResult
		err := json.Unmarshal([]byte(req.Content), &tmp)
		if err != nil {
			return err
		}
		basVulRes, basRuleIds, basTaskIds, err := basSrv.BasReceivResultt(ctx, tmp.Ip, tmp.Md5) //查询相关数据
		if err != nil {
			return err
		}
		basRuleMap := basSrv.GetBasRuleMapByIds(ctx, basRuleIds)          //查询规则数据
		err = basSrv.BasReceivUpdateBasVulRes(ctx, basVulRes, basRuleMap) //更改bas漏洞数据
		if err != nil {
			return err
		}
		//更改bas任务危险等级
		for i := 0; i < len(basTaskIds); i++ {
			err = basSrv.UpdateBasTaskRiskLevelById(ctx, basTaskIds[i])
			if err != nil {
				log.Info("BasReceivResult UpdateBasTaskRiskLevelById err:", err)
				continue
			}
		}
	}
	return nil
}

// BAS任务列表
func (b *Bas) BasGetTask(ctx context.Context, req *typespec.BasTaskGetReq, res *typespec.BasTaskGetRes) error {
	var basSrv services.Bas
	lists, total, _ := basSrv.BasGetTask(ctx, req.Page, req.Size, req.RiskLevel, req.Search)
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total

	// 获取中间数据用于获取其他表数据
	basTaskIds := make([]int, 0)
	templateIds := make([]int, 0)
	for _, item := range lists {
		basTaskIds = append(basTaskIds, item.ID)
		templateIds = append(templateIds, item.BasTemplateID)
	}

	// 获取评估方案信息
	basTemplateNames := basSrv.BasGetTemplateNameByIds(ctx, templateIds)

	// 获取用户数据
	var userMaps map[int]mysqls.User
	userIds := make([]int, 0)
	for _, v := range lists {
		if v.User > 0 {
			userIds = append(userIds, v.User)
		}
	}
	if len(userIds) > 0 {
		var srvUser services.User
		userMaps, _ = srvUser.AllForIds(ctx, userIds)
	}

	for _, item := range lists {
		username := ""
		if user, ok := userMaps[item.User]; ok {
			username = user.Username
		}
		res.List = append(res.List, typespec.BasTaskGetItem{
			Id:            item.ID,
			Name:          item.Name,
			TemplateId:    item.BasTemplateID,
			TemplateName:  basTemplateNames[item.BasTemplateID],
			RiskLevel:     item.RiskLevel,
			RiskLevelEnum: enums.BasEnum.BasRiskLevelEnum(item.RiskLevel),
			CreateTime:    item.CreateTime.Format(utils.DateTime),
			Status:        item.Status,
			StatusEnum:    enums.BasEnum.GetStatus(item.Status),
			UserId:        item.User,
			UserName:      username,
		})
	}
	return nil
}

// BAS任务结束
func (b *Bas) BasEndTaskById(ctx context.Context, req *typespec.BasTaskEndReq) error {
	//// 给缓存中写入任务结束，由执行的地方获取，进行结束操作
	//redisClient, err := redis.NewClient()
	//if err != nil {
	//	return err
	//}
	//redisClient.Set(ctx, "bas_end_task_"+strconv.Itoa(req.Id), "Y", 10*time.Second)

	// 更新任务状态为结束
	var basSrv services.Bas
	return basSrv.BasEndTaskById(ctx, req.Id)
}

// BAS任务删除
func (b *Bas) BasDelTask(ctx context.Context, req *typespec.BasTaskDelReq) error {
	ids := strings.Split(req.Id, ",")
	idInts := make([]int, 0)
	for _, id := range ids {
		idI, err := strconv.Atoi(id)
		if err == nil {
			idInts = append(idInts, idI)
		}
	}
	if len(idInts) > 0 {
		var basSrv services.Bas
		basTasks := basSrv.BasGetTaskByIds(ctx, idInts)
		for _, item := range basTasks {
			if item.Status != enums.BasTaskStatusDone {
				return errors.New("存在未结束的任务")
			}
		}
		tx := mysql.DB.Begin()
		dCtx := mysql.NewContext(ctx, tx)
		defer tx.Rollback()

		err := basSrv.BasDelTask(dCtx, idInts)
		if err != nil {
			return err
		}

		if err = tx.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

// BAS任务详情
func (b *Bas) BasGetTaskTargetPage(ctx context.Context, req *typespec.BasTaskDetailReq, res *typespec.BasTaskDetailRes) error {
	var basSrv services.Bas
	basTaskTargets, total, _ := basSrv.BasGetTaskTargetPageByTaskId(ctx, req.Id, req.Page, req.Size, req.Search)
	res.Page = req.Page
	res.Size = req.Size
	res.Total = total

	// 初始化风险等级
	targetRiskLevelMap := make(map[int]int)
	// 组合漏洞数量
	riskNumMap := make(map[int]map[int]int)
	// 目标ID
	targetIds := make([]int, 0)
	for _, item := range basTaskTargets {
		targetIds = append(targetIds, item.ID)
		targetRiskLevelMap[item.ID] = enums.BasRiskLevelSafe
		riskNumMap[item.ID] = make(map[int]int)
	}

	// 获取漏洞数据 组合漏洞数量 与风险等级
	vuls := basSrv.GetBasVulByTargetIds(ctx, targetIds, enums.BasVulStatusSuccess)
	for _, item := range vuls {
		numMap := riskNumMap[item.BasTargetID]
		if _, ok1 := numMap[item.RiskLevel]; ok1 {
			numMap[item.RiskLevel] += 1
		} else {
			numMap[item.RiskLevel] = 1
		}
		riskNumMap[item.BasTargetID] = numMap

		if targetRiskLevelMap[item.BasTargetID] > item.RiskLevel {
			targetRiskLevelMap[item.BasTargetID] = item.RiskLevel
		}
	}

	for _, item := range basTaskTargets {
		riskNum := riskNumMap[item.ID]
		res.List = append(res.List, typespec.BasTaskDetailItem{
			Id:            item.ID,
			Addr:          item.Addr,
			RiskLevel:     targetRiskLevelMap[item.ID],
			RiskLevelEnum: enums.BasEnum.BasRiskLevelEnum(targetRiskLevelMap[item.ID]),
			HighNum:       riskNum[enums.BasRiskLevelHigh],
			MidNum:        riskNum[enums.BasRiskLevelMiddle],
			LowNum:        riskNum[enums.BasRiskLevelLow],
			SafeNum:       riskNum[enums.BasRiskLevelSafe],
			Status:        item.Status,
			StatusEnum:    enums.BasEnum.GetStatus(item.Status),
			Create:        item.CreateTime.Format(utils.DateTime),
		})
	}
	return nil
}

// BAS任务目标日志
func (b *Bas) BasGetTargetLogs(ctx context.Context, req *typespec.BasTargetLogReq, res *typespec.BasTargetLogRes) error {
	var basSrv services.Bas
	logsData := basSrv.BasGetTargetLogs(ctx, req.Id)
	for _, item := range logsData {
		res.Content = append(res.Content, "["+item.CreateTime.Format(utils.DateTime)+"] "+item.Content)
	}
	return nil
}

// BAS任务删除
func (b *Bas) BasTargetDel(ctx context.Context, req *typespec.BasTargetDelReq) error {
	idsStr := strings.Split(req.Id, ",")
	if len(idsStr) == 0 {
		return nil
	}

	ids := make([]int, 0)
	for _, idS := range idsStr {
		id, err := strconv.Atoi(idS)
		if err != nil {
			return errors.New("请确认ID是否为number:" + err.Error())
		}
		ids = append(ids, id)
	}

	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	var basSrv services.Bas
	err := basSrv.BasTargetDel(dCtx, ids)
	if err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// BAS枚举
func (b *Bas) BasEnum(ctx context.Context, resp *typespec.BasEnumRes) error {
	var basEnums = enums.BasEnum
	resp.TaskStatus = basEnums.GetBasTaskStatusEnumArray()
	resp.RiskLevel = basEnums.GetBasRiskLevelEnumArray()
	resp.VulStatus = basEnums.GetBasVulStatusEnumArray()
	return nil
}

// 漏洞测试统计
func (b *Bas) BasVulStat(ctx context.Context, req *typespec.BasVulStatReq, resp *typespec.BasVulStatResp) error {
	var basSrv services.Bas
	//统计成功和失败的数量
	resp.Status = basSrv.GetBasVulCountByTaskIdAndStatus(ctx, req.BasTaskId)
	//统计成功状态的危险等级数量
	resp.RiskLevel = basSrv.GetBasVulRiskLevelCountByTaskIdAndStatus(ctx, req.BasTaskId, enums.BasVulStatusSuccess)
	return nil
}

// 漏洞测试列表
func (b *Bas) BasVulList(ctx context.Context, req *typespec.BasVulListReq, resp *typespec.BasVulListResp) error {
	var basSrv services.Bas
	basVulRes, total := basSrv.GetBasVulList(ctx, req.BasTaskId, req.Page, req.Size, req.RiskLevel, req.AttackStage, req.Status, req.Search, req.AttackMode)
	resp.Total = total
	resp.List = make([]typespec.BasVulListRespItems, 0)

	// 获取漏洞服务枚举值
	decisionOptions, err := httpclients.GetDecisionOptions(ctx)
	if err != nil {
		return errors.New("decision服务通信错误：" + err.Error())
	}
	// 攻击类型取漏洞的漏洞类型 - 启兵给的需求
	attackMode := decisionOptions.Data.VulLibrariesType
	attackModeMap := make(map[int]string)
	for _, item := range attackMode {
		attackModeMap[int(item.Value.(float64))] = item.Label.(string)
	}

	for i := 0; i < len(basVulRes); i++ {
		var tmpData = typespec.BasVulListRespItems{
			Id:              basVulRes[i].ID,
			RuleID:          basVulRes[i].RuleID,
			Addr:            basVulRes[i].Addr,
			RuleName:        basVulRes[i].RuleName,
			AttackMode:      basVulRes[i].AttackMode,
			AttackModeName:  attackModeMap[basVulRes[i].AttackMode],
			AttackStage:     basVulRes[i].AttackStage,
			AttackStageName: enums.BasRuleEnum.GetBasRuleAttackStage(basVulRes[i].AttackStage),
			RiskLevel:       basVulRes[i].RiskLevel,
			RiskLevelName:   enums.BasEnum.BasRiskLevelEnum(basVulRes[i].RiskLevel),
			Status:          basVulRes[i].Status,
			StatusName:      enums.BasEnum.BasVulStatusEnum(basVulRes[i].Status),
		}
		resp.List = append(resp.List, tmpData)
	}
	return nil
}

func (b *Bas) BasVulDel(ctx context.Context, req *typespec.BasVulDelReq) error {
	basVulArray := strings.Split(req.BasVulIds, ",")
	var basSrv services.Bas
	err := basSrv.DelBasVulByIds(ctx, basVulArray) //删除数据
	if err != nil {
		return err
	}
	err = basSrv.UpdateBasTaskRiskLevelById(ctx, req.BasTaskId) //更新任务的风险等级
	if err != nil {
		return err
	}
	return nil
}

// BAS Agent列表
func (b *Bas) BasAgentList(ctx context.Context, req *typespec.BasAgentListReq, res *typespec.BasAgentListRes) error {
	var basSrv services.Bas
	list, total, err := basSrv.BasAgentList(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		return err
	}
	res.Total = int(total)
	for _, item := range list {
		res.List = append(res.List, typespec.BasAgentListItem{
			Id:               item.ID,
			Name:             item.Name,
			Ip:               item.IP,
			OnlineStatus:     item.OnlineStatus,
			OnlineStatusEnum: enums.BasEnum.GetBaseNodeOnlineStatus(item.OnlineStatus),
			Status:           item.Status,
			StatusEnum:       enums.BasEnum.GetBaseNodeStatus(item.Status),
		})
	}
	return nil
}

// 可用节点列表
func (b *Bas) BasAgentLive(ctx context.Context, resp *typespec.BasAgentLiveResp) error {
	var basSrv services.Bas
	resp.List = make([]typespec.BasAgentLiveRespItem, 0)
	list := basSrv.GetBasAgentByStatusAndOnlinestatus(ctx, enums.BasNodeStatusEnable, enums.BasNodeOnlineStatusOnline)
	for i := 0; i < len(list); i++ {
		var tmp = typespec.BasAgentLiveRespItem{
			Id:   list[i].ID,
			Name: list[i].Name,
			Ip:   list[i].IP,
		}
		resp.List = append(resp.List, tmp)
	}
	return nil
}

// Bas Agent 状态修改 BasAgentStatusEdit
func (b *Bas) BasAgentStatusEdit(ctx context.Context, req *typespec.BasAgentStatusEditReq, res *typespec.BasAgentStatusEditRes) error {
	// 参数是否满足范围
	if req.Status != enums.BasNodeStatusEnable && req.Status != enums.BasNodeStatusDisable {
		return errors.New("状态仅支持1启用｜2禁用")
	}
	var basSrv services.Bas
	return basSrv.BasAgentStatusEdit(ctx, req.Id, req.Status)
}
