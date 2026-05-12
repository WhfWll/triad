package httpclients

import (
	"encoding/json"
	"errors"
	"gitlabee.4dogs.cn/common/httpclient"
)

type getUseScoreResp struct {
	Code int      `json:"code"`
	Data struct{} `json:"data"`
	Msg  string   `json:"msg"`
}

//请求计算目标可利用评分
func GetUseScore(param map[string]interface{}) error {
	h, err := httpclient.NewHttpSend("service_decision", "/decision/usescore/getscore")
	if err != nil {
		return err
	}
	h.SetBody(param)
	h.SetHeader(map[string]interface{}{"Authorization": "decision_open_api", "Connection": "close"})
	response, err := h.Post()
	if err != nil {
		return err
	}
	var tmpResp getUseScoreResp
	err = json.Unmarshal(response, &tmpResp)
	if err != nil {
		return err
	}
	if tmpResp.Code != 200 && len(tmpResp.Msg) > 0 {
		return errors.New(tmpResp.Msg)
	}
	return nil
}
