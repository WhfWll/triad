package enums

import "strconv"

var AssetScanEnum assetScan

type assetScan struct{}

func (a *assetScan) GetTaskExecCycleTypeMap() map[int]string {
	enum := map[int]string{
		TaskExecTypeCycleTypeWeek:  "每周一次",
		TaskExecTypeCycleTypeMonth: "每月一次",
	}
	return enum
}

func (a *assetScan) GetTaskExecCycleTypeName(execCycleType int) string {
	enum := a.GetTaskExecCycleTypeMap()
	if res, ok := enum[execCycleType]; ok {
		return res
	}
	return ""
}

func (a *assetScan) TaskExecCycleType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: TaskExecTypeCycleTypeWeek,
		Label: a.GetTaskExecCycleTypeName(TaskExecTypeCycleTypeWeek),
	}, {
		Value: TaskExecTypeCycleTypeMonth,
		Label: a.GetTaskExecCycleTypeName(TaskExecTypeCycleTypeMonth),
	}}
	return result
}

// 周期执行值 每月执行一次
func (a *assetScan) GetTaskExecCycleTypeMonthMap() map[int]string {
	enum := make(map[int]string)
	for i := 1; i <= 31; i++ {
		enum[i] = strconv.Itoa(i) + "号"
	}
	return enum
}
func (a *assetScan) GetTaskExecCycleTypeMonthName(execType int) string {
	enum := a.GetTaskExecCycleTypeMonthMap()
	if res, ok := enum[execType]; ok {
		return res
	}
	return ""
}
func (a *assetScan) TaskExecCycleTypeMonth() interface{} {
	var result []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}
	for i := 1; i <= 31; i++ {
		result = append(result, struct {
			Value int    `json:"value"`
			Label string `json:"label"`
		}{i, strconv.Itoa(i) + "号"})
	}
	return result
}

func (a *assetScan) GetTaskPortScanTypeMap() map[int]string {
	enum := map[int]string{
		TaskConfigurationPortScanTypeTop10:   "TOP10端口",
		TaskConfigurationPortScanTypeTop20:   "TOP20端口",
		TaskConfigurationPortScanTypeTop50:   "TOP50端口",
		TaskConfigurationPortScanTypeTop100:  "TOP100端口",
		TaskConfigurationPortScanTypeTop500:  "TOP500端口",
		TaskConfigurationPortScanTypeTop1000: "TOP1000端口",
		TaskConfigurationPortScanTypeAll:     "全部端口",
		TaskConfigurationPortScanTypeCustom:  "自定义端口",
	}
	return enum
}

func (a *assetScan) GetTaskPortScanTypeName(portScan int) string {
	enum := a.GetTaskPortScanTypeMap()
	if res, ok := enum[portScan]; ok {
		return res
	}
	return ""
}

func (a *assetScan) TaskPortScanType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: TaskConfigurationPortScanTypeTop10,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeTop10),
	}, {
		Value: TaskConfigurationPortScanTypeTop20,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeTop20),
	}, {
		Value: TaskConfigurationPortScanTypeTop50,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeTop50),
	}, {
		Value: TaskConfigurationPortScanTypeTop100,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeTop100),
	}, {
		Value: TaskConfigurationPortScanTypeTop500,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeTop500),
	}, {
		Value: TaskConfigurationPortScanTypeTop1000,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeTop1000),
	}, {
		Value: TaskConfigurationPortScanTypeAll,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeAll),
	}, {
		Value: TaskConfigurationPortScanTypeCustom,
		Label: a.GetTaskPortScanTypeName(TaskConfigurationPortScanTypeCustom),
	}}
	return result
}

// 端口范围选项后联动的端口数据
func (a *assetScan) TaskPortScanTypeValue() map[int]string {
	enum := map[int]string{
		TaskConfigurationPortScanTypeTop10:   "21,22,23,80,443,445,3306,8000,8080,8088",
		TaskConfigurationPortScanTypeTop20:   "21,22,23,80,443,445,3306,7000-7002,8000-8003,8080-8083,8088,9200",
		TaskConfigurationPortScanTypeTop50:   "21,22,23,80,88,106,110,111,113,119,135,139,143,144,179,199,389,427,1521,1630,1158,443,445,888,777,999,1070,1080,1090,3306,7000-7003,8000-8003,8008,8080-8083,8088,9000-9002,8090,9200,9300",
		TaskConfigurationPortScanTypeTop100:  "7,5555,9,13,21,22,23,25,26,37,53,79,80,81,88,106,110,111,113,119,135,139,143,144,179,199,389,427,443,444,445,465,513,514,515,543,544,548,554,587,631,646,873,888,990,993,995,1025,1026,1027,1028,1029,1080,1110,1433,1443,1720,1723,1755,1900,2000,2001,2049,2121,2181,2717,3000,3128,3306,3389,3986,4899,5000,5009,5051,5060,5101,5190,5357,5432,5631,5666,5800,5900,6000,6001,6646,7000,7001,7002,7003,7004,7005,7070,8000,8008,8009,8080,8081,8443,8888,9100,9999,10000,11211,32768,49152,49153,49154,49155,49156,49157,8088,9090,8090,8001,82,9080,8082,8089,9000,8002,89,8083,8200,90,8086,801,8011,8085,9001,9200,8100,8012,85,8084,8070,8091,8003,99,7777,8010,8028,8087,83,808,38888,8181,800,18080,8099,8899,86,8360,8300,8800,8180,3505,9002,8053,1000,7080,8989,28017,9060,8006,41516,880,8484,6677,8016,84,7200,9085,5555,8280,1980,8161,9091,7890,8060,6080,8880,8020,889,8881,9081,7007,8004,38501,1010,17,19,255,1024,1030,1041,1048,1049,1053,1054,1056,1064,1065,1801,2103,2107,2967,3001,3703,5001,5050,6004,8031,10010,10250,10255,6888,87,91,92,98,1081,1082,1118,1888,2008,2020,2100,2375,3008,6648,6868,7008,7071,7074,7078,7088,7680,7687,7688,8018,8030,8038,8042,8044,8046,8048,8069,8092,8093,8094,8095,8096,8097,8098,8101,8108,8118,8172,8222,8244,8258,8288,8448,8834,8838,8848,8858,8868,8879,8983,9008,9010,9043,9082,9083,9084,9086,9087,9088,9089,9092,9093,9094,9095,9096,9097,9098,9099,9443,9448,9800,9981,9986,9988,9998,10001,10002,10004,10008,12018,12443,14000,16080,18000,18001,18002,18004,18008,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018",
		TaskConfigurationPortScanTypeTop500:  "7,8,9,13,17,19,20,21,22,23,25,26,37,53,60,65,66,70,77,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,103,106,110,111,113,114,119,122,132,133,135,139,143,144,171,179,180,188,199,200,206,208,211,235,255,268,280,299,302,308,321,381,389,403,421,423,427,442,443,444,445,447,465,511,513,514,515,517,522,543,544,548,554,587,591,610,631,646,666,688,701,770,778,800,801,802,803,804,805,806,808,809,811,812,866,873,877,880,888,889,925,955,983,990,993,995,999,1000,1001,1005,1010,1024,1025,1026,1027,1028,1029,1030,1039,1041,1042,1048,1049,1053,1054,1056,1064,1065,1080,1081,1082,1085,1088,1100,1107,1108,1110,1118,1122,1123,1128,1158,1180,1182,1200,1212,1213,1234,1300,1301,1313,1356,1433,1443,1500,1550,1666,1680,1700,1720,1722,1723,1755,1790,1800,1801,1818,1863,1888,1900,1933,1949,1979,1980,1982,2000,2001,2005,2006,2007,2008,2009,2010,2011,2012,2013,2014,2015,2020,2046,2049,2051,2060,2070,2080,2093,2100,2103,2107,2110,2121,2125,2181,2222,2301,2348,2375,2382,2480,2490,2517,2521,2585,2717,2808,2886,2901,2967,3000,3001,3008,3012,3013,3030,3050,3080,3128,3216,3220,3306,3312,3333,3380,3389,3443,3456,3465,3503,3505,3535,3580,3588,3600,3606,3668,3690,3703,3721,3880,3938,3986,4000,4001,4016,4040,4300,4389,4430,4433,4440,5000,5001,5002,5003,5009,5013,5050,5600,5601,5631,5632,5644,5655,5656,5666,5678,5800,5811,5881,5887,5888,5898,5900,5902,5966,6000,6001,6002,6003,6004,6006,6010,6060,6080,6088,6090,6101,6118,6170,6180,6198,6226,6259,6379,6388,6886,6888,6889,6890,6920,6969,6988,7000,7001,7002,7003,7004,7005,7006,7007,7008,7009,7010,7011,7012,7017,7018,7020,7021,7022,7028,7031,7041,7048,7050,7060,7070,7071,7074,7078,7080,7081,7084,7086,7088,7094,7100,7101,7102,7108,7111,7117,7123,7129,7171,7180,7200,7201,7202,7215,7272,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8018,8019,8020,8021,8022,8023,8024,8025,8026,8027,8028,8029,8030,8031,8032,8033,8035,8036,8037,8038,8039,8040,8041,8042,8043,8044,8045,8046,8048,8050,8051,8053,8055,8056,8057,8058,8060,8061,8062,8064,8065,8066,8069,8070,8071,8073,8077,8078,8079,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8102,8103,8104,8172,8180,8181,8182,8183,8184,8186,8188,8288,8300,8308,8322,8333,8341,8343,8360,8380,8580,8582,8585,8600,8601,8610,8649,8660,9200,9201",
		TaskConfigurationPortScanTypeTop1000: "7,8,9,13,17,19,20,21,22,23,25,26,37,53,60,65,66,70,77,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,103,106,110,111,113,114,119,122,132,133,135,139,143,144,171,179,180,188,199,200,206,208,211,235,255,268,280,299,302,308,321,381,389,403,421,423,427,442,443,444,445,447,465,511,513,514,515,517,522,543,544,548,554,587,591,610,631,646,666,688,701,770,778,800,801,802,803,804,805,806,808,809,811,812,866,873,877,880,888,889,925,955,983,990,993,995,999,1000,1001,1005,1010,1024,1025,1026,1027,1028,1029,1030,1039,1041,1042,1048,1049,1053,1054,1056,1064,1065,1080,1081,1082,1085,1088,1100,1107,1108,1110,1118,1122,1123,1128,1158,1180,1182,1200,1212,1213,1234,1300,1301,1313,1356,1433,1443,1500,1550,1666,1680,1700,1720,1722,1723,1755,1790,1800,1801,1818,1863,1888,1900,1933,1949,1979,1980,1982,2000,2001,2005,2006,2007,2008,2009,2010,2011,2012,2013,2014,2015,2020,2046,2049,2051,2060,2070,2080,2093,2100,2103,2107,2110,2121,2125,2181,2222,2301,2348,2375,2382,2480,2490,2517,2521,2585,2717,2808,2886,2901,2967,3000,3001,3008,3012,3013,3030,3050,3080,3128,3216,3220,3306,3312,3333,3380,3389,3443,3456,3465,3503,3505,3535,3580,3588,3600,3606,3668,3690,3703,3721,3880,3938,3986,4000,4001,4016,4040,4300,4389,4430,4433,4440,4443,4567,4848,4850,4899,5000,5001,5002,5003,5009,5013,5050,5051,5060,5080,5081,5098,5100,5101,5155,5156,5190,5200,5201,5203,5233,5255,5256,5280,5357,5432,5544,5552,5555,5561,5600,5601,5631,5632,5644,5655,5656,5666,5678,5800,5811,5881,5887,5888,5898,5900,5902,5966,6000,6001,6002,6003,6004,6006,6010,6060,6080,6088,6090,6101,6118,6170,6180,6198,6226,6259,6379,6388,6443,6510,6543,6546,6565,6587,6600,6602,6603,6611,6646,6648,6666,6677,6680,6688,6699,6778,6789,6800,6801,6842,6868,6869,6879,6886,6888,6889,6890,6920,6969,6988,7000,7001,7002,7003,7004,7005,7006,7007,7008,7009,7010,7011,7012,7017,7018,7020,7021,7022,7028,7031,7041,7048,7050,7060,7070,7071,7074,7078,7080,7081,7084,7086,7088,7094,7100,7101,7102,7108,7111,7117,7123,7129,7171,7180,7200,7201,7202,7215,7272,7288,7321,7330,7380,7443,7500,7567,7680,7687,7688,7700,7702,7703,7709,7711,7713,7742,7776,7777,7778,7788,7791,7801,7856,7888,7890,7899,7900,7909,7915,7921,7925,7942,7943,7979,7999,8000,8001,8002,8003,8004,8005,8006,8007,8008,8009,8010,8011,8012,8013,8014,8015,8016,8018,8019,8020,8021,8022,8023,8024,8025,8026,8027,8028,8029,8030,8031,8032,8033,8035,8036,8037,8038,8039,8040,8041,8042,8043,8044,8045,8046,8048,8050,8051,8053,8055,8056,8057,8058,8060,8061,8062,8064,8065,8066,8069,8070,8071,8073,8077,8078,8079,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8102,8103,8104,8108,8111,8112,8118,8119,8122,8123,8130,8133,8136,8144,8161,8168,8172,8180,8181,8182,8183,8184,8186,8188,8189,8190,8191,8192,8193,8196,8197,8200,8213,8220,8222,8244,8258,8260,8280,8282,8283,8288,8300,8308,8322,8333,8341,8343,8360,8380,8381,8382,8383,8384,8390,8399,8400,8401,8402,8443,8445,8448,8477,8480,8481,8484,8500,8567,8580,8582,8585,8600,8601,8610,8649,8660,8666,8680,8686,8688,8700,8710,8720,8735,8780,8781,8787,8788,8799,8800,8801,8802,8806,8808,8809,8810,8813,8822,8834,8838,8839,8844,8848,8858,8860,8864,8866,8868,8877,8879,8880,8881,8885,8886,8887,8888,8889,8890,8891,8892,8895,8898,8899,8900,8902,8910,8912,8913,8955,8956,8972,8980,8983,8987,8988,8989,8990,8991,8997,8999,9000,9001,9002,9003,9004,9005,9006,9007,9008,9009,9010,9011,9012,9013,9014,9015,9020,9022,9025,9030,9031,9036,9039,9043,9050,9053,9060,9061,9070,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9101,9110,9111,9112,9131,9180,9182,9190,9191,9200,9201,9212,9231,9300,9301,9302,9437,9443,9448,9494,9500,9504,9507,9527,9595,9666,9696,9704,9800,9845,9876,9888,9889,9898,9900,9901,9909,9910,9912,9914,9918,9919,9980,9981,9986,9988,9990,9991,9992,9995,9997,9998,9999,10000,10001,10002,10003,10004,10007,10008,10009,10010,10016,10021,10025,10038,10040,10051,10066,10068,10080,10082,10086,10087,10088,10089,10099,10118,10154,10250,10255,10333,11000,11001,11080,11158,11211,11324,11347,11362,11366,11372,11381,12001,12018,12333,12345,12443,12881,13333,13382,13988,14000,14007,15000,15004,15018,15580,15672,15693,15801,15888,16080,16788,17000,17003,17095,17777,18000,18001,18002,18004,18008,18060,18080,18081,18082,18085,18088,18090,18098,18103,18264,18801,18803,18880,18881,18888,19001,19010,19045,19080,19101,19244,20000,20001,20021,20022,20046,20052,20140,20151,20153,20720,20806,20808,21000,21080,21245,21501,21502,22222,22343,22580,23352,23454,25006,25024,27000,27017,28017,28018,28080,28099,28214,28280,28780,30000,30001,30058,30082,30088,30551,31188,31945,32766,32768,34440,38000,38080,38086,38443,38501,38517,38888,40000,40069,40080,40310,41516,42424,43651,45149,45177,47078,47088,47583,48080,49152,49153,49154,49155,49156,49157,49705,49960,50000,50030,50045,50060,50070,50075,50080,50090,50240,51106,55351,55858,57880,58000,58031,58060,58080,58898,59009,59777,59999,60010,60022,60101,60465,61081,61999,65000,65001,65055,65129,65486,65493,65533,65535",
		TaskConfigurationPortScanTypeAll:     "0-65535",
		TaskConfigurationPortScanTypeCustom:  "",
	}
	return enum
}

// ip资产状态变化
const (
	TaskIpChangeTypeNo       = 1 //未变化
	TaskIpChangeTypeReduce   = 2 //已减少IP
	TaskIpChangeTypeAdd      = 3 //新增加IP
	TaskIpChangeTypePort     = 4 //端口变化IP
	TaskIpChangeTypeService  = 5 //服务变化IP
	TaskIpChangeTypeAssembly = 6 //组件变化IP
)

func (a *assetScan) GetTaskIpChangeTypeMap() map[int]string {
	enum := map[int]string{
		TaskIpChangeTypeNo:       "未变化",
		TaskIpChangeTypeReduce:   "已减少IP",
		TaskIpChangeTypeAdd:      "新增加IP",
		TaskIpChangeTypePort:     "端口变化IP",
		TaskIpChangeTypeService:  "服务变化IP",
		TaskIpChangeTypeAssembly: "组件变化IP",
	}
	return enum
}

func (a *assetScan) GetTaskIpChangeTypeName(changeType int) string {
	enum := a.GetTaskIpChangeTypeMap()
	if res, ok := enum[changeType]; ok {
		return res
	}
	return ""
}

func (a *assetScan) TaskIpChangeType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: TaskIpChangeTypeNo,
		Label: a.GetTaskIpChangeTypeName(TaskIpChangeTypeNo),
	}, {
		Value: TaskIpChangeTypeReduce,
		Label: a.GetTaskIpChangeTypeName(TaskIpChangeTypeReduce),
	}, {
		Value: TaskIpChangeTypeAdd,
		Label: a.GetTaskIpChangeTypeName(TaskIpChangeTypeAdd),
	}, {
		Value: TaskIpChangeTypePort,
		Label: a.GetTaskIpChangeTypeName(TaskIpChangeTypePort),
	}, {
		Value: TaskIpChangeTypeService,
		Label: a.GetTaskIpChangeTypeName(TaskIpChangeTypeService),
	}, {
		Value: TaskIpChangeTypeAssembly,
		Label: a.GetTaskIpChangeTypeName(TaskIpChangeTypeAssembly),
	}}
	return result
}

// ip存活状态
const (
	TaskIpIsLiveTypeYes = 1 //存活
	TaskIpIsLiveTypeNo  = 2 //不存活IP
)

func (a *assetScan) GetTaskIpIsLiveTypeMap() map[int]string {
	enum := map[int]string{
		TaskIpIsLiveTypeYes: "存活",
		TaskIpIsLiveTypeNo:  "不存活IP",
	}
	return enum
}

func (a *assetScan) GetTaskIpIsLiveTypeName(liveType int) string {
	enum := a.GetTaskIpIsLiveTypeMap()
	if res, ok := enum[liveType]; ok {
		return res
	}
	return ""
}

func (a *assetScan) TaskIpIsLiveType() interface{} {
	result := []struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}{{
		Value: TaskIpIsLiveTypeYes,
		Label: a.GetTaskIpIsLiveTypeName(TaskIpIsLiveTypeYes),
	}, {
		Value: TaskIpIsLiveTypeNo,
		Label: a.GetTaskIpIsLiveTypeName(TaskIpIsLiveTypeNo),
	}}
	return result
}

// ip探测状态
const (
	TaskIpStatusBegin            = 1 //待执行
	TaskIpStatusRunningBeginNmap = 2 //运行中-开始nmap扫描
	TaskIpStatusRunningEndNmap   = 3 //运行中-完成nmap扫描
	TaskIpStatusFinish           = 4 //已结束
)

func (a *assetScan) GetTaskIpStatusMap() map[int]string {
	enum := map[int]string{
		TaskIpStatusBegin:            "待执行",
		TaskIpStatusRunningBeginNmap: "运行中", //运行中-开始nmap扫描
		TaskIpStatusRunningEndNmap:   "运行中", //运行中-完成nmap扫描
		TaskIpStatusFinish:           "已完成",
	}
	return enum
}

func (a *assetScan) GetTaskIpStatusName(status int) string {
	enum := a.GetTaskIpStatusMap()
	if res, ok := enum[status]; ok {
		return res
	}
	return ""
}

func (a *assetScan) TaskIpStatus() interface{} {
	result := []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}{{
		Value: "1",
		Label: a.GetTaskIpStatusName(TaskIpStatusBegin),
	}, {
		Value: "2,3",
		Label: a.GetTaskIpStatusName(TaskIpStatusRunningBeginNmap),
	}, {
		Value: "4",
		Label: a.GetTaskIpStatusName(TaskIpStatusFinish),
	}}
	return result
}

const (
	TaskPortStatusBegin           = 1 //待执行
	TaskPortStatusRunningBeginWeb = 2 //运行中-开始web组件扫描
	TaskPortStatusRunningEndWeb   = 3 //运行中-完成web组件扫描
	TaskPortStatusFinish          = 4 //已结束
)

// 扫描端口状态变化
const (
	TaskPortChangeTypeNo        = 1 //未变化
	TaskPortChangeTypeReduce    = 2 //已减少端口
	TaskPortChangeTypeAdd       = 3 //新增加端口
	TaskPortChangeTypePortOpen  = 4 //关闭->开放端口
	TaskPortChangeTypeService   = 5 //服务变化端口
	TaskPortChangeTypeAssembly  = 6 //组件变化端口
	TaskPortChangeTypePortClose = 7 //开放->关闭端口
)

func (a *assetScan) GetScanPortChangeTypeMap() map[int]string {
	enum := map[int]string{
		TaskPortChangeTypeNo:        "未变化端口",
		TaskPortChangeTypeReduce:    "端口已关闭",
		TaskPortChangeTypeAdd:       "新开放端口",
		TaskPortChangeTypePortOpen:  "端口重新开放",
		TaskPortChangeTypeService:   "服务变化",
		TaskPortChangeTypeAssembly:  "组件变化",
		TaskPortChangeTypePortClose: "",
	}
	return enum
}

func (a *assetScan) GetScanPortChangeTypeName(changeType int) string {
	enum := a.GetScanPortChangeTypeMap()
	if res, ok := enum[changeType]; ok {
		return res
	}
	return ""
}

func (a *assetScan) GetAssetPortChangeTypeMap() map[int]string {
	enum := map[int]string{
		TaskPortChangeTypeNo:        "",
		TaskPortChangeTypeReduce:    "端口已关闭",
		TaskPortChangeTypeAdd:       "新开放端口",
		TaskPortChangeTypePortOpen:  "端口重新开放",
		TaskPortChangeTypeService:   "服务变化",
		TaskPortChangeTypeAssembly:  "组件变化",
		TaskPortChangeTypePortClose: "",
	}
	return enum
}

func (a *assetScan) GetAssetPortChangeTypeName(changeType int) string {
	enum := a.GetAssetPortChangeTypeMap()
	if res, ok := enum[changeType]; ok {
		return res
	}
	return ""
}
