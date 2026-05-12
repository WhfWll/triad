package gofishtypespec

type GetListInfoReq struct {
	Search string `json:"search" form:"search"`
	Page   int    `json:"page" form:"page" binding:"required"`
	Size   int    `json:"size" form:"size" binding:"required"`
}
