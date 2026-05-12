package mysqls

import (
	"context"
	"fmt"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"smart/tools/enums"
	"strings"
	"time"
)

type Assetgroup struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	Pid        int       `gorm:"column:pid" json:"pid"`                // 父级id
	Level      int       `gorm:"column:level" json:"level"`            // 级别，1-一级资产，2-二级资产...最多六级
	Name       string    `gorm:"column:name" json:"name"`              // 资产组名称
	Remark     string    `gorm:"column:remark" json:"remark"`          // 说明
	UserID     int       `gorm:"column:user_id" json:"userID"`         // 提交者id
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 修改时间
}

// TableName sets insert table name for this struct type
func (a *Assetgroup) TableName() string {
	return "asset_group"
}

// 查询所有资产组数据
func (a *Assetgroup) GetAllAssetgroup(ctx context.Context) []Assetgroup {
	var (
		assetgroupList []Assetgroup
		db             = mysql.FromContext(ctx).Model(&Assetgroup{})
	)
	db.Order("id desc").Find(&assetgroupList)
	return assetgroupList
}

// Get retrieves a single record of assetgroup from database
func (a *Assetgroup) GetAssetgroup(ctx context.Context) (Assetgroup, error) {
	var (
		assetgroup Assetgroup
		err        error
		db         = mysql.FromContext(ctx).Model(&Assetgroup{})
	)

	curErr := db.Where("id = ?", a.ID).First(&assetgroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return assetgroup, err
}

// Add persists assetgroup to database
func (a *Assetgroup) AddAssetgroup(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetgroup{})

	if err := db.Create(a).Error; err != nil {
		return err
	}

	return nil
}

// AddAndReturnID 创建资产组并返回id
func (a *Assetgroup) AddAndReturnID(ctx context.Context) (int, error) {
	db := mysql.FromContext(ctx).Model(&Assetgroup{})
	if err := db.Create(a).Error; err != nil {
		return 0, err
	}
	return a.ID, nil
}

// Update changes assetgroup by id
func (a *Assetgroup) UpdateAssetgroup(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Assetgroup{})

	if err := db.Where("id = ?", a.ID).Updates(a).Error; err != nil {
		return err
	}

	return nil
}

// DeleteById 通过主键ID删除资产组
func (a *Assetgroup) DeleteById(ctx context.Context, ids []int) error {
	return mysql.FromContext(ctx).Model(&Assetgroup{}).Where("id in ?", ids).Delete(Assetgroup{}).Error
}

// GetSubAssetGroup 获取子资产组信息
func (a *Assetgroup) GetSubAssetGroup(ctx context.Context) ([]Assetgroup, error) {
	var (
		assetgroup []Assetgroup
		err        error
		db         = mysql.FromContext(ctx).Model(&Assetgroup{})
	)
	curErr := db.Where("pid = ?", a.Pid).Find(&assetgroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetgroup, err
}

// GetAssetGroupInfo 获取默认资产组信息
func (a *Assetgroup) GetAssetGroupInfo(ctx context.Context, assetGroupID int) (Assetgroup, error) {
	var (
		assetgroup Assetgroup
		err        error
		db         = mysql.FromContext(ctx).Model(&Assetgroup{})
	)
	curErr := db.Where("id = ?", assetGroupID).First(&assetgroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetgroup, err
}

// GetDefaultAssetGroupInfo 获取默认资产组信息
func (a *Assetgroup) GetDefaultAssetGroupInfo(ctx context.Context) (Assetgroup, error) {
	var (
		assetgroup Assetgroup
		err        error
		db         = mysql.FromContext(ctx).Model(&Assetgroup{})
	)
	curErr := db.Where("name = ?", "默认组").First(&assetgroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetgroup, err
}

// GetAssetGroupInfoByPid 获取默认资产组信息
func (a *Assetgroup) GetAssetGroupInfoByPid(ctx context.Context, pid int) ([]Assetgroup, error) {
	var (
		assetgroup []Assetgroup
		err        error
		db         = mysql.FromContext(ctx).Model(&Assetgroup{})
	)
	curErr := db.Where("pid = ?", pid).Find(&assetgroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetgroup, err
}

// GetAssetGroupByName 通过资产组名查找
func (a *Assetgroup) GetAssetGroupByName(ctx context.Context) (Assetgroup, error) {
	var (
		assetgroup Assetgroup
		err        error
		db         = mysql.FromContext(ctx).Model(&Assetgroup{})
	)

	curErr := db.Where("name = ?", a.Name).First(&assetgroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return assetgroup, err
}

// GetAssetGroupByNamePid 通过资产组名查找
func (a *Assetgroup) GetAssetGroupByNamePid(ctx context.Context) (Assetgroup, error) {
	var (
		assetgroup Assetgroup
		err        error
		db         = mysql.FromContext(ctx).Model(&Assetgroup{})
	)

	curErr := db.Where("name = ? and pid = ? ", a.Name, a.Pid).First(&assetgroup).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return assetgroup, err
}

// GetAssetGroupList 获取资产组列表
func (a *Assetgroup) GetAssetGroupList(ctx context.Context) ([]Assetgroup, error) {
	var (
		assetGroupList []Assetgroup
		err            error
		db             = mysql.FromContext(ctx).Model(&Assetgroup{})
	)
	curErr := db.Where("name = ?", a.Name).Find(&assetGroupList).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return assetGroupList, err
}

// GetAssetGroupBreadcrumb 根据资产组ID返回层级路径（面包屑）
func (a *Assetgroup) GetAssetGroupBreadcrumb(ctx context.Context, id int) (string, error) {
	var breadcrumb []string
	currentID := id
	for {
		var (
			db = mysql.FromContext(ctx).Model(&Assetgroup{})
			ag Assetgroup
		)
		err := db.Where("id = ?", currentID).First(&ag).Error
		if err != nil {
			return "", err
		}
		breadcrumb = append([]string{ag.Name}, breadcrumb...)
		if ag.Pid == 0 || ag.ID == ag.Pid {
			break
		}
		currentID = ag.Pid
	}
	return strings.Join(breadcrumb, " - "), nil
}

// GetOrCreateAssetGroupByPath 根据路径获取或创建资产组，返回最末级资产组的 ID
func (a *Assetgroup) GetOrCreateAssetGroupByPath(ctx context.Context, path string) (int, error) {
	if path == "" {
		return enums.DefaultAssetGroup, nil
	}
	names := strings.Split(path, "/")
	pid := 0
	var (
		currentID int
	)
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var group Assetgroup // 每次新建变量，避免复用
		err := mysql.FromContext(ctx).Model(&Assetgroup{}).
			Where("name = ? AND pid = ?", name, pid).
			First(&group).Error
		if err == gorm.ErrRecordNotFound {
			fmt.Println("add groupname" + name)
			newGroup := Assetgroup{
				Pid:        pid,
				Name:       name,
				Level:      a.getAssetGroupLevelByPid(ctx, pid),
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			if err := mysql.FromContext(ctx).Model(&Assetgroup{}).Create(&newGroup).Error; err != nil {
				return enums.DefaultAssetGroup, err
			}
			currentID = newGroup.ID
		} else if err != nil {
			return enums.DefaultAssetGroup, err
		} else {
			currentID = group.ID
		}
		pid = currentID
	}
	return pid, nil
}

func (a *Assetgroup) getAssetGroupLevelByPid(ctx context.Context, pid int) int {
	if pid == 0 {
		return 1
	}
	var parent Assetgroup
	err := mysql.FromContext(ctx).Model(&Assetgroup{}).Where("id = ?", pid).First(&parent).Error
	if err != nil {
		return 1
	}
	return parent.Level + 1
}
