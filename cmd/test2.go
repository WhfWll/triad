package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	//dialer := &net.Dialer{
	//	Timeout: time.Second * 5, // 设置连接超时为 5 秒
	//}
	//conn, err := dialer.Dial("tcp", "192.168.0.68:88") // 连接到本地telnet服务器，默认端口23
	//if err != nil {
	//	fmt.Println("Failed to connect:", err)
	//	return
	//}
	//defer conn.Close()
	//fmt.Println(conn)
	//
	//urlRaw := "http://192.168.0.1:8888/abc.php?id=1"
	//fmt.Println(urlRaw)
	//url := network.ParseUrl2(urlRaw)
	//fmt.Println(url.Path)
	//fmt.Println(url.Query())

	url := "http://172.16.103.8:8084/yan/autologin/checkssocode"
	payload := strings.NewReader(`{"sso_token":"e57f896f-e644-4011-a316-1f2a0b62e078"}`)
	req, _ := http.NewRequest("POST", url, payload)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
