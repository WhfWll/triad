package services

import (
	"context"
	"errors"
	"smart/api/typespec"
	"smart/models/mysqls"
	"time"
)

// AiScenario AI应用场景服务层
type AiScenario struct{}

// GetAiScenarioList 获取AI应用场景列表
func (a *AiScenario) GetAiScenarioList(ctx context.Context) ([]typespec.AiScenarioItem, error) {
	var aiScenarioModel mysqls.AiScenario
	scenarios, err := aiScenarioModel.GetAiScenarioList(ctx)
	if err != nil {
		return nil, err
	}

	var result []typespec.AiScenarioItem
	for _, scenario := range scenarios {
		// 获取关联的大模型名称
		var llmModelName string
		if scenario.LlmModelID > 0 {
			var llmModelModel mysqls.LlmModel
			llmModel, err := llmModelModel.GetLlmModelByID(ctx, scenario.LlmModelID)
			if err == nil {
				llmModelName = llmModel.ModelName
			}
		}

		item := typespec.AiScenarioItem{
			ID:           scenario.ID,
			Name:         scenario.Name,
			Description:  scenario.Description,
			Icon:         scenario.Icon,
			IsEnabled:    scenario.IsEnabled,
			LlmModelID:   scenario.LlmModelID,
			LlmModelName: llmModelName,
			Prompt:       scenario.Prompt,
			CreateTime:   scenario.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:   scenario.UpdateTime.Format("2006-01-02 15:04:05"),
		}
		result = append(result, item)
	}

	return result, nil
}

// UpdateAiScenarioConfig 更新AI应用场景配置
func (a *AiScenario) UpdateAiScenarioConfig(ctx context.Context, id int, description string, isEnabled, llmModelID int, prompt string) error {
	var aiScenarioModel mysqls.AiScenario
	
	// 检查场景是否存在
	scenario, err := aiScenarioModel.GetAiScenarioByID(ctx, id)
	if err != nil {
		return errors.New("AI应用场景不存在")
	}

	// 如果指定了大模型ID，检查大模型是否存在
	if llmModelID > 0 {
		var llmModelModel mysqls.LlmModel
		_, err := llmModelModel.GetLlmModelByID(ctx, llmModelID)
		if err != nil {
			return errors.New("指定的大模型不存在")
		}
	}

	// 更新场景配置
	if description != "" {
		scenario.Description = description
	}
	scenario.IsEnabled = isEnabled
	scenario.LlmModelID = llmModelID
	scenario.Prompt = prompt
	scenario.UpdateTime = time.Now()

	return aiScenarioModel.UpdateAiScenario(ctx, scenario)
}