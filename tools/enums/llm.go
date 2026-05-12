// Package enums
// @Author bcy2007  2025/12/18 15:50
package enums

type llm struct{}

var LlmEnums llm

const (
	LlmModelTypeText = iota + 1
	LlmModelTypeImg
)

const (
	LlmModelIsDefault = iota + 1
	LlmModelNotDefault
)

func (l *llm) LlmModelIsDefaultEnums() interface{} {
	return []struct {
		Label string `json:"label"`
		Value int    `json:"value"`
	}{
		{Label: "默认", Value: LlmModelIsDefault},
		{Label: "非默认", Value: LlmModelNotDefault},
	}
}

const (
	LlmModelStatusEnable = iota + 1
	LlmModelStatusDisable
)

const (
	LlmModelEnabledYes = iota + 1
	LlmModelEnabledNo
)

const (
	LlmModelPlatformOpenAi = "openai"
)

var LlmModelEnableNameMap = map[int]string{
	LlmModelEnabledYes: "连接成功",
	LlmModelEnabledNo:  "连接失败",
}

var LlmModelPlatformNameMap = map[string]string{
	LlmModelPlatformOpenAi: "OpenAI兼容模型",
}

var LlmModelTypeNameMap = map[int]string{
	LlmModelTypeText: "文本模型",
	LlmModelTypeImg:  "图像模型",
}

func (l *llm) GetLlmModelEnableName(enabled int) string {
	if desc, ok := LlmModelEnableNameMap[enabled]; ok {
		return desc
	}
	return "未知状态"
}

func (l *llm) GetLlmModelTypeName(modelType int) string {
	if desc, ok := LlmModelTypeNameMap[modelType]; ok {
		return desc
	}
	return "未知模型类型"
}

func (l *llm) LlmModelTypeEnums() interface{} {
	return []struct {
		Label string `json:"label"`
		Value int    `json:"value"`
	}{
		{Label: "文本模型", Value: LlmModelTypeText},
		{Label: "图像模型", Value: LlmModelTypeImg},
	}
}

func (l *llm) GetLlmModelPlatformName(platform string) string {
	if desc, ok := LlmModelPlatformNameMap[platform]; ok {
		return desc
	}
	return "未知平台"
}

func (l *llm) LlmModelPlatformEnums() interface{} {
	return []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}{
		{Label: "OpenAI兼容模型", Value: LlmModelPlatformOpenAi},
	}
}

const (
	LlmModelEnhancementOpen = iota + 1
	LlmModelEnhancementClose
)

var LlmModelEnhancementNameMap = map[int]string{
	LlmModelEnhancementOpen:  "开启",
	LlmModelEnhancementClose: "关闭",
}

func (l *llm) GetLlmModelEnhancementName(status int) string {
	if desc, ok := LlmModelEnhancementNameMap[status]; ok {
		return desc
	}
	return "未知状态"
}

func (l *llm) LlmModelEnhancementEnums() interface{} {
	return []struct {
		Label string `json:"label"`
		Value int    `json:"value"`
	}{
		{Label: "开启", Value: LlmModelEnhancementOpen},
		{Label: "关闭", Value: LlmModelEnhancementClose},
	}
}
