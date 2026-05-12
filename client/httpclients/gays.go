package httpclients

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/httpclient"
)

type GetUserGaysReq struct {
	Authorization string `form:"authorization" json:"authorization" binding:"required"`
}

type GetUserGaysResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		User struct {
			Id      string `json:"id"`
			Account string `json:"account"`
		}
	} `json:"data"`
}

// GetUserByGays 获取公安一所用户信息
func GetUserByGays(req GetUserGaysReq) (string, error) {
	h, err := httpclient.NewHttpSend("service_gays", "/power-api/base/user/context/getUser")
	if err != nil {
		return "", err
	}
	//param := map[string]interface{}{
	//	"authorization": req.Authorization,
	//}
	//log.Println("param", param)
	//log.Println("token", req.Authorization)
	h.SetHeader(map[string]interface{}{
		"authorization": req.Authorization,
		"Content-Type":  "application/json",
		"accept":        "*/*",
	})
	log.Println("authorization", req.Authorization)
	log.Println("header", h.Header)
	//h.SetBody(param)
	response, err := h.Get() //发送get请求
	if err != nil {
		return "", err
	}
	var res GetUserGaysResp
	// 解析响应
	err = json.Unmarshal(response, &res)
	if err != nil {
		return "", err
	}
	log.Println("res", res.Data)
	if res.Code != 200 {
		log.Println("failed", res)
		return "", errors.New("认证失败")
	}
	return res.Data.User.Account, nil
}
