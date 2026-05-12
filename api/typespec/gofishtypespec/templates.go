package gofishtypespec

type TemplateIDReq struct {
	ID int64 `json:"id" form:"id" common:"模板id" v:"required|min:1#模板id最小值为1"`
}

type CreateTemplateReq struct {
	Name           string       `json:"name"`
	Subject        string       `json:"subject"`
	Text           string       `json:"text"`
	HTML           string       `json:"html"`
	EnvelopeSender string       `json:"envelope_sender"`
	Attachments    []Attachment `json:"attachments"`
}

// Attachment 表示邮件附件
type Attachment struct {
	Content string `json:"content"` // base64 编码的内容
	Type    string `json:"type"`
	Name    string `json:"name"`
}

type UpdateTemplateReq struct {
	ID             int64        `json:"id" form:"id" common:"模板id" v:"required|min:1#模板id最小值为1"`
	Name           string       `json:"name"`
	Subject        string       `json:"subject"`
	Text           string       `json:"text"`
	HTML           string       `json:"html"`
	EnvelopeSender string       `json:"envelope_sender"`
	Attachments    []Attachment `json:"attachments"`
}

type ImportEmailReq struct {
	Content      string `json:"content"`
	ConvertLinks bool   `json:"convert_links"`
}
