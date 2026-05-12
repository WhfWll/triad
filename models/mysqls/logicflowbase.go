package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

type LogicFlowBase struct {
	ID              int       `gorm:"column:id;primary_key" json:"id"`                 // 主键
	TaskID          int       `gorm:"column:task_id" json:"taskID"`                    // 所属流量分析任务id
	TargetID        int       `gorm:"column:target_id" json:"targetID"`                // 所属流量分析目标id
	Hash            string    `gorm:"column:hash" json:"hash"`                         // 唯一标识，用于区分数据是否重复
	YakID           int       `gorm:"column:yak_id" json:"yakID"`                      // yak数据ID，用于增量数据标识
	Url             string    `gorm:"column:url" json:"url"`                           // 请求url
	Host            string    `gorm:"column:host" json:"host"`                         // 请求地址
	IP              string    `gorm:"column:ip" json:"ip"`                             // IP地址
	Method          string    `gorm:"column:method" json:"method"`                     // 请求方法
	Protocol        int       `gorm:"column:protocol" json:"protocol"`                 // 网络协议，1-http，2-https
	RespTitle       string    `gorm:"column:resp_title" json:"respTitle"`              // 响应title
	Tags            string    `gorm:"column:tags" json:"tags"`                         // 标签
	RespCode        string    `gorm:"column:resp_code" json:"respCode"`                // 响应码
	RespContentType string    `gorm:"column:resp_content_type" json:"respContentType"` // 响应类型
	ReqHeader       string    `gorm:"column:req_header" json:"reqHeader"`              // 源请求头
	RespHeader      string    `gorm:"column:resp_header" json:"respHeader"`            // 源响应头
	RespContent     string    `gorm:"column:resp_content" json:"respContent"`          // 源响应数据
	LikeField       string    `gorm:"column:like_field" json:"likeField"`              // 冗余like字段
	CreateTime      time.Time `gorm:"column:create_time" json:"createTime"`            // 创建时间
	UpdateTime      time.Time `gorm:"column:update_time" json:"updateTime"`            // 修改时间
}

// TableName sets insert table name for this struct type
func (l *LogicFlowBase) TableName() string {
	return "logic_flow_base"
}

// Get retrieves a list of flowbase from database
func (l *LogicFlowBase) GetFlowbaseList(ctx context.Context, search string, flowTaskId int, page, limit int) ([]LogicFlowBase, int64, error) {
	var (
		flowbaseList []LogicFlowBase
		count        int64
		db           = mysql.FromContext(ctx).Model(&LogicFlowBase{})
		query        string
		args         []interface{}
	)
	query += "task_id = ?"
	args = append(args, flowTaskId)
	if len(search) > 0 {
		query += " and like_field LIKE ?"
		args = append(args, "%"+search+"%")
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&flowbaseList)
	return flowbaseList, count, nil
}

// Get retrieves a single record of flowbase from database
func (l *LogicFlowBase) GetFlowbase(ctx context.Context, id int) (LogicFlowBase, error) {
	var (
		flowbase LogicFlowBase
		err      error
		db       = mysql.FromContext(ctx).Model(&LogicFlowBase{})
	)

	curErr := db.Where("id = ?", id).First(&flowbase).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return flowbase, err
}

// Add persists flowbase to database
func (l *LogicFlowBase) AddFlowbase(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&LogicFlowBase{})

	if err := db.Create(l).Error; err != nil {
		return err
	}

	return nil
}

// Update changes flowbase by id
func (l *LogicFlowBase) UpdateFlowbase(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&LogicFlowBase{})

	if err := db.Where("id = ?", l.ID).Updates(l).Error; err != nil {
		return err
	}

	return nil
}

// DeleteFlowbase 根据id删除
// ids必须为[]int/[]string
func (l *LogicFlowBase) DeleteFlowbase(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&LogicFlowBase{})
	if err := db.Where("id IN ?", ids).Delete(&LogicFlowBase{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFlowbase 根据id删除
// ids必须为[]int/[]string
func (l *LogicFlowBase) DeleteFlowBaseByTaskIds(ctx context.Context, flowTaskIds any) error {
	var db = mysql.FromContext(ctx).Model(&LogicFlowBase{})
	if err := db.Where("task_id IN ?", flowTaskIds).Delete(&LogicFlowBase{}).Error; err != nil {
		return err
	}
	return nil
}

func (l *LogicFlowBase) GetFlowBaseByTaskIdsTargetIdUrlMethod(ctx context.Context, flowTaskId int, hash, method, url string, tags []string) error {
	var db = mysql.FromContext(ctx).Model(&LogicFlowBase{})
	var param = map[string]interface{}{
		"tags": strings.Join(tags, ","),
	}
	if err := db.Where("hash = ? and task_id = ? and url = ? and method = ?", hash, flowTaskId, url, method).Updates(param).Error; err != nil {
		return err
	}
	return nil
}
