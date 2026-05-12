package enums

const (
	XrayStatusWait    = 1
	XrayStatusRunning = 2
	XrayStatusDone    = 3
)

var XrayEnum xray

type xray struct {
}

func (x *xray) AllStatusEnum() map[int]string {
	res := map[int]string{
		XrayStatusWait:    "待运行",
		XrayStatusRunning: "运行中",
		XrayStatusDone:    "已完成",
	}
	return res
}
func (x *xray) GetStatusEnum(status int) string {
	enum := x.AllStatusEnum()
	if res, ok := enum[status]; ok {
		return res
	}
	return ""
}
