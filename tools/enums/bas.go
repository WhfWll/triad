package enums

const (
	BasTaskStatusWait    = 1
	BasTaskStatusRunning = 2
	BasTaskStatusDone    = 3
)

var BasEnum bas

type bas struct {
}

func (b *bas) AllStatusEnum() map[int]string {
	res := map[int]string{
		BasTaskStatusWait:    "待开始",
		BasTaskStatusRunning: "运行中",
		BasTaskStatusDone:    "已完成",
	}
	return res
}

func (b *bas) BasTaskStatusEnum(status int) string {
	enum := b.AllStatusEnum()
	value, ok := enum[status]
	if ok {
		return value
	}
	return ""
}

func (b *bas) GetStatus(status int) string {
	return b.AllStatusEnum()[status]
}

func (b *bas) GetBasTaskStatusEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: BasTaskStatusWait,
		Label: b.BasTaskStatusEnum(BasTaskStatusWait),
	}, {
		Value: BasTaskStatusRunning,
		Label: b.BasTaskStatusEnum(BasTaskStatusRunning),
	}, {
		Value: BasTaskStatusDone,
		Label: b.BasTaskStatusEnum(BasTaskStatusDone),
	}}
	return result
}

const (
	BasRiskLevelHigh   = 1 //高危
	BasRiskLevelMiddle = 2 //中危
	BasRiskLevelLow    = 3 //低危
	BasRiskLevelSafe   = 4 //安全
)

func (b *bas) BasRiskLevelEnum(level int) string {
	enum := map[int]string{
		BasRiskLevelHigh:   "高危",
		BasRiskLevelMiddle: "中危",
		BasRiskLevelLow:    "低危",
		BasRiskLevelSafe:   "安全",
	}
	value, ok := enum[level]
	if ok {
		return value
	}
	return ""
}

func (b *bas) GetBasRiskLevelEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: BasRiskLevelHigh,
		Label: b.BasRiskLevelEnum(BasRiskLevelHigh),
	}, {
		Value: BasRiskLevelMiddle,
		Label: b.BasRiskLevelEnum(BasRiskLevelMiddle),
	}, {
		Value: BasRiskLevelLow,
		Label: b.BasRiskLevelEnum(BasRiskLevelLow),
	}, {
		Value: BasRiskLevelSafe,
		Label: b.BasRiskLevelEnum(BasRiskLevelSafe),
	}}
	return result
}

const (
	BasVulStatusFail    = 1 //失败
	BasVulStatusSuccess = 2 //成功
)

func (b *bas) BasVulStatusEnum(status int) string {
	enum := map[int]string{
		BasVulStatusFail:    "失败",
		BasVulStatusSuccess: "成功",
	}
	value, ok := enum[status]
	if ok {
		return value
	}
	return ""
}

func (b *bas) GetBasVulStatusEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: BasVulStatusFail,
		Label: b.BasVulStatusEnum(BasVulStatusFail),
	}, {
		Value: BasVulStatusSuccess,
		Label: b.BasVulStatusEnum(BasVulStatusSuccess),
	}}
	return result
}

const (
	BasNodeOnlineStatusOnline  = 1 //在线
	BasNodeOnlineStatusOffline = 2 //离线
	BasNodeStatusEnable        = 1 //启用
	BasNodeStatusDisable       = 2 //禁用
	BasRuleFileName            = "bas_rule_*"
)

// 节点在线状态
func (b *bas) AllBaseNodeOnlineStatus() map[int]string {
	return map[int]string{
		BasNodeOnlineStatusOnline:  "在线",
		BasNodeOnlineStatusOffline: "离线",
	}
}
func (b *bas) GetBaseNodeOnlineStatus(onlineStatus int) string {
	return b.AllBaseNodeOnlineStatus()[onlineStatus]
}

// 节点状态
func (b *bas) AllBaseNodeStatus() map[int]string {
	return map[int]string{
		BasNodeStatusEnable:  "启用",
		BasNodeStatusDisable: "禁用",
	}
}
func (b *bas) GetBaseNodeStatus(status int) string {
	return b.AllBaseNodeStatus()[status]
}
