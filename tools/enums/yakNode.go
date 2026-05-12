package enums

const (
	YakNodeStatusOffline = 0
	YakNodeStatusOnline  = 1
)

const (
	YakNodeIsDisableN = 0
	YakNodeIsDisableY = 1
)

var YakNode yakNode

type yakNode struct {
}

// 在线状态
func (y *yakNode) AllStatusEnum() map[int]string {
	return map[int]string{
		YakNodeStatusOffline: "离线",
		YakNodeStatusOnline:  "在线",
	}
}

func (y *yakNode) GetStatusEnum(status int) string {
	return y.AllStatusEnum()[status]
}

// 禁用状态
func (y *yakNode) AllIsDisableEnum() map[int]string {
	return map[int]string{
		YakNodeIsDisableN: "启用",
		YakNodeIsDisableY: "禁用",
	}
}

func (y *yakNode) GetIsDisableEnum(disable int) string {
	return y.AllIsDisableEnum()[disable]
}
