package utils

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hyperjumptech/beda"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"smart/tools/enums"
	"sort"
	"strconv"
	"strings"
	"time"
)

var UrlHandleLogic urlHandleLogic

type urlHandleLogic struct {
}

// GetFullUrl 获取完整的链接
func (u urlHandleLogic) GetFullUrl(target string) string {
	if !strings.HasPrefix(target, "http") {
		target = "http" + "://" + target
	}
	return target
}

// ParseUrl 解析url
func (u urlHandleLogic) ParseUrl(targetUrl string) (string, string, string, string, string, error) {
	targetUrl = u.GetFullUrl(targetUrl)
	urlPtr, err := url.Parse(targetUrl)
	if err != nil {
		return "", "", "", "", "", err
	}
	scheme := urlPtr.Scheme
	hostname := urlPtr.Hostname()
	port := urlPtr.Port()
	path := urlPtr.Path
	query := urlPtr.RawQuery
	return scheme, hostname, port, path, query, nil
}

/*
FormatHeadersNoHttpHeader 把原始的请求头格式化为字典（无HTTP头的RequestHeader用此方法）
HTTP头 如："GET / HTTP/1.1\r\n"
*/
func (u urlHandleLogic) FormatHeadersNoHttpHeader(rawHeaders string) (map[string]string, error) {
	header := make(map[string]string, 0)
	rawHeaders = strings.ReplaceAll(rawHeaders, "\r\n", "\n")
	rawHeaders = strings.ReplaceAll(rawHeaders, "\n", "\r\n")
	if !strings.HasSuffix(rawHeaders, "\r\n") {
		rawHeaders += "\r\n"
	}

	// don't forget to make certain the headers end with a second "\r\n"
	reader := bufio.NewReader(strings.NewReader(rawHeaders + "\r\n"))
	tp := textproto.NewReader(reader)

	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	httpHeader := http.Header(mimeHeader)

	//把httpHeader中的value由切片转为字符串
	for key, value := range httpHeader {
		header[key] = strings.Join(value, ",")
	}

	return header, nil
}

/*
FormatHeaders 把原始的请求头格式化为字典(有HTTP头的RequestHeader用此方法，默认用此方法)
HTTP头 如："GET / HTTP/1.1\r\n"
*/
func (u urlHandleLogic) FormatHeaders(rawHeaders string) (map[string]string, error) {
	header := make(map[string]string, 0)
	rawHeaders = strings.ReplaceAll(rawHeaders, "\r\n", "\n")
	rawHeaders = strings.ReplaceAll(rawHeaders, "\n", "\r\n")
	if !strings.HasSuffix(rawHeaders, "\r\n") {
		rawHeaders += "\r\n"
	}
	tempValue := strings.Split(rawHeaders, ":")
	if len(tempValue) != 2 {
		return header, errors.New("header格式错误")
	}

	// we need to make sure to add a fake HTTP header here to make a valid request.
	//reader := bufio.NewReader(strings.NewReader(rawHeaders + "\r\n"))
	//req, err := http.ReadRequest(reader)

	//if err != nil {
	//	return nil, err
	//}
	//httpHeader := req.Header

	//把httpHeader中的value由切片转为字符串
	//for key, value := range httpHeader {
	//	header[key] = strings.Join(value, ",")
	//}
	return header, nil
}

// FormatCookie 把cookie的值由字符串转为字典
func (u urlHandleLogic) FormatCookie(rawCookie string) map[string]string {
	cookie := make(map[string]string)
	cookieList := strings.Split(strings.ReplaceAll(strings.TrimSpace(rawCookie), "; ", ";"), ";")
	for _, value := range cookieList {
		cookieSubStrList := strings.Split(value, "=")
		if len(cookieSubStrList) == 2 {
			cookie[cookieSubStrList[0]] = cookieSubStrList[1]
		}
	}
	return cookie
}

// RedirectFunc 重定向禁止
func (u urlHandleLogic) RedirectFunc(req *http.Request, via []*http.Request) error {
	fmt.Println(req.RequestURI)
	// 如果返回 非nil 则禁止向下重定向 返回nil 则 一直向下请求 10 次 重定向
	return http.ErrUseLastResponse
}

func (u urlHandleLogic) HttpGet(url string, header map[string]string, cookie map[string]string) (int, []byte, error) {
	//禁止跳转，超时时间为10s
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
	}
	client := http.Client{CheckRedirect: u.RedirectFunc, Timeout: 10 * time.Second, Transport: tr}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}
	// 添加请求头
	for key, value := range header {
		req.Header.Add(key, value)
	}

	// 添加cookie
	for key, value := range cookie {
		cookie := &http.Cookie{
			Name:  key,
			Value: value,
		}
		req.AddCookie(cookie)
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	respByteSlice, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return statusCode, respByteSlice, nil
}

func (u urlHandleLogic) VerifyStatus(url string, verifyType int, verifyValue string) (int, error) {
	defaultUserHeaders := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:68.0) Gecko/20100101 Firefox/68.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,application/json,text/plain,*/*;q=0.8",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Connection":      "close",
	}
	headersRight := make(map[string]string, 0)
	headersError := defaultUserHeaders
	cookieRight := make(map[string]string, 0)
	cookieError := make(map[string]string, 0)
	if verifyType == enums.TaskConfigurationWebsiteLoginCookie {
		cookieRight = u.FormatCookie(verifyValue)
		headersRight = defaultUserHeaders
	} else if verifyType == enums.TaskConfigurationWebsiteLoginHeader {
		headers, err := u.FormatHeaders(verifyValue)
		if err != nil {
			return 0, err
		}
		headersRight = headers
	}
	statusCode, respByteSlice, err := u.HttpGet(url, headersRight, cookieRight)
	if err != nil {
		return 0, err
	}
	if statusCode == 404 {
		return enums.TaskConfigurationWebsiteLoginFail, nil
	}
	if u.In(strconv.Itoa(statusCode), []string{"301", "302", "303"}) {
		return enums.TaskConfigurationWebsiteLoginSuccess, nil
	}
	statusCode2, respByteSlice2, err := u.HttpGet(url, headersError, cookieError)
	if statusCode == 200 && u.In(strconv.Itoa(statusCode2), []string{"301", "302", "303"}) {
		return enums.TaskConfigurationWebsiteLoginVerify, nil
	} else if statusCode == 200 && statusCode2 == 200 {
		if beda.JaroDistance(string(respByteSlice), string(respByteSlice2)) >= 0.9 {
			return enums.TaskConfigurationWebsiteLoginSuccess, nil
		} else {
			return enums.TaskConfigurationWebsiteLoginVerify, nil
		}
	} else {
		return enums.TaskConfigurationWebsiteLoginFail, nil
	}
}

func (u urlHandleLogic) In(target string, targetList []string) bool {
	sort.Strings(targetList)
	index := sort.SearchStrings(targetList, target)
	if index < len(targetList) && targetList[index] == target {
		return true
	}
	return false
}
