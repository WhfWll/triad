package enums

var WifiApInfoEnum wifiApInfoEnum

type wifiApInfoEnum struct {
}

// wifi 加密方式枚举获取
func (*wifiApInfoEnum) GetSsidCryptsetEnum(encrypt int) string {
	if encrypt == 0 {
		return "NONE"
	} else if (encrypt&(1<<20)) != 0 && (encrypt&(1<<21)) != 0 {
		return "WPA2"
	} else if (encrypt&(1<<7)) != 0 && (encrypt&(1<<20)) != 0 {
		return "WPA"
	} else if (encrypt & (1 << 1)) != 0 {
		return "WEP"
	} else {
		return "UNKNOWN"
	}
}

// wifi 信号通道枚举
func (*wifiApInfoEnum) GetHandlerEnum(channel int) string {
	if channel < 36 {
		return "2.4GHz"
	}
	return "5GHz"
}
