package mysqls

import (
	"context"
	"encoding/json"
	"smart/tools/enums"
	"strconv"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

// RemoteSession 远程会话
type RemoteSession struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`                // 主键
	Estate          string    `gorm:"column:estate" json:"estate"`                    // 数据状态valid/deleted
	TaskID          int       `gorm:"column:task_id" json:"taskID"`                   // 所属任务id
	TargetID        int       `gorm:"column:target_id" json:"targetID"`               // 所属目标id
	SessionID       string    `gorm:"column:session_id" json:"sessionID"`             // 会话编号
	TargetURL       string    `gorm:"column:target_url" json:"targetURL"`             // 目标地址
	Route           string    `gorm:"column:route" json:"route"`                      // 路由
	RemoteControl   string    `gorm:"column:remote_control" json:"remoteControl"`     // 远程控制
	Detail          string    `gorm:"column:detail" json:"detail"`                    // 详情
	Status          int       `gorm:"column:status" json:"status"`                    // 状态
	DownloadedFiles string    `gorm:"column:downloaded_files" json:"downloadedFiles"` // 已下载文件
	UserID          int       `gorm:"column:user_id" json:"userID"`                   // 操作人id
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`           // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`           // 修改时间
	SessionKey      string    `gorm:"column:session_key" json:"session_key"`          // 已下载文件
}

// TableName sets insert table name for this struct type
func (r *RemoteSession) TableName() string {
	return "remote_session"
}

// GetRemoteSessionList 获取会话证据列表
func (r *RemoteSession) GetRemoteSessionList(ctx context.Context, taskID, targetId, page, limit int, search string) ([]RemoteSession, int64, error) {
	var (
		remoteSessionList []RemoteSession
		query             string
		args              []interface{}
		count             int64
		db                = mysql.FromContext(ctx).Model(&RemoteSession{})
	)
	if targetId > 0 {
		query += "target_id = " + strconv.Itoa(targetId) + " and "
	}
	query += "estate = 'valid' and task_id = ?"
	args = append(args, taskID)
	if search != "" {
		query += " and vul_name LIKE ?"
		args = append(args, "%"+search+"%")
	}
	db.Where(query, args...).Limit(limit).Offset(limit * (page - 1)).Find(&remoteSessionList)
	db.Count(&count)
	return remoteSessionList, count, nil
}

// GetRemoteSession 获取远端会话详情
func (r *RemoteSession) GetRemoteSession(ctx context.Context) (RemoteSession, error) {
	var (
		taskEvidence RemoteSession
		err          error
		db           = mysql.FromContext(ctx).Model(&RemoteSession{})
	)
	curErr := db.Where("id = ?", r.ID).First(&taskEvidence).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskEvidence, err
}

// AddRemoteSession 添加远程会话
func (r *RemoteSession) AddRemoteSession(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&RemoteSession{})
	if err := db.Debug().Create(r).Error; err != nil {
		return err
	}
	return nil
}

// AddRemoteSessionInfo 添加远程会话信息
func (r *RemoteSession) AddRemoteSessionInfo(ctx context.Context, taskID, targetID, sessionType int, targetURL, sessionID, downloadedFiles string, detail map[string]interface{}) error {
	detailBytes, _ := json.Marshal(detail)
	remoteSession := RemoteSession{
		TaskID:          taskID,
		TargetID:        targetID,
		SessionID:       sessionID,
		TargetURL:       targetURL,
		Estate:          "valid",
		Status:          sessionType,
		Detail:          string(detailBytes),
		DownloadedFiles: downloadedFiles,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	return remoteSession.AddRemoteSession(ctx)
}

// UpdateRemoteSession 更新会话
func (r *RemoteSession) UpdateRemoteSession(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&RemoteSession{})
	if err := db.Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}
	return nil
}

// UpdateRemoteSession 更新会话
func (r *RemoteSession) UpdateRemoteSessionById(ctx context.Context, id int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&RemoteSession{})
	if err := db.Where("id = ?", id).Updates(param).Error; err != nil {
		return err
	}
	return nil
}

// DeleteRemoteSession 删除会话
func (r *RemoteSession) DeleteRemoteSession(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&RemoteSession{})
	r.Estate = "deleted"
	r.UpdateTime = time.Now()
	if err := db.Where("id = ?", r.ID).Updates(r).Error; err != nil {
		return err
	}
	return nil
}

// DeleteRemoteSessionIds 批量删除
func (r *RemoteSession) DeleteRemoteSessionIds(ctx context.Context, ids []string) error {
	var db = mysql.FromContext(ctx).Model(&RemoteSession{})
	r.Estate = "deleted"
	r.UpdateTime = time.Now()
	if err := db.Where(" id IN ?", ids).Updates(r).Error; err != nil {
		return err
	}
	return nil
}

// GetRemoteSessionCount 获取远程会话总数
func (r *RemoteSession) GetRemoteSessionCount(ctx context.Context, uid int, role int) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&RemoteSession{})
	)
	if role == enums.UserRoleOrdinary {
		db = db.Where("user_id = ?", uid)
	}
	db.Where("estate = ?", "valid")
	db.Count(&count)
	return count
}

// All 获取所有会话
func (t *RemoteSession) All(ctx context.Context, filter string) ([]RemoteSession, error) {
	var (
		remoteSessionList []RemoteSession
		db                = mysql.FromContext(ctx).Model(&RemoteSession{})
	)
	if filter != "" {
		db.Where(filter).Find(&remoteSessionList)
	} else {
		db.Find(&remoteSessionList)
	}
	return remoteSessionList, nil
}

func (r *RemoteSession) GetRemoteSessionById(ctx context.Context, id int) (RemoteSession, error) {
	var (
		remoteSession RemoteSession
		db            = mysql.FromContext(ctx).Model(&RemoteSession{})
		err           error
	)
	err = db.Where("estate = 'valid' and id = ?", id).Find(&remoteSession).Error
	return remoteSession, err
}
