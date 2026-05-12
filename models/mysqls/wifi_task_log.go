package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type WifiTaskLog struct {
	ID           int    `gorm:"column:id;primary_key" json:"id"`          //
	ApMac        string `gorm:"column:ap_mac" json:"apMac"`               //
	TaskID       int    `gorm:"column:task_id" json:"taskID"`             //
	Content      string `gorm:"column:content" json:"content"`            //
	GenerateTime int64  `gorm:"column:generate_time" json:"generateTime"` //
}

// TableName sets insert table name for this struct type
func (w *WifiTaskLog) TableName() string {
	return "wifi_task_log"
}

// Get retrieves a list of wifiTaskLog from database
func (w *WifiTaskLog) GetWifiTaskLogList(ctx context.Context, page, limit int) ([]WifiTaskLog, int64, error) {
	var (
		wifiTaskLogList []WifiTaskLog
		count           int64
		db              = mysql.FromContext(ctx).Model(&WifiTaskLog{})
	)

	db.Limit(limit).Offset(limit * (page - 1)).Find(&wifiTaskLogList)
	db.Count(&count)

	return wifiTaskLogList, count, nil
}

// Get retrieves a single record of wifiTaskLog from database
func (w *WifiTaskLog) GetWifiTaskLog(ctx context.Context) (WifiTaskLog, error) {
	var (
		wifiTaskLog WifiTaskLog
		err         error
		db          = mysql.FromContext(ctx).Model(&WifiTaskLog{})
	)

	curErr := db.Where("id = ?", w.ID).First(&wifiTaskLog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return wifiTaskLog, err
}

func (w *WifiTaskLog) GetWifiTaskLogByTaskId(ctx context.Context, taskIds []int) ([]WifiTaskLog, error) {
	var (
		wifiTaskLog []WifiTaskLog
		err         error
		db          = mysql.FromContext(ctx).Model(&WifiTaskLog{})
	)

	curErr := db.Where("task_id in ?", taskIds).Find(&wifiTaskLog).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return wifiTaskLog, err
}

// Add persists wifiTaskLog to database
func (w *WifiTaskLog) AddWifiTaskLog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&WifiTaskLog{})

	if err := db.Create(w).Error; err != nil {
		return err
	}

	return nil
}

// Update changes wifiTaskLog by id
func (w *WifiTaskLog) UpdateWifiTaskLog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&WifiTaskLog{})

	if err := db.Where("id = ?", w.ID).Updates(w).Error; err != nil {
		return err
	}

	return nil
}

// Delete wifiTaskLog by id
func (w *WifiTaskLog) DeleteWifiTaskLog(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&WifiTaskLog{})

	//w.Estate = "deleted"
	//w.UpdateTime = time.Now()
	if err := db.Where("id = ?", w.ID).Updates(w).Error; err != nil {
		return err
	}

	return nil
}

// Delete wifiTaskLog by taskId
func (w *WifiTaskLog) DeleteByTaskId(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&WifiTaskLog{})

	//w.Estate = "deleted"
	//w.UpdateTime = time.Now()
	if err := db.Where("task_id in ?", taskIds).Delete(w).Error; err != nil {
		return err
	}

	return nil
}
