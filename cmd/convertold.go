//go:build ignore

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strings"
	"time"
)

type TaskChecktask struct {
	Id                  int       `json:"id"                  ` //
	TaskName            string    `json:"taskName"            ` //
	TaskType            int       `json:"taskType"            ` //
	TargetNumber        int       `json:"targetNumber"        ` //
	RiskLevel           int       `json:"riskLevel"           ` //
	UpdateTime          time.Time `json:"updateTime"          ` //
	Progress            int       `json:"progress"            ` //
	CheckTarget         string    `json:"checkTarget"         ` //
	VulnNumber          int       `json:"vulnNumber"          ` //
	Hc                  int       `json:"hc"                  ` //
	TaskConfigId        int       `json:"taskConfigId"        ` //
	TaskTemplateId      int       `json:"taskTemplateId"      ` //
	UserId              int       `json:"userId"              ` //
	TaskStatus          int       `json:"taskStatus"          ` //
	CreateTime          time.Time `json:"createTime"          ` //
	RelationTask        string    `json:"relationTask"        ` //
	MonitorTaskId       string    `json:"monitorTaskId"       ` // 监控任务id（自定义）
	IsLatest            int       `json:"isLatest"            ` // 是否是同一任务体系的最后一个
	TaskConfigurationId int       `json:"taskConfigurationId" ` // 任务配置id
}

// TaskChecktarget is the golang structure for table task_checktarget.
type TaskChecktarget struct {
	Id              int       `json:"id"              ` //
	TargetId        string    `json:"targetId"        ` //
	TargetUrl       string    `json:"targetUrl"       ` //
	HighNumber      int       `json:"highNumber"      ` //
	MiddleNumber    int       `json:"middleNumber"    ` //
	LowNumber       int       `json:"lowNumber"       ` //
	Progress        int       `json:"progress"        ` //
	IsAlive         int       `json:"isAlive"         ` //
	TargetType      int       `json:"targetType"      ` //
	StartTime       time.Time `json:"startTime"       ` //
	EndTime         time.Time `json:"endTime"         ` //
	CheckTaskId     int       `json:"checkTaskId"     ` //
	TaskConfigId    int       `json:"taskConfigId"    ` //
	HostId          string    `json:"hostId"          ` //
	TargetStatus    int       `json:"targetStatus"    ` //
	InfoNumber      int       `json:"infoNumber"      ` //
	OperatingSystem string    `json:"operatingSystem" ` //
	AttackFaceId    int       `json:"attackFaceId"    ` //
	AssetBaseId     uint      `json:"assetBaseId"     ` // 资产收集asset_collect_base表ID
	RealTargetUrl   string    `json:"realTargetUrl"   ` //
	IsRemote        int       `json:"isRemote"        ` //
	CloudHost       string    `json:"cloudHost"       ` // 云主机
	Area            string    `json:"area"            ` // ip归属区域
	CDN             string    `json:"cDN"             ` //
	AS              string    `json:"aS"              ` //
	OpenPort        string    `json:"openPort"        ` // 开放端口
	AssetIp         string    `json:"assetIp"         ` // 资产ip
}

type TaskCheckvulnerability struct {
	Id                   int       `json:"id"                   ` //
	Name                 string    `json:"name"                 ` //
	Object               string    `json:"object"               ` //
	Type                 string    `json:"type"                 ` //
	RiskLevel            int       `json:"riskLevel"            ` //
	UseImpactValue       string    `json:"useImpactValue"       ` //
	VulId                string    `json:"vulId"                ` //
	Detail               string    `json:"detail"               ` //
	FixSuggest           string    `json:"fixSuggest"           ` //
	RefUrl               string    `json:"refUrl"               ` //
	VulnLocation         string    `json:"vulnLocation"         ` //
	Pocname              string    `json:"pocname"              ` //
	TargetResultId       string    `json:"targetResultId"       ` //
	TargetId             string    `json:"targetId"             ` //
	Status               string    `json:"status"               ` //
	IsRepeat             int       `json:"isRepeat"             ` //
	UpdateTime           time.Time `json:"updateTime"           ` //
	ToolsVulnerabilityId int       `json:"toolsVulnerabilityId" ` //
	AssetIp              string    `json:"assetIp"              ` // 资产ip
}

type ToolsVulnerability struct {
	Id              int       `json:"id"              ` //
	VulId           string    `json:"vulId"           ` // 漏洞唯一标识符
	Source          string    `json:"source"          ` // 漏洞来源
	Cvss            string    `json:"cvss"            ` // cvss级别
	Priority        string    `json:"priority"        ` // 优先级
	VulNum          string    `json:"vulNum"          ` // 漏洞编号（cve/cnvd/cnnvd）
	PublishedTime   string    `json:"publishedTime"   ` // 纰漏时间
	UseImpact       string    `json:"useImpact"       ` // 利用影响
	Disable         int       `json:"disable"         ` // 是否禁用
	VulName         string    `json:"vulName"         ` // 漏洞名称
	VulClass        string    `json:"vulClass"        ` // 漏洞分类
	VulType         string    `json:"vulType"         ` // 漏洞类型
	VulRisk         int       `json:"vulRisk"         ` // 风险等级
	FixSuggest      string    `json:"fixSuggest"      ` // 修复建议
	AffectRange     string    `json:"affectRange"     ` // 影响范围
	Description     string    `json:"description"     ` // 漏洞描述
	VulUid          string    `json:"vulUid"          ` // uuid（漏洞验证时使用）
	AddTime         time.Time `json:"addTime"         ` // 漏洞添加时间
	Hide            int       `json:"hide"            ` // 方案中添加漏洞列表页是否隐藏该漏洞
	VulAnalysis     string    `json:"vulAnalysis"     ` // 漏洞分析
	VulExample      string    `json:"vulExample"      ` // 示例说明
	Component       string    `json:"component"       ` // 受影响的组建名称
	OperatingSystem int       `json:"operatingSystem" ` // 1:windows, 2:linux, 3:unix, 4:其他
}

type TaskCheckresult struct {
	Id           int       `json:"id"           ` //
	FatherId     string    `json:"fatherId"     ` //
	NodeId       string    `json:"nodeId"       ` //
	TargetId     string    `json:"targetId"     ` //
	Location     string    `json:"location"     ` //
	Pocname      string    `json:"pocname"      ` //
	CreateTime   time.Time `json:"createTime"   ` //
	Success      int       `json:"success"      ` //
	Result       string    `json:"result"       ` //
	AdditionalId int       `json:"additionalId" ` //
	AssetIp      string    `json:"assetIp"      ` // 资产ip
}

type TaskAdditional struct {
	Id           int    `json:"id"           ` //
	RequestInfo  string `json:"requestInfo"  ` //
	ResponseInfo string `json:"responseInfo" ` //
}

func main() {
	dsn := "xiaozhi:xiaozhi.4dogs.cn@tcp(192.168.4.131:33306)/qiming_db2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	ctx := context.Background()
	convertTask(ctx, db)
	convertTarget(ctx, db)
	convertTaskVul(ctx, db)
}

func convertTask(ctx context.Context, db *gorm.DB) {
	var checkTaskList []TaskChecktask
	db.Table("task_checktask").Find(&checkTaskList)
	newTaskList := make([]mysqls.TaskTask, 0)
	for _, checkTask := range checkTaskList {
		var checkTargetList []TaskChecktarget
		db.Table("task_checktarget").Where("check_task_id = ?", checkTask.Id).Find(&checkTargetList)
		highNum, middleNum, lowNum, safeNum := 0, 0, 0, 0
		for _, target := range checkTargetList {
			if target.HighNumber > 0 {
				highNum += 1
			} else if target.MiddleNumber > 0 {
				middleNum += 1
			} else if target.LowNumber > 0 {
				lowNum += 1
			} else {
				safeNum += 1
			}
		}

		newTaskList = append(newTaskList, mysqls.TaskTask{
			ID:             checkTask.Id,
			TaskName:       checkTask.TaskName,
			RiskLevel:      checkTask.RiskLevel,
			Status:         checkTask.TaskStatus,
			Weight:         enums.TaskWeightLow,
			IsStats:        enums.TaskIsStatsNo,
			TaskType:       enums.TaskTypeMultipleTask,
			ExecuteType:    enums.TaskExecTypeImmediate,
			TaskTemplateID: 0,
			TargetNum:      checkTask.TargetNumber,
			HigeNum:        highNum,
			MiddleNum:      middleNum,
			LowNum:         lowNum,
			SafeNum:        safeNum,
			UserID:         checkTask.UserId,
			Pid:            0,
			CreateTime:     checkTask.CreateTime,
			UpdateTime:     checkTask.UpdateTime,
		})
	}
	bytes, err := json.Marshal(newTaskList)
	if err != nil {
		panic(err)
	}
	filename := "task.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func convertTarget(ctx context.Context, db *gorm.DB) {
	var checkTargetList []TaskChecktarget
	db.Table("task_checktarget").Find(&checkTargetList)
	fmt.Println(len(checkTargetList))
	newTargetList := make([]mysqls.TaskTarget, 0)
	for _, checkTarget := range checkTargetList {
		riskLevel := enums.TargetRiskLowNoFound
		if checkTarget.HighNumber > 0 {
			riskLevel = enums.TargetRiskHigh
		} else if checkTarget.MiddleNumber > 0 {
			riskLevel = enums.TargetRiskMid
		} else if checkTarget.LowNumber > 0 {
			riskLevel = enums.TargetRiskLow
		}
		newTargetList = append(newTargetList, mysqls.TaskTarget{
			ID:               checkTarget.Id,
			TaskID:           checkTarget.CheckTaskId,
			TargetURL:        checkTarget.TargetUrl,
			Status:           checkTarget.TargetStatus,
			Weight:           enums.TaskWeightLow,
			RiskLevel:        riskLevel,
			OpSys:            checkTarget.OperatingSystem,
			IsAlive:          checkTarget.IsAlive,
			TargetType:       checkTarget.TargetType,
			TaskTemplateID:   0,
			TaskTemplateJSON: "",
			IsRemoteSession:  0,
			UserID:           1,
			CreateTime:       checkTarget.StartTime,
			UpdateTime:       checkTarget.EndTime,
			EndTime:          checkTarget.EndTime,
			UseScore:         0,
			IsScore:          0,
			ExtendField:      "",
		})
	}
	bytes, err := json.Marshal(newTargetList)
	if err != nil {
		panic(err)
	}
	filename := "target.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func convertTaskVul(ctx context.Context, db *gorm.DB) {
	var checkVulList []TaskCheckvulnerability
	db.Table("task_checkvulnerability").Where(`target_id != ""`).Find(&checkVulList)
	newVulList := make([]mysqls.TaskVul, 0)
	for index, checkVul := range checkVulList {
		if index%100 == 0 {
			fmt.Println("获取漏洞结果个数: ", index)
		}

		var (
			checkTarget        TaskChecktarget
			toolsVulnerability ToolsVulnerability
			taskCheckresult    TaskCheckresult
		)
		db.Table("task_checktarget").Where("target_id = ?", checkVul.TargetId).First(&checkTarget)
		if checkTarget.Id == 0 {
			continue
		}
		db.Table("tools_vulnerability").Where("id = ?", checkVul.ToolsVulnerabilityId).First(&toolsVulnerability)
		if toolsVulnerability.Id == 0 {
			continue
		}
		db.Table("task_checkresult").Where("id = ?", checkVul.TargetResultId).First(&taskCheckresult)
		var (
			result string
			verMsg string
		)
		if taskCheckresult.Id != 0 {
			result = taskCheckresult.Result

			var taskAdditional TaskAdditional
			verMsgList := make([]map[string]string, 0)
			db.Table("task_additional").Where("id = ?", taskCheckresult.AdditionalId).First(&taskAdditional)
			fmt.Println(taskAdditional.Id)
			if taskAdditional.Id != 0 {
				verMsgList = getPackage(taskAdditional)
			}
			//fmt.Println(verMsgList)
			//time.Sleep(1 * time.Second)
			verMsgByte, _ := json.Marshal(verMsgList)
			verMsg = string(verMsgByte)
		}

		newVulList = append(newVulList, mysqls.TaskVul{
			ID:             checkVul.Id,
			DataType:       enums.VulDataTypOne,
			TaskID:         checkTarget.CheckTaskId,
			TargetID:       checkTarget.Id,
			TargetUrl:      checkTarget.TargetUrl,
			Pocname:        checkVul.Pocname,
			Name:           toolsVulnerability.VulName,
			Class:          enums.VulLibrariesClassOther,
			Type:           enums.VulLibrariesClassOther,
			Risk:           checkVul.RiskLevel,
			Location:       checkVul.VulnLocation,
			Status:         enums.VulStatusVerifySuccess,
			TestStatus:     enums.VulTestStatusNotTest,
			ExploitImpact:  checkVul.UseImpactValue,
			VulID:          "",
			Description:    toolsVulnerability.Description,
			FixSuggest:     toolsVulnerability.FixSuggest,
			PublishedTime:  toolsVulnerability.PublishedTime,
			AffectRange:    toolsVulnerability.AffectRange,
			TargetResultID: 0,
			VulNumber:      toolsVulnerability.VulNum,
			VulAddress:     checkVul.VulnLocation,
			RefUrl:         "",
			Cvss:           toolsVulnerability.Cvss,
			VulResult:      result,
			VulParam:       "",
			VerMsg:         verMsg,
			DecisionVulId:  "",
			Snapshot:       "",
			CreateTime:     checkVul.UpdateTime,
			UpdateTime:     checkVul.UpdateTime,
		})
	}
	bytes, err := json.Marshal(newVulList)
	if err != nil {
		panic(err)
	}
	filename := "taskVul.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func getPackage(taskAdditional TaskAdditional) []map[string]string {
	requestInfo := taskAdditional.RequestInfo
	responseInfo := taskAdditional.ResponseInfo
	requestInfo = strings.ReplaceAll(requestInfo, "\\r\\n", "\\n")
	responseInfo = strings.ReplaceAll(responseInfo, "\\r\\n", "\\n")
	requestInfoSlice := strings.Split(requestInfo, "\\n")
	responseInfoSlice := strings.Split(responseInfo, "\\n")
	reqTempSlice := make([]string, 0)
	for _, req := range requestInfoSlice {
		reqTempSlice = append(reqTempSlice, strings.TrimSpace(req))
	}
	reqTempStr := strings.Join(reqTempSlice, "\n")
	resTempSlice := make([]string, 0)
	for _, res := range responseInfoSlice {
		resTempSlice = append(resTempSlice, strings.TrimSpace(res))
	}
	resTempStr := strings.Join(resTempSlice, "\n")

	requestInfoSlice2 := make([]string, 0)
	responseInfoSlice2 := make([]string, 0)
	if reqTempStr != "" && !strings.HasPrefix(reqTempStr, "[") {
		requestInfoSlice2 = append(requestInfoSlice2, reqTempStr)
	}
	if resTempStr != "" {
		responseInfoSlice2 = append(responseInfoSlice2, resTempStr)
	}

	packetSlice := make([]map[string]string, 0)
	index := 0
	// 如果request为空，response不为空时，单独返回response
	if len(requestInfoSlice2) == 0 && len(responseInfoSlice2) != 0 {
		// 单独返回response相关信息
		for _, res := range responseInfoSlice2 {
			if res == "" {
				continue
			}
			index += 1
			resBs64Encode := base64.StdEncoding.EncodeToString([]byte(res))
			resBs64Decode, err := base64.StdEncoding.DecodeString(resBs64Encode)
			if err != nil {
				continue
			}
			resBs64DecodeStr := string(resBs64Decode)
			resBs64DecodeStr = strings.ReplaceAll(resBs64DecodeStr, "            ", "")
			packetSlice = append(packetSlice, map[string]string{
				"response": resBs64DecodeStr,
			})
		}
	} else if len(requestInfoSlice2) != 0 { // request不为空，response可能为空，可能不为空，为空时需判断一下
		for i, v := range requestInfoSlice2 {
			if v == "" {
				continue
			}
			index += 1
			req := v
			reqBs64Encode := base64.StdEncoding.EncodeToString([]byte(req))
			reqBs64Decode, err := base64.StdEncoding.DecodeString(reqBs64Encode)
			if err != nil {
				continue
			}
			reqBs64DecodeStr := string(reqBs64Decode)

			var res string
			var resBs64DecodeStr string
			if len(responseInfoSlice2) != 0 {
				res = responseInfoSlice2[i]
				resBs64Encode := base64.StdEncoding.EncodeToString([]byte(res))
				resBs64Decode, err := base64.StdEncoding.DecodeString(resBs64Encode)
				if err != nil {
					continue
				}
				resBs64DecodeStr = string(resBs64Decode)
			}
			reqBs64DecodeStr = strings.ReplaceAll(reqBs64DecodeStr, "            ", "")
			resBs64DecodeStr = strings.ReplaceAll(resBs64DecodeStr, "            ", "")
			packetSlice = append(packetSlice, map[string]string{
				"request":  reqBs64DecodeStr,
				"response": resBs64DecodeStr,
			})
		}
	}
	return packetSlice
}
