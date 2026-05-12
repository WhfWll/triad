package mysqls

import (
	"context"
	"smart/tools/enums"
	"strconv"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

// TaskEvidence 任务证据表
type TaskEvidence struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`                // 主键
	Estate          string    `gorm:"column:estate" json:"estate"`                    // 数据状态valid/deleted
	TaskID          int       `gorm:"column:task_id" json:"taskID"`                   // 所属任务id
	TargetID        int       `gorm:"column:target_id" json:"targetID"`               // 所属目标id
	TargetURL       string    `gorm:"column:target_url" json:"targetURL"`             // 关联目标地址
	VulName         string    `gorm:"column:vul_name" json:"vulName"`                 // 关联漏洞名称
	RiskType        int       `gorm:"column:risk_type" json:"riskType"`               // 风险类型
	RiskDetail      string    `gorm:"column:risk_detail" json:"riskDetail"`           // 风险详情
	DownloadedFiles string    `gorm:"column:downloaded_files" json:"downloadedFiles"` // 已下载文件
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`           // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`           // 修改时间
}

// TableName 任务证据表
func (t *TaskEvidence) TableName() string {
	return "task_evidence"
}

// GetTaskEvidenceList 获取任务证据列表
func (t *TaskEvidence) GetTaskEvidenceList(ctx context.Context, taskID, targetId, page, limit, riskType int, search string) ([]TaskEvidence, int64, error) {
	var (
		taskEvidenceList []TaskEvidence
		query            string
		args             []interface{}
		count            int64
		db               = mysql.FromContext(ctx).Model(&TaskEvidence{})
	)
	if targetId != 0 {
		query += "target_id = " + strconv.Itoa(targetId) + " and "
	}
	query += "estate = 'valid' and task_id = ?"
	args = append(args, taskID)
	if search != "" {
		query += " and vul_name LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if riskType != 0 {
		query += " and risk_type = ? "
		args = append(args, riskType)
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Find(&taskEvidenceList)
	return taskEvidenceList, count, nil
}

// GetTaskEvidence 获取任务证据详情
func (t *TaskEvidence) GetTaskEvidence(ctx context.Context) (TaskEvidence, error) {
	var (
		taskEvidence TaskEvidence
		err          error
		db           = mysql.FromContext(ctx).Model(&TaskEvidence{})
	)
	curErr := db.Where("id = ?", t.ID).First(&taskEvidence).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return taskEvidence, err
}

// AddTaskEvidence 添加任务证据
func (t *TaskEvidence) AddTaskEvidence(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskEvidence{})
	if err := db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTaskEvidence 更新任务证据
func (t *TaskEvidence) UpdateTaskEvidence(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskEvidence{})
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTaskEvidence 删除任务证据
func (t *TaskEvidence) DeleteTaskEvidence(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskEvidence{})
	t.Estate = "deleted"
	t.UpdateTime = time.Now()
	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTaskEvidenceIds 批量删除
func (t *TaskEvidence) DeleteTaskEvidenceIds(ctx context.Context, ids []string) error {
	var db = mysql.FromContext(ctx).Model(&TaskEvidence{})
	t.Estate = "deleted"
	t.UpdateTime = time.Now()
	if err := db.Where(" id IN ?", ids).Updates(t).Error; err != nil {
		return err
	}
	return nil
}

type VulEvidenceStat struct {
	RiskType  int `json:"risk_type"`
	RiskCount int `json:"risk_count"`
}

// GetVulEvidenceStat 获取漏洞取证统计
func (t *TaskEvidence) GetVulEvidenceStat(ctx context.Context, uid int, role int) []VulEvidenceStat {
	var (
		VulEvidenceStatList []VulEvidenceStat
		db                  = mysql.FromContext(ctx).Model(&TaskEvidence{})
	)

	if role == enums.UserRoleOrdinary {
		db = db.Joins("JOIN task_task ON task_evidence.task_id = task_task.id").Where("task_task.user_id = ?", uid)
	}

	db.Where("task_evidence.estate", "valid").
		Select("task_evidence.risk_type, Count(*) as risk_count").
		Group("task_evidence.risk_type").
		Find(&VulEvidenceStatList)

	return VulEvidenceStatList
}

// All 获取所有目标
func (t *TaskEvidence) All(ctx context.Context, filter string) ([]TaskEvidence, error) {
	var (
		taskEvidenceList []TaskEvidence
		db               = mysql.FromContext(ctx).Model(&TaskEvidence{})
	)
	if filter != "" {
		db.Where(filter).Find(&taskEvidenceList)
	} else {
		db.Find(&taskEvidenceList)
	}
	return taskEvidenceList, nil
}

// GetTaskEvidenceByTaskIds 通过任务id获取证据
func (t *TaskEvidence) GetTaskEvidenceByTaskIds(ctx context.Context, taskId []int) []TaskEvidence {
	var (
		taskEvidence []TaskEvidence
		db           = mysql.FromContext(ctx).Model(&TaskEvidence{})
	)
	db.Where("task_id in ?", taskId).Find(&taskEvidence)
	return taskEvidence
}
