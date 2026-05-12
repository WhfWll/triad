package application

import (
	"context"
	"smart/api/typespec"
	"smart/services"
)

// AiScenario AI应用场景应用层
type AiScenario struct{}

// AiScenarioList AI应用场景列表
func (a *AiScenario) AiScenarioList(ctx context.Context) ([]typespec.AiScenarioItem, error) {
	var aiScenarioService services.AiScenario
	return aiScenarioService.GetAiScenarioList(ctx)
}

// AiScenarioConfig AI应用场景配置
func (a *AiScenario) AiScenarioConfig(ctx context.Context, id int, description string, isEnabled, llmModelID int, prompt string) error {
	var aiScenarioService services.AiScenario
	return aiScenarioService.UpdateAiScenarioConfig(ctx, id, description, isEnabled, llmModelID, prompt)
}