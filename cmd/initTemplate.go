//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/gookit/color"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	"smart/client/httpclients"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"strings"
	"time"
)

/*
*
初始化场景

	-name 场景名称 必填
	-import_type 导入类型 必填
		[
			eq_pocnames 完整匹配pocname |
			like_pocnames 模糊匹配pocnames |
			script_type 脚本类型
		]
	-data 依据import_type导入的数据 多个逗号分割 必填
	-force 是否强制覆盖 默认false 选填
		[
			true：强制覆盖场景名一样的漏洞数据
			false：仅更新
		]
*/

// 入参
var (
	configFile string
	name       string
	importType string
	data       string
	force      bool
	waitSecond int // 每一步操作间隔时间，默认n秒
)

// 枚举
// 可选值枚举
const importTypeEqPocnames = "eq_pocnames"
const importTypeLikePocnames = "like_pocnames"
const importTypeScriptType = "script_type"

// 导入类型可选值
var importTypeInArrayMap = map[string]struct{}{importTypeEqPocnames: {}, importTypeLikePocnames: {}, importTypeScriptType: {}}

// 默认配置项
const defaultConfig = `{
    "portScanConfig":{
        "isOpen": true,
        "intelligencePort":false,
        "tcpScanType":1,
        "tcpScanTypeZh":"",
        "portScanType":100,
        "portScanTypeZh":"",
        "scanPort":"7,5555,9,13,21,22,23,25,26,37,53,79,80,81,88,106,110,111,113,119,135,139,143,144,179,199,389,427,443,444,445,465,513,514,515,543,544,548,554,587,631,646,873,888,990,993,995,1025,1026,1027,1028,1029,1080,1110,1433,1443,1720,1723,1755,1900,2000,2001,2049,2121,2181,2717,3000,3128,3306,3389,3986,4899,5000,5009,5051,5060,5101,5190,5357,5432,5631,5666,5800,5900,6000,6001,6646,7000,7001,7002,7003,7004,7005,7070,8000,8008,8009,8080,8081,8443,8888,9100,9999,10000,11211,32768,49152,49153,49154,49155,49156,49157,8088,9090,8090,8001,82,9080,8082,8089,9000,8002,89,8083,8200,90,8086,801,8011,8085,9001,9200,8100,8012,85,8084,8070,8091,8003,99,7777,8010,8028,8087,83,808,38888,8181,800,18080,8099,8899,86,8360,8300,8800,8180,3505,9002,8053,1000,7080,8989,28017,9060,8006,41516,880,8484,6677,8016,84,7200,9085,5555,8280,1980,8161,9091,7890,8060,6080,8880,8020,889,8881,9081,7007,8004,38501,1010,17,19,255,1024,1030,1041,1048,1049,1053,1054,1056,1064,1065,1801,2103,2107,2967,3001,3703,5001,5050,6004,8031,10010,10250,10255,6888,87,91,92,98,1081,1082,1118,1888,2008,2020,2100,2375,3008,6648,6868,7008,7071,7074,7078,7088,7680,7687,7688,8018,8030,8038,8042,8044,8046,8048,8069,8092,8093,8094,8095,8096,8097,8098,8101,8108,8118,8172,8222,8244,8258,8288,8448,8834,8838,8848,8858,8868,8879,8983,9008,9010,9043,9082,9083,9084,9086,9087,9088,9089,9092,9093,9094,9095,9096,9097,9098,9099,9443,9448,9800,9981,9986,9988,9998,10001,10002,10004,10008,12018,12443,14000,16080,18000,18001,18002,18004,18008,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018"
    },
    "webCrawlerConfig":{
        "isOpen":false,
        "maxDepth":0,
        "maxDepthZh":"",
        "maxUrl":0,
        "maxUrlZh":"",
        "scanRange":0,
        "scanRangeZh":"",
        "timeout":0,
        "timeoutZh":"",
        "fullTimeout":0,
        "fullTimeoutZh":"",
        "scanRepeat":0,
        "scanRepeatZh":"",
        "blackList":"",
        "whiteList":"",
        "headers":null,
        "suffixFilter":"",
        "localStorage":{
            "isOpen":false,
            "list":null
        }
    },
    "webPathScanConfig":{
        "isOpen":false,
        "guessRate":0,
        "guessRateZh":"",
        "guessTimeout":0,
        "guessTimeoutZh":"",
        "titleBlack":"",
        "scanDict":null,
        "dickNames":null,
        "isIntelligent":false
    },
    "weakPassConfig":{
        "isOpen":false,
        "services":null,
        "servicesZh":null,
        "dictType":0,
        "dictTypeZh":"",
        "commonUserDict":0,
        "commonUserDictZh":"",
        "commonPassDict":0,
        "commonPassDictZh":"",
        "addAccount":"",
        "addPass":"",
        "onlyUseAdd":false,
        "guessNum":0,
        "guessNumZh":"",
        "guessTimeout":0,
        "guessTimeoutZh":"",
        "guessRate":0,
        "guessRateZh":"",
        "captchaMode":""
    },
    "subdomainCollectConfig":{
        "isOpen":false,
        "subdomainDict":0,
        "subdomainDictZh":""
    },
    "websiteLoginConfig":{
        "isOpen":false,
        "list":null
    },
    "vulIdsConfig":null,
    "testIntensity":0,
    "vulExploit":false,
    "safeTest":false,
    "lateralMove":{
        "targetNum":0,
        "jumpNum":0,
        "range":"",
        "isOpen":false
    },
    "mode":{
        "mode":0,
        "modeZh":"",
        "distributeNodeId":0
    },
    "proxyConfig":{
        "isOpen":false,
        "addr":"",
        "port":"",
        "proto":0,
        "isAuth":false,
        "username":"",
        "password":""
    }
}`

func main() {
	flag.StringVar(&configFile, "config", "", "脚本使用到了配置文件中的decision远程调用地址，所以默认只允许在config.json同级目录下使用此脚本，否则必须传配置文件地址")
	flag.StringVar(&name, "name", "", "场景名称 必填")
	flag.StringVar(&importType, "import_type", "", "导入类型 可选值【eq_pocnames 完整匹配pocname | like_pocnames 模糊匹配pocnames | script_type 脚本类型匹配】 必填")
	flag.StringVar(&data, "data", "", "依据import_type导入的数据 多个逗号分割 必填")
	flag.BoolVar(&force, "force", false, "是否强制覆盖 默认false 选填")
	flag.IntVar(&waitSecond, "wait_second", 5, "每一步操作间隔时间，默认5秒")
	flag.Parse()

	// 校验参数
	err := checkParam()
	if err != nil {
		color.Redln("[error] " + err.Error())
		return
	}

	// 配置文件初始化
	if configFile != "" {
		err = config.NewConfig(configFile)
	} else {
		err = config.NewConfig()
	}
	if err != nil {
		color.Redln("[error] new config err:" + err.Error())
		return
	}

	// 数据库初始化
	mysql.Setup()

	// 处理结果
	err = process()
	if err != nil {
		color.Redln("[error] " + err.Error())
		return
	}

	// 结束
	color.Greenln("Success")
}

// 校验必填参数
func checkParam() (err error) {
	errMsg := make([]string, 0)
	if name == "" {
		errMsg = append(errMsg, "name场景名称必填")
	}
	if importType == "" {
		errMsg = append(errMsg, "import_type导入类型必填")
	}
	if data == "" {
		errMsg = append(errMsg, "依据import_type导入的数据必填")
	}
	// 导入类型可选值校验
	if _, ok := importTypeInArrayMap[importType]; !ok {
		errMsg = append(errMsg, "import_type可选值错误")
	}

	if len(errMsg) > 0 {
		err = errors.New(strings.Join(errMsg, " && "))
	}

	return
}

// 处理数据
func process() error {
	ctx := context.Background()

	libIds, err := getLibIds()
	if err != nil {
		return err
	}
	if len(libIds) == 0 {
		return errors.New("未查询到任何漏洞，请确认参数是否正确")
	}
	color.Info.Println("检测到" + strconv.Itoa(len(libIds)) + "个漏洞信息，" + strconv.Itoa(waitSecond) + "秒后开始同步")
	// 睡眠n秒用于给用户提供反悔时间
	time.Sleep(time.Duration(waitSecond) * time.Second)

	// 检测操作类型：新增 / 更新 / 强制更新
	var templateSrv services.SceneTaskTemplate
	taskTemplate, err := templateSrv.GetByName(ctx, name)
	if err != nil {
		return err
	}

	// 更新操作
	if taskTemplate.ID > 0 {
		color.Info.Println("模板名为【" + name + "】已存在，此次将进行更新，" + strconv.Itoa(waitSecond) + "秒后继续")
		// 睡眠n秒用于给用户提供反悔时间
		time.Sleep(time.Duration(waitSecond) * time.Second)

		// 获取当前模板下的配置
		var srvTemplate services.SceneTaskTemplate
		_, configData, err := srvTemplate.Detail(ctx, taskTemplate.ID, taskTemplate.Estate)
		if err != nil {
			return err
		}

		if force {
			color.Info.Println("本次操作将覆盖更新，" + strconv.Itoa(waitSecond) + "秒后继续")
			// 睡眠n秒用于给用户提供反悔时间
			time.Sleep(time.Duration(waitSecond) * time.Second)

			// 强制替换
			configData.VulIdsConfig = libIds

			_, err = templateSrv.Update(ctx, taskTemplate.ID, configData, name, "系统添加", 0)
			if err != nil {
				return err
			}
		} else {
			color.Info.Println("本次操作仅添加不存在的漏洞，" + strconv.Itoa(waitSecond) + "秒后继续")
			// 睡眠n秒用于给用户提供反悔时间
			time.Sleep(time.Duration(waitSecond) * time.Second)

			//
			alreadyVulIdMap := make(map[int]struct{})
			for _, vulId := range configData.VulIdsConfig {
				alreadyVulIdMap[vulId] = struct{}{}
			}
			for _, vulId := range libIds {
				if _, ok := alreadyVulIdMap[vulId]; !ok {
					configData.VulIdsConfig = append(configData.VulIdsConfig, vulId)
				}
			}

			_, err = templateSrv.Update(ctx, taskTemplate.ID, configData, name, "系统添加", 0)
			if err != nil {
				return err
			}
		}
	} else {
		// 新增操作
		color.Warn.Println("模板名为【" + name + "】不存在，此次将进行添加，" + strconv.Itoa(waitSecond) + "秒后继续")
		// 睡眠n秒用于给用户提供反悔时间
		time.Sleep(time.Duration(waitSecond) * time.Second)

		// 组织添加的配置
		var templateConfig enums.ConfigJson
		err = json.Unmarshal([]byte(defaultConfig), &templateConfig)
		if err != nil {
			return errors.New("场景配置格式化失败：" + err.Error())
		}
		templateConfig.VulIdsConfig = libIds

		_, err = templateSrv.Create(ctx, templateConfig, name, "系统添加", 0)
		if err != nil {
			return errors.New("添加失败：" + err.Error())
		}
	}

	return nil
}

// 依据传的参数获取脚本ID
func getLibIds() (libIds []int, err error) {
	switch importType {
	case importTypeEqPocnames, importTypeLikePocnames: // [完整匹配 或 like] pocname
		var req httpclients.OpenVulLibIdByPocnamesReq
		req.MatchType = importType
		req.Pocnames = data
		res, err := httpclients.OpenVulLibIdByPocnames(req)
		if err != nil {
			return nil, err
		}
		libIds = res.Data.LibIds
	case importTypeScriptType: // 脚本类型匹配
		var req httpclients.OpenVulLibIdByScriptReq
		req.ScriptType = data
		res, err := httpclients.OpenVulLibIdByScript(req)
		if err != nil {
			return nil, err
		}
		libIds = res.Data.LibIds
		break
	}
	return
}
