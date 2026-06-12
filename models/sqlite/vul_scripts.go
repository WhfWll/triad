package sqlite

import (
	"context"
)

// VulScript 对应 scanner.db 中的 vul_scripts 表
type VulScript struct {
	ID           int    `gorm:"column:id;primaryKey" json:"id"`
	UserID       int    `gorm:"column:user_id" json:"user_id"`
	ScriptName   string `gorm:"column:script_name" json:"script_name"`
	Type         string `gorm:"column:type" json:"type"`
	LibName      string `gorm:"column:lib_name" json:"lib_name"`
	Content      string `gorm:"column:content" json:"content"`
	VulID        string `gorm:"column:vul_id" json:"vul_id"`
	VerifyType   string `gorm:"column:verify_type" json:"verify_type"`
	Params       string `gorm:"column:params" json:"params"`
	CreateTime   string `gorm:"column:create_time" json:"create_time"`
	UpdateTime   string `gorm:"column:update_time" json:"update_time"`
	EvidenceType string `gorm:"column:evidence_type" json:"evidence_type"`
	Component    string `gorm:"column:component" json:"component"`
	Product      string `gorm:"column:product" json:"product"`
}

// TableName 指定表名
func (v *VulScript) TableName() string {
	return "vul_scripts"
}

// GetScriptsByFuzzyName 根据名称模糊匹配脚本名
// name通常是组件名或产品名
func GetScriptsByFuzzyName(ctx context.Context, name string) ([]VulScript, error) {
	db, err := GetScannerDB()
	if err != nil {
		return nil, err
	}

	var scripts []VulScript
	// 模糊匹配 script_name
	if err := db.WithContext(ctx).Where("script_name LIKE ?", "%"+name+"%").Find(&scripts).Error; err != nil {
		return nil, err
	}
	return scripts, nil
}

// GetScriptByScriptName 根据脚本名称精确匹配获取单个脚本
func GetScriptByScriptName(ctx context.Context, scriptName string) (*VulScript, error) {
	db, err := GetScannerDB()
	if err != nil {
		return nil, err
	}
	var script VulScript
	if err := db.WithContext(ctx).Where("script_name = ?", scriptName).First(&script).Error; err != nil {
		return nil, err
	}
	return &script, nil
}
