package httpclients

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/httpclient"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GetUser15suoReq struct {
	AccessToken string `form:"access_token" json:"access_token" binding:"required"`
}

type UserInfo15Suo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id         int    `json:"id"`
		UserName   string `json:"username"`
		RealName   string `json:"realName"`
		Phone      string `json:"phone"`
		Email      string `json:"email"`
		UserGroups []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"userGroups"`
		Roles []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
	} `json:"data"`
}

type GetAccessTokenReq struct {
	GrantType    string `form:"grant_type" json:"grant_type" binding:"required"`
	ClientId     string `form:"client_id" json:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" json:"client_secret" binding:"required"`
	Code         string `form:"code" json:"code" binding:"required"`
}

type GetAccessTokenResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		ClientId         string `json:"client_id"`
		Scope            string `json:"scope"`
	}
}

func GetAccessToken(req GetAccessTokenReq) (string, error) {
	// 第一步 根据授权码获取 access_token
	var clientConfig map[string]map[string]interface{}
	if err := config.Load("client", &clientConfig); err != nil {
		return "", err
	}
	if clientConfig["service_15suo"] == nil {
		return "", errors.New("免密登录配置缺失")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	DataUrlVal := url.Values{}
	DataUrlVal.Add("client_id", clientConfig["service_15suo"]["client_id"].(string))
	DataUrlVal.Add("client_secret", clientConfig["service_15suo"]["client_secret"].(string))
	DataUrlVal.Add("grant_type", "authorization_code")
	DataUrlVal.Add("code", req.Code)

	baseUri := clientConfig["service_15suo"]["base_uri"].(string)
	fmt.Println("11111111111", clientConfig["service_15suo"])
	fmt.Println("11111111111", baseUri)
	request, err := http.NewRequest("POST", baseUri+"/api/iam/oauth2/token", strings.NewReader(DataUrlVal.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	reqBody, err := io.ReadAll(request.Body)
	fmt.Println("2222222222222222", reqBody)
	if err != nil {
		fmt.Println("Request failed:", err)
		return "", err
	}
	defer resp.Body.Close()

	var res GetAccessTokenResp
	// 解析响应
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("333333333333333", string(body))
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	fmt.Println("44444444444", res.Data.AccessToken)
	if res.Data.AccessToken == "" {
		return "", errors.New("认证报错, code码: " + req.Code + ", 错误消息: " + res.Msg)
	}
	return res.Data.AccessToken, nil
}

type GetTokenValidLianTongReq struct {
	SsoToken string `form:"sso_token" json:"sso_token" binding:"required"`
}

type GetTokenValidLianTongResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct{}
}

// GetTokenValidLianTong
func GetTokenValidLianTong(req GetTokenValidLianTongReq) error {
	h, err := httpclient.NewHttpSend("service_liantong", "/yan/autologin/checkssocode")
	if err != nil {
		return err
	}
	param := map[string]interface{}{
		"sso_token": req.SsoToken,
	}
	h.SetHeader(map[string]interface{}{
		"Content-Type": "application/json",
		"accept":       "*/*",
	})
	h.SetBody(param)
	response, err := h.Post() //发送get请求
	if err != nil {
		return err
	}
	var res GetTokenValidLianTongResp
	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return err
	}
	if res.Code != 200 {
		return errors.New("认证失败")
	}
	return nil
}

func GetUserByAccessToken(req GetUser15suoReq) (res UserInfo15Suo, err error) {
	h, err := httpclient.NewHttpSend("service_15suo", "/api/iam/oauth2/userinfo")
	if err != nil {
		return res, err
	}
	h.SetHeader(map[string]interface{}{"Authorization": "Bearer " + req.AccessToken})
	response, err := h.Get() //发送get请求
	if err != nil {
		return res, err
	}
	fmt.Println(string(response))
	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
