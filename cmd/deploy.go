//go:build ignore

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"smart/client/httpclients"
	"smart/models/mysqls"
	"smart/services"
	dbutils "smart/tools/data"
	"smart/tools/enums"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/urfave/cli"
	mysqlsEnv "gitlabee.4dogs.cn/common/mysql"
	"gitlabee.4dogs.cn/common/redis"
)

//一键部署脚本
// 添加 rabbitmq 配置
// 部署 mysql 初始化 数据
// 部署 neo4j 初始化 数据
// 部署 supervisor 配置文件
// 进行 授权 设置
// 一键 启动 docker 容器
// 数据结构升级比较
// 二进制文件升级比较
// 生成 产品序列号

var InitSubCommand = cli.Command{
	Name:  "init",
	Usage: "进行系统初始化管理",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "productId"},
		&cli.StringFlag{Name: "template"},
		&cli.BoolFlag{Name: "node"},
		&cli.BoolFlag{Name: "dbStruct"},
	},
	Action: func(c *cli.Context) error {
		mysqlsEnv.Setup()
		redis.Setup()

		ctx := context.Background()
		productId := c.Bool("productId")
		templateName := c.String("template")
		node := c.Bool("node")
		dbStruct := c.Bool("dbStruct")
		if productId {
			initProductId(ctx)
		}
		if templateName != "" {
			initTemplate(ctx, templateName)
		}
		if node {
			initNodeData(ctx)
		}
		if dbStruct {
			initDbStruct(ctx)
		}
		return nil
	},
}

var ChangeSubCommand = cli.Command{
	Name:  "change",
	Usage: "进行系统更改管理",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "sysName"},
		&cli.StringFlag{Name: "userPass"},
		&cli.StringFlag{Name: "sysVersion"},
		&cli.StringFlag{Name: "vulVersion"},
		&cli.BoolFlag{Name: "closeWhite"},
	},
	Action: func(c *cli.Context) error {
		ctx := context.Background()
		sysName := c.String("sysName")
		userPass := c.String("userPass")
		sysVersion := c.String("sysVersion")
		vulVersion := c.String("vulVersion")
		closeWhite := c.Bool("closeWhite")
		if sysName != "" {
			changeSysName(ctx, sysName)
		}
		if userPass != "" {
			changeUserPass(ctx, sysName)
		}
		if sysVersion != "" {
			changeSysVersion(ctx, sysVersion)
		}
		if vulVersion != "" {
			changeVulVersion(ctx, vulVersion)
		}
		if closeWhite {
			closeWhiteFunc(ctx)
		}
		return nil
	},
}

var CleanSubCommand = cli.Command{
	Name:  "clean",
	Usage: "进行系统数据清理",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "decision"},
		&cli.BoolFlag{Name: "smart"},
		&cli.BoolFlag{Name: "yakitDb"},
	},
	Action: func(c *cli.Context) error {
		ctx := context.Background()
		decision := c.Bool("decision")
		smart := c.Bool("smart")
		yakitDb := c.Bool("yakitDb")
		if decision {
			cleanDecisionData(ctx)
		}
		if smart {
			cleanSmartData(ctx)
		}
		if yakitDb {
			cleanYakitDb(ctx)
		}
		return nil
	},
}

var CheckSubCommand = cli.Command{
	Name:  "check",
	Usage: "进行系统检查",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "auth"},
		&cli.BoolFlag{Name: "env"},
		&cli.BoolFlag{Name: "blockEnv"},
		&cli.BoolFlag{Name: "sysVersion"},
		&cli.BoolFlag{Name: "vulVersion"},
		&cli.StringFlag{Name: "dbStruct"},
	},
	Action: func(c *cli.Context) error {
		ctx := context.Background()
		isAuth := c.Bool("auth")
		env := c.Bool("env")
		blockEnv := c.Bool("blockEnv")
		sysVersion := c.Bool("sysVersion")
		vulVersion := c.Bool("vulVersion")
		dbStructFile := c.String("dbStruct")
		if isAuth {
			checkSystemAuth(ctx)
		}
		if env {
			checkEnv(ctx)
		}
		if blockEnv {
			blockCheckEnv(ctx)
		}
		if sysVersion {
			getSysVersion(ctx)
		}
		if vulVersion {
			getVulVersion(ctx)
		}
		if dbStructFile != "" {
			dbutils.CheckDbStruct(ctx, dbStructFile)
		}
		return nil
	},
}

var FixSubCommand = cli.Command{
	Name:  "fix",
	Usage: "进行系统环境修复",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "docker"},
	},
	Action: func(c *cli.Context) error {
		ctx := context.Background()
		docker := c.Bool("docker")
		if docker {
			err := fixDocker(ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil
	},
}

var ExportSubCommand = cli.Command{
	Name:  "export",
	Usage: "进行数据导出",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "task"},
		&cli.BoolFlag{Name: "target"},
		&cli.BoolFlag{Name: "taskVul"},
		&cli.BoolFlag{Name: "taskInfo"},
		&cli.BoolFlag{Name: "taskResult"},
		&cli.BoolFlag{Name: "taskLog"},
		&cli.BoolFlag{Name: "taskEvidence"},
		&cli.BoolFlag{Name: "remoteSession"},
		&cli.BoolFlag{Name: "taskData"},
	},
	Action: func(c *cli.Context) error {
		ctx := context.Background()
		task := c.Bool("task")
		target := c.Bool("target")
		taskVul := c.Bool("taskVul")
		taskInfo := c.Bool("taskInfo")
		taskResult := c.Bool("taskResult")
		taskLog := c.Bool("taskLog")
		taskEvidence := c.Bool("taskEvidence")
		remoteSession := c.Bool("remoteSession")
		taskData := c.Bool("taskData")

		mysqlsEnv.Setup()

		if task {
			_, err := exportTask(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if target {
			_, err := exportTarget(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if taskVul {
			_, err := exportTaskVul(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if taskInfo {
			_, err := exportTaskInfo(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if taskResult {
			_, err := exportTaskResult(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if taskLog {
			_, err := exportTaskLog(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if taskEvidence {
			_, err := exportTaskEvidence(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if remoteSession {
			_, err := exportRemoteSession(ctx, "")
			if err != nil {
				fmt.Println(err)
			}
		}
		if taskData {
			err := exportTaskData(ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil
	},
}

var ImportSubCommand = cli.Command{
	Name:  "import",
	Usage: "进行数据导入",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "task"},
		&cli.StringFlag{Name: "target"},
		&cli.StringFlag{Name: "taskVul"},
		&cli.StringFlag{Name: "taskInfo"},
		&cli.StringFlag{Name: "taskResult"},
	},
	Action: func(c *cli.Context) error {
		ctx := context.Background()
		taskPath := c.String("task")
		targetPath := c.String("target")
		taskVulPath := c.String("taskVul")
		taskInfoPath := c.String("taskInfo")
		taskResultPath := c.String("taskResult")

		mysqlsEnv.Setup()
		taskId, targetId := getLatestTaskId(ctx)
		if taskId <= 0 || targetId <= 0 {
			return errors.New("请输入基准任务id和基准目标id")
		}
		fmt.Println("任务基准id: ", taskId, "目标基准id: ", targetId)

		if taskPath != "" {
			importTask(ctx, taskPath, taskId)
		}
		if targetPath != "" {
			importTarget(ctx, targetPath, taskId, targetId)
		}
		if taskVulPath != "" {
			importTaskVul(ctx, taskVulPath, taskId, targetId)
		}
		if taskInfoPath != "" {
			importTaskInfo(ctx, taskInfoPath, taskId, targetId)
		}
		if taskResultPath != "" {
			importTaskResult(ctx, taskResultPath, taskId, targetId)
		}
		return nil
	},
}

func main() {
	dbutils.InitConfigInfo()
	app := cli.NewApp()
	app.Usage = "smart manage tools"
	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, InitSubCommand)
	app.Commands = append(app.Commands, CleanSubCommand)
	app.Commands = append(app.Commands, ChangeSubCommand)
	app.Commands = append(app.Commands, CheckSubCommand)
	app.Commands = append(app.Commands, FixSubCommand)
	app.Commands = append(app.Commands, ExportSubCommand)
	app.Commands = append(app.Commands, ImportSubCommand)
	err := app.Run(os.Args)
	if err != nil {
		println(err.Error())
		return
	}
}

func initProductId(ctx context.Context) {
	var srv services.Auth
	serialNumber, err := srv.GenerateSystemSerialNumber(ctx)
	if err != nil {
		panic(err)
	}
	tempData, err := buildAuthData2(serialNumber, "0")
	if err != nil {
		panic(err)
	}
	var auth services.Auth
	ciphertext := auth.RsaEncrypt(tempData, []byte(enums.SWPubKey))
	err = srv.UpdateAuthInfo(ctx, map[string]string{
		"productID":   serialNumber,
		"productName": "自动化渗透测试系统",
		"authCode":    hex.EncodeToString(ciphertext),
		"authTime":    "未授权",
	})
	err = srv.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
	err = srv.CleanAuthRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("初始化 序列号 成功")
}

func buildAuthData2(serialNumber string, days string) ([]byte, error) {
	tempMap := make(map[string]string, 0)
	tempMap["productID"] = serialNumber
	tempMap["generateTime"] = time.Now().Format("2006-01-02")
	tempMap["authDays"] = days
	tempData, err := json.Marshal(tempMap)
	if err != nil {
		return nil, nil
	}
	return tempData, nil
}

// 清理 决策引擎 临时数据
func cleanDecisionData(ctx context.Context) {
	scriptParamList := []string{"exec", "laozhi_db", "mysql", "-uroot", "-proot.4dogs.cn", "decision", "--default-character-set=utf8", "-e", "truncate", "history_tasks"}
	fmt.Println(scriptParamList)
	cmd := exec.Command("docker", scriptParamList...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("exec error: ", err)
	}
}

// 清理 smart 临时数据
func cleanSmartData(ctx context.Context) {
	command := `docker exec laozhi_db mysql -uroot -proot.4dogs.cn smart --default-character-set=utf8 -e 'truncate bas_log;truncate bas_node;truncate bas_target;truncate bas_task;truncate bas_vul;truncate burpsuite_task;truncate burpsuite_task_result;truncate flow_base;truncate flow_log;truncate flow_risk;truncate flow_target;truncate flow_task;truncate log_audit;truncate log_backup;truncate remote_session;truncate report_record;truncate report_verify_port;truncate report_verify_target;truncate report_verify_task;truncate report_verify_vul;truncate system_config_backup;truncate system_message;truncate task_evidence;truncate task_log;truncate task_log_info;truncate task_result;truncate task_target;truncate task_target_result;truncate task_task;truncate task_task_info;truncate task_task_result;truncate task_vul;truncate wifi_ap_info;truncate wifi_cli_info;truncate wifi_task;truncate wifi_task_log;truncate xray_task;truncate xray_task_result;truncate logic_task;truncate logic_target;truncate logic_vul;truncate logic_log;truncate logic_log_info;'`
	cmd := exec.Command("bash", "-c", command)
	fmt.Println("[command]", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("exec error: ", err)
	}
}

// 清理 default-yakit.db
func cleanYakitDb(ctx context.Context) {
	rootDisk, err := disk.Usage("/")
	if err != nil {
		fmt.Println(err)
	}
	freeSize := int(math.Round(float64(rootDisk.Free) / 1024 / 1024 / 1024))
	if freeSize > 2 {
		return
	}
	// 清理磁盘
	u, err := user.Current()
	homeList := []string{u.HomeDir, "/root", "/home/dogs", "/home/admin1/"}
	for _, path := range homeList {
		yakitDbPath := filepath.Join(path, "yakit-projects", "default-yakit.db")
		fmt.Println(yakitDbPath)
		err = os.Remove(yakitDbPath)
		if err != nil {
			continue
		} else {
			fmt.Println("File removed successfully: ", yakitDbPath)
			break
		}
	}
}

// 更新系统名称
func changeSysName(ctx context.Context, sysName string) {
	mysqlsEnv.Setup()
	redis.Setup()
	var srv services.Auth
	authInfo, err := srv.GetAuthInfo(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	authInfo["productName"] = sysName
	err = srv.UpdateAuthInfo(ctx, authInfo)
	if err != nil {
		fmt.Println("系统名称更新报错: ", err)
	}
	fmt.Println("系统名称更新成功")
}

// 更新用户名称
func changeUserPass(ctx context.Context, sysName string) {
	fmt.Println("系统名称更新成功")
}

// 更新系统版本
func changeSysVersion(ctx context.Context, sysVersion string) {
	mysqlsEnv.Setup()
	redis.Setup()
	var (
		srv                 services.MapSet
		systemVersionMapSet services.SystemVersionMapSet
	)

	objValueStr, err := srv.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if objValueStr == "" {
		systemVersionMapSet.CurrentVersion = enums.SystemBaseVersion
		systemVersionMapSet.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
		objValueByte, err := json.Marshal(systemVersionMapSet)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = srv.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent); err != nil {
			fmt.Println(err)
			return
		}
	}
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &systemVersionMapSet); err != nil {
			fmt.Println(err)
			return
		}
	}
	systemVersionMapSet.CurrentVersion = sysVersion
	systemVersionMapSet.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	systemVersionData, err := json.Marshal(systemVersionMapSet)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = srv.Create(ctx, enums.SystemVersionMapSetObjKey, string(systemVersionData), "系统版本信息")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("系统版本号更新成功")
}

// 更新漏洞版本
func changeVulVersion(ctx context.Context, vulVersion string) {
	mysqlsEnv.Setup()
	redis.Setup()
	var (
		srv                 services.MapSet
		systemVersionMapSet services.SystemVersionMapSet
	)

	objValueStr, err := srv.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if objValueStr == "" {
		systemVersionMapSet.CurrentVersion = enums.SystemBaseVersion
		systemVersionMapSet.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
		objValueByte, err := json.Marshal(systemVersionMapSet)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = srv.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent); err != nil {
			fmt.Println(err)
			return
		}
	}
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &systemVersionMapSet); err != nil {
			fmt.Println(err)
			return
		}
	}
	systemVersionMapSet.VulVersion = vulVersion
	systemVersionMapSet.VulUpdateTime = time.Now().Format("2006-01-02 15:04:05")
	systemVersionData, err := json.Marshal(systemVersionMapSet)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = srv.Create(ctx, enums.SystemVersionMapSetObjKey, string(systemVersionData), "系统版本信息")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("漏洞版本号更新成功")
}

// 关闭系统白名单
func closeWhiteFunc(ctx context.Context) {
	mysqlsEnv.Setup()
	redis.Setup()
	var (
		srv           services.MapSet
		ipWhiteMapSet services.SystemSettingIpWhiteMapSet
	)
	//如果不存在获取系统白名单key,就新建一个map key
	objValueStr, err := srv.GetMapValue(ctx, enums.SystemAccessIpWhiteMapSetObjKey)
	if objValueStr == "" {
		ipWhiteMapSet.IsOpen = enums.ConfigClose
		objValueByte, err := json.Marshal(ipWhiteMapSet)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = srv.Create(ctx, enums.SystemAccessIpWhiteMapSetObjKey, string(objValueByte), enums.SystemAccessIpWhiteMapSetContent); err != nil {
			fmt.Println(err)
			return
		}
	}
	//如果存在系统白名单key，就更新进去
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &ipWhiteMapSet); err != nil {
			fmt.Println(err)
			return
		}
	}
	ipWhiteMapSet.IsOpen = enums.ConfigClose
	ipWhiteMapSetByte, err := json.Marshal(ipWhiteMapSet)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = srv.Create(ctx, enums.SystemAccessIpWhiteMapSetObjKey, string(ipWhiteMapSetByte), enums.SystemAccessIpWhiteMapSetContent)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("系统白名单恢复成功")
}

// 检测系统授权
func checkSystemAuth(ctx context.Context) {
	mysqlsEnv.Setup()
	redis.Setup()
	var srv services.Auth
	//第一步  获取授权信息
	authInfoMap, err := srv.GetAuthInfo(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	//第二步 检查产品id是否一致
	serialNumber, err := srv.GenerateSystemSerialNumber(ctx)
	if err != nil {
		fmt.Println("授权状态异常： ", err)
		return
	}
	if serialNumber != authInfoMap["productID"] {
		fmt.Println("授权状态异常： ", "系统序列号与本机特征码不一致")
		return
	}
	//第三步 解码授权码，并比较 产品id和授权码中的产品id
	authCode := authInfoMap["authCode"]
	decryptMap, err := srv.RsaDecrypt(ctx, authCode)
	if err != nil {
		fmt.Println("授权状态异常: ", "授权码解码出错")
	}
	if decryptMap["productID"] != serialNumber {
		fmt.Println("授权状态异常: ", "授权码不匹配")
		return
	}

	//第四步 计算授权时间
	authTimeString := authInfoMap["authTime"]
	authDays, err := strconv.Atoi(decryptMap["authDays"])
	if err != nil {
		fmt.Println("授权状态异常: ", "时间解码出错", err.Error())
		return
	}
	authTime, err := time.Parse(enums.ResTimeDayLayout, authTimeString)
	authExpireTime := authTime.AddDate(0, 0, authDays)
	leftDays := int(math.Ceil(authExpireTime.Sub(time.Now()).Hours() / 24))
	if leftDays <= 0 {
		fmt.Println("授权状态异常: ", "授权已过期")
		return
	}
	fmt.Println("当前授权状态正常")
}

// 检测系统环境
func checkEnv(ctx context.Context) {
	isNormal := true
	err := dbutils.CheckMysql(ctx)
	if err != nil {
		isNormal = false
		fmt.Println("运行环境异常: ", err.Error())
	}
	err = dbutils.CheckRedis(ctx)
	if err != nil {
		isNormal = false
		fmt.Println("运行环境异常: ", err.Error())
	}
	if isNormal {
		fmt.Println("运行环境正常")
	}
}

// 检测系统环境
func blockCheckEnv(ctx context.Context) {
	for {
		err := dbutils.CheckMysql(ctx)
		if err != nil {
			fmt.Println("运行环境异常: ", err.Error())
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	for {
		err := dbutils.CheckRedis(ctx)
		if err != nil {
			fmt.Println("运行环境异常: ", err.Error())
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	fmt.Println("系统环境正常")
}

// 更新系统版本
func getSysVersion(ctx context.Context) {
	mysqlsEnv.Setup()
	redis.Setup()
	var (
		srv                 services.MapSet
		systemVersionMapSet services.SystemVersionMapSet
	)

	objValueStr, err := srv.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if objValueStr == "" {
		return
	}
	if err = json.Unmarshal([]byte(objValueStr), &systemVersionMapSet); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(systemVersionMapSet.CurrentVersion)
}

// 更新系统版本
func getVulVersion(ctx context.Context) {
	mysqlsEnv.Setup()
	redis.Setup()
	var (
		srv                 services.MapSet
		systemVersionMapSet services.SystemVersionMapSet
	)

	objValueStr, err := srv.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if objValueStr == "" {
		return
	}
	if err = json.Unmarshal([]byte(objValueStr), &systemVersionMapSet); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(systemVersionMapSet.VulVersion)
}

// 修复neo4j运行环境
func fixDocker(ctx context.Context) error {
	cmd1 := exec.Command("bash", "-c", "docker-compose down")
	cmd1.Dir = "/opt/laozhi"
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stdout
	err := cmd1.Run()

	cmd2 := exec.Command("bash", "-c", "docker-compose up -d")
	cmd2.Dir = "/opt/laozhi"
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stdout
	err = cmd2.Run()
	if err != nil {
		return errors.New("docker 环境恢复异常: " + err.Error())
	}
	err = restartSupervisord(ctx)
	if err != nil {
		return errors.New("supervisord 环境恢复异常: " + err.Error())
	}
	fmt.Println("supervisord 环境恢复成功")
	return nil
}

// 重启进程
func restartSupervisord(ctx context.Context) error {
	cmd1 := exec.Command("bash", "-c", "supervisord ctl restart all ")
	cmd1.Dir = "/opt/laozhi"
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stdout
	err := cmd1.Run()

	if err != nil {
		return errors.New("docker 命令执行异常: " + err.Error())
	}
	return nil
}

func initTemplate(ctx context.Context, templateName string) {
	scriptList := []string{
		"mitm_xss_check",
		"mitm_sql_check",
		"mitm_ssrf_check",
		"mitm_xxe_check",
		"mitm_file_download_check",
		"mitm_fileupload_check",
		"mitm_command_injection_check",
		"mitm_sensinfo_check",
		"mitm_urlredirect_check",
		"mitm_sqlstatement_heck",
		"mitm_contenttype_check",
		"mitm_csp_check",
		"mitm_skipVerificationStep_check",
		"mitm_cleartext_check",
		"mitm_csrf_check",
		"mitm_http_check",
		"mitm_code_leak_check",
		"mitm_fastjson_rce_check",
		"mitm_cors_check",
		"Xunyi_Technology_74CMS_security_vulnerability_poc",
		"Activemq_weak_password_poc",
		"Activemq_CVE_2016_3088_poc",
		"ActiveMQ_information_leakage_vulnerability_poc",
		"airflow_unauth_vul",
		"bash_cve_2014_6271_rce",
		"coldfusion_cve_2010_2861_lfi",
		"Adobe ColdFusion 11 LDAP 反序列化",
		"Confluence_CVE_2019_3396",
		"Confluence_CVE_2021_26084",
		"Confluence_CVE_2022_26134",
		"CouchDB_CVE_2017_12635_poc",
		"CouchDB_CVE_2017_12636_exp",
		"discuz_v7.2_sqli",
		"discuz_wooyun_2010_080723",
		"default-django-page",
		"Docker_unauthorized_access_vulnerability_poc",
		"Drupal_SQL_injection_vulnerability_(CVE-2014-3704)._poc",
		"CVE-2018-7600",
		"Dubbo_default_password",
		"ecshop_collection_list_sqli_poc",
		"ecshop_xianzhi_2017_02_82239600_poc",
		"elasticsearch_CVE_2014_3120_poc",
		"elasticsearch_CVE_2015_1427_poc",
		"elasticsearch_CVE_2015_3337_poc",
		"elasticsearch_CVE_2015_5531_poc",
		"elasticsearch_wooyun_2015_110216_poc",
		"mitm_fastjson_rce_check",
		"Apache_Flink_CVE_2020_17519_fileread",
		"CVE-2021-22205",
		"goahead_CVE_2017_17562",
		"Grafana welcome 任意文件读取漏洞",
		"hadoop-unauth-rce",
		"apache_http_parsing_CVE_2017_15715_poc",
		"apache_ssi_rce",
		"influxdb_Unauthorized_access",
		"Jboss_2017_12149_RCE",
		"Jboss_2017_7504_RCE",
		"Jboss_jmxinvokerServlet_rce_exp_poc",
		"jenkins_version_detect",
		"Jenkins_brute",
		"jenkins_CVE_2018_1000861_poc",
		"CVE_2021_28164",
		"CVE_2021_28169",
		"CVE-2021-34429",
		"Atlassian_Jira_cfx_CVE_2021_26086_fileread",
		"Joomla_3.7.0_(CVE-2017-8917)_SQL_injection_vulnerability_poc",
		"Jupyter_Notebook_Unauthorized_access",
		"Kibana_CVE_2018_17246_file_inclusion",
		"Kibana_Unauthorized_access",
		"laravel_CVE_2021_3129",
		"Metabase_geojson_file_read",
		"mongo_express_CVE_2019_10758_rce",
		"MySQL_CVE_2012_2122_sql_inject",
		"nacos_default_password",
		"nacos_QVD_2023_6271_console_bypass",
		"nexus_cve_2019_7238_rce",
		"nginx_cve_2017_7529",
		//"phpstudy_nginx_wrong_resolve",
		"CVE-2017-14849",
		"CVE-2021-28073",
		"CVE-2020-9496",
		"CVE-2020-35476",
		"PHP_CGI_cve_2012_1823_rce",
		"phpmyadmin_cve_2018_12613_file_read",
		"PHPUnit_eval_stdin_CVE_2017_9841_rce",
		"rabbitmq_default_password",
		"Rabbitmq_unauthorized_access_vulnerability_poc",
		"Ruby_on_Rails_Path_Crossing_Vulnerability_(CVE-2018-3760)__poc",
		"CVE-2019-5418",
		"redis_unauth",
		"Rocket.Chat_sql_CVE_2021_22911",
		"rsync_Unauthorized_access_poc",
		"CVE-2020-16846",
		"CVE-2016-4437",
		"yunshikong_erp_shiro_RCE",
		"Apache_Skywalking_8.3.0_SQL",
		"CVE-2017-12629",
		"apache solr_Unauthorized_access",
		"spark_api_unauth",
		"CVE-2016-4977",
		"Spring_Data_Rest_CVE_2017_8046_poc",
		"Spring_Data_Commons_CVE_2018_1273_rce_poc",
		"Spring_Cloud_Gateway_cve_2022_22947_poc",
		"CVE-2022-22963",
		"SpringFramework_CVE_2022_22965_poc",
		"Supervisord_CVE_2017_11610",
		"Thinkphp_5.0.22_5.1.29_RCE",
		"ThinkPHP_5.0.23_2_RCE",
		"ThinkPHP5_SQL_Injection_exp",
		"ThinkPHP_2.x~3.0Beta_RCE",
		"Thinkphp_lang_RCE",
		"tomcat_CVE_2017_12615_fileupload",
		"Apache_Tomcat_CVE_2020_1938",
		"tomcat_weak_password_vulnerabilites",
		"CVE-2020-13942",
		"CVE-2018-7490",
		"weblogic_CVE_2017_10271",
		"weblogic_CVE_2018_2894_poc",
		"weblogic_CVE_2018_14882_poc",
		"weblogic_CVE_2014_4210_poc",
		"weblogic_brute",
		"Webmin_rce",
		"CVE-2021-21351",
		"CVE-2021-29505",
		"xxljob-default-login",
		"pi-holidayapi",
		"zabbix_CVE_2016_10134_sql",
		"zabbix_console_default_passwd",
		"zabbix_authentication_bypass",
		"mitm_htaccess_File_Readable",
		"mitm_js_library_check",
		"mitm_open_redirect_check",
		"arbitrary_file_download",
		"mitm_cookie_httponly_check",
		"cleartext_password_check",
		"mitm_accessdb_check",
		"mitm_dir_index_check",
		"http_to_https_redirect_check",
		"mitm_file_inclusion_check",
		"mitm_phpinfo_check",
		"mitm_devfile_check",
		"mitm_readmeinfo_check",
		"mitm_gitinfo_check",
	}
	libIds := make([]int, 0)
	for _, script := range scriptList {
		var req httpclients.OpenVulLibIdByPocnamesReq
		req.MatchType = "eq_pocnames"
		req.Pocnames = script
		res, err := httpclients.OpenVulLibIdByPocnames(req)
		if err != nil {
			continue
		}
		fmt.Println(script, res.Data.LibIds)
		libIds = append(libIds, res.Data.LibIds...)
	}

	// 检测操作类型：新增 / 更新 / 强制更新
	var (
		templateSrv services.SceneTaskTemplate
		srvTemplate services.SceneTaskTemplate
	)
	taskTemplate, err := templateSrv.GetByName(ctx, templateName)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, configData, err := srvTemplate.Detail(ctx, taskTemplate.ID, taskTemplate.Estate)
	if err != nil {
		fmt.Println(err)
	}
	configData.VulIdsConfig = libIds
	fmt.Println(configData)
	_, err = templateSrv.Update(ctx, taskTemplate.ID, configData, templateName, "系统添加", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 导出任务数据
func exportTask(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出任务文件")
	var taskTask mysqls.TaskTask
	tasks, err := taskTask.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}
	filename := "task.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportTarget(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出目标文件")
	var taskTarget mysqls.TaskTarget
	targets, err := taskTarget.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(targets)
	if err != nil {
		panic(err)
	}
	filename := "target.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportTaskVul(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出任务漏洞文件")
	var taskVul mysqls.TaskVul
	taskVuls, err := taskVul.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(taskVuls)
	if err != nil {
		panic(err)
	}
	filename := "taskVul.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportTaskInfo(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出任务信息文件")
	var taskInfo mysqls.TaskTaskInfo
	taskInfos, err := taskInfo.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(taskInfos)
	if err != nil {
		panic(err)
	}
	filename := "taskInfo.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportTaskResult(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出任务扫描结果文件")
	var taskResult mysqls.TaskTaskResult
	taskResults, err := taskResult.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(taskResults)
	if err != nil {
		panic(err)
	}
	filename := "taskResult.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportTaskLog(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出任务日志文件")
	var taskLog mysqls.Tasklog
	taskLogs, err := taskLog.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(taskLogs)
	if err != nil {
		panic(err)
	}
	filename := "taskLog.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportTaskEvidence(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出任务取证文件")
	var taskEvidence mysqls.TaskEvidence
	taskEvidences, err := taskEvidence.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(taskEvidences)
	if err != nil {
		panic(err)
	}
	filename := "taskEvidence.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportRemoteSession(ctx context.Context, filter string) (string, error) {
	fmt.Println("开始导出远程会话文件")
	var remoteSession mysqls.RemoteSession
	remoteSessions, err := remoteSession.All(ctx, filter)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(remoteSessions)
	if err != nil {
		panic(err)
	}
	filename := "remoteSession.json"
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
	return filename, nil
}

// 导出任务数据
func exportTaskData(ctx context.Context) error {
	return nil
}

func getLatestTaskId(ctx context.Context) (int, int) {
	filename := "id.json"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Printf("文件不存在：%s\n", filename)
		saveLatestTaskId(ctx)
	}
	jsonByte, _ := ioutil.ReadFile(filename)
	var jsonData map[string]int
	err = json.Unmarshal(jsonByte, &jsonData)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	return jsonData["taskId"], jsonData["targetId"]

}

func saveLatestTaskId(ctx context.Context) {
	var taskTask mysqls.TaskTask
	task, _ := taskTask.GetLatestTask(ctx)
	var taskTarget mysqls.TaskTarget
	target, _ := taskTarget.GetLatestTarget(ctx)
	jsonData := map[string]int{
		"taskId":   task.ID,
		"targetId": target.ID,
	}
	jsonByte, _ := json.Marshal(jsonData)
	ioutil.WriteFile("id.json", jsonByte, 0644)
}

// importTask 导入任务表数据
func importTask(ctx context.Context, filepath string, id int) {
	fmt.Println("开始导入任务信息")
	jsonData, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var taskList []mysqls.TaskTask
	err = json.Unmarshal(jsonData, &taskList)
	if err != nil {
		panic(err)
	}
	for index, task := range taskList {
		taskEntity := mysqls.TaskTask{
			ID:             id + task.ID,
			TaskName:       task.TaskName,
			RiskLevel:      task.RiskLevel,
			Status:         task.Status,
			Weight:         task.Weight,
			IsStats:        task.IsStats,
			TaskType:       task.TaskType,
			ExecuteType:    task.ExecuteType,
			TaskTemplateID: task.TaskTemplateID,
			TargetNum:      task.TargetNum,
			HigeNum:        task.HigeNum,
			MiddleNum:      task.MiddleNum,
			LowNum:         task.LowNum,
			SafeNum:        task.SafeNum,
			UserID:         task.UserID,
			Pid:            task.Pid,
			CreateTime:     task.CreateTime,
			UpdateTime:     task.UpdateTime,
		}
		if index%1000 == 0 {
			fmt.Println("更新任务信息个数: ", index)
		}
		err = taskEntity.AddTaskCheckTask(ctx)
	}
	fmt.Println("更新任务信息个数: ", len(taskList))
}

// importTarget 导入目标数据
func importTarget(ctx context.Context, filepath string, taskId, targetId int) {
	fmt.Println("开始导入目标信息")
	jsonData, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var targetList []mysqls.TaskTarget
	err = json.Unmarshal(jsonData, &targetList)
	if err != nil {
		panic(err)
	}
	for index, target := range targetList {
		targetEntity := mysqls.TaskTarget{
			ID:               targetId + target.ID,
			TaskID:           taskId + target.TaskID,
			TargetURL:        target.TargetURL,
			Status:           target.Status,
			Weight:           target.Weight,
			RiskLevel:        target.RiskLevel,
			OpSys:            target.OpSys,
			IsAlive:          target.IsAlive,
			TargetType:       target.TargetType,
			TaskTemplateID:   target.TaskTemplateID,
			TaskTemplateJSON: target.TaskTemplateJSON,
			IsRemoteSession:  target.IsRemoteSession,
			UserID:           target.UserID,
			CreateTime:       target.CreateTime,
			UpdateTime:       target.UpdateTime,
			EndTime:          target.EndTime,
			UseScore:         target.UseScore,
			IsScore:          target.IsScore,
			ExtendField:      target.ExtendField,
		}
		if index%1000 == 0 {
			fmt.Println("更新目标信息个数: ", index)
		}
		err = targetEntity.AddTaskTarget(ctx)
	}
	fmt.Println("更新目标信息个数: ", len(targetList))
}

// importTaskVul 导入任务漏洞
func importTaskVul(ctx context.Context, filepath string, taskId, targetId int) {
	fmt.Println("开始导入任务漏洞信息")
	jsonData, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var taskVulList []mysqls.TaskVul
	err = json.Unmarshal(jsonData, &taskVulList)
	if err != nil {
		panic(err)
	}
	for index, taskVul := range taskVulList {
		taskVulEntity := mysqls.TaskVul{
			DataType:       taskVul.DataType,
			TaskID:         taskId + taskVul.TaskID,
			TargetID:       targetId + taskVul.TargetID,
			TargetUrl:      taskVul.TargetUrl,
			Pocname:        taskVul.Pocname,
			Name:           taskVul.Name,
			Class:          taskVul.Class,
			Type:           taskVul.Type,
			Risk:           taskVul.Risk,
			Location:       taskVul.Location,
			Status:         taskVul.Status,
			TestStatus:     taskVul.TestStatus,
			ExploitImpact:  taskVul.ExploitImpact,
			VulID:          taskVul.VulID,
			Description:    taskVul.Description,
			FixSuggest:     taskVul.FixSuggest,
			PublishedTime:  taskVul.PublishedTime,
			AffectRange:    taskVul.AffectRange,
			TargetResultID: taskVul.TargetResultID,
			VulNumber:      taskVul.VulNumber,
			VulAddress:     taskVul.VulAddress,
			RefUrl:         taskVul.RefUrl,
			Cvss:           taskVul.Cvss,
			VulResult:      taskVul.VulResult,
			VulParam:       taskVul.VulParam,
			VerMsg:         taskVul.VerMsg,
			DecisionVulId:  taskVul.DecisionVulId,
			Snapshot:       taskVul.Snapshot,
			CreateTime:     taskVul.CreateTime,
			UpdateTime:     taskVul.UpdateTime,
		}
		if index%1000 == 0 {
			fmt.Println("更新漏洞信息个数: ", index)
		}
		err = taskVulEntity.AddTaskVul(ctx)
	}
	fmt.Println("更新漏洞信息个数: ", len(taskVulList))
}

// importTaskInfo 导入任务信息
func importTaskInfo(ctx context.Context, filepath string, taskId, targetId int) {
	fmt.Println("开始导入任务信息表数据")
	jsonData, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var taskInfoList []mysqls.TaskTaskInfo
	err = json.Unmarshal(jsonData, &taskInfoList)
	if err != nil {
		panic(err)
	}
	for index, taskInfo := range taskInfoList {
		taskInfoEntity := mysqls.TaskTaskInfo{
			TaskID:           taskId + taskInfo.TaskID,
			TaskName:         taskInfo.TaskName,
			Status:           taskInfo.Status,
			Weight:           taskInfo.Weight,
			TaskType:         taskInfo.TaskType,
			CheckTarget:      taskInfo.CheckTarget,
			ExecuteType:      taskInfo.ExecuteType,
			ExecuteLastTime:  taskInfo.ExecuteLastTime,
			ExecuteNextTime:  taskInfo.ExecuteNextTime,
			ExecuteJSON:      taskInfo.ExecuteJSON,
			TaskTemplateID:   taskInfo.TaskTemplateID,
			TaskTemplateJSON: taskInfo.TaskTemplateJSON,
			Overview:         taskInfo.Overview,
			UserID:           taskInfo.UserID,
			CreateTime:       taskInfo.CreateTime,
			UpdateTime:       taskInfo.UpdateTime,
		}
		if index%1000 == 0 {
			fmt.Println("更新任务信息个数: ", index)
		}
		err = taskInfoEntity.AddTaskTaskInfo(ctx)
	}
	fmt.Println("更新任务信息个数: ", len(taskInfoList))
}

// importTaskInfo 导入任务结果信息
func importTaskResult(ctx context.Context, filepath string, taskId, targetId int) {
	fmt.Println("开始导入任务结果数据")
	jsonData, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var taskResultList []mysqls.TaskTaskResult
	err = json.Unmarshal(jsonData, &taskResultList)
	if err != nil {
		panic(err)
	}
	for index, taskResult := range taskResultList {
		objId, err := strconv.Atoi(taskResult.ObjID)
		if err != nil {
			continue
		}
		subObjID, err := strconv.Atoi(taskResult.SubObjID)
		if err != nil {
			continue
		}

		taskResultEntity := mysqls.TaskTaskResult{
			ObjType:    taskResult.ObjType,
			SubObjType: taskResult.SubObjType,
			ObjID:      strconv.Itoa(taskId + objId),
			SubObjID:   strconv.Itoa(targetId + subObjID),
			Identify:   taskResult.Identify,
			Field1:     taskResult.Field1,
			Field2:     taskResult.Field2,
			Field3:     taskResult.Field3,
			Field4:     taskResult.Field4,
			JSONResult: taskResult.JSONResult,
			CreateTime: taskResult.CreateTime,
		}
		if index%1000 == 0 {
			fmt.Println("更新任务结果个数: ", index)
		}
		err = taskResultEntity.AddTaskTaskResult(ctx)
	}
}

func initNodeData(ctx context.Context) {
	commandLine1 := `docker exec laozhi_db mysql -uroot -root.4dogs.cn decision --default-character-set=utf8 -e 'delete from node'`
	cmd1 := exec.Command("/bin/bash", "-c", commandLine1)
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err := cmd1.Run()
	if err != nil {
		fmt.Println("exec error: ", err)
	}
	commandLine2 := `docker exec laozhi_db mysql -uroot -proot.4dogs.cn decision --default-character-set=utf8 -e "INSERT INTO decision.node (id, node_id, node_type, token, main_ip_addr, flow_task_id, flow_server_id, node_status, create_time, update_time) VALUES (1, 'node3', 'scanner-agent', ' aba9ec64-d143-433e-b2d2-ceb4de79d653', '0.0.0.0', '', '2693724', '', '2023-08-03 08:00:01', '2024-01-11 11:26:30')"`
	cmd2 := exec.Command("/bin/bash", "-c", commandLine2)
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		fmt.Println("exec error: ", err)
	}
	fmt.Println("")
}

// 初始化数据库结构
func initDbStruct(ctx context.Context) {
	var sourceDecisionDBUrl = "xiaozhi:xiaozhi.4dogs.cn@tcp(192.168.0.68:33306)/decision"
	var sourceSmartDBUrl = "xiaozhi:xiaozhi.4dogs.cn@tcp(192.168.0.68:33306)/smart"
	dbutils.GetNewestDbStruct(sourceDecisionDBUrl, sourceSmartDBUrl)
}
