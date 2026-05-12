package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type Basvul struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`         // 主键
	BasTaskID   int       `gorm:"column:bas_task_id" json:"basTaskID"`     // bas任务ID
	BasTargetID int       `gorm:"column:bas_target_id" json:"basTargetID"` // bas目标ID
	Addr        string    `gorm:"column:addr" json:"addr"`                 // 攻击目标
	RuleID      int       `gorm:"column:rule_id" json:"ruleID"`            // bas目标ID
	RuleName    string    `gorm:"column:rule_name" json:"ruleName"`        // 剧本名称
	AttackMode  int       `gorm:"column:attack_mode" json:"attackMode"`    // 攻击方式
	AttackStage int       `gorm:"column:attack_stage" json:"attackStage"`  // 攻击阶段
	RiskLevel   int       `gorm:"column:risk_level" json:"riskLevel"`      // 任务风险等级，1-高危、2-中危、3-低危、4-安全
	Md5Code     string    `gorm:"column:md5_code" json:"md5Code"`          // md5值，逗号隔开
	Status      int       `gorm:"column:status" json:"status"`             // 攻击结果：1失败 2成功
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`    // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`    // 修改时间
}

func (b *Basvul) TableName() string {
	return "bas_vul"
}

// 根据ip、状态和md5查询数据
func (b *Basvul) GetBasvulByIPAndMd5(ctx context.Context, ip string, status int, md5 []string) ([]Basvul, error) {
	var (
		basvulList []Basvul
		db         = mysql.FromContext(ctx).Model(&Basvul{})
		query      string
		args       []interface{}
	)
	query += "addr = ? and status = ? "
	args = append(args, ip, status)

	for i := 0; i < len(md5); i++ {
		if i == 0 {
			query += " and (md5_code LIKE ?"
			args = append(args, "%"+md5[i]+"%")

		} else {
			query += " or md5_code LIKE ?"
			args = append(args, "%"+md5[i]+"%")
		}
		if i == len(md5)-1 {
			query += ")"
		}
	}

	db.Where(query, args...).Find(&basvulList)
	return basvulList, nil
}

// Get retrieves a single record of basvul from database
func (b *Basvul) GetBasvul(ctx context.Context) (Basvul, error) {
	var (
		basvul Basvul
		err    error
		db     = mysql.FromContext(ctx).Model(&Basvul{})
	)

	curErr := db.Where("id = ?", b.ID).First(&basvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return basvul, err
}

//根据任务id查询风险等级最高的数据
func (b *Basvul) GetBasvulRiskLevelByTaskId(ctx context.Context, basTaskId int, status int) (Basvul, error) {
	var (
		basvul Basvul
		err    error
		db     = mysql.FromContext(ctx).Model(&Basvul{})
	)
	curErr := db.Where("bas_task_id = ? and status = ?", basTaskId, status).Order("risk_level").First(&basvul).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}
	return basvul, err
}

// Add persists basvul to database
func (b *Basvul) AddBasvul(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&Basvul{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// 同时更新多条数据
func (b *Basvul) UpdateBasvul(ctx context.Context, ids []int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Basvul{})
	if err := db.Where("id IN ?", ids).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 更新一条数据
func (b *Basvul) UpdateBasvulById(ctx context.Context, id int, params map[string]interface{}) error {
	var db = mysql.FromContext(ctx).Model(&Basvul{})
	if err := db.Where("id = ?", id).Updates(params).Error; err != nil {
		return err
	}
	return nil
}

// 插入多条bas漏洞检测数据
func (b *Basvul) AddBasvulMany(ctx context.Context, datas []Basvul) error {
	var db = mysql.FromContext(ctx).Model(&Basvul{})
	if err := db.Create(datas).Error; err != nil {
		return err
	}
	return nil
}

// 根据任务id获取漏洞测试数据
func (b *Basvul) GetBasvulByTaskId(ctx context.Context, basTaskID int) []Basvul {
	var (
		basvulList []Basvul
		db         = mysql.FromContext(ctx).Model(&Basvul{})
	)
	db.Where("bas_task_id = ?", basTaskID).Find(&basvulList)
	return basvulList
}

func (b *Basvul) GetBasvulByTargetIds(ctx context.Context, basTargetIDs []int, status int) []Basvul {
	var (
		basvulList []Basvul
		db         = mysql.FromContext(ctx).Model(&Basvul{})
	)
	db.Where("bas_target_id in ? and status = ?", basTargetIDs, status).Find(&basvulList)
	return basvulList
}

// 根据任务id和状态获取漏洞数量
func (b *Basvul) GetBasvulCountByTaskIdAndStatus(ctx context.Context, basTaskID int, status int) int64 {
	var (
		count int64
		db    = mysql.FromContext(ctx).Model(&Basvul{})
	)
	db.Where("bas_task_id = ? and status = ?", basTaskID, status).Count(&count)
	return count
}

type BasVulRisk struct {
	RiskLevel int   `gorm:"column:risk_level" json:"riskLevel"`
	Total     int64 `gorm:"column:total" json:"total"`
}

// 根据任务id和状态获取危险等级的数量
func (b *Basvul) GetBasVulRiskLevelCountByTaskIdAndStatus(ctx context.Context, basTaskID int, status int) []BasVulRisk {
	var (
		basVulRiskResult []BasVulRisk
		db               = mysql.FromContext(ctx).Model(&Basvul{})
	)
	db.Select("risk_level, count(id) as total").Where("bas_task_id = ? and status = ?", basTaskID, status).Group("risk_level").Find(&basVulRiskResult)
	return basVulRiskResult
}

func (b *Basvul) GetBasvulList(ctx context.Context, basTaskID int, page, limit, riskLevel, attackStage, status int, search string, attackMode string) ([]Basvul, int64) {
	var (
		basvulList []Basvul
		db         = mysql.FromContext(ctx).Model(&Basvul{})
		count      int64
		query      string
		args       []interface{}
	)
	query += "bas_task_id = ?"
	args = append(args, basTaskID)
	if riskLevel > 0 {
		query += " and risk_level = ?"
		args = append(args, riskLevel)
	}
	if attackStage > 0 {
		query += " and attack_stage = ?"
		args = append(args, attackStage)
	}
	if len(attackMode) > 0 {
		query += " and attack_mode = ?"
		args = append(args, attackMode)
	}
	if len(search) > 0 {
		query += " and rule_name LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if status != 0 {
		query += " and status = ?"
		args = append(args, status)
	}
	dbs := db.Where(query, args...)
	dbs.Count(&count)
	dbs.Debug().Limit(limit).Offset(limit * (page - 1)).Order("id desc").Find(&basvulList)
	return basvulList, count
}

// 根据id删除漏洞测试数据
func (b *Basvul) DelBasvulByIds(ctx context.Context, ids any) error {
	var db = mysql.FromContext(ctx).Model(&Basvul{})
	if err := db.Where("id IN ?", ids).Delete(&Basvul{}).Error; err != nil {
		return err
	}
	return nil
}

// 根据任务id删除漏洞测试数据
func (b *Basvul) DelBasvulByTaskIds(ctx context.Context, taskIds any) error {
	var db = mysql.FromContext(ctx).Model(&Basvul{})
	if err := db.Where("bas_task_id IN ?", taskIds).Delete(&Basvul{}).Error; err != nil {
		return err
	}
	return nil
}
