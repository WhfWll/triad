package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

// TaskResult 任务检测结果表
type TaskResult struct {
	ID            int       `gorm:"column:id;primary_key" json:"id"`             // 主键
	TaskID        int       `gorm:"column:task_id" json:"taskID"`                // 所属任务id
	TargetID      int       `gorm:"column:target_id" json:"targetID"`            // 目标id
	NodeID        string    `gorm:"column:node_id" json:"nodeID"`                // 节点id
	FatherID      string    `gorm:"column:father_id" json:"fatherID"`            // 父节点id
	Pocname       string    `gorm:"column:pocname" json:"pocname"`               // 漏洞标识
	DecisionLibId int       `gorm:"column:decision_lib_id" json:"decisionLibId"` // 决策引擎漏洞库ID
	VulName       string    `gorm:"column:vul_name" json:"vulName"`              // 漏洞名称 冗余字段用于构建路径图
	Risk          int       `gorm:"column:risk" json:"risk"`                     // 漏洞风险等级
	Result        string    `gorm:"column:result" json:"result"`                 // 漏洞结果
	Request       string    `gorm:"column:request" json:"request"`               // 请求报文
	Response      string    `gorm:"column:response" json:"response"`             // 响应报文
	Checked       int       `gorm:"column:checked" json:"Checked"`               // 检测结果 0未检测成功 1检测成功
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`        // 创建时间
}

// TableName 任务检测结果表
func (t *TaskResult) TableName() string {
	return "task_result"
}

// GetTaskResultList 返回任务结果列表
func (t *TaskResult) GetTaskResultList(ctx context.Context, taskID, riskType, page, limit int, search, targetIDs string) ([]TaskResult, int64, error) {
	var (
		taskResultList []TaskResult
		count          int64
		db             = mysql.FromContext(ctx).Model(&TaskResult{})
		query          string
		args           []interface{}
	)
	query += "1 = 1"
	//query += "task_id = ?"
	//args = append(args, taskID)
	if riskType != 0 {
	}
	if targetIDs != "" {
		query += "target_id in ( ? )"
		args = append(args, targetIDs)
	}
	if len(search) > 0 {
		query += " and result LIKE ?"
		args = append(args, "%"+search+"%")
	}
	query += " and ( pocname LIKE '[DATA]' or pocname LIKE '[RCE]' or pocname LIKE '[BF]' or pocname LIKE '[IL]' or pocname LIKE '[FL]' )"
	db.Where(query, args...).Limit(limit).Offset(limit * (page - 1)).Find(&taskResultList)
	db.Count(&count)
	return taskResultList, count, nil
}

// Get retrieves a single record of taskResult from database
func (t *TaskResult) GetTaskResult(ctx context.Context) (TaskResult, error) {
	var (
		taskResult TaskResult
		err        error
		db         = mysql.FromContext(ctx).Model(&TaskResult{})
	)

	curErr := db.Where("id = ?", t.ID).First(&taskResult).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return taskResult, err
}

// Add persists taskResult to database
func (t *TaskResult) AddTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskResult{})

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// Update changes taskResult by id
func (t *TaskResult) UpdateTaskResult(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&TaskResult{})

	if err := db.Where("id = ?", t.ID).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

// UpdateRiskByNodeID updates the risk of a taskResult by its node_id
func (t *TaskResult) UpdateRiskByNodeID(ctx context.Context, nodeID string, risk int) error {
	var db = mysql.FromContext(ctx).Model(&TaskResult{})
	if err := db.Where("node_id = ?", nodeID).Update("risk", risk).Error; err != nil {
		return err
	}
	return nil
}

// AllByTargetId retrieves all taskResult records by targetId
func (t *TaskResult) AllByTargetId(ctx context.Context, targetId int, pocname string) []TaskResult {
	var (
		taskResultList []TaskResult
		db             = mysql.FromContext(ctx).Model(&TaskResult{})
	)
	query := db.Where("target_id = ?", targetId)
	if pocname != "" {
		query = query.Where("pocname = ?", pocname)
	}
	query.Find(&taskResultList)
	return taskResultList
}

// DeleteByTaskIds deletes taskResult records by taskIds
func (t *TaskResult) DeleteByTaskIds(ctx context.Context, taskIds []int) error {
	var db = mysql.FromContext(ctx).Model(&TaskResult{})
	if err := db.Where("task_id IN ?", taskIds).Delete(&TaskResult{}).Error; err != nil {
		return err
	}
	return nil
}

// GetByTargetIdAndNodeId retrieves a single record of taskResult from database by targetId and nodeId
func (t *TaskResult) GetByTargetIdAndNodeId(ctx context.Context, targetId int, nodeId string) (TaskResult, error) {
	var (
		taskResult TaskResult
		err        error
		db         = mysql.FromContext(ctx).Model(&TaskResult{})
	)
	err = db.Where("target_id = ? AND node_id = ?", targetId, nodeId).First(&taskResult).Error
	return taskResult, err
}

// DeleteTaskResultByTargetIds deletes taskResult records by targetIds
// targetIds must be []int/[]string slice
func (t *TaskResult) DeleteTaskResultByTargetIds(ctx context.Context, targetIds any) error {
	var db = mysql.FromContext(ctx).Model(&TaskResult{})
	if err := db.Where("target_id IN ?", targetIds).Delete(&TaskResult{}).Error; err != nil {
		return err
	}
	return nil
}
