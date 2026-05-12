package gofishtypespec

type ProfileIDReq struct {
	ID int64 `json:"id" form:"id" common:"发送配置id" v:"required|min:1#发送配置id最小值为1"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateProfileReq struct {
	Name             string   `json:"name"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	Host             string   `json:"host"`
	InterfaceType    string   `json:"interface_type"`
	FromAddress      string   `json:"from_address"`
	IgnoreCertErrors bool     `json:"ignore_cert_errors"`
	Headers          []Header `json:"headers,omitempty"`
}

type UpdateProfileReq struct {
	ID               int64    `json:"id" form:"id" common:"发送配置id" v:"required|min:1#发送配置id最小值为1"`
	Name             string   `json:"name"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	Host             string   `json:"host"`
	InterfaceType    string   `json:"interface_type"`
	FromAddress      string   `json:"from_address"`
	IgnoreCertErrors bool     `json:"ignore_cert_errors"`
	Headers          []Header `json:"headers,omitempty"`
}

type SendTestEmailReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Position  string `json:"position"`
	Url       string `json:"url"`

	Template struct {
		Name string `json:"name"`
	} `json:"template"`

	Page struct {
		Name string `json:"name"`
	} `json:"page"`

	Smtp struct {
		Name             string   `json:"name,omitempty"`
		FromAddress      string   `json:"from_address"`
		Host             string   `json:"host"`
		Username         string   `json:"username"`
		Password         string   `json:"password"`
		IgnoreCertErrors bool     `json:"ignore_cert_errors"`
		Headers          []Header `json:"headers"`
	} `json:"smtp"`
}
