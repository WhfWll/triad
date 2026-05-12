package mysqls

import (
	"context"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/tools/enums"
	"time"
)

type LlmModel struct {
	ID         int       `gorm:"column:id;primary_key" json:"id"`      // 主键
	ModelName  string    `gorm:"column:model_name" json:"modelName"`   // 模型名称
	Platform   string    `gorm:"column:platform" json:"platform"`      // 平台类型：openai、baidu、ali等
	ApiUrl     string    `gorm:"column:api_url" json:"apiUrl"`         // API地址
	ApiKey     string    `gorm:"column:api_key" json:"apiKey"`         // API密钥
	ModelID    string    `gorm:"column:model_id" json:"modelID"`       // 模型ID
	ModelType  int       `gorm:"column:model_type" json:"modelType"`   // 模型类型 图像模型或者文本模型
	IsDefault  int       `gorm:"column:is_default" json:"isDefault"`   // 是否默认：1-是，2-否
	Status     int       `gorm:"column:status" json:"status"`          // 状态：1-启用，2-禁用
	Enabled    int       `gorm:"column:enabled" json:"enabled"`        // 是否可用
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
}

// TableName 设置表名
func (l *LlmModel) TableName() string {
	return "llm_models"
}

// GetLlmModelList 获取大模型列表
func (l *LlmModel) GetLlmModelList(ctx context.Context, page, size int, search string) ([]LlmModel, int64, error) {
	var (
		llmModels []LlmModel
		total     int64
		offset    = (page - 1) * size
	)

	db := mysql.GetDB().WithContext(ctx).Model(&LlmModel{})

	// 搜索条件
	if search != "" {
		db = db.Where("model_name LIKE ? OR platform LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表数据
	if err := db.Order("create_time DESC").Offset(offset).Limit(size).Find(&llmModels).Error; err != nil {
		return nil, 0, err
	}

	return llmModels, total, nil
}

// GetLlmModelByID 根据ID获取大模型详情
func (l *LlmModel) GetLlmModelByID(ctx context.Context, id int) (LlmModel, error) {
	var llmModel LlmModel
	err := mysql.GetDB().WithContext(ctx).Where("id = ?", id).First(&llmModel).Error
	return llmModel, err
}

// GetLlmModelByName 根据模型名称获取大模型
func (l *LlmModel) GetLlmModelByName(ctx context.Context, modelName string) (LlmModel, error) {
	var llmModel LlmModel
	err := mysql.GetDB().WithContext(ctx).Where("model_name = ?", modelName).First(&llmModel).Error
	return llmModel, err
}

// GetDefaultLlmModel 获取默认大模型
func (l *LlmModel) GetDefaultLlmModel(ctx context.Context, modelType int) (LlmModel, error) {
	var llmModel LlmModel
	err := mysql.GetDB().WithContext(ctx).
		Where("is_default = ? AND status = ? AND model_type = ?",
			enums.LlmModelIsDefault,
			enums.LlmModelStatusEnable,
			modelType).Find(&llmModel).Error
	if err != nil {
		return llmModel, err
	}
	if llmModel.ID == 0 {
		err = mysql.GetDB().WithContext(ctx).Where("status = ? AND model_type = ?",
			enums.LlmModelStatusEnable,
			modelType).Find(&llmModel).Error
	}
	return llmModel, err
}

// AddLlmModel 添加大模型
func (l *LlmModel) AddLlmModel(ctx context.Context) error {
	return mysql.GetDB().WithContext(ctx).Create(l).Error
}

// UpdateLlmModel 更新大模型
func (l *LlmModel) UpdateLlmModel(ctx context.Context) error {
	return mysql.GetDB().WithContext(ctx).Model(l).Where("id = ?", l.ID).Updates(l).Error
}

// DeleteLlmModels 批量删除大模型
func (l *LlmModel) DeleteLlmModels(ctx context.Context, ids []int) error {
	return mysql.GetDB().WithContext(ctx).Where("id IN ?", ids).Delete(&LlmModel{}).Error
}

// UpdateDefaultStatus 更新默认状态
func (l *LlmModel) UpdateDefaultStatus(ctx context.Context, id int, isDefault int) error {
	return mysql.GetDB().WithContext(ctx).Model(&LlmModel{}).Where("id = ?", id).Update("is_default", isDefault).Error
}

// ClearAllDefault 清除所有默认状态
func (l *LlmModel) ClearAllDefault(ctx context.Context, modelType int) error {
	return mysql.GetDB().WithContext(ctx).
		Model(&LlmModel{}).Where("is_default = ? AND model_type = ?", enums.LlmModelIsDefault, modelType).
		Update("is_default", enums.LlmModelNotDefault).Error
}
