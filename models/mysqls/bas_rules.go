package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
	"time"
)

type BasRules struct {
	ID                             int       `gorm:"column:id;primary_key" json:"id"`                                                  // 主键
	Content                        string    `gorm:"column:content" json:"content"`                                                    // 规则内容
	Name                           string    `gorm:"column:name" json:"name"`                                                          // 规则名称
	NameZh                         string    `gorm:"column:name_zh" json:"nameZh"`                                                     // 规则名称
	Hash                           string    `gorm:"column:hash" json:"hash"`                                                          // hash值
	ClassType                      string    `gorm:"column:class_type" json:"classType"`                                               // 类型
	Protocol                       string    `gorm:"column:protocol" json:"protocol"`                                                  // 协议名称
	Keywords                       string    `gorm:"column:keywords" json:"keywords"`                                                  // 关键字
	KeywordsZh                     string    `gorm:"column:keywords_zh" json:"keywordsZh"`                                             // 关键字中文描述
	Description                    string    `gorm:"column:description" json:"description"`                                            // 描述
	DescriptionZh                  string    `gorm:"column:description_zh" json:"descriptionZh"`                                       // 描述中文
	Cve                            string    `gorm:"column:cve" json:"cve"`                                                            // cve编号
	RawTrafficBeyondIpPacketBase64 string    `gorm:"column:raw_traffic_beyond_ip_packet_base64" json:"rawTrafficBeyondIpPacketBase64"` // 基于ip数据包的原始流量
	RawTrafficBeyondHttpBase64     string    `gorm:"column:raw_traffic_beyond_http_base64" json:"rawTrafficBeyondHttpBase64"`          // http 之外的原始流量
	CreateTime                     time.Time `gorm:"column:create_time" json:"createTime"`                                             // 创建时间
	UpdateTime                     time.Time `gorm:"column:update_time" json:"updateTime"`                                             // 修改时间
	AttackMode                     int       `gorm:"column:attack_mode" json:"attackMode"`                                             // 攻击方式
	AttackStage                    int       `gorm:"column:attack_stage" json:"attackStage"`                                           // 攻击阶段
	RiskLevel                      int       `gorm:"column:risk_level" json:"riskLevel"`                                               // 任务风险等级，1-高危、2-中危、3-低危、4-安全
	AffectTarget                   string    `gorm:"column:effect_target" json:"affectTarget"`                                         // 影响目标
	AffectScore                    string    `gorm:"column:effect_score" json:"affectSource"`                                          // 影响评分
	RelationAttackMethod           string    `gorm:"column:relation_attack_method" json:"relationAttackMethod"`                        // 关联攻击方式
	RefUrl                         string    `gorm:"column:ref_url" json:"refUrl"`                                                     // 参考链接
	FixSuggest                     string    `gorm:"column:fix_suggest" json:"fixSuggest"`                                             // 修复建议
}

// TableName sets insert table name for this struct type
func (b *BasRules) TableName() string {
	return "bas_rules"
}

// Get retrieves a list of basRules from database
func (b *BasRules) GetBasRulesList(ctx context.Context, page, limit, attackStage int, attackType string, risk int, classType string, search string, ids []string) ([]BasRules, int64, error) {
	var (
		basRulesList []BasRules
		count        int64
		db           = mysql.FromContext(ctx).Model(&BasRules{})
	)
	if len(ids) > 0 {
		db = db.Where("id in ?", ids)
	}
	if search != "" {
		db = db.Where("name_zh like ?", "%"+search+"%")
	}
	if classType != "" {
		db = db.Where("class_type = ?", classType)
	}
	if attackStage != 0 {
		db = db.Where("attack_stage = ?", attackStage)
	}
	if len(attackType) > 0 {
		db = db.Where("attack_mode = ?", attackType)
	}
	if risk != 0 {
		db = db.Where("risk_level = ?", risk)
	}

	db.Count(&count)
	db.Limit(limit).Offset(limit * (page - 1)).Find(&basRulesList)

	return basRulesList, count, nil
}

// Get retrieves a single record of basRules from database
func (b *BasRules) GetBasRulesById(ctx context.Context, id int) (BasRules, error) {
	var (
		basRules BasRules
		err      error
		db       = mysql.FromContext(ctx).Model(&BasRules{})
	)

	curErr := db.Where("id = ?", id).First(&basRules).Error
	if curErr != nil && curErr != gorm.ErrRecordNotFound {
		err = curErr
	}

	return basRules, err
}

// Add persists basRules to database
func (b *BasRules) AddBasRules(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasRules{})

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

// Update changes basRules by id
func (b *BasRules) UpdateBasRules(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasRules{})

	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

// Delete basRules by id
func (b *BasRules) DeleteBasRules(ctx context.Context) error {
	var db = mysql.FromContext(ctx).Model(&BasRules{})

	//b.Estate = "deleted"
	b.UpdateTime = time.Now()
	if err := db.Where("id = ?", b.ID).Updates(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *BasRules) AddAll(ctx context.Context, datas []BasRules) error {
	var db = mysql.FromContext(ctx).Model(&BasRules{})

	if err := db.CreateInBatches(datas, 1000).Error; err != nil {
		return err
	}

	return nil
}

// Get retrieves a list of basRules from database
func (b *BasRules) GetListByNames(ctx context.Context, names []string) []BasRules {
	var (
		basRulesList []BasRules
		db           = mysql.FromContext(ctx).Model(&BasRules{})
	)

	db.Where("name in ?", names).Find(&basRulesList)

	return basRulesList
}

func (b *BasRules) GetByIds(ctx context.Context, ids []int) []BasRules {
	var (
		basRules []BasRules
		db       = mysql.FromContext(ctx).Model(&BasRules{})
	)
	db.Where("id in ?", ids).Find(&basRules)
	return basRules
}
