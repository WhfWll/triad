package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type VulLifecycle struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	PocName    string    `gorm:"column:poc_name" json:"pocName"`  // 漏洞pocname
	VulName    string    `gorm:"column:vul_name" json:"vulName"`  // 漏洞名称
	Location   string    `gorm:"column:location" json:"location"` // 漏洞位置
	Content    string    `gorm:"column:content" json:"content"`   // 周期信息
	FindNum    int       `gorm:"column:find_num" json:"findNum"`  // 发现次数
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName sets insert table name for this struct type
func (v *VulLifecycle) TableName() string {
	return "vul_lifecycle"
}

// Add 添加生命周期记录
func (v *VulLifecycle) Add(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&VulLifecycle{})
	v.CreateTime = time.Now()
	v.UpdateTime = time.Now()
	if err := db.Create(v).Error; err != nil {
		return 0, err
	}
	return v.ID, nil
}

// GetVulLifecycleDetail 通过 poc_name, vul_name, location 获取详情信息
func (v *VulLifecycle) GetVulLifecycleDetail(ctx context.Context, pocName, name, location string) (VulLifecycle, error) {
	var (
		info VulLifecycle
		err  error
		db   = mysql.FromContext(ctx).Model(&VulLifecycle{})
	)
	curErr := db.Where("poc_name = ? AND vul_name = ? AND location = ?", pocName, name, location).
		First(&info).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return info, err
}

// UpdateContentByUniqueKey 更新周期任务 content
func (v *VulLifecycle) UpdateContentByUniqueKey(ctx context.Context, pocName, name, location, content string) error {
	var db = mysql.FromContext(ctx).Model(&VulLifecycle{})
	// 构建更新数据
	updateData := map[string]interface{}{
		"content":     content,
		"update_time": time.Now(),
	}
	err := db.Where("poc_name = ? AND vul_name = ? AND location = ?", pocName, name, location).
		Updates(updateData).Error
	return err
}

// IncrementFindNum 增加发现次数
// 场景：每次扫描发现同一个漏洞时调用，find_num + 1，同时更新 content
func (v *VulLifecycle) IncrementFindNum(ctx context.Context, pocName, name, location, content string) error {
	var db = mysql.FromContext(ctx).Model(&VulLifecycle{})
	updateData := map[string]interface{}{
		"find_num":    gorm.Expr("find_num + ?", 1),
		"content":     content,
		"update_time": time.Now(),
	}
	err := db.Where("poc_name = ? AND vul_name = ? AND location = ?", pocName, name, location).
		Updates(updateData).Error
	return err
}

// UpdateContentById 如果你知道 ID，也可以通过 ID 直接更新 Content
func (v *VulLifecycle) UpdateContentById(ctx context.Context, id int, content string) error {
	var db = mysql.FromContext(ctx).Model(&VulLifecycle{})
	err := db.Where("id = ?", id).
		Updates(map[string]interface{}{
			"content":     content,
			"update_time": time.Now(),
		}).Error

	return err
}

// GetList 获取生命周期列表（带分页和搜索）
// 这是一个通用的查询方法，参考了 VulScripts 的写法
func (v *VulLifecycle) GetList(ctx context.Context, page, limit int, pocName, name string) ([]VulLifecycle, int64, error) {
	var (
		list  []VulLifecycle
		count int64
		db    = mysql.FromContext(ctx).Model(&VulLifecycle{})
	)
	if pocName != "" {
		db = db.Where("poc_name like ?", "%"+pocName+"%")
	}
	if name != "" {
		db = db.Where("vul_name like ?", "%"+name+"%")
	}

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(limit * (page - 1)).
		Order("update_time desc"). // 通常按更新时间倒序
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}
