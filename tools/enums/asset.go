package enums

var AssetEnum asset

type asset struct{}

// 备案等级
const (
	AssetFilingLevelZero  int = 1 //备案等级-无
	AssetFilingLevelOne   int = 2 //备案等级-等保一级
	AssetFilingLevelTwo   int = 3 //备案等级-等保二级
	AssetFilingLevelThree int = 4 //备案等级-等保三级
	AssetFilingLevelFour  int = 5 //备案等级-等保四级
	AssetFilingLevelFive  int = 6 //备案等级-等保五级
)

func (a *asset) AllAssetFilingLevelEnum() map[int]string {
	res := map[int]string{
		AssetFilingLevelZero:  "无",
		AssetFilingLevelOne:   "一级",
		AssetFilingLevelTwo:   "二级",
		AssetFilingLevelThree: "三级",
		AssetFilingLevelFour:  "四级",
		AssetFilingLevelFive:  "五级",
	}
	return res
}

func (a *asset) GetAssetFilingLevel(filingLevel int) string {
	return a.AllAssetFilingLevelEnum()[filingLevel]
}

func (a *asset) AllAssetFilingLevelEnumReversal() map[string]int {
	res := map[string]int{
		"无":  AssetFilingLevelZero,
		"一级": AssetFilingLevelOne,
		"二级": AssetFilingLevelTwo,
		"三级": AssetFilingLevelThree,
		"四级": AssetFilingLevelFour,
		"五级": AssetFilingLevelFive,
	}
	return res
}

func (a *asset) GetAssetFilingLevelID(name string) int {
	return a.AllAssetFilingLevelEnumReversal()[name]
}
func (a *asset) GetAssetFilingLevelEnumArray() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: AssetFilingLevelZero,
		Label: a.GetAssetFilingLevel(AssetFilingLevelZero),
	}, {
		Value: AssetFilingLevelOne,
		Label: a.GetAssetFilingLevel(AssetFilingLevelOne),
	}, {
		Value: AssetFilingLevelTwo,
		Label: a.GetAssetFilingLevel(AssetFilingLevelTwo),
	}, {
		Value: AssetFilingLevelThree,
		Label: a.GetAssetFilingLevel(AssetFilingLevelThree),
	}, {
		Value: AssetFilingLevelFour,
		Label: a.GetAssetFilingLevel(AssetFilingLevelFour),
	}, {
		Value: AssetFilingLevelFive,
		Label: a.GetAssetFilingLevel(AssetFilingLevelFive),
	}}
	return result
}

// 资产状态变化
const (
	AssetChangeTypeNo       = 1 //未变化
	AssetChangeTypeReduce   = 2 //已减少IP
	AssetChangeTypeAdd      = 3 //新增加IP
	AssetChangeTypePort     = 4 //端口变化IP
	AssetChangeTypeService  = 5 //服务变化IP
	AssetChangeTypeAssembly = 6 //组件变化IP
)

func (a *asset) GetAssetChangeTypeMap() map[int]string {
	enum := map[int]string{
		AssetChangeTypeNo:       "未变化",
		AssetChangeTypeReduce:   "已减少IP",
		AssetChangeTypeAdd:      "新增加IP",
		AssetChangeTypePort:     "端口变化IP",
		AssetChangeTypeService:  "服务变化IP",
		AssetChangeTypeAssembly: "组件变化IP",
	}
	return enum
}

func (a *asset) GetAssetChangeTypeName(changeType int) string {
	enum := a.GetAssetChangeTypeMap()
	if res, ok := enum[changeType]; ok {
		return res
	}
	return ""
}

func (a *asset) AssetChangeType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: AssetChangeTypeNo,
		Label: a.GetAssetChangeTypeName(AssetChangeTypeNo),
	}, {
		Value: AssetChangeTypeReduce,
		Label: a.GetAssetChangeTypeName(AssetChangeTypeReduce),
	}, {
		Value: AssetChangeTypeAdd,
		Label: a.GetAssetChangeTypeName(AssetChangeTypeAdd),
	}, {
		Value: AssetChangeTypePort,
		Label: a.GetAssetChangeTypeName(AssetChangeTypePort),
	}, {
		Value: AssetChangeTypeService,
		Label: a.GetAssetChangeTypeName(AssetChangeTypeService),
	}, {
		Value: AssetChangeTypeAssembly,
		Label: a.GetAssetChangeTypeName(AssetChangeTypeAssembly),
	}}
	return result
}

// 资产存活状态
const (
	AssetIsLiveYes = 1 //存活
	AssetIsLiveNo  = 2 //不存活
)

func (a *asset) GetAssetIsLiveTypeMap() map[int]string {
	enum := map[int]string{
		AssetIsLiveYes: "存活",
		AssetIsLiveNo:  "不存活IP",
	}
	return enum
}

func (a *asset) GetAssetIsLiveTypeName(liveType int) string {
	enum := a.GetAssetIsLiveTypeMap()
	if res, ok := enum[liveType]; ok {
		return res
	}
	return ""
}

func (a *asset) AssetIsLiveType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: AssetIsLiveYes,
		Label: a.GetAssetIsLiveTypeName(AssetIsLiveYes),
	}, {
		Value: AssetIsLiveNo,
		Label: a.GetAssetIsLiveTypeName(AssetIsLiveNo),
	}}
	return result
}

// 近期变化资产是否忽略
const (
	AssetIsIgnoreNo  int = 1 //资产变化-不忽略
	AssetIsIgnoreYes int = 2 //资产变化-忽略
)

// 资产组默认值
const (
	AssetGroupUngroupedId   int    = 1     //资产组-未分组id
	AssetGroupUngroupedName string = "未分组" //资产组-未分组名称

)

// 资产危险等级
const (
	SafeAsset       = 0 // 安全资产
	FatalRiskAsset  = 1 // 严重资产
	HighRiskAsset   = 2 // 高危资产
	MiddleRiskAsset = 3 // 中危资产
	LowRiskAsset    = 4 // 低危资产
)

func (a *asset) GetAssetRiskTypeName(liveType int) string {
	enum := a.GetAssetRiskTypeMap()
	if res, ok := enum[liveType]; ok {
		return res
	}
	return ""
}

func (a *asset) GetAssetRiskTypeMap() map[int]string {
	enum := map[int]string{
		FatalRiskAsset:  "严重",
		HighRiskAsset:   "高危",
		MiddleRiskAsset: "中危",
		LowRiskAsset:    "低危",
		SafeAsset:       "安全",
	}
	return enum
}

func (a *asset) AssetRiskType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: HighRiskAsset,
		Label: a.GetAssetRiskTypeName(HighRiskAsset),
	}, {
		Value: MiddleRiskAsset,
		Label: a.GetAssetRiskTypeName(MiddleRiskAsset),
	}, {
		Value: LowRiskAsset,
		Label: a.GetAssetRiskTypeName(LowRiskAsset),
	}, {
		Value: SafeAsset,
		Label: a.GetAssetRiskTypeName(SafeAsset),
	}}
	return result
}

// 资产-类型
const (
	OASystem                   = 1 // OA系统
	OfficeEquipment            = 2 // 办公设备
	NetworkEquipment           = 3 // 网络设备
	SafeEquipment              = 4 // 安全设备
	Video                      = 5 // 视频监控
	IndustrialControlEquipment = 6 // 工控设备
	IOT                        = 7 // IOT
	BusinessEquipment          = 8 // 业务系统
	NoAssetType                = 0 // 未知资产
)

func (a *asset) GetAssetTypeName(liveType int) string {
	enum := a.GetAssetTypeMap()
	if res, ok := enum[liveType]; ok {
		return res
	}
	return "未知"
}

func (a *asset) GetAssetTypeMap() map[int]string {
	enum := map[int]string{
		OASystem:                   "OA系统",
		OfficeEquipment:            "办公设备",
		NetworkEquipment:           "网络设备",
		SafeEquipment:              "安全设备",
		Video:                      "视频监控",
		IndustrialControlEquipment: "工控设备",
		IOT:                        "IOT",
		BusinessEquipment:          "业务系统",
		NoAssetType:                "未知",
	}
	return enum
}

func (a *asset) GetAssetTypeID(name string) int {
	enum := a.GetAssetTypeMapReversal()
	if res, ok := enum[name]; ok {
		return res
	}
	return 0
}

func (a *asset) GetAssetTypeMapReversal() map[string]int {
	enum := map[string]int{
		"OA系统": OASystem,
		"办公设备": OfficeEquipment,
		"网络设备": NetworkEquipment,
		"安全设备": SafeEquipment,
		"视频监控": Video,
		"工控设备": IndustrialControlEquipment,
		"IOT":  IOT,
		"业务系统": BusinessEquipment,
		"未知":   NoAssetType,
	}
	return enum
}

func (a *asset) AssetType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: OASystem,
		Label: a.GetAssetTypeName(OASystem),
	}, {
		Value: OfficeEquipment,
		Label: a.GetAssetTypeName(OfficeEquipment),
	}, {
		Value: NetworkEquipment,
		Label: a.GetAssetTypeName(NetworkEquipment),
	}, {
		Value: SafeEquipment,
		Label: a.GetAssetTypeName(SafeEquipment),
	}, {
		Value: Video,
		Label: a.GetAssetTypeName(Video),
	}, {
		Value: IndustrialControlEquipment,
		Label: a.GetAssetTypeName(IndustrialControlEquipment),
	}, {
		Value: IOT,
		Label: a.GetAssetTypeName(IOT),
	}, {
		Value: BusinessEquipment,
		Label: a.GetAssetTypeName(BusinessEquipment),
	}, {
		Value: NoAssetType,
		Label: a.GetAssetTypeName(NoAssetType),
	}}
	return result
}

// 可信设备
const (
	TrustLevelTrusted   = 1 // 可信
	TrustLevelUntrusted = 2 // 未登记
)

func (a *asset) GetTrustLevelMap() map[int]string {
	return map[int]string{
		TrustLevelTrusted:   "可信",
		TrustLevelUntrusted: "未登记",
	}
}

func (a *asset) GetTrustLevelName(level int) string {
	if res, ok := a.GetTrustLevelMap()[level]; ok {
		return res
	}
	return ""
}

func (a *asset) GetTrustLevelID(name string) int {
	return a.GetTrustLevelReversal()[name]
}

func (a *asset) GetTrustLevelReversal() map[string]int {
	return map[string]int{
		"可信":  TrustLevelTrusted,
		"未登记": TrustLevelUntrusted,
	}
}

func (a *asset) TrustLevel() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{Value: TrustLevelTrusted, Label: a.GetTrustLevelName(TrustLevelTrusted)},
		{Value: TrustLevelUntrusted, Label: a.GetTrustLevelName(TrustLevelUntrusted)},
	}
	return result
}

// 设备权重
const (
	DeviceWeightHigh     = 1 // 高
	DeviceWeightMedium   = 2 // 中
	DeviceWeightLow      = 3 // 低
	DeviceWeightVeryHigh = 4 // 极高
	DeviceWeightVeryLow  = 5 // 极低
)

func (a *asset) GetDeviceWeightMap() map[int]string {
	return map[int]string{
		DeviceWeightHigh:     "高",
		DeviceWeightMedium:   "中",
		DeviceWeightLow:      "低",
		DeviceWeightVeryHigh: "极高",
		DeviceWeightVeryLow:  "极低",
	}
}

func (a *asset) GetDeviceWeightName(weight int) string {
	if res, ok := a.GetDeviceWeightMap()[weight]; ok {
		return res
	}
	return ""
}

func (a *asset) GetDeviceWeightNameID(name string) int {
	return a.GetDeviceWeightMapReversal()[name]
}

func (a *asset) GetDeviceWeightMapReversal() map[string]int {
	return map[string]int{
		"高":  DeviceWeightHigh,
		"中":  DeviceWeightMedium,
		"低":  DeviceWeightLow,
		"极高": DeviceWeightVeryHigh,
		"极低": DeviceWeightVeryLow,
	}
}

func (a *asset) DeviceWeight() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{
		{Value: DeviceWeightHigh, Label: a.GetDeviceWeightName(DeviceWeightHigh)},
		{Value: DeviceWeightMedium, Label: a.GetDeviceWeightName(DeviceWeightMedium)},
		{Value: DeviceWeightLow, Label: a.GetDeviceWeightName(DeviceWeightLow)},
		{Value: DeviceWeightVeryHigh, Label: a.GetDeviceWeightName(DeviceWeightVeryHigh)},
		{Value: DeviceWeightVeryLow, Label: a.GetDeviceWeightName(DeviceWeightVeryLow)},
	}
	return result
}

func GetAssetRisk(risk int) string {
	enum := GetAssetRiskEnum()
	if v, ok := enum[risk]; ok {
		return v
	}
	return ""
}

func GetAssetRiskEnum() map[int]string {
	enum := map[int]string{
		FatalRiskAsset:  "严重",
		HighRiskAsset:   "高危",
		MiddleRiskAsset: "中危",
		LowRiskAsset:    "低危",
		SafeAsset:       "安全",
	}
	return enum
}

type LoginProtocol string

const (
	ProtocolSSH    LoginProtocol = "ssh"
	ProtocolRDP    LoginProtocol = "rdp"
	ProtocolTelnet LoginProtocol = "telnet"
	ProtocolVNC    LoginProtocol = "vnc"
	ProtocolFTP    LoginProtocol = "ftp"
	ProtocolSFTP   LoginProtocol = "sftp"
	ProtocolMySQL  LoginProtocol = "mysql"
	ProtocolHTTP   LoginProtocol = "http"
	ProtocolHTTPS  LoginProtocol = "https"
	ProtocolSNMP   LoginProtocol = "snmp"
)

var LoginProtocolTools loginProtocolEnum

type loginProtocolEnum struct{}

func (l *loginProtocolEnum) GetProtocolMap() map[int]string {
	return map[int]string{
		ConnMethodSSH: "ssh",
	}
}

func (l *loginProtocolEnum) GetProtocolLabel(protocol int) string {
	if label, ok := l.GetProtocolMap()[protocol]; ok {
		return label
	}
	return ""
}

func (l *loginProtocolEnum) ProtocolList() interface{} {
	result := []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}{}
	for val, label := range l.GetProtocolMap() {
		result = append(result, struct {
			Value string `json:"value"`
			Label string `json:"label"`
		}{
			Value: string(val),
			Label: label,
		})
	}

	return result
}
