package services

import (
	"context"
	"errors"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
)

type LlmModel struct {
}

// GetLlmModelList 获取大模型列表
func (l *LlmModel) GetLlmModelList(ctx context.Context, page, size int, search string) ([]mysqls.LlmModel, int64, error) {
	var llmModelModel mysqls.LlmModel
	return llmModelModel.GetLlmModelList(ctx, page, size, search)
}

// GetLlmModelDetail 获取大模型详情
func (l *LlmModel) GetLlmModelDetail(ctx context.Context, id int) (mysqls.LlmModel, error) {
	var llmModelModel mysqls.LlmModel
	return llmModelModel.GetLlmModelByID(ctx, id)
}

// AddLlmModel 添加大模型
func (l *LlmModel) AddLlmModel(ctx context.Context, modelName, platform, apiUrl, apiKey, modelID string, modelType, isDefault int) error {
	// 检查模型名称是否已存在
	var llmModelModel mysqls.LlmModel
	existModel, err := llmModelModel.GetLlmModelByName(ctx, modelName)
	if err == nil && existModel.ID != 0 {
		return errors.New("模型名称已存在")
	}

	// 如果设置为默认，先清除其他默认状态
	if isDefault == enums.LlmModelIsDefault {
		if err := llmModelModel.ClearAllDefault(ctx, modelType); err != nil {
			return err
		}
	} else {
		defaultModel, _ := llmModelModel.GetDefaultLlmModel(ctx, modelType)
		if defaultModel.ID == 0 {
			isDefault = enums.LlmModelIsDefault
		}
	}

	// 创建新模型
	newModel := mysqls.LlmModel{
		ModelName:  modelName,
		Platform:   platform,
		ApiUrl:     apiUrl,
		ApiKey:     apiKey,
		ModelID:    modelID,
		ModelType:  modelType,
		IsDefault:  isDefault,
		Status:     enums.LlmModelStatusEnable, // 默认启用
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	return newModel.AddLlmModel(ctx)
}

// UpdateLlmModel 更新大模型
func (l *LlmModel) UpdateLlmModel(ctx context.Context, id int, modelName, platform, apiUrl, apiKey, modelID string, modelType, isDefault, status int) error {
	// 检查模型是否存在
	var llmModelModel mysqls.LlmModel
	existModel, err := llmModelModel.GetLlmModelByID(ctx, id)
	if err != nil {
		return errors.New("模型不存在")
	}

	// 检查模型名称是否与其他模型重复
	if modelName != existModel.ModelName {
		nameModel, err := llmModelModel.GetLlmModelByName(ctx, modelName)
		if err == nil && nameModel.ID != 0 && nameModel.ID != id {
			return errors.New("模型名称已存在")
		}
	}

	// 如果设置为默认，先清除其他默认状态
	if isDefault == enums.LlmModelIsDefault {
		if err = llmModelModel.ClearAllDefault(ctx, modelType); err != nil {
			return err
		}
	} else {
		defaultModel, _ := llmModelModel.GetDefaultLlmModel(ctx, modelType)
		if defaultModel.ID == 0 {
			isDefault = enums.LlmModelIsDefault
		}
	}

	// 更新模型信息
	updateModel := mysqls.LlmModel{
		ID:         id,
		ModelName:  modelName,
		Platform:   platform,
		ApiUrl:     apiUrl,
		ApiKey:     apiKey,
		ModelID:    modelID,
		ModelType:  modelType,
		IsDefault:  isDefault,
		Status:     status,
		UpdateTime: time.Now(),
	}

	return updateModel.UpdateLlmModel(ctx)
}

// DeleteLlmModels 批量删除大模型
func (l *LlmModel) DeleteLlmModels(ctx context.Context, idsStr string) error {
	// 解析ID字符串
	idStrs := strings.Split(strings.TrimSpace(idsStr), ",")
	if len(idStrs) == 0 {
		return errors.New("请提供要删除的模型ID")
	}

	var ids []int
	for _, idStr := range idStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return errors.New("ID格式错误：" + idStr)
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return errors.New("请提供有效的模型ID")
	}

	var llmModelModel mysqls.LlmModel

	// 直接删除，忽略不存在的ID
	// 数据库的 WHERE id IN ? 语句会自动忽略不存在的ID
	return llmModelModel.DeleteLlmModels(ctx, ids)
}

// SetDefaultLlmModel 设置默认大模型
func (l *LlmModel) SetDefaultLlmModel(ctx context.Context, id int) error {
	var llmModelModel mysqls.LlmModel

	// 检查模型是否存在且启用
	existModel, err := llmModelModel.GetLlmModelByID(ctx, id)
	if err != nil {
		return errors.New("模型不存在")
	}
	if existModel.Status != 1 {
		return errors.New("只能设置启用状态的模型为默认")
	}

	// 先清除所有默认状态
	if err = llmModelModel.ClearAllDefault(ctx, existModel.ModelType); err != nil {
		return err
	}

	// 设置新的默认模型
	return llmModelModel.UpdateDefaultStatus(ctx, id, enums.LlmModelIsDefault)
}

// GetDefaultLlmModel 获取默认大模型
func (l *LlmModel) GetDefaultLlmModel(ctx context.Context, modelType int) (mysqls.LlmModel, error) {
	var llmModelModel mysqls.LlmModel
	return llmModelModel.GetDefaultLlmModel(ctx, modelType)
}
