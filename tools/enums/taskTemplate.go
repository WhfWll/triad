package enums

// 状态
const (
	TaskTemplateStatusSuccess = "valid"   //数据状态正常
	TaskTemplateStatusDel     = "deleted" //数据状态不正常
)

// 来源
const (
	TaskTemplateSourceUser   = 1 //用户创建
	TaskTemplateSourceSystem = 2 //系统创建
)

// 是否默认
const (
	TaskTemplateIsDefaultY = 1
	TaskTemplateIsDefaultN = 2
)

var TaskTemplate taskTemplate

type taskTemplate struct {
}

func (a *taskTemplate) GetSourceEnum(source uint) string {
	sourceEnum := map[uint]string{
		TaskTemplateSourceUser:   "用户创建",
		TaskTemplateSourceSystem: "系统创建",
	}

	if res, ok := sourceEnum[source]; ok {
		return res
	}
	return ""
}

func (t *taskTemplate) IsDefaultStr(k uint) string {
	enum := map[uint]string{
		TaskTemplateIsDefaultY: "是",
		TaskTemplateIsDefaultN: "否",
	}
	value, ok := enum[k]
	if ok {
		return value
	}
	return ""
}
