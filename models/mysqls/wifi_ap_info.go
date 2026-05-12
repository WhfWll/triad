package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"time"
)

type WifiApInfo struct {
	SourceMac            string    `gorm:"column:source_mac;primary_key" json:"sourceMac"`            //
	AttackType           int       `gorm:"column:attack_type" json:"attackType"`                      //
	Manuf                string    `gorm:"column:manuf" json:"manuf"`                                 //
	FirstTime            int64     `gorm:"column:first_time" json:"firstTime"`                        //
	LastTime             int64     `gorm:"column:last_time" json:"lastTime"`                          //
	Duration             int       `gorm:"column:duration" json:"duration"`                           //
	Channel              int       `gorm:"column:channel" json:"channel"`                             //
	Maxseenrate          int       `gorm:"column:maxseenrate" json:"maxseenrate"`                     //
	Carrierset           int       `gorm:"column:carrierset" json:"carrierset"`                       //
	Encodingset          string    `gorm:"column:encodingset" json:"encodingset"`                     //
	LlcPackets           int64     `gorm:"column:llc_packets" json:"llcPackets"`                      //
	DataPackets          int64     `gorm:"column:data_packets" json:"dataPackets"`                    //
	Datasize             int64     `gorm:"column:datasize" json:"datasize"`                           //
	LastSignalRssi       int       `gorm:"column:last_signal_rssi" json:"lastSignalRssi"`             //
	MinSignalRssi        int       `gorm:"column:min_signal_rssi" json:"minSignalRssi"`               //
	MaxSignalRssi        int       `gorm:"column:max_signal_rssi" json:"maxSignalRssi"`               //
	SsidCloaked          int       `gorm:"column:ssid_cloaked" json:"ssidCloaked"`                    // 0：广播，1：隐藏
	Ssid                 string    `gorm:"column:ssid" json:"ssid"`                                   //
	SsidBeaconInfo       string    `gorm:"column:ssid_beacon_info" json:"ssidBeaconInfo"`             //
	SsidFirstTime        int64     `gorm:"column:ssid_first_time" json:"ssidFirstTime"`               //
	SsidLastTime         int64     `gorm:"column:ssid_last_time" json:"ssidLastTime"`                 //
	SsidMaxrate          float64   `gorm:"column:ssid_maxrate" json:"ssidMaxrate"`                    //
	SsidBeaconrate       int       `gorm:"column:ssid_beaconrate" json:"ssidBeaconrate"`              //
	SsidCountry          string    `gorm:"column:ssid_country" json:"ssidCountry"`                    //
	SsidStartchan        int       `gorm:"column:ssid_startchan" json:"ssidStartchan"`                //
	SsidEndchan          int       `gorm:"column:ssid_endchan" json:"ssidEndchan"`                    //
	SsidTxpower          int       `gorm:"column:ssid_txpower" json:"ssidTxpower"`                    //
	SsidCryptset         int       `gorm:"column:ssid_cryptset" json:"ssidCryptset"`                  //
	LocationFlag         int       `gorm:"column:location_flag" json:"locationFlag"`                  //
	Lat                  float64   `gorm:"column:lat" json:"lat"`                                     //
	Lon                  float64   `gorm:"column:lon" json:"lon"`                                     //
	Alt                  float64   `gorm:"column:alt" json:"alt"`                                     //
	DroneNum             int       `gorm:"column:drone_num" json:"droneNum"`                          //
	NearestDrone         string    `gorm:"column:nearest_drone" json:"nearestDrone"`                  //
	NearestDroneRssi     int       `gorm:"column:nearest_drone_rssi" json:"nearestDroneRssi"`         //
	NearestDroneDistance float64   `gorm:"column:nearest_drone_distance" json:"nearestDroneDistance"` //
	FreqCounter          string    `gorm:"column:freq_counter" json:"freqCounter"`                    //
	List                 int       `gorm:"column:list" json:"list"`                                   //
	Status               int       `gorm:"column:status" json:"status"`                               // 0：无，1：强制下线
	ConnectCliCount      int       `gorm:"column:connect_cli_count" json:"connectCliCount"`           //
	AreaID               int       `gorm:"column:area_id" json:"areaID"`                              //
	LabelID              int       `gorm:"column:label_id" json:"labelID"`                            //
	IsMydevice           int       `gorm:"column:is_mydevice" json:"isMydevice"`                      // 0:非自有设备，1：自有设备
	NetworkType          int       `gorm:"column:network_type" json:"networkType"`                    // 0：内网，1：外网
	Tag                  int       `gorm:"column:tag" json:"tag"`                                     // 1:标记为非私搭乱建
	DeviceType           int       `gorm:"column:device_type" json:"deviceType"`                      // 0:AP，  1:终端
	OperateUser          string    `gorm:"column:operate_user" json:"operateUser"`                    //
	AddTime              int64     `gorm:"column:add_time" json:"addTime"`                            //
	CreateTime           time.Time `gorm:"column:create_time" json:"createTime"`                      // 创建时间
	UpdateTime           time.Time `gorm:"column:update_time" json:"updateTime"`                      // 修改时间
}

// TableName sets insert table name for this struct type
func (w *WifiApInfo) TableName() string {
	return "wifi_ap_info"
}

// Get retrieves a list of wifiApInfo from database
func (w *WifiApInfo) GetWifiApInfoList(ctx context.Context, page, limit int) ([]WifiApInfo, int64, error) {
	var (
		wifiApInfoList []WifiApInfo
		count          int64
		db             = mysql.FromContext(ctx).Model(&WifiApInfo{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&wifiApInfoList)
	db.Count(&count)

	return wifiApInfoList, count, nil
}

// Add persists wifiApInfo to database
func (w *WifiApInfo) AddWifiApInfo(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&WifiApInfo{})

	if err := db.Create(w).Error; err != nil {
		return err
	}

	return nil
}

// ApList wifi 所有在线wifi列表
func (w *WifiApInfo) ApList(ctx context.Context, search, startDate, endDate string) ([]WifiApInfo, int64) {
	var (
		wifiApInfoList []WifiApInfo
		count          int64
		db             = mysql.FromContext(ctx).Model(&WifiApInfo{})
	)
	if search != "" {
		db = db.Where("ssid like ? or source_mac = ?", "%"+search+"%", search)
	}

	if startDate != "" {
		db = db.Where("last_time >= ? ", startDate).Where("last_time <= ? ", endDate)
	} else {
		db = db.Where("last_time >= ? ", time.Now().Unix()-180)
	}

	db = db.Where("last_signal_rssi > -100")
	db.Count(&count)
	db.Order("last_time DESC").Find(&wifiApInfoList)

	return wifiApInfoList, count
}

// GetBySourceMac 通过mac地址获取单个wifi信息
func (w *WifiApInfo) GetBySourceMac(ctx context.Context, sourceMac string) WifiApInfo {
	var db = mysql.FromContext(ctx).Model(&WifiApInfo{})
	var data WifiApInfo
	db.Where("source_mac = ?", sourceMac).Find(&data)

	return data
}
