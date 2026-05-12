package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// UserResponse 用户响应结构体
type UserResponse struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data UserData    `json:"data"`
}

// UserData 用户数据结构体
type UserData struct {
	User     User     `json:"user"`
	Employee Employee `json:"employee"`
}

// User 用户信息结构体
type User struct {
	ID           string      `json:"id"`
	CreateBy     string      `json:"createby"`
	CreateTime   string      `json:"createTime"`
	UpdateBy     string      `json:"updateby"`
	UpdateTime   string      `json:"updateTime"`
	Deleted      string      `json:"deleted"`
	Account      string      `json:"account"`
	Password     string      `json:"password"`
	IsSuper      string      `json:"isSuper"`
	FixedFlag    interface{} `json:"fixedFlag"`
	DataCheck    string      `json:"dataCheck"`
	OrgId        interface{} `json:"orgId"`
	OrgName      interface{} `json:"orgName"`
	Name         interface{} `json:"name"`
	Mobile       interface{} `json:"mobile"`
	FullName     interface{} `json:"fullName"`
	LoginNumbers int         `json:"loginNumbers"`
	LoginTime    string      `json:"loginTime"`
	Ids          interface{} `json:"ids"`
}

// Employee 员工信息结构体
type Employee struct {
	ID         string `json:"id"`
	CreateBy   string `json:"createby"`
	CreateTime string `json:"createTime"`
	UpdateBy   string `json:"updateby"`
	UpdateTime string `json:"updateTime"`
}

// getMockUserData 获取模拟用户数据
func getMockUserData() UserResponse {
	return UserResponse{
		Code: 200,
		Msg:  nil,
		Data: UserData{
			User: User{
				ID:           "135693120164174236",
				CreateBy:     "1",
				CreateTime:   "2025-04-02 10:00:09",
				UpdateBy:     "",
				UpdateTime:   "2025-08-27 14:28:13",
				Deleted:      "Y",
				Account:      "test",
				Password:     "7180d959d6383d590d92b69f8f032517b6151ae1073e36",
				IsSuper:      "",
				FixedFlag:    nil,
				DataCheck:    "75c289798413d012f89050ec9ad493d16868cb3b04b89",
				OrgId:        nil,
				OrgName:      nil,
				Name:         nil,
				Mobile:       nil,
				FullName:     nil,
				LoginNumbers: 0,
				LoginTime:    "2025-06-13 10:05:24",
				Ids:          nil,
			},
			Employee: Employee{
				ID:         "135693120164174236",
				CreateBy:   "1",
				CreateTime: "2025-04-02 10:00:09",
				UpdateBy:   "1",
				UpdateTime: "2025-08-13 10:05:32",
			},
		},
	}
}

// getUserHandler 处理获取用户信息的请求
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 处理OPTIONS预检请求
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 只处理GET请求
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取模拟数据
	userData := getMockUserData()

	// 将数据转换为JSON
	jsonData, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 返回JSON响应
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// main 主函数
func main() {
	// 注册路由
	http.HandleFunc("/power-api/base/user/context/getUser", getUserHandler)

	// 启动服务器
	port := ":8080"
	fmt.Printf("Mock服务器启动在端口%s\n", port)
	fmt.Println("访问地址: http://localhost:8080/power-api/base/user/context/getUser")
	log.Fatal(http.ListenAndServe(port, nil))
}
