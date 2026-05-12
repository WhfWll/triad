package enums

var TaskTaskLogEnum taskTaskLogEnum

type taskTaskLogEnum struct {
}

const (
	TaskLogTypeTarget          = 1
	TaskLogTypeAddService      = 2
	TaskLogTypeAddFileUpload   = 3
	TaskLogTypeAddLoginEntry   = 4
	TaskLogTypeAddLoginVoucher = 5
	TaskLogTypeAddVul          = 6
	TaskLogTypePassiveFlow     = 7
	TaskLogTypeSensitivePath   = 8
	TaskLogTypeTargetVerify    = 9
)

func (t taskTaskLogEnum) TypeStr(tasklog_type int) string {
	enum := map[int]string{
		TaskLogTypeTarget:          "目标测试",
		TaskLogTypeAddService:      "新增服务测试",
		TaskLogTypeAddFileUpload:   "新增文件上传测试",
		TaskLogTypeAddLoginEntry:   "新增登录入口测试",
		TaskLogTypeAddLoginVoucher: "新增登录凭证测试",
		TaskLogTypeAddVul:          "新增漏洞测试",
		TaskLogTypePassiveFlow:     "被动流量测试",
		TaskLogTypeSensitivePath:   "敏感路径",
		TaskLogTypeTargetVerify:    "目标验证",
	}
	value, ok := enum[tasklog_type]
	if ok {
		return value
	}
	return ""
}
