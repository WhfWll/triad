package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type WifiTask struct {
	TaskID   int    `gorm:"column:task_id;AUTO_INCREMENT;primary_key" json:"taskID"` //
	Mac      string `gorm:"column:mac;primary_key" json:"mac"`                       //
	TaskName string `gorm:"column:task_name" json:"taskName"`
	Status   int    `gorm:"column:status" json:"status"` // 0：初始状态，等待开始检测
	// 1：正在进行密码爆破
	// 2：密码爆破失败
	// 3：密码爆破成功，等待WiFi模拟
	// 4：正在WiFi模拟
	// 5：完成WiFi模拟
	// 2020-10-13 增加一种状态
	// 0：初始状态，等待开始检测
	// 1：开始进行报文收集
	// 2：正在进行密码爆破
	// 3：密码爆破失败
	// 4：密码爆破成功，等待AP渗透
	// 5：AP渗透完成，等待WiFi模拟
	// 6：正在进行WiFi模块
	// 7：完成WiFi模拟
	PasswdSource     int       `gorm:"column:passwd_source" json:"passwdSource"`         //
	PasswdDict       string    `gorm:"column:passwd_dict" json:"passwdDict"`             //
	Channel          int       `gorm:"column:channel" json:"channel"`                    //
	Encrypt          int       `gorm:"column:encrypt" json:"encrypt"`                    //
	Carrier          int       `gorm:"column:carrier" json:"carrier"`                    //
	Passwd           string    `gorm:"column:passwd" json:"passwd"`                      //
	StartTime        int64     `gorm:"column:start_time" json:"startTime"`               //
	EndTime          int64     `gorm:"column:end_time" json:"endTime"`                   //
	Ssid             string    `gorm:"column:ssid" json:"ssid"`                          //
	IsSimulate       int       `gorm:"column:is_simulate" json:"isSimulate"`             // 是否模拟：0：否，1：是
	SimulateDuration int       `gorm:"column:simulate_duration" json:"simulateDuration"` // 模拟时长，单位秒
	Model            string    `gorm:"column:model" json:"model"`                        //
	IsCrack          int       `gorm:"column:is_crack" json:"isCrack"`                   // 是否爆破，0：否，1：是
	ID               int       `gorm:"column:id" json:"id"`                              //
	IsEmbed          int       `gorm:"column:is_embed" json:"isEmbed"`                   // 是否植入，0：否，1：是
	ReasonCode       int       `gorm:"column:reason_code" json:"reasonCode"`             //
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`             // 创建时间
	UpdateTime       time.Time `gorm:"column:update_time" json:"updateTime"`             // 修改时间
}

// TableName sets insert table name for this struct type
func (w *WifiTask) TableName() string {
	return "wifi_task"
}

// Get retrieves a list of WifiTask from database
func (w *WifiTask) GetWifiTaskList(ctx context.Context, page, limit int, search string) ([]WifiTask, int64) {
	var (
		WifiTaskList []WifiTask
		count        int64
		db           = mysql.FromContext(ctx).Model(&WifiTask{})
	)

	if search != "" {
		db = db.Where("ssid like ?", "%"+search+"%")
	}
	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Order("create_time DESC").Find(&WifiTaskList)

	return WifiTaskList, count
}

// Get retrieves a single record of WifiTask from database
func (w *WifiTask) GetWifiTask(ctx context.Context) (WifiTask, error) {
	var (
		wifiTask WifiTask
		err      error
		db       = mysql.FromContext(ctx).Model(&WifiTask{})
	)

	curErr := db.Where("id = ?", w.ID).First(&wifiTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return wifiTask, err
}

// Add persists WifiTask to database
func (w *WifiTask) AddWifiTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&WifiTask{})

	if err := db.Create(&w).Error; err != nil {
		return err
	}

	return nil
}

// Update changes WifiTask by id
func (w *WifiTask) UpdateWifiTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&WifiTask{})

	if err := db.Where("id = ?", w.ID).Updates(w).Error; err != nil {
		return err
	}

	return nil
}

// Delete WifiTask by id
func (w *WifiTask) DeleteWifiTask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&WifiTask{})

	//w.Estate = "deleted"
	w.UpdateTime = time.Now()
	if err := db.Where("id = ?", w.ID).Updates(w).Error; err != nil {
		return err
	}

	return nil
}

// Delete WifiTask by taskId
func (w *WifiTask) DeleteByTaskId(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&WifiTask{})

	//w.Estate = "deleted"
	w.UpdateTime = time.Now()
	if err := db.Where("task_id in ?", taskIds).Delete(w).Error; err != nil {
		return err
	}

	return nil
}

// 获取最后一个taskID 这张表结构暂时不要动，因为后面的C程序也在用
func (w *WifiTask) GetLastTaskID(ctx context.Context) (int, error) {
	var (
		wifiTask WifiTask
		err      error
		db       = mysql.FromContext(ctx).Model(&WifiTask{})
	)

	curErr := db.Order("create_time DESC").First(&wifiTask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return wifiTask.TaskID, err
}
