package enums

const (
	BurpsuiteStatusWait    = 1
	BurpsuiteStatusRunning = 2
	BurpsuiteStatusDone    = 3

	BurpsuiteRiskHigh int = 2 //高危
	BurpsuiteRiskMid  int = 3 //中危
	BurpsuiteRiskLow  int = 4 //低危
	BurpsuiteRiskInfo int = 5 //信息
)

var BurpsuiteEnum Burpsuite

type Burpsuite struct {
}

func (x *Burpsuite) AllStatusEnum() map[int]string {
	res := map[int]string{
		BurpsuiteStatusWait:    "待运行",
		BurpsuiteStatusRunning: "运行中",
		BurpsuiteStatusDone:    "已完成",
	}
	return res
}
func (x *Burpsuite) GetStatusEnum(status int) string {
	enum := x.AllStatusEnum()
	if res, ok := enum[status]; ok {
		return res
	}
	return ""
}

func (x *Burpsuite) AllRiskEnum() map[int]string {
	enum := map[int]string{
		BurpsuiteRiskHigh: "高危",
		BurpsuiteRiskMid:  "中危",
		BurpsuiteRiskLow:  "低危",
		BurpsuiteRiskInfo: "信息",
	}
	return enum
}
func (x *Burpsuite) GetRiskEnum(risk int) string {
	enum := x.AllRiskEnum()
	if v, ok := enum[risk]; ok {
		return v
	}
	return ""
}
