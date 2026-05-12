package typespec

import "mime/multipart"

type DictionaryEnumRes struct {
	Types   interface{}              `json:"types"`
	Service DictionaryEnumServiceRes `json:"service"`
}

type DictionaryEnumServiceRes struct {
	WeakPass      interface{} `json:"weakPass"`
	Wifi          interface{} `json:"wifi"`
	WebPathScan   interface{} `json:"webPathScan"`
	SubdomainScan interface{} `json:"subdomainScan"`
}

type DictionaryListReq struct {
	Page   int    `form:"page" json:"page" binding:"required"`
	Size   int    `form:"size" json:"size" binding:"required"`
	Search string `form:"search" json:"search"`
	Types  int    `form:"types" json:"types"`
}

type DictionaryListRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type DictionaryDetailReq struct {
	DictId int `json:"dictId" form:"dictId" binding:"required"`
}

type DictionaryDetailRes struct {
	Id         int    `json:"id"`
	Sources    int    `json:"sources"`
	Types      int    `json:"types"`
	Service    int    `json:"service"`
	IsDefault  int    `json:"isDefault"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type DictionaryDeleteReq struct {
	DictIds string `json:"dictIds" form:"dictIds" binding:"required"`
}

type DictionaryDeleteRes struct {
}

type DictionarySetDefaultReq struct {
	DictId  int `json:"dictId" form:"dictId" binding:"required"`
	Service int `json:"service" form:"service" binding:"required"`
	Types   int `json:"types" form:"types" binding:"required"`
}

type DictionarySetDefaultRes struct {
}

type DictionaryAddOrEditReq struct {
	Id      int                   `json:"id" form:"id"` // 0 - 新增，非0 - 修改
	Sources int                   `json:"sources" form:"sources" binding:"required"`
	Types   int                   `json:"types" form:"types" binding:"required"`
	Service int                   `json:"service" form:"service" binding:"required"`
	Name    string                `json:"name" form:"name" binding:"required"`
	Content string                `json:"content" form:"content"`
	File    *multipart.FileHeader `json:"file" form:"file"`
}

type DictionaryAddOrEditRes struct {
}
