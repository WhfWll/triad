package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
)

//具体用法参考gorm官网文档：https://v1.gorm.io/zh_CN/docs/index.html
type Class struct {
	ID     int    `gorm:"primarykey" json:"id"`
	Estate string `gorm:"column:estate" json:"estate"`
	Name   string `gorm:"column:name" json:"name"`
	Nums   int    `gorm:"nums" json:"nums"`
}

func (s *Class) TableName() string {
	return "class"
}

//新增
func (s *Class) AddClass(ctx context.Context) (int, error) {
	var db = mysql.FromContext(ctx).Model(&Class{})
	if err := db.Create(s).Error; err != nil {
		return 0, err
	}
	return s.ID, nil
}
