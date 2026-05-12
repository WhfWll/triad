package typespec

// AiScenarioListReq AI应用场景列表请求
type AiScenarioListReq struct {
}

// AiScenarioListRes AI应用场景列表响应
type AiScenarioListRes struct {
	List []AiScenarioItem `json:"list"`
}

// AiScenarioItem AI应用场景项
type AiScenarioItem struct {
	ID          int    `json:"id"`          // 场景ID
	Name        string `json:"name"`        // 场景名称
	Description string `json:"description"` // 场景描述
	Icon        string `json:"icon"`        // 场景图标
	IsEnabled   int    `json:"is_enabled"`  // 是否启用：1-启用，2-禁用
	LlmModelID  int    `json:"llm_model_id"` // 关联的大模型ID
	LlmModelName string `json:"llm_model_name"` // 关联的大模型名称
	Prompt      string `json:"prompt"`      // 场景提示词
	CreateTime  string `json:"create_time"` // 创建时间
	UpdateTime  string `json:"update_time"` // 更新时间
}

// AiScenarioConfigReq AI应用场景配置请求
type AiScenarioConfigReq struct {
	ID          int    `form:"id" json:"id" binding:"required"`                 // 场景ID
	Description string `form:"description" json:"description"`                 // 场景描述
	IsEnabled   int    `form:"is_enabled" json:"is_enabled" binding:"required"` // 是否启用：1-启用，2-禁用
	LlmModelID  int    `form:"llm_model_id" json:"llm_model_id"`               // 关联的大模型ID
	Prompt      string `form:"prompt" json:"prompt"`                           // 场景提示词
}

// AiScenarioConfigRes AI应用场景配置响应
type AiScenarioConfigRes struct {
}