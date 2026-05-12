package invoke

import (
	"context"
	"os"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type CallInfo struct {
	CallId        string      `json:"callId"`        //调用id
	TaskId        string      `json:"taskId"`        //任务id
	ToolName      string      `json:"toolName"`      //工具名称
	ToolParamList []ToolParam `json:"toolParam"`     //工具参数
	Status        string      `json:"status"`        //调用状态
	Result        []string    `json:"result"`        //调用结果
	Previous      []string    `json:"previous"`      //前置脚本
	Following     []string    `json:"following"`     //后置关系
	FrontToolList []string    `json:"frontToolList"` //前置脚本名称列表
	Pid           int         `json:"pid"`           //进程运行id
	Source        string      `json:"source"`        //数据来源，vultest-待测漏洞手动触发
	SafeTestId    int         `json:"safeTestId"`    //执行安全测试的漏洞id
	Proxy         string      `json:"proxy"`         // 使用的代理，参数如：socks5://127.0.0.1:1080
}

type CallResult struct {
	CallID   string            `json:"callId"`
	ToolName string            `json:"toolName"`
	Status   string            `json:"status"`
	Details  map[string]string `json:"details"`
}

// 工具参数数据
type ToolParam struct {
	ParamName    string `json:"paramName"`    //参数名称
	ParamValue   string `json:"paramValue"`   //参数值
	Label        string `json:"label"`        //公开信息或者 私有信息
	RelationName string `json:"relationName"` //所需关系
	BelongTo     string `json:"belongTo"`     //所属脚本
}

type Risk struct {
	Hash string `json:"hash"`

	// essential
	IP        string `json:"ip"`
	IPInteger int64  `json:"ip_integer"`

	// extraTargets
	Url  string `json:"url"`
	Port int    `json:"port"`
	Host string `json:"host"`

	//
	Title           string `json:"title"`
	TitleVerbose    string `json:"title_verbose"`
	Description     string `json:"description"`
	Solution        string `json:"solution"`
	RiskType        string `json:"RiskType"`
	RiskTypeVerbose string `json:"RiskTypeVerbose"`
	Parameter       string `json:"parameter"`
	Payload         string `json:"payload"`
	Details         string `json:"details"`
	Severity        string `json:"severity"`
	Request         []byte `json:"request"`
	Response        []byte `json:"response"`

	// 新增字段
	ScriptName string `json:"ScriptName,omitempty"` // 脚本名称
	ScriptType string `json:"ScriptType,omitempty"` // 脚本类型
}

type ScannerLog struct {
	CreatedAt int64  `json:"CreatedAt"`
	Hash      string `json:"Hash"`
	Content   string `json:"content"`
	Pocname   string `json:"pocname"`

	// 攻击路径相关字段
	CallChainID       string                 `json:"call_chain_id,omitempty"`       // 调用链ID
	ParentScriptID    string                 `json:"parent_script_id,omitempty"`    // 父脚本ID
	ScriptExecutionID string                 `json:"script_execution_id,omitempty"` // 脚本执行ID
	TriggerReason     string                 `json:"trigger_reason,omitempty"`      // 触发原因
	AttackPath        []string               `json:"attack_path,omitempty"`         // 攻击路径
	RelatedFindings   []string               `json:"related_findings,omitempty"`    // 相关发现
	LogType           string                 `json:"log_type,omitempty"`            // 日志类型：scan_start, script_execution, vulnerability_found等
	Result            map[string]interface{} `json:"result,omitempty"`              // 扫描结果
}

// Stop 结束yak脚本运行
func (c *CallInfo) Stop(ctx context.Context) error {
	log.Println("stop:", c.ToolName, c.Pid)
	if c.Pid != 0 {
		cmd := exec.Command("kill", "-9", strconv.Itoa(c.Pid))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Error("kill process error: ", err)
			return err
		}
		log.Info("kill process success: ", c.ToolName, c.Pid)
	}
	return nil
}

// Stop 结束web_scanner脚本运行
func (c *CallInfo) StopSoft(ctx context.Context) error {
	log.Println("stop soft:", c.ToolName, c.Pid)
	if c.Pid != 0 {
		cmd := exec.Command("kill", strconv.Itoa(c.Pid))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Error("kill process error: ", err)
			return err
		}
		log.Info("kill process success: ", c.ToolName, c.Pid)
	}
	return nil
}
