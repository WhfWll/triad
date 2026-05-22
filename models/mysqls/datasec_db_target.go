package mysqls

import (
	"context"
	"time"

	"gitlabee.4dogs.cn/common/mysql"
)

type DatasecDBTarget struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	UserID     int       `gorm:"column:user_id" json:"userId"`
	Name       string    `gorm:"column:name" json:"name"`
	GroupName  string    `gorm:"column:group_name" json:"groupName"`
	DBType     int       `gorm:"column:db_type" json:"dbType"`
	DBHost     string    `gorm:"column:db_host" json:"dbHost"`
	DBPort     int       `gorm:"column:db_port" json:"dbPort"`
	DBName     string    `gorm:"column:db_name" json:"dbName"`
	DBUser     string    `gorm:"column:db_user" json:"dbUser"`
	DBPassword string    `gorm:"column:db_password" json:"dbPassword"`
	Remark     string    `gorm:"column:remark" json:"remark"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (DatasecDBTarget) TableName() string {
	return "datasec_db_target"
}

func (m *DatasecDBTarget) Create(ctx context.Context, row *DatasecDBTarget) error {
	now := time.Now()
	row.CreateTime = now
	row.UpdateTime = now
	return mysql.FromContext(ctx).Model(m).Create(row).Error
}

func (m *DatasecDBTarget) Update(ctx context.Context, row *DatasecDBTarget) error {
	row.UpdateTime = time.Now()
	return mysql.FromContext(ctx).Model(m).Where("id = ?", row.ID).Select("*").Omit("id", "create_time", "user_id").Updates(row).Error
}

func (m *DatasecDBTarget) Delete(ctx context.Context, id, userID int) error {
	db := mysql.FromContext(ctx).Model(m).Where("id = ?", id)
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	return db.Delete(&DatasecDBTarget{}).Error
}

func (m *DatasecDBTarget) GetByID(ctx context.Context, id, userID int) (*DatasecDBTarget, error) {
	var row DatasecDBTarget
	db := mysql.FromContext(ctx).Model(m).Where("id = ?", id)
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if err := db.First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *DatasecDBTarget) ListByUser(ctx context.Context, userID, dbType int, search string, page, size int) ([]DatasecDBTarget, int64, error) {
	var list []DatasecDBTarget
	var total int64
	db := mysql.FromContext(ctx).Model(m)
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if dbType > 0 {
		db = db.Where("db_type = ?", dbType)
	}
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("name LIKE ? OR db_host LIKE ? OR group_name LIKE ? OR db_name LIKE ?", like, like, like, like)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	err := db.Order("id desc").Limit(size).Offset((page - 1) * size).Find(&list).Error
	return list, total, err
}

func (m *DatasecDBTarget) ListByIDs(ctx context.Context, userID int, ids []int) ([]DatasecDBTarget, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []DatasecDBTarget
	db := mysql.FromContext(ctx).Model(m).Where("id IN ?", ids)
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	err := db.Find(&list).Error
	return list, err
}
