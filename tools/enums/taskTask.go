package enums

import "strconv"

var TaskTaskEnum task

type task struct {
}

const (
	TaskExChangeName   string = "decision-task"
	TaskRoutingKeyName string = "decision-task.#"
	TaskCorrelationID  string = "1"
	TaskReplyTo        string = "3"
)

/************ 执行方式 **************/
const (
	TaskExecTypeImmediate = 1 //即时执行
	TaskExecTypeTiming    = 2 //定时执行
	TaskExecTypeCycle     = 3 //周期执行
)

func (t task) AllExecTypeEnum() map[int]string {
	enum := map[int]string{
		TaskExecTypeImmediate: "即时任务",
		TaskExecTypeTiming:    "定时任务",
		TaskExecTypeCycle:     "周期任务",
	}
	return enum
}
func (t task) ExecTypeEnum(execType int) string {
	enum := t.AllExecTypeEnum()
	if res, ok := enum[execType]; ok {
		return res
	}
	return ""
}

/************ 周期执行方式 **************/
const (
	TaskExecTypeCycleTypeWeek  = 1 // 周期执行方式
	TaskExecTypeCycleTypeMonth = 2 // 周期执行方式
)

func (t task) AllExecCycleTypeEnum() map[int]string {
	enum := map[int]string{
		TaskExecTypeCycleTypeWeek:  "每周一次",
		TaskExecTypeCycleTypeMonth: "每月一次",
	}
	return enum
}
func (t task) ExecCycleTypeEnum(execType int) string {
	enum := t.AllExecCycleTypeEnum()
	if res, ok := enum[execType]; ok {
		return res
	}
	return ""
}

/************ 周期执行值 每周执行一次 **************/
const (
	TaskExecTypeCycleTypeWeek1 = 1
	TaskExecTypeCycleTypeWeek2 = 2
	TaskExecTypeCycleTypeWeek3 = 3
	TaskExecTypeCycleTypeWeek4 = 4
	TaskExecTypeCycleTypeWeek5 = 5
	TaskExecTypeCycleTypeWeek6 = 6
	TaskExecTypeCycleTypeWeek7 = 7
)

func (t task) AllExecCycleTypeWeekValueEnum() map[int]string {
	enum := map[int]string{
		TaskExecTypeCycleTypeWeek1: "周一",
		TaskExecTypeCycleTypeWeek2: "周二",
		TaskExecTypeCycleTypeWeek3: "周三",
		TaskExecTypeCycleTypeWeek4: "周四",
		TaskExecTypeCycleTypeWeek5: "周五",
		TaskExecTypeCycleTypeWeek6: "周六",
		TaskExecTypeCycleTypeWeek7: "周日",
	}
	return enum
}
func (t task) ExecCycleTypeWeekValueEnum(execType int) string {
	enum := t.AllExecCycleTypeWeekValueEnum()
	if res, ok := enum[execType]; ok {
		return res
	}
	return ""
}

/************ 周期执行值 每月执行一次 **************/
func (t task) AllExecCycleTypeMonthValueEnum() map[int]string {
	enum := make(map[int]string)
	for i := 1; i <= 31; i++ {
		enum[i] = strconv.Itoa(i) + "号"
	}
	return enum
}
func (t task) ExecCycleTypeMonthValueEnum(execType int) string {
	enum := t.AllExecCycleTypeMonthValueEnum()
	if res, ok := enum[execType]; ok {
		return res
	}
	return ""
}

/************ 运行时间段 **************/
const (
	_ = iota
	TaskRuntimePeriod1
	TaskRuntimePeriod2
	TaskRuntimePeriod3
	TaskRuntimePeriod4
	TaskRuntimePeriod5
	TaskRuntimePeriod6
	TaskRuntimePeriod7
	TaskRuntimePeriod8
	TaskRuntimePeriod9
	TaskRuntimePeriod10
	TaskRuntimePeriod11
	TaskRuntimePeriod12
	TaskRuntimePeriod13
	TaskRuntimePeriod14
	TaskRuntimePeriod15
	TaskRuntimePeriod16
	TaskRuntimePeriod17
	TaskRuntimePeriod18
	TaskRuntimePeriod19
	TaskRuntimePeriod20
	TaskRuntimePeriod21
	TaskRuntimePeriod22
	TaskRuntimePeriod23
	TaskRuntimePeriod24
)

func (t task) AllRuntimePeriodEnum() map[int]string {
	enum := map[int]string{
		TaskRuntimePeriod1:  "00:00-01:00",
		TaskRuntimePeriod2:  "01:00-02:00",
		TaskRuntimePeriod3:  "02:00-03:00",
		TaskRuntimePeriod4:  "03:00-04:00",
		TaskRuntimePeriod5:  "04:00-05:00",
		TaskRuntimePeriod6:  "05:00-06:00",
		TaskRuntimePeriod7:  "06:00-07:00",
		TaskRuntimePeriod8:  "07:00-08:00",
		TaskRuntimePeriod9:  "08:00-09:00",
		TaskRuntimePeriod10: "09:00-10:00",
		TaskRuntimePeriod11: "10:00-11:00",
		TaskRuntimePeriod12: "11:00-12:00",
		TaskRuntimePeriod13: "12:00-13:00",
		TaskRuntimePeriod14: "13:00-14:00",
		TaskRuntimePeriod15: "14:00-15:00",
		TaskRuntimePeriod16: "15:00-16:00",
		TaskRuntimePeriod17: "16:00-17:00",
		TaskRuntimePeriod18: "17:00-18:00",
		TaskRuntimePeriod19: "18:00-19:00",
		TaskRuntimePeriod20: "19:00-20:00",
		TaskRuntimePeriod21: "20:00-21:00",
		TaskRuntimePeriod22: "21:00-22:00",
		TaskRuntimePeriod23: "22:00-23:00",
		TaskRuntimePeriod24: "23:00-00:00",
	}

	return enum
}
func (t task) RuntimePeriodEnum(runtimePeriod int) string {
	enum := t.AllRuntimePeriodEnum()

	if res, ok := enum[runtimePeriod]; ok {
		return res
	}
	return ""
}

/************ 状态 **************/
const (
	_                 = iota
	TaskStatusTrigger //待触发
	TaskStatusBegin   //待执行
	TaskStatusRunning //运行中
	TaskStatusFinish  //已结束
	TaskStatusPausing //暂停中
)

func (t task) AllStatusEnum() map[int]string {
	enum := map[int]string{
		TaskStatusTrigger: "待触发",
		TaskStatusBegin:   "待执行",
		TaskStatusRunning: "运行中",
		TaskStatusFinish:  "已完成",
		TaskStatusPausing: "暂停中",
	}
	return enum
}
func (t task) StatusEnum(status int) string {
	enum := t.AllStatusEnum()
	if res, ok := enum[status]; ok {
		return res
	}
	return ""
}

/************* 任务类型 **************/
const (
	TaskTypeMultipleTask        = 1 // 综合任务
	TaskTypeAppSecDynamic       = 2 // 应用安全-动态扫描
	TaskTypeAppSecApp           = 3 // 应用安全-专项应用检测
	TaskTypeDataSecDB           = 4 // 数据安全-数据库基线检查
	TaskTypeDataSecSensitive    = 5 // 数据安全-敏感数据发现
)

func (t task) AllTypeEnum() map[int]string {
	enum := map[int]string{
		TaskTypeMultipleTask:     "综合任务",
		TaskTypeAppSecDynamic:    "应用安全-动态扫描",
		TaskTypeAppSecApp:        "应用安全-专项应用检测",
		TaskTypeDataSecDB:        "数据安全-数据库检查",
		TaskTypeDataSecSensitive: "数据安全-敏感数据发现",
	}
	return enum
}
func (t task) TypeEnum(types int) string {
	enum := t.AllTypeEnum()

	if res, ok := enum[types]; ok {
		return res
	}
	return ""
}

const (
	TaskPermeateModeCommon = "common" // 通用渗透
)

/************** 风险等级 ***************/
const (
	TaskRiskHigh = 1 // 高危
	TaskRiskMid  = 2 // 中危
	TaskRiskLow  = 3 // 低危
	TaskRiskSafe = 4 // 安全
)

func (t task) AllRiskEnum() map[int]string {
	enum := map[int]string{
		TaskRiskHigh: "高危",
		TaskRiskMid:  "中危",
		TaskRiskLow:  "低危",
		TaskRiskSafe: "安全",
	}
	return enum
}
func (t task) RiskEnum(risk int) string {
	enum := t.AllRiskEnum()
	if res, ok := enum[risk]; ok {
		return res
	}
	return ""
}

const (
	_              = iota
	TaskIsStatsNo  //不进行统计
	TaskIsStatsYes //进行统计
)

const (
	TaskWeightLow  = 1 // 低
	TaskWeightMid  = 2 // 中
	TaskWeightHigh = 3 // 高
)

func (t task) AllTaskWeightEnum() map[int]string {
	enum := map[int]string{
		TaskWeightLow:  "低",
		TaskWeightMid:  "中",
		TaskWeightHigh: "高",
	}
	return enum
}

func (t task) GetTaskWeightEnum(weight int) string {
	enum := t.AllTaskWeightEnum()
	if res, ok := enum[weight]; ok {
		return res
	}
	return ""
}

func (t task) GetTaskWeightEnumArray() interface{} {
	result := []struct {
		Value   int    `json:"value"`
		Label   string `json:"label"`
		Default int    `json:"default"`
	}{{
		Value:   TaskWeightLow,
		Label:   t.GetTaskWeightEnum(TaskWeightLow),
		Default: 0,
	}, {
		Value:   TaskWeightMid,
		Label:   t.GetTaskWeightEnum(TaskWeightMid),
		Default: 1,
	}, {
		Value:   TaskWeightHigh,
		Label:   t.GetTaskWeightEnum(TaskWeightHigh),
		Default: 0,
	}}
	return result
}

const (
	TaskTestIntensityOne   = 1
	TaskTestIntensityTwo   = 2
	TaskTestIntensityThree = 3
	TaskTestIntensityFour  = 4
	TaskTestIntensityFive  = 5
)

func (t task) AllTaskTestIntensityEnum() map[int]string {
	enum := map[int]string{
		TaskTestIntensityOne:   "1",
		TaskTestIntensityTwo:   "2",
		TaskTestIntensityThree: "3",
		TaskTestIntensityFour:  "4",
		TaskTestIntensityFive:  "5",
	}
	return enum
}

func (t task) GetTaskTestIntensityEnum(intensity int) string {
	enum := t.AllTaskTestIntensityEnum()
	if res, ok := enum[intensity]; ok {
		return res
	}
	return ""
}

func (t task) GetTaskTestIntensityEnumArray() interface{} {
	result := []struct {
		Value   int    `json:"value"`
		Label   string `json:"label"`
		Default int    `json:"default"`
	}{{
		Value:   TaskTestIntensityOne,
		Label:   t.GetTaskTestIntensityEnum(TaskTestIntensityOne),
		Default: 1,
	}, {
		Value:   TaskTestIntensityTwo,
		Label:   t.GetTaskTestIntensityEnum(TaskTestIntensityTwo),
		Default: 0,
	}, {
		Value:   TaskTestIntensityThree,
		Label:   t.GetTaskTestIntensityEnum(TaskTestIntensityThree),
		Default: 0,
	}, {
		Value:   TaskTestIntensityFour,
		Label:   t.GetTaskTestIntensityEnum(TaskTestIntensityFour),
		Default: 0,
	}, {
		Value:   TaskTestIntensityFive,
		Label:   t.GetTaskTestIntensityEnum(TaskTestIntensityFive),
		Default: 0,
	}}
	return result
}

const (
	VulScriptExploitImpactNoImpact         = 2 //无影响
	VulScriptExploitImpactServerSlow       = 3 //运行缓慢
	VulScriptExploitImpactRefuseServer     = 4 //拒绝服务漏洞
	VulScriptExploitImpactServiceBreakdown = 5 //服务崩溃
	VulScriptExploitImpactOutAge           = 6 //系统奔溃
)

func (t task) TestIntensityToExploitImpact(testIntensity int) string {
	var optionsScriptExploitImpact = map[int]string{
		TaskTestIntensityOne:   "0,2,3,4,5,6", //执行利用影响中的所有漏洞、
		TaskTestIntensityTwo:   "0,2,3,4,5",   //不执行系统崩溃漏洞
		TaskTestIntensityThree: "0,2,3,4",     //不执行系统崩溃与服务崩溃漏洞
		TaskTestIntensityFour:  "0,2,3",       //不执行系统崩溃/服务崩溃/拒绝服务漏洞
		TaskTestIntensityFive:  "0,2",         //不执行系统崩溃/服务崩溃/拒绝服务漏洞/运行缓慢漏洞
	}
	if res, ok := optionsScriptExploitImpact[testIntensity]; ok {
		return res
	}
	return ""
}
