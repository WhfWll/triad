package typespec

type GetStudentListReq struct {
	Name  string `form:"name" json:"name"`
	Page  int    `form:"page" json:"page" binding:"required"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
}

type GetStudentListResp struct {
	Total int64                     `json:"total"`
	List  []GetStudentListRespItems `json:"list"`
}

type GetStudentListRespItems struct {
	Id   int    `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	Sex  int    `form:"sex" json:"sex"`
}

type AddStudentReq struct {
	Name string `form:"name" json:"name"`
	Sex  int    `form:"sex" json:"sex"`
}
type AddStudentResp struct {
	Id int `form:"id" json:"id"`
}

type UpdateStudentReq struct {
	Id   int    `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	Sex  int    `form:"sex" json:"sex"`
}
type UpdateStudentResp struct {
}

type DelStudentReq struct {
	Id int `form:"id" json:"id"`
}
type DelStudentResp struct {
}

type AddStudentAndClassReq struct {
	Name string `form:"name" json:"name"`
	Sex  int    `form:"sex" json:"sex"`
}
type AddStudentAndClassResp struct {
}

type ChangeDatabaseReq struct {
}
type ChangeDatabaseResp struct {
}

type RedisUseReq struct {
}
type RedisUseResp struct {
}

type Config2 struct {
	Name string `json:"name"`
	Age  []int  `json:"age"`
}

type GetConfigResp struct {
	Demon1 string `json:"demon1"`
	Demon2 string `json:"demon2"`
}

type ReceiveMqMessageReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
