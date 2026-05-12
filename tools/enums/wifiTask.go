package enums

const (
	WifiTaskStatusIni = 0
	//WifiTaskStatusCollect          = 1
	WifiTaskStatusPassExplode      = 1
	WifiTaskStatusPassExplodeFalse = 2
	WifiTaskStatusPassExplodeTrue  = 3

	WifiTaskCryptsetNone    = 1
	WifiTaskCryptsetWp2     = 4
	WifiTaskCryptsetWpa     = 3
	WifiTaskCryptsetWep     = 2
	WifiTaskCryptsetUnknown = 0

	WifiTaskCarrier1 = 1
	WifiTaskCarrier2 = 2
	WifiTaskCarrier3 = 3
	WifiTaskCarrier4 = 4
	WifiTaskCarrier5 = 5
	WifiTaskCarrier6 = 6
	WifiTaskCarrier7 = 7
	WifiTaskCarrier8 = 8
)

var WifiTaskEnum wifiTask

type wifiTask struct {
}

func (w *wifiTask) AllStatus() map[int]string {
	res := map[int]string{
		WifiTaskStatusIni: "初始状态，等待开始检测",
		//WifiTaskStatusCollect:          "开始进行报文收集",
		WifiTaskStatusPassExplode:      "正在进行密码爆破",
		WifiTaskStatusPassExplodeFalse: "密码爆破失败",
		WifiTaskStatusPassExplodeTrue:  "密码爆破成功",
	}
	return res
}

func (w *wifiTask) GetStatus(status int) string {
	res := w.AllStatus()
	return res[status]
}

func (w *wifiTask) RecoverCryptset(encrypt int) int {
	if encrypt == 0 {
		return WifiTaskCryptsetNone
	} else if (encrypt&(1<<20)) != 0 && (encrypt&(1<<21)) != 0 {
		return WifiTaskCryptsetWp2
	} else if (encrypt&(1<<7)) != 0 && (encrypt&(1<<20)) != 0 {
		return WifiTaskCryptsetWpa
	} else if (encrypt & (1 << 1)) != 0 {
		return WifiTaskCryptsetWep
	} else {
		return WifiTaskCryptsetUnknown
	}
}

func (w *wifiTask) AllCryptset() map[int]string {
	res := map[int]string{
		WifiTaskCryptsetNone:    "NONE",
		WifiTaskCryptsetWp2:     "WPA2",
		WifiTaskCryptsetWpa:     "WPA",
		WifiTaskCryptsetWep:     "WEP",
		WifiTaskCryptsetUnknown: "UNKNOWN",
	}
	return res
}

func (w *wifiTask) GetCryptset(cryptset int) string {
	res := w.AllCryptset()
	return res[cryptset]
}

func (w *wifiTask) AllCarrier() map[int]string {
	res := map[int]string{
		WifiTaskCarrier1: "802.11b",
		WifiTaskCarrier2: "802.11bplus",
		WifiTaskCarrier3: "802.11a",
		WifiTaskCarrier4: "802.11g",
		WifiTaskCarrier5: "802.11fhss",
		WifiTaskCarrier6: "802.11dsss",
		WifiTaskCarrier7: "802.11n",
		WifiTaskCarrier8: "802.11n",
	}
	return res
}
func (w *wifiTask) GetCarrier(carrier int) string {
	return w.AllCarrier()[carrier]
}
