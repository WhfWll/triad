package enums

type Flow struct{}

const (
	FlowTaskStatusDisable = 1 //禁用
	FlowTaskStatusWait    = 2 //待开始
	FlowTaskStatusExec    = 3 //运行中
	FlowTaskStatusDone    = 4 //已结束
)

func (f Flow) FlowTaskStatusEnum(status int) string {
	enum := map[int]string{
		FlowTaskStatusDisable: "禁用",
		FlowTaskStatusWait:    "待开始",
		FlowTaskStatusExec:    "运行中",
		FlowTaskStatusDone:    "已结束",
	}
	value, ok := enum[status]
	if ok {
		return value
	}
	return ""
}

func (f Flow) GetFlowTaskStatusEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: FlowTaskStatusDisable,
		Label: f.FlowTaskStatusEnum(FlowTaskStatusDisable),
	}, {
		Value: FlowTaskStatusWait,
		Label: f.FlowTaskStatusEnum(FlowTaskStatusWait),
	}, {
		Value: FlowTaskStatusExec,
		Label: f.FlowTaskStatusEnum(FlowTaskStatusExec),
	}, {
		Value: FlowTaskStatusDone,
		Label: f.FlowTaskStatusEnum(FlowTaskStatusDone),
	}}
	return result
}

const (
	FlowTargetStatusOne = 1 //禁用
	FlowTargetStatusTwo = 2 //正常
)
const (
	FlowRiskLevelFatal = 1 //致命
	FlowRiskLevelHigh  = 2 //高危
	FlowRiskLevelMid   = 3 //中危
	FlowRiskLevelLow   = 4 //低危
	FlowRiskLevelSafe  = 5 //安全
)

func (f Flow) FlowRiskLevelEnum(risk int) string {
	enum := map[int]string{
		FlowRiskLevelFatal: "致命",
		FlowRiskLevelHigh:  "高危",
		FlowRiskLevelMid:   "中危",
		FlowRiskLevelLow:   "低危",
		FlowRiskLevelSafe:  "安全",
	}
	value, ok := enum[risk]
	if ok {
		return value
	}
	return ""
}

func (f Flow) GetFlowRiskLevelEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: FlowRiskLevelFatal,
		Label: f.FlowRiskLevelEnum(FlowRiskLevelFatal),
	}, {
		Value: FlowRiskLevelHigh,
		Label: f.FlowRiskLevelEnum(FlowRiskLevelHigh),
	}, {
		Value: FlowRiskLevelMid,
		Label: f.FlowRiskLevelEnum(FlowRiskLevelMid),
	}, {
		Value: FlowRiskLevelLow,
		Label: f.FlowRiskLevelEnum(FlowRiskLevelLow),
	}, {
		Value: FlowRiskLevelSafe,
		Label: f.FlowRiskLevelEnum(FlowRiskLevelSafe),
	}}
	return result
}

func (f Flow) FlowRiskLevelToTargetRiskLevelName(risk int) string {
	enum := map[int]string{
		FlowRiskLevelFatal: "高危",
		FlowRiskLevelHigh:  "高危",
		FlowRiskLevelMid:   "中危",
		FlowRiskLevelLow:   "低危",
		FlowRiskLevelSafe:  "安全",
	}
	value, ok := enum[risk]
	if ok {
		return value
	}
	return ""
}

const (
	FlowTaskExpireTimeOneHour        = 3600      //1小时
	FlowTaskExpireTimeTwoHour        = 7200      //2小时
	FlowTaskExpireTimeThreeHour      = 10800     //3小时
	FlowTaskExpireTimeSixHour        = 21600     //6小时
	FlowTaskExpireTimeTwelveHour     = 43200     //12小时
	FlowTaskExpireTimeTwentyFourHour = 86400     //24小时
	FlowTaskExpireTimeFortyHour      = 144000    //40小时
	FlowTaskExpireTimeSeventyTwoHour = 259200    //72小时
	FlowTaskExpireTimeNoLimit        = 315360000 //不限,默认10年
)

func (f Flow) FlowTaskExpireTimeEnum(expireTime int) string {
	enum := map[int]string{
		FlowTaskExpireTimeOneHour:        "1小时",
		FlowTaskExpireTimeTwoHour:        "2小时",
		FlowTaskExpireTimeThreeHour:      "3小时",
		FlowTaskExpireTimeSixHour:        "6小时",
		FlowTaskExpireTimeTwelveHour:     "12小时",
		FlowTaskExpireTimeTwentyFourHour: "24小时",
		FlowTaskExpireTimeFortyHour:      "40小时",
		FlowTaskExpireTimeSeventyTwoHour: "72小时",
		FlowTaskExpireTimeNoLimit:        "不限",
	}
	value, ok := enum[expireTime]
	if ok {
		return value
	}
	return ""
}

func (f Flow) GetFlowTaskExpireTimeEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: FlowTaskExpireTimeOneHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeOneHour),
	}, {
		Value: FlowTaskExpireTimeTwoHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeTwoHour),
	}, {
		Value: FlowTaskExpireTimeThreeHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeThreeHour),
	}, {
		Value: FlowTaskExpireTimeSixHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeSixHour),
	}, {
		Value: FlowTaskExpireTimeTwelveHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeTwelveHour),
	}, {
		Value: FlowTaskExpireTimeTwentyFourHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeTwentyFourHour),
	}, {
		Value: FlowTaskExpireTimeFortyHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeFortyHour),
	}, {
		Value: FlowTaskExpireTimeSeventyTwoHour,
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeSeventyTwoHour),
	}, {
		Value: FlowTaskExpireTimeNoLimit, // 默认10年
		Label: f.FlowTaskExpireTimeEnum(FlowTaskExpireTimeNoLimit),
	}}
	return result
}
