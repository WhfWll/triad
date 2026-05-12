package typespec

import "time"

// LlmModelItem 大模型项目结构
type LlmModelItem struct {
	ID            int       `json:"id"`        // 主键
	ModelName     string    `json:"modelName"` // 模型名称
	Platform      string    `json:"platform"`  // 平台类型
	PlatformDesc  string    `json:"platformDesc"`
	ApiUrl        string    `json:"apiUrl"`  // API地址
	ApiKey        string    `json:"apiKey"`  // API密钥
	ModelID       string    `json:"modelID"` // 模型ID
	ModelType     int       `json:"modelType"`
	ModelTypeDesc string    `json:"modelTypeDesc"`
	IsDefault     int       `json:"isDefault"`  // 是否默认
	Status        int       `json:"status"`     // 状态
	CreateTime    time.Time `json:"createTime"` // 创建时间
	UpdateTime    time.Time `json:"updateTime"` // 更新时间
}

// LlmModelListReq 大模型列表请求
type LlmModelListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"` // 页码
	Size   int    `form:"size" json:"size" binding:"required"` // 每页数量
	Search string `form:"search" json:"search"`                // 搜索关键词
}

// LlmModelListRes 大模型列表响应
type LlmModelListRes struct {
	List  []LlmModelItem `json:"list"`  // 大模型列表
	Total int64          `json:"total"` // 总数
}

// LlmModelDetailReq 大模型详情请求
type LlmModelDetailReq struct {
	ID int `form:"id" json:"id" binding:"required"` // 大模型ID
}

// LlmModelDetailRes 大模型详情响应
type LlmModelDetailRes struct {
	ID            int       `json:"id"`        // 主键
	ModelName     string    `json:"modelName"` // 模型名称
	Platform      string    `json:"platform"`  // 平台类型
	PlatformDesc  string    `json:"platformDesc"`
	ApiUrl        string    `json:"apiUrl"`    // API地址
	ApiKey        string    `json:"apiKey"`    // API密钥
	ModelID       string    `json:"modelID"`   // 模型ID
	ModelType     int       `json:"modelType"` // 模型类型
	ModelTypeDesc string    `json:"modelTypeDesc"`
	IsDefault     int       `json:"isDefault"`  // 是否默认
	Status        int       `json:"status"`     // 状态
	CreateTime    time.Time `json:"createTime"` // 创建时间
	UpdateTime    time.Time `json:"updateTime"` // 更新时间
}

// LlmModelSaveReq 保存大模型请求（添加/编辑）
type LlmModelSaveReq struct {
	ID        int    `form:"id" json:"id"`                                  // 大模型ID（编辑时必填，添加时不填）
	ModelName string `form:"modelName" json:"modelName" binding:"required"` // 模型名称
	Platform  string `form:"platform" json:"platform" binding:"required"`   // 平台类型
	ApiUrl    string `form:"apiUrl" json:"apiUrl" binding:"required"`       // API地址
	ApiKey    string `form:"apiKey" json:"apiKey" binding:"required"`       // API密钥
	ModelID   string `form:"modelID" json:"modelID" binding:"required"`     // 模型ID
	ModelType int    `form:"modelType" json:"modelType" binding:"required"` // 模型类型 文本或者图像
	IsDefault int    `form:"isDefault" json:"isDefault"`                    // 是否默认：1-是，2-否
	Status    int    `form:"status" json:"status"`                          // 状态：1-启用，2-禁用
}

// LlmModelSaveRes 保存大模型响应
type LlmModelSaveRes struct {
	Success bool `json:"success"` // 是否成功
}

// LlmModelDeleteReq 删除大模型请求
type LlmModelDeleteReq struct {
	IDs string `form:"ids" json:"ids" binding:"required"` // 大模型ID，支持单个ID或逗号分隔的多个ID，如：1 或 1,2,3
}

// LlmModelDeleteRes 删除大模型响应
type LlmModelDeleteRes struct {
	Success bool `json:"success"` // 是否成功
}

// LlmModelSetDefaultReq 设置默认大模型请求
type LlmModelSetDefaultReq struct {
	ID int `form:"id" json:"id" binding:"required"` // 大模型ID
}

// LlmModelSetDefaultRes 设置默认大模型响应
type LlmModelSetDefaultRes struct {
	Success bool `json:"success"` // 是否成功
}

type LlmModelEnabledTestReq struct {
	Id int `form:"id" json:"id" binding:"required"` // 大模型ID
}

type LlmModelEnabledTestRes struct {
	Enabled     int    `json:"enabled"`
	EnabledDesc string `json:"enabledDesc"`
	ErrorMsg    string `json:"errorMsg"`
}

type LlmModelEnumsResp struct {
	ModelType  interface{} `json:"modelType"`
	Platform   interface{} `json:"platform"`
	Enhance    interface{} `json:"enhance"`
	DefaultLlm interface{} `json:"defaultLlm"`
}

type LlmModelEnhancementDetailResp struct {
	Status     int    `json:"status"`
	StatusDesc string `json:"statusDesc"`
}

type LlmModelEnhancementEditReq struct {
	Status int `form:"status" json:"status" binding:"required"`
}
