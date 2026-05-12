package services

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/tools/data"
	"smart/tools/enums"
	"smart/tools/file"
	"strings"
	"time"
)

type Bas struct {
}

// 规则 - 规则添加
func (b *Bas) BasInsertRule(ctx context.Context, datas []BasImportRules) error {
	inserts := make([]mysqls.BasRules, 0)
	for _, item := range datas {
		inserts = append(inserts, mysqls.BasRules{
			Content:                        item.Content,
			Name:                           item.Name,
			NameZh:                         item.NameZh,
			ClassType:                      item.ClassType,
			Protocol:                       item.Protocal,
			Keywords:                       item.Keywords,
			KeywordsZh:                     item.KeywordsZH,
			Description:                    item.Description,
			DescriptionZh:                  item.DescriptionZH,
			Cve:                            item.Cve,
			RawTrafficBeyondIpPacketBase64: item.RawTrafficBeyondIpPacketBase64Name,
			RawTrafficBeyondHttpBase64:     item.RawTrafficBeyondHttpBase64Name,
			CreateTime:                     item.CreateTime,
			UpdateTime:                     item.UpdateTime,
			Hash:                           item.Hash,
		})
	}
	if len(inserts) > 0 {
		var basRuleModel mysqls.BasRules
		return basRuleModel.AddAll(ctx, inserts)
	}
	return nil
}

// 规则 - 规则导入
func (b *Bas) BasUpdateRule(ctx context.Context, datas []BasImportRules) error {
	for _, item1 := range datas {
		go func(item BasImportRules) {
			var basRuleModel mysqls.BasRules
			basRuleModel.ID = item.Id
			basRuleModel.NameZh = item.NameZh
			basRuleModel.ClassType = item.ClassType
			basRuleModel.Content = item.Content
			basRuleModel.Protocol = item.Protocal
			basRuleModel.Keywords = item.Keywords
			basRuleModel.KeywordsZh = item.KeywordsZH
			basRuleModel.Description = item.Description
			basRuleModel.DescriptionZh = item.DescriptionZH
			basRuleModel.Cve = item.Cve
			basRuleModel.RawTrafficBeyondIpPacketBase64 = item.RawTrafficBeyondIpPacketBase64Name
			basRuleModel.RawTrafficBeyondHttpBase64 = item.RawTrafficBeyondHttpBase64Name
			basRuleModel.UpdateTime = item.UpdateTime
			basRuleModel.Hash = item.Hash
			basRuleModel.UpdateBasRules(ctx)
		}(item1)

	}
	return nil
}

// 规则 - 分类枚举获取 BasRuleClassEnum
func (b *Bas) BasRuleClassEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.BasRuleEnum.AllClassEnum())
}

// 规则 - 风险等级枚举获取 BasRuleLevelRiskEnum
func (b *Bas) BasRuleLevelRiskEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.BasRuleEnum.AllBasRuleRisk())
}

// 规则 - 检测阶段枚举获取 BasRuleAttackStageEnum
func (b *Bas) BasRuleAttackStageEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.BasRuleEnum.AllBasRuleAttackStage())
}

// 规则 - 是否成功状态枚举 BasRuleStatusEnum
func (b *Bas) BasRuleStatusEnum() []typespec.GlobalOptionsItemRes {
	return toolsSort(enums.BasRuleEnum.AllBasResultStatus())
}

// 规则 - 通过names获取 [hash=>{}]
func (b *Bas) BasGetRuleByNames(ctx context.Context, names []string) map[string]mysqls.BasRules {
	var basRulesModel mysqls.BasRules
	data := basRulesModel.GetListByNames(ctx, names)
	returnData := make(map[string]mysqls.BasRules)
	for _, item := range data {
		returnData[item.Hash] = item
	}
	return returnData
}

// 剧本集 - 通过name获取
func (b *Bas) BasGetTemplateByName(ctx context.Context, name string) (mysqls.BasTemplate, error) {
	var basRulesModel mysqls.BasTemplate
	return basRulesModel.GetByName(ctx, name)
}

// 规则查询列表
func (b *Bas) BasRuleGet(ctx context.Context, page, size, attackStage int, attackType string, risk int, class string, search string, ids []string) ([]mysqls.BasRules, int64, error) {
	var basRuleModel mysqls.BasRules
	return basRuleModel.GetBasRulesList(ctx, page, size, attackStage, attackType, risk, class, search, ids)
}

// 规则查询 依据规则ID
func (b *Bas) GetBasRuleGetByIds(ctx context.Context, ids []int) []mysqls.BasRules {
	var basRuleModel mysqls.BasRules
	return basRuleModel.GetByIds(ctx, ids)
}

// 规则查询 依据规则ID 返回map
func (b *Bas) GetBasRuleMapByIds(ctx context.Context, ids []int) map[int]mysqls.BasRules {
	var (
		basRuleModel mysqls.BasRules
		result       = make(map[int]mysqls.BasRules, 0)
	)
	basRuleRes := basRuleModel.GetByIds(ctx, ids)
	for i := 0; i < len(basRuleRes); i++ {
		result[basRuleRes[i].ID] = basRuleRes[i]
	}
	return result
}

// 规则查询 依据规则ID
func (b *Bas) GetBasRuleById(ctx context.Context, id int) (mysqls.BasRules, error) {
	var basRuleModel mysqls.BasRules
	return basRuleModel.GetBasRulesById(ctx, id)
}

// 规则编辑
func (b *Bas) BasRuleEdit(ctx context.Context, id int, affectTarget string, attackMode, attackStage, riskLevel int, affectScore, relationAttackMethod, fixSuggest, refUrl string) error {
	var basRuleModel mysqls.BasRules
	basRule, err := basRuleModel.GetBasRulesById(ctx, id)
	if err != nil {
		return err
	}
	if basRule.ID == 0 {
		return errors.New("未知的规则")
	}
	basRule.AffectTarget = affectTarget
	basRule.AttackMode = attackMode
	basRule.AttackStage = attackStage
	basRule.RiskLevel = riskLevel
	basRule.AffectScore = affectScore
	basRule.RelationAttackMethod = relationAttackMethod
	basRule.FixSuggest = fixSuggest
	basRule.RefUrl = refUrl
	return basRule.UpdateBasRules(ctx)
}

// 剧本集创建
func (b *Bas) BasTemplateCreate(ctx context.Context, name, desc, ruleIds string) error {
	var basTemplateModel mysqls.BasTemplate
	basTemplateModel.Name = name
	basTemplateModel.Desc = desc
	basTemplateModel.RuleIds = ruleIds
	basTemplateModel.CreateTime = time.Now()
	basTemplateModel.UpdateTime = time.Now()
	return basTemplateModel.AddBasTemplate(ctx)
}

// 剧本集更新
func (b *Bas) BasTemplateUpdate(ctx context.Context, name, desc, ruleIds string, id int) error {
	var basTemplateModel mysqls.BasTemplate
	basTemplateModel.ID = id
	basTemplateModel.Name = name
	basTemplateModel.Desc = desc
	basTemplateModel.RuleIds = ruleIds
	basTemplateModel.UpdateTime = time.Now()
	return basTemplateModel.UpdateBasTemplate(ctx)
}

// 剧本集列表
func (b *Bas) BasTemplateGet(ctx context.Context, page, size int, search string) ([]mysqls.BasTemplate, int64, error) {
	var basTemplateModel mysqls.BasTemplate
	return basTemplateModel.GetBasTemplateList(ctx, page, size, search)
}

// 剧本集详情
func (b *Bas) BasGetTemplateById(ctx context.Context, id int) (mysqls.BasTemplate, error) {
	var basTemplateModel mysqls.BasTemplate
	basTemplateModel.ID = id
	return basTemplateModel.GetBasTemplate(ctx)
}

// 剧本集 依据IDs获取
func (b *Bas) BasGetTemplateNameByIds(ctx context.Context, ids []int) map[int]string {
	var basTemplateModel mysqls.BasTemplate
	basTemplate := basTemplateModel.GetByIds(ctx, ids)
	fmt.Println(basTemplate, ids)
	basTemplateNameMaps := make(map[int]string)
	for _, item := range basTemplate {
		basTemplateNameMaps[item.ID] = item.Name
	}
	return basTemplateNameMaps
}

// 剧本集删除
func (b *Bas) BasDelTemplateById(ctx context.Context, ids []int) error {
	var basTemplateModel mysqls.BasTemplate
	return basTemplateModel.DeleteByIds(ctx, ids)
}

// 剧本集设置是否默认
func (b *Bas) BasTemplateSetDefault(ctx context.Context, id int) error {
	var basTemplateModel mysqls.BasTemplate
	err := basTemplateModel.CancelDefault(ctx)
	if err != nil {
		return err
	}
	err = basTemplateModel.SetDefault(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// BAS任务列表
func (b *Bas) BasGetTask(ctx context.Context, page, size, riskLevel int, search string) ([]mysqls.BasTask, int64, error) {
	var basTaskModel mysqls.BasTask

	return basTaskModel.GetBasTaskList(ctx, page, size, riskLevel, search)
}

// 根据状态查询
func (b *Bas) GetBasTargetByTaskIdsAndStatus(ctx context.Context, taskIds any, status int) []mysqls.BasTarget {
	var basTargetModel mysqls.BasTarget
	return basTargetModel.GetByTaskIdsAndStatus(ctx, taskIds, status)
}

// BAS 任务及目标都设置为完成
func (b *Bas) BasFinishTask(ctx context.Context, taskId int) {
	// 任务已完成
	var basTaskModel mysqls.BasTask
	basTaskModel.UpdateStatusById(ctx, taskId, enums.BasTaskStatusDone)
	// 目标已完成
	var basTaskTargetModel mysqls.BasTarget
	basTaskTargetModel.UpdateStatusByTaskId(ctx, taskId, enums.BasTaskStatusDone)
}

// 修改任务和目标数据的状态
func (b *Bas) UpdateBasTaskTargetStatus(ctx context.Context, taskIds any, targetIds any, status int) error {
	var (
		basTaskMysql   mysqls.BasTask
		basTargetMysql mysqls.BasTarget
		params         = map[string]interface{}{"status": status, "update_time": time.Now()}
	)
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	err := basTaskMysql.UpdateByIds(dCtx, taskIds, params)
	if err != nil {
		return err
	}
	err = basTargetMysql.UpdateByIds(dCtx, targetIds, params)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// BAS 目标设置为运行中
func (b *Bas) BasRunningTarget(ctx context.Context, taskId, targetId int) {
	// 目标运行中
	var basTaskTargetModel mysqls.BasTarget
	basTaskTargetModel.UpdateStatusById(ctx, targetId, enums.BasTaskStatusRunning)

	// 任何一个目标为运行中，那么任务也必须为运行中
	var basTaskModel mysqls.BasTask
	basTaskModel.UpdateStatusById(ctx, taskId, enums.BasTaskStatusRunning)
}

// BAS 目标设置为完成
func (b *Bas) BasFinishTarget(ctx context.Context, taskId, targetId int) {
	// 目标已完成
	var basTaskTargetModel mysqls.BasTarget
	basTaskTargetModel.UpdateStatusById(ctx, targetId, enums.BasTaskStatusDone)
	// 任何一个目标为已完成，那么都需要校验该任务下是否都已完成，如已完成则需要将任务改为已完成
	notFinish := basTaskTargetModel.GetByTaskIdAndInStatus(ctx, taskId, []int{enums.BasTaskStatusWait, enums.BasTaskStatusRunning})
	if len(notFinish) == 0 {
		var basTaskModel mysqls.BasTask
		basTaskModel.UpdateStatusById(ctx, taskId, enums.BasTaskStatusDone)
	}
}

// BAS 目标添加日志
func (b *Bas) BasTargetLogCreates(ctx context.Context, data []BasTaskLogCreate) error {
	if len(data) == 0 {
		return nil
	}

	datas := make([]mysqls.BasLog, 0)
	for _, item := range data {
		datas = append(datas, mysqls.BasLog{
			BasTaskID:   item.BasTaskId,
			BasTargetId: item.BasTaskTargetId,
			Content:     item.Content,
			CreateTime:  time.Now(),
		})
	}
	var basTaskTargetLogModel mysqls.BasLog
	return basTaskTargetLogModel.AddAll(ctx, datas)
}

// BAS任务结束
func (b *Bas) BasEndTaskById(ctx context.Context, id int) error {
	var basTaskModel mysqls.BasTask
	return basTaskModel.UpdateStatusById(ctx, id, enums.BasTaskStatusDone)
}

// BAS任务删除
func (b *Bas) BasDelTask(ctx context.Context, ids []int) error {
	// 任务删除
	var basTaskModel mysqls.BasTask
	err := basTaskModel.DeleteById(ctx, ids)
	if err != nil {
		return err
	}

	// 目标删除
	var basTaskTargetModel mysqls.BasTarget
	err = basTaskTargetModel.DeleteByTaskId(ctx, ids)
	if err != nil {
		return err
	}

	// 日志删除
	var basTaskTargetLogModel mysqls.BasLog
	err = basTaskTargetLogModel.DeleteByTaskIds(ctx, ids)
	if err != nil {
		return err
	}

	//bas漏洞删除
	var basVulModel mysqls.Basvul
	err = basVulModel.DelBasvulByTaskIds(ctx, ids)
	if err != nil {
		return err
	}

	return nil
}

// BAS 依据ID获取任务
func (b *Bas) BasGetTaskById(ctx context.Context, id int) (mysqls.BasTask, error) {
	var basTaskModel mysqls.BasTask
	return basTaskModel.GetBasTaskById(ctx, id)
}

// BAS 依据ID获取任务
func (b *Bas) GetBasTaskIdsByStatus(ctx context.Context, status int) ([]int, error) {
	var (
		basTaskModel mysqls.BasTask
		result       = make([]int, 0)
	)
	taskRes, err := basTaskModel.GetBasTaskByStatus(ctx, status)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(taskRes); i++ {
		result = append(result, taskRes[i].ID)
	}
	return result, nil
}

// BAS 依据ids获取任务
func (b *Bas) BasGetTaskByIds(ctx context.Context, ids []int) []mysqls.BasTask {
	var basTaskModel mysqls.BasTask
	return basTaskModel.GetByIds(ctx, ids)
}

// BAS 依据ids获取任务 返回map
func (b *Bas) GetBasTaskMapByIds(ctx context.Context, ids []int) map[int]mysqls.BasTask {
	var (
		basTaskModel mysqls.BasTask
		result       = make(map[int]mysqls.BasTask, 0)
	)
	basTaskRes := basTaskModel.GetByIds(ctx, ids)
	for i := 0; i < len(basTaskRes); i++ {
		result[basTaskRes[i].ID] = basTaskRes[i]
	}
	return result
}

// BAS任务目标 依据任务ID获取
func (b *Bas) BasGetTaskTargetPageByTaskId(ctx context.Context, taskId, page, size int, search string) ([]mysqls.BasTarget, int64, error) {
	var basTaskTargetModel mysqls.BasTarget
	return basTaskTargetModel.GetBasTaskTargetList(ctx, taskId, page, size, search)
}

// BAS任务目标列表 批量获取目标的

// BAS任务目标 删除
func (b *Bas) BasTargetDel(ctx context.Context, targetIds []int) error {
	// 目标删除
	var basTaskTargetModel mysqls.BasTarget
	err := basTaskTargetModel.DeleteById(ctx, targetIds)
	if err != nil {
		return err
	}
	// 日志删除
	var basTaskTargetLogModel mysqls.BasLog
	err = basTaskTargetLogModel.DeleteByTargetIds(ctx, targetIds)
	if err != nil {
		return err
	}
	return nil
}

// BAS 目标日志
func (b *Bas) BasGetTargetLogs(ctx context.Context, targetId int) []mysqls.BasLog {
	var basTargetLogModel mysqls.BasLog
	return basTargetLogModel.AllByTargetId(ctx, targetId)
}

// BAS 目标日志
func (b *Bas) GetBasNodeByIds(ctx context.Context, targetIds any, onlinStatus, status int) ([]string, error) {
	var (
		basNodeMysql mysqls.Basnode
		result       []string
	)
	nodeRes, err := basNodeMysql.GetBasnodesByIds(ctx, targetIds, onlinStatus, status)
	if err != nil {
		return result, err
	}
	for i := 0; i < len(nodeRes); i++ {
		result = append(result, nodeRes[i].IP)
	}
	return result, nil
}

// BAS任务创建
func (b *Bas) BasCreateTask(ctx context.Context, name string, templateId int, ruleIds string, agentIps []string, uid int) error {
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()
	// 创建任务
	var basTaskModel = mysqls.BasTask{
		Name:            name,
		BasTemplateID:   templateId,
		BasTemplateJSON: ruleIds,
		RiskLevel:       enums.BasRiskLevelSafe,
		Status:          enums.BasTaskStatusWait,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		User:            uid,
	}
	if err := basTaskModel.AddBasTask(dCtx); err != nil {
		return err
	}
	// 创建target
	basTaskTargets := make([]mysqls.BasTarget, 0)
	for i := 0; i < len(agentIps); i++ {
		basTaskTargets = append(basTaskTargets, mysqls.BasTarget{
			BasTaskID:       basTaskModel.ID,
			Addr:            agentIps[i],
			Status:          enums.BasTaskStatusWait,
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
			BasTemplateID:   templateId,
			BasTemplateJSON: ruleIds,
		})
	}
	var basTaskTargetModel mysqls.BasTarget
	if err := basTaskTargetModel.AddAll(dCtx, basTaskTargets); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (b *Bas) SaveBasRuleFilePath(ruleRes []mysqls.BasRules) (string, error) {
	var fileContent []BasRuleFile
	for i := 0; i < len(ruleRes); i++ {
		var tmpData = BasRuleFile{
			RuleId:  ruleRes[i].ID,
			Content: ruleRes[i].Content,
		}
		fileContent = append(fileContent, tmpData)
	}
	fileContentByte, err := json.Marshal(fileContent)
	if err != nil {
		return "", err
	}
	filePatah, err := file.SaveTmpFile(enums.BasRuleFileName, string(fileContentByte))
	if err != nil {
		return "", err
	}
	return filePatah, nil
}

func (b *Bas) SendBasTask(ctx context.Context, basTaskId int, basTargetId int, ip string, filPath string, callBackFunc func(context.Context, int, int, string, string)) {
	paramList := []string{"--ip", ip, "--ruleFilePath", filPath}
	log.Info("SendBasTask cmd paramList:", paramList)
	cmd := exec.CommandContext(ctx, "./bas-send", paramList...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Info("SendBasTask cmd StdoutPipe err:", err)
		return
	}
	defer stdout.Close()
	if err = cmd.Start(); err != nil {
		log.Info("SendBasTask cmd Start err:", err)
		return
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		data, err := simplifiedchinese.GB18030.NewDecoder().Bytes(scanner.Bytes()) //处理乱码
		if err != nil {
			log.Info("SendBasTask simplifiedchinese NewDecoder err:", err)
			continue
		}
		log.Info("SendBasTask receive msg:", string(data))
		go callBackFunc(ctx, basTaskId, basTargetId, ip, string(data))
	}
	if err = cmd.Wait(); err != nil {
		log.Info("SendBasTask Wait err:", err)
	}
}

// 插入一条bas日志
func (b *Bas) AddBasLog(ctx context.Context, basTaskId, basTargetId int, data string) error {
	var basLogMysql = mysqls.BasLog{
		BasTaskID:   basTaskId,
		BasTargetId: basTargetId,
		Content:     data,
		CreateTime:  time.Now(),
	}
	return basLogMysql.AddBasLog(ctx)
}

// 根据规则id查询规则信息
func (b *Bas) GetBasRulesByIds(ctx context.Context, Ids []int) map[int]mysqls.BasRules {
	var (
		basRuleMysql mysqls.BasRules
		result       = make(map[int]mysqls.BasRules, 0)
	)

	basRuleRes := basRuleMysql.GetByIds(ctx, Ids)
	for i := 0; i < len(basRuleRes); i++ {
		result[basRuleRes[i].ID] = basRuleRes[i]
	}
	return result
}

// 整理发送bas返回的数据
func (b *Bas) OrderBasSendData(basTaskId, basTargetId int, ip string, contentMsg []BaSendMsgContent, basRulesMap map[int]mysqls.BasRules) []mysqls.Basvul {
	var result = make([]mysqls.Basvul, 0)
	for i := 0; i < len(contentMsg); i++ {
		if v, ok := basRulesMap[contentMsg[i].RuleId]; ok {
			tmpName := v.NameZh
			if len(tmpName) == 0 {
				tmpName = v.Name
			}
			var tmpData = mysqls.Basvul{
				BasTaskID:   basTaskId,
				BasTargetID: basTargetId,
				Addr:        ip,
				RuleID:      contentMsg[i].RuleId,
				RuleName:    tmpName,
				AttackMode:  v.AttackMode,
				AttackStage: v.AttackStage,
				RiskLevel:   enums.BasRiskLevelSafe, //默认安全
				Md5Code:     strings.Join(contentMsg[i].Md5, ","),
				Status:      enums.BasVulStatusFail, //默认失败
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			}
			result = append(result, tmpData)
		}
	}
	return result
}

// 批量新增漏洞测试数据
func (b *Bas) SaveBasVul(ctx context.Context, basVuls []mysqls.Basvul) error {
	var basVulMysql mysqls.Basvul
	return basVulMysql.AddBasvulMany(ctx, basVuls)
}

// 根据任务id获取漏洞测试数据
func (b *Bas) GetBasVulByTaskId(ctx context.Context, basTaskID int) []mysqls.Basvul {
	var basVulMysql mysqls.Basvul
	return basVulMysql.GetBasvulByTaskId(ctx, basTaskID)
}

// 根据目标id获取漏洞测试数据
func (b *Bas) GetBasVulByTargetIds(ctx context.Context, basTargetIDs []int, status int) []mysqls.Basvul {
	var basVulMysql mysqls.Basvul
	return basVulMysql.GetBasvulByTargetIds(ctx, basTargetIDs, status)
}

// 根据任务id获取漏洞测试（成功/失败）数量
func (b *Bas) GetBasVulCountByTaskIdAndStatus(ctx context.Context, basTaskID int) [2]int64 {
	var (
		basVulMysql mysqls.Basvul
		result      [2]int64
	)
	result[0] = basVulMysql.GetBasvulCountByTaskIdAndStatus(ctx, basTaskID, enums.BasVulStatusSuccess)
	result[1] = basVulMysql.GetBasvulCountByTaskIdAndStatus(ctx, basTaskID, enums.BasVulStatusFail)
	return result
}

// 根据任务id和状态获取危险等级数量
func (b *Bas) GetBasVulRiskLevelCountByTaskIdAndStatus(ctx context.Context, basTaskID int, status int) [4]int64 {
	var (
		basVulMysql mysqls.Basvul
		result      [4]int64
	)
	basVulRes := basVulMysql.GetBasVulRiskLevelCountByTaskIdAndStatus(ctx, basTaskID, status)
	for i := 0; i < len(basVulRes); i++ {
		if basVulRes[i].RiskLevel > 0 && basVulRes[i].RiskLevel < 5 && basVulRes[i].Total > 0 {
			result[basVulRes[i].RiskLevel-1] = basVulRes[i].Total
		}
	}
	return result
}

func (b *Bas) GetBasVulList(ctx context.Context, basTaskID int, page, size, riskLevel, attackStage, status int, search string, attackMode string) ([]mysqls.Basvul, int64) {
	var basVulMysql mysqls.Basvul
	return basVulMysql.GetBasvulList(ctx, basTaskID, page, size, riskLevel, attackStage, status, search, attackMode)
}

// 根据id删除漏洞测试数据
func (b *Bas) DelBasVulByIds(ctx context.Context, basVulIds any) error {
	var basVulMysql mysqls.Basvul
	return basVulMysql.DelBasvulByIds(ctx, basVulIds)
}

//更新BAS任务的风险等级
func (b *Bas) UpdateBasTaskRiskLevelById(ctx context.Context, basTaskId int) error {
	var (
		basTaskMysql mysqls.BasTask
		basVulMysql  mysqls.Basvul
	)
	basTaskRes, err := basTaskMysql.GetBasTaskById(ctx, basTaskId) //查询任务
	if err != nil {
		return err
	}
	basVulRes, err := basVulMysql.GetBasvulRiskLevelByTaskId(ctx, basTaskId, enums.BasVulStatusSuccess) //查询BAS漏洞风险等级高的数据
	if err != nil {
		return err
	}
	if basVulRes.ID > 0 && (basTaskRes.RiskLevel != basVulRes.RiskLevel) { //更新任务
		var params = map[string]interface{}{"risk_level": basVulRes.RiskLevel, "update_time": time.Now()}
		err = basTaskMysql.UpdateById(ctx, basTaskId, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// 处理agent心跳检测
func (b *Bas) BasReceivHeartBeat(ctx context.Context, ip string) error {
	var baseNodeMysql mysqls.Basnode
	//查询bas_node表
	basNodeRes, err := baseNodeMysql.GetBasNodeByIp(ctx, ip)
	if err != nil {
		return err
	}
	if basNodeRes.ID == 0 { //新增
		baseNodeMysql.Name = "渗透节点TX_" + ip
		baseNodeMysql.IP = ip
		baseNodeMysql.OnlineStatus = enums.BasNodeOnlineStatusOnline
		baseNodeMysql.Status = enums.BasNodeStatusEnable
		baseNodeMysql.CreateTime = time.Now()
		baseNodeMysql.UpdateTime = time.Now()
		err = baseNodeMysql.AddBasnode(ctx)
		if err != nil {
			return err
		}
	} else { //更新
		if basNodeRes.Status != enums.BasNodeStatusEnable {
			return nil
		}
		var params = map[string]interface{}{"online_status": enums.BasNodeOnlineStatusOnline, "update_time": time.Now()}
		err = basNodeRes.UpdateBasnode(ctx, basNodeRes.ID, params)
		if err != nil {
			return err
		}
	}
	return nil
}

//查询接收结果相关的数据
func (b *Bas) BasReceivResultt(ctx context.Context, ip string, md5 []string) ([]mysqls.Basvul, []int, []int, error) {
	var (
		basVulMysql mysqls.Basvul
		basRuleIds  = make([]int, 0)
		basTaskIds  = make([]int, 0)
	)
	//查询数据
	basVulRes, err := basVulMysql.GetBasvulByIPAndMd5(ctx, ip, enums.BasVulStatusFail, md5)
	if err != nil {
		return nil, basRuleIds, basTaskIds, err
	}
	//更改状态
	for i := 0; i < len(basVulRes); i++ {
		basRuleIds = append(basRuleIds, basVulRes[i].RuleID)
		basTaskIds = append(basTaskIds, basVulRes[i].BasTaskID)
	}
	basTaskIds = data.ArrayIntUnique(basTaskIds) //basTaskId去重
	basRuleIds = data.ArrayIntUnique(basRuleIds) //basRuleId去重
	return basVulRes, basRuleIds, basTaskIds, nil
}

//查询接收结果相关的数据
func (b *Bas) BasReceivUpdateBasVulRes(ctx context.Context, basVulData []mysqls.Basvul, basRuleMapData map[int]mysqls.BasRules) error {
	var basVulMysql mysqls.Basvul
	for i := 0; i < len(basVulData); i++ {
		riskLevel := 4
		if v, ok := basRuleMapData[basVulData[i].RuleID]; ok {
			riskLevel = v.RiskLevel
		}
		var params = map[string]interface{}{"status": enums.BasVulStatusSuccess, "risk_level": riskLevel, "update_time": time.Now()}
		err := basVulMysql.UpdateBasvulById(ctx, basVulData[i].ID, params)
		if err != nil {
			log.Info("BasReceivUpdateBasVulRes UpdateBasvulById err", err)
			continue
		}
	}
	return nil
}

// BAS Agent列表
func (b *Bas) BasAgentList(ctx context.Context, page, size int, search string) ([]mysqls.Basnode, int64, error) {
	var basNodeModel mysqls.Basnode
	return basNodeModel.GetBasNodeList(ctx, page, size, search)
}

// 根据状态和在线状态查询agent数据
func (b *Bas) GetBasAgentByStatusAndOnlinestatus(ctx context.Context, status, onlineStatus int) []mysqls.Basnode {
	var basNodeModel mysqls.Basnode
	return basNodeModel.GetBasNodeListByStatusAndOnlinestatus(ctx, status, onlineStatus)
}

// Bas Agent 状态修改
func (b *Bas) BasAgentStatusEdit(ctx context.Context, id, status int) error {
	var basNodeModel mysqls.Basnode
	return basNodeModel.UpdateBasnode(ctx, id, map[string]interface{}{"status": status})
}

// 根据在线状态和更新时间查询agent数据
func (b *Bas) GetBasAgentByOnlineStatusAndUpdatetime(ctx context.Context, onlineStatus int, updateTime string) []int {
	var (
		basNodeModel mysqls.Basnode
		result       = make([]int, 0)
	)
	basNodeRes := basNodeModel.GetBasAgentByOnlineStatusAndUpdatetime(ctx, onlineStatus, updateTime)
	for i := 0; i < len(basNodeRes); i++ {
		result = append(result, basNodeRes[i].ID)
	}
	return result
}

//批量更新basnode
func (b *Bas) UpdateBasNodeByIds(ctx context.Context, ids any, params map[string]interface{}) error {
	var basNodeModel mysqls.Basnode
	return basNodeModel.UpdateBasnodeByIds(ctx, ids, params)
}
