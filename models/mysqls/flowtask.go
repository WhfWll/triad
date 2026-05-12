package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Flowtask struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`         // 主键
	TaskName    string    `gorm:"column:task_name" json:"taskName"`        // 任务名称
	NodeID      string    `gorm:"column:node_id" json:"nodeID"`            // 节点ID
	NetwordCard string    `gorm:"column:netword_card" json:"networdCard"`  // 节点网卡
	Port        string    `gorm:"column:port" json:"port"`                 // 节点端口
	ExpireTime  int       `gorm:"column:expire_time" json:"expireTime"`    // 分析结束时间
	Status      int       `gorm:"column:status" json:"status"`             // 任务运行状态，1-已禁用，2-待执行 3-已开始 4-已结束
	UserID      int       `gorm:"column:user_id" json:"userID"`            // 操作人id
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`    // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`    // 修改时间
	OtherConfig string    `gorm:"column:other_config" json:"other_config"` // 其他配置
	VulConfig   string    `gorm:"column:vul_config" json:"vul_config"`     // 漏洞配置
}

// TableName sets insert table name for this struct type
func (f *Flowtask) TableName() string {
	return "flow_task"
}

// GetFlowtaskList 流量分析任务列表
func (f *Flowtask) GetFlowtaskList(ctx context.Context, search string, page, limit int, userIdList []int) ([]Flowtask, int64, error) {
	var (
		flowtaskList []Flowtask
		count        int64
		db           = mysql.FromContext(ctx).Model(&Flowtask{})
		query        string
		args         []interface{}
	)
	query += "1 = 1"
	if len(search) != 0 {
		query += " and task_name LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if len(userIdList) != 0 {
		query += " and user_id in ?"
		args = append(args, userIdList)
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&flowtaskList)
	return flowtaskList, count, nil
}

// GetFlowtaskListByIdsAndStatus 根据id和状态查询数据
func (f *Flowtask) GetFlowtaskListByIdsAndStatus(ctx context.Context, status int, ids any) ([]Flowtask, error) {
	var (
		flowtaskList []Flowtask
		db           = mysql.FromContext(ctx).Model(&Flowtask{})
	)
	db.Where("status = ? and id in ?", status, ids).Find(&flowtaskList)
	return flowtaskList, nil
}

// GetFlowtask 查询任务详情
func (f *Flowtask) GetFlowtask(ctx context.Context, id int) (Flowtask, error) {
	var (
		flowtask Flowtask
		err      error
		db       = mysql.FromContext(ctx).Model(&Flowtask{})
	)
	curErr := db.Where("id = ?", id).First(&flowtask).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return flowtask, err
}

// AddFlowtask 新增一条数据
func (f *Flowtask) AddFlowtask(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Flowtask{})

	if err := db.Create(f).Error; err != nil {
		return err
	}

	return nil
}

// UpdateFlowtask 根据id修改数据
func (f *Flowtask) UpdateFlowtask(ctx context.Context, flowTaskId int, param map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Flowtask{})
	if err := db.Where("id = ?", flowTaskId).Updates(param).Error; err != nil {
		return err
	}

	return nil
}

// DeleteFlowtskByIds 批量删除
func (f *Flowtask) DeleteFlowtaskByIds(ctx context.Context, flowTaskIds any) error {
	var db = mysql.FromContext(ctx).Model(&Flowtask{})
	if err := db.Where("id in ?", flowTaskIds).Delete(f).Error; err != nil {
		return err
	}
	return nil
}
