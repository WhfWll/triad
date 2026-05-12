package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Flowrisk struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`                 // 主键
	FlowTaskID      int       `gorm:"column:flow_task_id" json:"flowTaskID"`           // 所属流量分析任务id
	FlowTargetID    int       `gorm:"column:flow_target_id" json:"flowTargetID"`       // 所属流量分析目标id
	Hash            string    `gorm:"column:hash" json:"hash"`                         // 唯一标识，用于区分数据是否重复
	YakID           int       `gorm:"column:yak_id" json:"yakID"`                      // yak数据ID，用于增量数据标识
	Host            string    `gorm:"column:host" json:"host"`                         // 请求地址
	IP              string    `gorm:"column:ip" json:"ip"`                             // IP地址
	IPInteger       string    `gorm:"column:ip_integer" json:"ipinteger"`              // IP？Yak的记录，不知道是啥
	Port            string    `gorm:"column:port" json:"port"`                         // 端口
	Title           string    `gorm:"column:title" json:"title"`                       // 漏洞标题
	RiskType        string    `gorm:"column:risk_type" json:"riskType"`                // 漏洞类型
	RiskTypeVerbose string    `gorm:"column:risk_type_verbose" json:"riskTypeVerbose"` // 漏洞类型冗长
	Payload         string    `gorm:"column:payload" json:"payload"`                   // 内容
	Detail          string    `gorm:"column:detail" json:"detail"`                     // 详情
	RiskLevel       int       `gorm:"column:risk_level" json:"riskLevel"`              // 风险等级
	Request         string    `gorm:"column:request" json:"request"`                   // 请求
	Response        string    `gorm:"column:response" json:"response"`                 // 响应
	Parameter       string    `gorm:"column:parameter" json:"parameter"`               // 参数
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`            // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`            // 修改时间
}

// TableName sets insert table name for this struct type
func (f *Flowrisk) TableName() string {
	return "flow_risk"
}

//GetFlowriskList 列表查询
func (f *Flowrisk) GetFlowriskList(ctx context.Context, search string, riskLevel int, flowTaskId int, page, limit int) ([]Flowrisk, int64, error) {
	var (
		flowriskList []Flowrisk
		count        int64
		db           = mysql.FromContext(ctx).Model(&Flowrisk{})
		query        string
		args         []interface{}
	)
	query += "flow_task_id = ?"
	args = append(args, flowTaskId)
	if len(search) > 0 {
		query += " and (title LIKE ? or host LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if riskLevel > 0 {
		query += " and risk_level = ?"
		args = append(args, riskLevel)
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&flowriskList)
	return flowriskList, count, nil
}

//GetFlowriskListByTaskId 根据任务id查询数据
func (f *Flowrisk) GetFlowriskListByTaskId(ctx context.Context, flowTaskId int) ([]Flowrisk, error) {
	var (
		flowriskList []Flowrisk
		db           = mysql.FromContext(ctx).Model(&Flowrisk{})
	)
	db.Where("flow_task_id = ?", flowTaskId).Find(&flowriskList)
	return flowriskList, nil
}

//GetFlowriskListByTaskIds 根据任务ids查询数据
func (f *Flowrisk) GetFlowriskListByTaskIds(ctx context.Context, flowTaskIds any) ([]Flowrisk, error) {
	var (
		flowriskList []Flowrisk
		db           = mysql.FromContext(ctx).Model(&Flowrisk{})
	)
	db.Where("flow_task_id IN ?", flowTaskIds).Find(&flowriskList)
	return flowriskList, nil
}

//GetFlowrisk 查询漏洞详情
func (f *Flowrisk) GetFlowrisk(ctx context.Context, id int) (Flowrisk, error) {
	var (
		flowrisk Flowrisk
		err      error
		db       = mysql.FromContext(ctx).Model(&Flowrisk{})
	)

	curErr := db.Where("id = ?", id).First(&flowrisk).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return flowrisk, err
}

// Add persists flowrisk to database
func (f *Flowrisk) AddFlowrisk(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Flowrisk{})

	if err := db.Create(f).Error; err != nil {
		return err
	}

	return nil
}

// Update changes flowrisk by id
func (f *Flowrisk) UpdateFlowrisk(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Flowrisk{})

	if err := db.Where("id = ?", f.ID).Updates(f).Error; err != nil {
		return err
	}

	return nil
}

// DeleteFlowrisk 根据id删除
// ids必须为[]int/[]string
func (f *Flowrisk) DeleteFlowrisk(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&Flowrisk{})
	if err := db.Where("id IN ?", ids).Delete(&Flowrisk{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFlowriskByTaskIds 根据任务id删除
// ids必须为[]int/[]string
func (f *Flowrisk) DeleteFlowriskByTaskIds(ctx context.Context, flowTaskIds any) error {
	var db = mysql.FromContext(ctx).Model(&Flowrisk{})
	if err := db.Where("flow_task_id IN ?", flowTaskIds).Delete(&Flowrisk{}).Error; err != nil {
		return err
	}
	return nil
}
