package httpclients

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/httpclient"
)

type GetTokenValidSiYuanReq struct {
	Token        string `form:"token" json:"token" binding:"required"`
	PlatformCode string `form:"platformCode" json:"platformCode" binding:"required"`
}

type GetTokenValidSiYuanResp struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	SubCode    string `json:"subCode"`
	SubMessage string `json:"subMessage"`
	ErrorLevel string `json:"errorLevel"`
	Data       struct {
		UserId string `json:"userId"`
	} `json:"data"`
}

// GetTokenValidSiYuan
func GetTokenValidSiYuan(req GetTokenValidSiYuanReq) (string, error) {
	h, err := httpclient.NewHttpSend("service_siyuan", "/api/user/validate")
	if err != nil {
		return "", err
	}
	platformCode := "LDWJJSMXT"
	if req.PlatformCode != "" {
		platformCode = req.PlatformCode
	}
	param := map[string]interface{}{
		"request":      "JumpValidateRequest",
		"platformCode": platformCode,
		"token":        req.Token,
	}
	log.Println("param", param)
	log.Println("token", req.Token)
	h.SetHeader(map[string]interface{}{
		"Content-Type": "application/json",
		"accept":       "*/*",
	})
	h.SetBody(param)
	response, err := h.Post() //发送get请求
	if err != nil {
		return "", err
	}
	var res GetTokenValidSiYuanResp
	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return "", err
	}
	log.Println("res", res.Data)
	if res.Code != "SUCCESS" {
		log.Println("failed", res)
		return "", errors.New("认证失败")
	}
	return res.Data.UserId, nil
}

// GetTokenByHTYZ 航天运载项目从云平台获取token
func GetTokenByHTYZ(req GetTokenValidSiYuanReq) (string, error) {
	h, err := httpclient.NewHttpSend("service_siyuan", "/api/user/validate")
	if err != nil {
		return "", err
	}
	platformCode := "LDWJJSMXT"
	if req.PlatformCode != "" {
		platformCode = req.PlatformCode
	}
	param := map[string]interface{}{
		"request":      "JumpValidateRequest",
		"platformCode": platformCode,
		"token":        req.Token,
	}
	log.Println("param", param)
	log.Println("token", req.Token)
	h.SetHeader(map[string]interface{}{
		"Content-Type": "application/json",
		"accept":       "*/*",
	})
	h.SetBody(param)
	response, err := h.Post() //发送get请求
	if err != nil {
		return "", err
	}
	var res GetTokenValidSiYuanResp
	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return "", err
	}
	log.Println("res", res.Data)
	if res.Code != "SUCCESS" {
		log.Println("failed", res)
		return "", errors.New("认证失败")
	}
	return res.Data.UserId, nil
}
