package httpclients

import (
	"gitlabee.4dogs.cn/common/httpclient"
	"io"
)

//VerMsg 请求报文
func VerMsg(method, url string, headerMap, body map[string]interface{}) (string, string, map[string]string, string, error) {
	h := httpclient.NewHttpSendByHost(url)
	if len(body) > 0 {
		if method == "GET" {
			h.GetUrlBuild(body)
		} else {
			h.SetBody(body)
		}
	}
	h.SetHeader(headerMap)
	resp, err := h.Send(method)
	if err != nil {
		return "", "", nil, "", err
	}
	defer resp.Body.Close()
	bodystring, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", nil, "", err
	}
	headermap := h.GetHeader(resp)
	return resp.Proto, resp.Status, headermap, string(bodystring), nil
}
