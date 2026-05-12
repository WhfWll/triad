package network

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
)

func UseProxyURL(u, proxy string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	}}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil { // 使用传入代理
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	r, err := (&http.Client{Transport: tr}).Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
