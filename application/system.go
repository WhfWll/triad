package application

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/data"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/network"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/mysql"
	"gorm.io/gorm"
)

type SystemManage struct {
}

// AuthInfo 授权信息
func (sm *SystemManage) AuthInfo(ctx context.Context, resp *typespec.AuthInfoReq, res *typespec.AuthInfoRes) (err error) {
	var mapSetSer services.MapSet
	// 判断授权状态
	authStatus := mapSetSer.GetProductAuthState(ctx)
	authInfo, err := mapSetSer.GetMapValue(ctx, "productAuthInfo")
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(authInfo), &res); err != nil {
		return err
	}
	res.Status = authStatus
	res.SoftwareVersion = enums.ProductSoftwareDisplayVersion
	if runtime.GOOS != "linux" {
		if res.AuthTime == "未授权" || res.AuthCode == "" {
			res.Status = false
		}
	}
	if strings.Trim(res.LeftDays, " ") == "--" || res.LeftDays == "0" {
		res.LeftDays = "授权已过期"
	} else if res.AuthTime == "未授权" {
		res.LeftDays = "-- " + "天"
		res.AuthDays = "-- " + "天"
	} else {
		res.LeftDays = res.LeftDays + "天"
	}
	return
}

// AuthSave 授权
func (sm *SystemManage) AuthSave(ctx context.Context, req *typespec.AuthSaveReq) (err error) {
	// 1 获取系统授权状态和授权信息
	var (
		authOldInfo typespec.AuthInfoRes
		mapSet      services.MapSet
		authStatus  bool
		authErr     error
	)
	sm.AuthInfo(ctx, nil, &authOldInfo)
	authStatus = mapSet.GetProductAuthState(ctx)
	var authService services.Auth
	authInfoRes, err := authService.RsaDecrypt(ctx, req.AuthCode)
	if err != nil {
		return errors.New("授权码错误")
	}
	// 1 判断产品ID是否相符
	proID, err := authService.GetProductID(ctx)
	// 增加 TrimSpace 比较，防止数据库里存的有空格，而授权码里没空格（或反之）导致的误判
	if err != nil || (authInfoRes["productID"] != proID && authInfoRes["productID"] != strings.TrimSpace(proID)) {
		return errors.New("产品ID错误 授权失败")
	}
	// 2 判断验证码是否处于有效期内
	generateTimeStr := authInfoRes["generateTime"]
	generateTime, _ := time.Parse(enums.TimeYMDBarLayout, generateTimeStr)
	// 计算后的过期时间 默认10年有效期
	generateTimeStr = generateTime.AddDate(0, 0, 3650).Format(enums.TimeYMDBarLayout)
	// if time.Now().Format(enums.TimeYMDBarLayout) > generateTimeStr {
	// 	return errors.New("授权码已过期")
	// }
	// 判断是否重复使用
	if authService.CheckAuthRecord(ctx, req.AuthCode) {
		return errors.New("不要重复授权")
	}
	// 3 计算验证码有效期
	if authStatus == true {
		//  判断一下老的授权码是否无效/被修改
		if _, checkErr := authService.RsaDecrypt(ctx, authOldInfo.AuthCode); checkErr != nil {
			return errors.New("老的授权码错误 请检查是否篡改")
		}
		// 如果授权过了且验证码是新的 就更新时间
		oldLeftDay, _ := strconv.Atoi(strings.Trim(authOldInfo.LeftDays, "天"))
		oldAuthDay, _ := strconv.Atoi(strings.Trim(authOldInfo.AuthDays, "天"))
		authDays, _ := strconv.Atoi(strings.Trim(authInfoRes["authDays"], "天"))
		leftDays := oldLeftDay + authDays
		newAuthDays := oldAuthDay + authDays
		authTime, _ := time.Parse(enums.TimeYMDBarLayout, authOldInfo.AuthTime)

		tempMap := make(map[string]string, 0)
		tempMap["productID"] = proID
		tempMap["generateTime"] = time.Now().Format("2006-01-02")
		tempMap["authDays"] = strconv.Itoa(newAuthDays)
		tempData, _ := json.Marshal(tempMap)
		ciphertext := authService.RsaEncrypt(tempData, []byte(enums.SWPubKey))

		// 跟新一下过期时间和剩余时间
		authErr = authService.UpdateAuthInfo(ctx, map[string]string{
			"authCode":    hex.EncodeToString(ciphertext),
			"authTime":    authOldInfo.AuthTime,
			"productID":   authOldInfo.ProductID,
			"productName": authOldInfo.ProductName,
			"authExpTime": authTime.AddDate(0, 0, newAuthDays).Format(enums.TimeYMDBarLayout),
			"authDays":    strconv.Itoa(newAuthDays) + "天",
			"leftDays":    strconv.Itoa(leftDays),
		})
	} else {
		authDays, _ := strconv.Atoi(authInfoRes["authDays"])
		authExpTimeStr := time.Now().AddDate(0, 0, authDays).Format(enums.TimeYMDBarLayout)
		// 修改授权状态
		authService.UpdateAuthState(ctx, enums.ProductAuthStateSuccess)
		// 录入授权信息
		authErr = authService.UpdateAuthInfo(ctx, map[string]string{
			"authCode":    req.AuthCode,
			"authTime":    time.Now().Format(enums.TimeYMDBarLayout),
			"productID":   authInfoRes["productID"],
			"productName": "自动化渗透测试平台",
			"authExpTime": authExpTimeStr,
			"authDays":    authInfoRes["authDays"] + "天",
			"leftDays":    authInfoRes["authDays"],
		})
	}
	// 4 如果授权成功 则记录到记录表里
	if authErr == nil {
		authService.UpdateAuthRecord(ctx, req.AuthCode)
	}
	return
}

// GenerateProductID 生成产品id
func (sm *SystemManage) GenerateProductID(ctx context.Context) error {
	var auth services.Auth
	authInfo, err := auth.GetAuthInfo(ctx)
	if authInfo["authCode"] != "" {
		return errors.New("系统已授权")
	}
	productID, err := auth.GenerateSystemSerialNumber(ctx)
	if err != nil {
		return err
	}
	authInfoMap := make(map[string]string, 0)
	authInfoMap["productID"] = productID
	authInfoMap["authCode"] = ""
	authInfoMap["authExpTime"] = ""
	authInfoMap["authTime"] = ""
	authInfoMap["leftDays"] = "0"
	authInfoMap["productName"] = "自动化渗透测试平台"
	err = auth.UpdateAuthInfo(ctx, authInfoMap)
	if err != nil {
		return err
	}
	return nil
}

// TargetIpSave 测试目标黑白名单更新
func (sm *SystemManage) TargetIpSave(ctx context.Context, req *typespec.TargetIpSaveReq) error {
	//ip解析成数组
	var targetList = make([]string, 0)
	if req.IsOpen == enums.TargetIpWhiteBlackIsOpenOn {
		var ipListString string
		if req.Type == enums.TargetIpWhiteBlackTypeWhite { //选择的是白名单
			ipListString = req.WhiteList
		} else { //选择的是黑名单
			ipListString = req.BlackList
		}
		var analysisTarget data.TaskCheckTaskAnalysisTarget
		analysisTarget.AnalysisTarget(ipListString, "")
		errorTargetList := analysisTarget.ErrorTargetList
		if len(errorTargetList) > 0 {
			return errors.New(strings.Join(errorTargetList, ","))
		}
		targetList = analysisTarget.TargetList
	}
	fmt.Println(req)

	//插入mapset表
	var objValue = services.TargetIpWhiteBlackMapSet{
		IsOpen:      req.IsOpen,
		Type:        req.Type,
		WhiteList:   req.WhiteList,
		BlackList:   req.BlackList,
		IpListArray: targetList,
	}
	objValueByte, err := json.Marshal(objValue)
	if err != nil {
		return err
	}
	var mapSet services.MapSet
	err = mapSet.Create(ctx, enums.TargetIpWhiteBlackMapSetObjKey, string(objValueByte), enums.TargetIpWhiteBlackMapSetContent)
	if err != nil {
		return err
	}
	return nil
}

// TargetIpList 测试目标黑白名单查询
func (sm *SystemManage) TargetIpList(ctx context.Context, resp *typespec.TargetIpListResp) error {
	//查询mapset
	var mapSet services.MapSet
	objValueString, err := mapSet.GetMapValue(ctx, enums.TargetIpWhiteBlackMapSetObjKey)
	if err != nil {
		return err
	}
	var objValue services.TargetIpWhiteBlackMapSet
	err = json.Unmarshal([]byte(objValueString), &objValue)
	if err != nil {
		return err
	}
	//复制到返回值
	resp.IsOpen = objValue.IsOpen
	resp.Type = objValue.Type
	resp.WhiteList = objValue.WhiteList
	resp.BlackList = objValue.BlackList
	return nil
}

// TargetIpList 测试目标黑白名单查询
func (sm *SystemManage) GetReverseIpHost(ctx context.Context, resp *typespec.GetReverseIpHostResp) error {
	//查询mapset
	var mapSet services.MapSet
	objValueString, err := mapSet.GetMapValue(ctx, enums.ReverseIpHostMapSetObjKey)
	if err != nil {
		return err
	}
	var objValue services.ReverseIpHost
	err = json.Unmarshal([]byte(objValueString), &objValue)
	if err != nil {
		return err
	}
	//复制到返回值
	resp.ReverseType = objValue.ReverseType
	resp.ReverseHost = objValue.ReverseHost
	resp.ReversePort = objValue.ReversePort
	return nil
}

// TargetIpSave 测试目标黑白名单更新
func (sm *SystemManage) ReverseIpHostSave(ctx context.Context, req *typespec.ReverseIpHostSaveReq) error {
	//插入mapset表
	var objValue = services.ReverseIpHost{
		ReverseType: req.ReverseType,
		ReverseHost: req.ReverseHost,
		ReversePort: req.ReversePort,
	}
	objValueByte, err := json.Marshal(objValue)
	if err != nil {
		return err
	}
	var mapSet services.MapSet
	err = mapSet.Create(ctx, enums.ReverseIpHostMapSetObjKey, string(objValueByte), enums.ReverseIpHostMapSetContent)
	if err != nil {
		return err
	}
	return nil
}

// SystemConfigBackupConfigSave 系统管理-配置备份-保存配置
func (sm *SystemManage) SystemConfigBackupConfigSave(ctx context.Context, req *typespec.SystemConfigBackupConfigReq) error {
	var mapSetService services.MapSet
	currentTime := time.Now()
	if req.SaveTime.IsZero() {
		req.SaveTime = currentTime
	}
	if req.RunTime.IsZero() {
		req.RunTime = currentTime.AddDate(0, req.Cycle, 0)
	}
	objValueStruct := services.SystemConfigBackupConfigMapSet{
		IsOpen:   req.IsOpen,
		Cycle:    req.Cycle,
		SaveTime: req.SaveTime,
		RunTime:  req.RunTime,
	}
	dataByte, err := json.Marshal(objValueStruct)
	if err != nil {
		return err
	}
	objValue := string(dataByte)
	return mapSetService.Create(ctx, enums.SystemConfigBackupConfigMapSetObjKey, objValue, enums.SystemConfigBackupConfigMapSetContent)
}

// SystemConfigBackupConfigInfo 系统管理-配置备份-配置信息
func (sm *SystemManage) SystemConfigBackupConfigInfo(ctx context.Context, res *typespec.SystemConfigBackupConfigInfoRes) error {
	var mapSetService services.MapSet
	mapSetValue, err := mapSetService.GetMapValue(ctx, enums.SystemConfigBackupConfigMapSetObjKey)
	if err != nil {
		return err
	}
	if mapSetValue == "" {
		return nil
	}
	var systemConfigBackupConfigMapSet services.SystemConfigBackupConfigMapSet
	if err = json.Unmarshal([]byte(mapSetValue), &systemConfigBackupConfigMapSet); err != nil {
		return err
	}
	res.Cycle = systemConfigBackupConfigMapSet.Cycle
	res.IsOpen = systemConfigBackupConfigMapSet.IsOpen
	return nil
}

// SystemConfigBackupList 系统管理-配置备份-列表
func (sm *SystemManage) SystemConfigBackupList(ctx context.Context, req *typespec.SystemConfigBackupListReq, res *typespec.SystemConfigBackupListRes) error {
	var systemConfigBackupService services.SystemConfigBackup
	systemConfigBackupList, total, err := systemConfigBackupService.SystemConfigBackupList(ctx, req.Page, req.Size)
	if err != nil {
		return err
	}
	res.Total = total
	for _, systemConfigBackup := range systemConfigBackupList {
		res.List = append(res.List, typespec.SystemConfigBackupListItemRes{
			Id:         systemConfigBackup.ID,
			Name:       systemConfigBackup.Name,
			Path:       systemConfigBackup.Path,
			CreateTime: systemConfigBackup.CreateTime.Format(enums.TimeLayout),
		})
	}
	return nil
}

// SystemConfigBackupDownload 系统管理-配置备份-下载
func (sm *SystemManage) SystemConfigBackupDownload(ctx context.Context, req *typespec.SystemConfigBackupDownloadReq, res *typespec.SystemConfigBackupDownloadRes) error {
	var systemConfigBackupService services.SystemConfigBackup
	systemConfigBackup, err := systemConfigBackupService.SystemConfigBackupRecord(ctx, req.Id)
	if err != nil {
		return err
	}
	if systemConfigBackup.Path == "" {
		return errors.New("下载文件不存在")
	}

	// 校验文件路径是否合法
	cleanPath := filepath.Clean(systemConfigBackup.Path)
	cleanBase := filepath.Clean(enums.SystemConfigBackupDir)
	if !strings.HasPrefix(cleanPath, cleanBase) {
		return errors.New("非法文件路径")
	}

	res.Path = systemConfigBackup.Path
	return nil
}

// SystemConfigBackupDelete 系统管理-配置备份-删除
func (sm *SystemManage) SystemConfigBackupDelete(ctx context.Context, req *typespec.SystemConfigBackupDeleteReq) error {
	var systemConfigBackupService services.SystemConfigBackup
	return systemConfigBackupService.SystemConfigBackupDelete(ctx, req.Id)
}

// SystemConfigBackupNow 系统管理-配置备份-立即备份
func (sm *SystemManage) SystemConfigBackupNow(ctx context.Context) error {
	var systemConfigBackupService services.SystemConfigBackup
	return systemConfigBackupService.SystemConfigBackupNow(ctx)
}

// SystemConfigBackupRestore 系统管理-配置备份-恢复
func (sm *SystemManage) SystemConfigBackupRestore(ctx context.Context, req *typespec.SystemConfigBackupRestoreReq) error {
	//获取备份文件路径
	var systemConfigBackupService services.SystemConfigBackup
	systemConfigBackup, err := systemConfigBackupService.SystemConfigBackupRecord(ctx, req.Id)
	if err != nil {
		return err
	}
	path := systemConfigBackup.Path

	//获取备份文件内容
	dataList, err := file.ReadCsv(path)
	if err != nil {
		return err
	}
	if len(dataList) == 0 {
		return errors.New("数据不存在")
	}

	//开启事务
	tx := mysql.DB.Begin()
	dCtx := mysql.NewContext(ctx, tx)
	defer tx.Rollback()

	//恢复数据
	var mapSetService services.MapSet
	for index, line := range dataList[1:] {
		if len(line) == 0 {
			log.Infof("第%d行数据为空", index+1)
			continue
		}
		err = mapSetService.Create(dCtx, line[2], line[3], line[4])
		if err != nil {
			return err
		}
	}

	//提交事务
	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// SystemSettingIpWhiteSave 系统管理 - 系统设置 - 系统访问白名单 - 保存
func (sm *SystemManage) SystemSettingIpWhiteSave(ctx context.Context, req *typespec.SystemSettingIpWhiteSaveReq) error {
	//组合objValue数据
	ipWhiteMapSet := services.SystemSettingIpWhiteMapSet{
		IsOpen: req.IsOpen,
		Ip:     req.Ip,
	}
	objValueByte, err := json.Marshal(ipWhiteMapSet)
	if err != nil {
		return err
	}
	//保存数据
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.SystemAccessIpWhiteMapSetObjKey, string(objValueByte), enums.SystemAccessIpWhiteMapSetContent)
}

// SystemSettingIpWhiteInfo 系统管理 - 系统设置 - 系统访问白名单 - 信息
func (sm *SystemManage) SystemSettingIpWhiteInfo(ctx context.Context, res *typespec.SystemSettingIpWhiteInfoRes) error {
	//获取objValue
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.SystemAccessIpWhiteMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var ipWhiteMapSet services.SystemSettingIpWhiteMapSet
	if err = json.Unmarshal([]byte(objValueStr), &ipWhiteMapSet); err != nil {
		return err
	}
	res.IsOpen = ipWhiteMapSet.IsOpen
	res.Ip = ipWhiteMapSet.Ip
	return nil
}

// SystemSettingSyslogSave 系统管理 - 系统设置 - syslog服务 - 保存
func (sm *SystemManage) SystemSettingSyslogSave(ctx context.Context, req *typespec.SystemSettingSyslogSaveReq) error {
	//组合objValue数据
	syslogMapSet := services.SystemSettingSyslogMapSet{
		IsOpen: req.IsOpen,
		Ip:     req.Ip,
		Port:   req.Port,
		Types:  req.Types,
	}
	objValueByte, err := json.Marshal(syslogMapSet)
	if err != nil {
		return err
	}
	//保存数据
	var mapSetService services.MapSet
	if err = mapSetService.Create(ctx, enums.SyslogConfigMapSetObjKey, string(objValueByte), enums.SyslogConfigMapSetContent); err != nil {
		return err
	}
	//处理syslog后台服务
	var syslogService services.SyslogServer
	return syslogService.HandleSyslog(req.Ip, req.Types, req.IsOpen, req.Port)
}

// SystemSettingSyslogInfo 系统管理 - 系统设置 - syslog服务 - 信息
func (sm *SystemManage) SystemSettingSyslogInfo(ctx context.Context, res *typespec.SystemSettingSyslogInfoRes) error {
	//获取objValue
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.SyslogConfigMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var syslogMapSet services.SystemSettingSyslogMapSet
	if err = json.Unmarshal([]byte(objValueStr), &syslogMapSet); err != nil {
		return err
	}
	res.IsOpen = syslogMapSet.IsOpen
	res.Ip = syslogMapSet.Ip
	res.Port = syslogMapSet.Port
	res.Types = syslogMapSet.Types
	return nil
}

// SystemSettingMailSave 系统设置 - 邮箱配置 - 保存
func (sm *SystemManage) SystemSettingMailSave(ctx context.Context, req *typespec.SystemSettingMailSaveReq) error {
	//邮箱密码进行16进制解码
	rawPassByte, err := hex.DecodeString(req.Password)
	if err != nil {
		return errors.New("密码加密协商不一致")
	}
	//邮箱密码加密，并进行base64编码
	encryptedPassBase64 := services.CbcPasswordEncryption(string(rawPassByte))
	//组合objValue数据
	mailMapSet := services.SystemSettingMailMapSet{
		Address:  req.Address,
		Port:     req.Port,
		Username: req.Username,
		Password: encryptedPassBase64,
		Encrypt:  req.Encrypt,
	}
	objValueByte, err := json.Marshal(mailMapSet)
	if err != nil {
		return err
	}
	//保存数据
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.MailConfigMapSetObjKey, string(objValueByte), enums.MailConfigMapSetContent)
}

// SystemSettingMailInfo 系统设置 - 邮箱配置 - 信息
func (sm *SystemManage) SystemSettingMailInfo(ctx context.Context, res *typespec.SystemSettingMailInfoRes) error {
	//获取objValue
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.MailConfigMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	//解析数据
	var mailMapSet services.SystemSettingMailMapSet
	if err = json.Unmarshal([]byte(objValueStr), &mailMapSet); err != nil {
		return err
	}
	//base64解码，邮箱密码解密
	rawPass := services.CbcPasswordDecodeString(mailMapSet.Password)
	//16进制编码
	hexPass := hex.EncodeToString([]byte(rawPass))
	//组装返回数据
	res.Address = mailMapSet.Address
	res.Port = mailMapSet.Port
	res.Username = mailMapSet.Username
	res.Password = hexPass

	return nil
}

// SystemSettingMailVerify 系统设置 - 邮箱配置 - 验证
func (sm *SystemManage) SystemSettingMailVerify(ctx context.Context, req *typespec.SystemSettingMailSaveReq) error {
	//邮箱密码进行16进制解码
	rawPassByte, err := hex.DecodeString(req.Password)
	if err != nil {
		return errors.New("密码加密协商不一致")
	}
	var mailService services.MailManage
	attachmentPaths := make([]string, 0)
	recipients := []string{req.Username}
	return mailService.SendEmail(recipients, attachmentPaths, req.Username, string(rawPassByte), req.Address, req.Port, 5*time.Second)
}

// SystemSettingNetworkConfigSave 系统管理 - 系统设置 - 网络配置 - 保存
func (sm *SystemManage) SystemSettingNetworkConfigSave(ctx context.Context, req *typespec.SystemSettingNetworkConfigSaveReq) error {
	networkConfigMapSet := services.SystemSettingNetworkConfigMapSet{
		Ip:               req.Ip,
		Mask:             req.Mask,
		Gateway:          req.Gateway,
		DnsServer:        req.DnsServer,
		StandbyDnsServer: req.StandbyDnsServer,
		WebPort:          req.WebPort,
	}
	objValueByte, err := json.Marshal(networkConfigMapSet)
	if err != nil {
		return err
	}
	var mapSetService services.MapSet
	if err = mapSetService.Create(ctx, enums.NetworkConfigMapSetObjKey, string(objValueByte), enums.NetworkConfigMapSetContent); err != nil {
		return err
	}

	if runtime.GOOS != "linux" {
		return errors.New("当前系统不支持")
	}
	output, err := exec.Command("hostnamectl").Output()
	if err != nil {
		return errors.New("get system platform err: " + err.Error())
	}
	var networkConfig services.NetWorkConfig
	if strings.Contains(string(output), "Ubuntu") {
		go func() {
			//配置网卡
			if err = networkConfig.UbuntuNetworkConfigSave(req.Ip, req.Mask, req.Gateway, req.DnsServer, req.StandbyDnsServer); err != nil {
				log.Errorf("network save err: %s", err.Error())
				return
			}
			//配置docker-compose中nginx的port
			if err = networkConfig.DockerComposeConfigUpdate(strconv.Itoa(req.WebPort)); err != nil {
				log.Errorf("docker-compose config update err: %s", err.Error())
				return
			}
			//先重启docker中的nginx
			if err = networkConfig.DockerComposeUpNginx(); err != nil {
				log.Errorf("restart nginx err: %s", err.Error())
				return
			}
			//再重启网卡
			if err = networkConfig.RestartUbuntuNetwork(); err != nil {
				log.Errorf("restart network err: %s", err.Error())
				return
			}
		}()
	} else if strings.Contains(string(output), "Centos Linux 7") {
		go func() {
			//配置网卡
			if err = networkConfig.CentosNetworkConfigSave(req.Ip, req.Mask, req.Gateway, req.DnsServer, req.StandbyDnsServer); err != nil {
				log.Errorf("network save err: %s", err.Error())
				return
			}
			//配置docker-compose中nginx的port
			if err = networkConfig.DockerComposeConfigUpdate(strconv.Itoa(req.WebPort)); err != nil {
				log.Errorf("docker-compose config update err: %s", err.Error())
				return
			}
			//先重启docker中的nginx
			if err = networkConfig.DockerComposeUpNginx(); err != nil {
				log.Errorf("restart nginx err: %s", err.Error())
				return
			}
			//再重启网卡
			if err = networkConfig.RestartCentosNetwork(); err != nil {
				log.Errorf("restart network err: %s", err.Error())
				return
			}
		}()
	} else {
		return errors.New("当前系统不支持")
	}
	if err != nil {
		return err
	}
	return nil
}

// SystemSettingNetworkConfigInfo 系统管理 - 系统设置 - 网络配置 - 信息
func (sm *SystemManage) SystemSettingNetworkConfigInfo(ctx context.Context, res *typespec.SystemSettingNetworkConfigInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.NetworkConfigMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var networkConfigMapSet services.SystemSettingNetworkConfigMapSet
	if err = json.Unmarshal([]byte(objValueStr), &networkConfigMapSet); err != nil {
		return err
	}
	res.Ip = networkConfigMapSet.Ip
	res.Mask = networkConfigMapSet.Mask
	res.Gateway = networkConfigMapSet.Gateway
	res.DnsServer = networkConfigMapSet.DnsServer
	res.StandbyDnsServer = networkConfigMapSet.StandbyDnsServer
	res.WebPort = networkConfigMapSet.WebPort
	return nil
}

// SystemSettingMonitorWarnInfo 系统管理 - 系统设置 - 系统监控告警 - 信息
func (sm *SystemManage) SystemSettingMonitorWarnInfo(ctx context.Context, res *typespec.SystemMonitorWarnInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.MonitorWarnMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var monitorWarnMapSet services.MonitorWarnMapSet
	if err = json.Unmarshal([]byte(objValueStr), &monitorWarnMapSet); err != nil {
		return err
	}
	res.IsOpen = monitorWarnMapSet.IsOpen
	res.CpuWarn = monitorWarnMapSet.CpuWarn
	res.MemoryWarn = monitorWarnMapSet.MemoryWarn
	res.DiskWarn = monitorWarnMapSet.DiskWarn
	res.FlowWarn = monitorWarnMapSet.FlowWarn
	return nil
}

// SystemSettingMonitorWarnSave 系统管理 - 系统设置 - 系统监控告警 - 保存
func (sm *SystemManage) SystemSettingMonitorWarnSave(ctx context.Context, req *typespec.SystemMonitorWarnSaveReq) error {
	monitorWarnMapSet := services.MonitorWarnMapSet{
		IsOpen:     req.IsOpen,
		CpuWarn:    req.CpuWarn,
		MemoryWarn: req.MemoryWarn,
		DiskWarn:   req.DiskWarn,
		FlowWarn:   req.FlowWarn,
	}
	objValueByte, err := json.Marshal(monitorWarnMapSet)
	if err != nil {
		return err
	}
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.MonitorWarnMapSetObjKey, string(objValueByte), enums.MonitorWarnMapSetContent)
}

// BusinessSettingTcpBlindTestSave 系统管理 - 业务设置 - Tcp盲测平台 - 保存
func (sm *SystemManage) BusinessSettingTcpBlindTestSave(ctx context.Context, req *typespec.TcpBlindTestSaveReq) error {
	tcpBlindTestMapSet := services.TcpBlindTestMapSet{
		Type: req.Type,
		Host: req.Host,
		Port: req.Port,
	}
	objValueByte, err := json.Marshal(tcpBlindTestMapSet)
	if err != nil {
		return err
	}
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.TcpBlindTestMapSetObjKey, string(objValueByte), enums.TcpBlindTestContent)
}

// BusinessSettingTcpBlindTestInfo 系统管理 - 业务设置 - Tcp盲测平台 - 信息
func (sm *SystemManage) BusinessSettingTcpBlindTestInfo(ctx context.Context, res *typespec.TcpBlindTestInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.TcpBlindTestMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var tcpBlindTestMapSet services.TcpBlindTestMapSet
	if err = json.Unmarshal([]byte(objValueStr), &tcpBlindTestMapSet); err != nil {
		return err
	}
	res.Type = tcpBlindTestMapSet.Type
	res.Host = tcpBlindTestMapSet.Host
	res.Port = tcpBlindTestMapSet.Port
	return nil
}

// BusinessSettingHttpBlindTestSave 系统管理 - 业务设置 - http盲测平台 - 保存
func (sm *SystemManage) BusinessSettingHttpBlindTestSave(ctx context.Context, req *typespec.HttpBlindTestSaveReq) error {
	httpBlindTestMapSet := services.HttpBlindTestMapSet{
		Type: req.Type,
		Host: req.Host,
		Port: req.Port,
	}
	objValueByte, err := json.Marshal(httpBlindTestMapSet)
	if err != nil {
		return err
	}
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.HttpBlindTestMapSetObjKey, string(objValueByte), enums.HttpBlindTestContent)
}

// BusinessSettingHttpBlindTestInfo 系统管理 - 业务设置 - http盲测平台 - 信息
func (sm *SystemManage) BusinessSettingHttpBlindTestInfo(ctx context.Context, res *typespec.HttpBlindTestInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.HttpBlindTestMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var httpBlindTestMapSet services.HttpBlindTestMapSet
	if err = json.Unmarshal([]byte(objValueStr), &httpBlindTestMapSet); err != nil {
		return err
	}
	res.Type = httpBlindTestMapSet.Type
	res.Host = httpBlindTestMapSet.Host
	res.Port = httpBlindTestMapSet.Port
	return nil
}

// BusinessSettingDnsBlindTestSave 系统管理 - 业务设置 - Dns盲测平台 - 保存
func (sm *SystemManage) BusinessSettingDnsBlindTestSave(ctx context.Context, req *typespec.DnsBlindTestSaveReq) error {
	dnsBlindTestMapSet := services.DnsBlindTestMapSet{
		Type:   req.Type,
		Domain: req.Domain,
	}
	objValueByte, err := json.Marshal(dnsBlindTestMapSet)
	if err != nil {
		return err
	}
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.DnsBlindTestMapSetObjKey, string(objValueByte), enums.DnsBlindTestContent)
}

// BusinessSettingDnsBlindTestInfo 系统管理 - 业务设置 - Dns盲测平台 - 信息
func (sm *SystemManage) BusinessSettingDnsBlindTestInfo(ctx context.Context, res *typespec.DnsBlindTestInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.DnsBlindTestMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var dnsBlindTestMapSet services.DnsBlindTestMapSet
	if err = json.Unmarshal([]byte(objValueStr), &dnsBlindTestMapSet); err != nil {
		return err
	}
	res.Type = dnsBlindTestMapSet.Type
	res.Domain = dnsBlindTestMapSet.Domain
	return nil
}

// BusinessSettingIcmpBlindTestSave 系统管理 - 业务设置 - Icmp盲测平台 - 保存
func (sm *SystemManage) BusinessSettingIcmpBlindTestSave(ctx context.Context, req *typespec.IcmpBlindTestSaveReq) error {
	icmpBlindTestMapSet := services.IcmpBlindTestMapSet{
		Type: req.Type,
		Host: req.Host,
	}
	objValueByte, err := json.Marshal(icmpBlindTestMapSet)
	if err != nil {
		return err
	}
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.IcmpBlindTestMapSetObjKey, string(objValueByte), enums.IcmpBlindTestContent)
}

// BusinessSettingIcmpBlindTestInfo 系统管理 - 业务设置 - Icmp盲测平台 - 信息
func (sm *SystemManage) BusinessSettingIcmpBlindTestInfo(ctx context.Context, res *typespec.IcmpBlindTestInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.IcmpBlindTestMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var icmpBlindTestMapSet services.IcmpBlindTestMapSet
	if err = json.Unmarshal([]byte(objValueStr), &icmpBlindTestMapSet); err != nil {
		return err
	}
	res.Type = icmpBlindTestMapSet.Type
	res.Host = icmpBlindTestMapSet.Host
	return nil
}

// CurTasksInfo  业务设置 - 任务并发配置 - 信息
func (sm *SystemManage) CurTasksInfo(ctx context.Context, res *typespec.CurTasksInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.CurTasksInfoMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	if err = json.Unmarshal([]byte(objValueStr), res); err != nil {
		return err
	}
	return nil
}

// CurTasksSave 业务设置 - 任务并发配置 - 保存
func (sm *SystemManage) CurTasksSave(ctx context.Context, req *typespec.CurTasksSaveReq) error {
	var mapSetService services.MapSet
	objValueByte, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return mapSetService.Create(ctx, enums.CurTasksInfoMapSetObjKey, string(objValueByte), enums.CurTasksInfoMapSetContent)
}

// SecurityScanConcurrencyInfo 系统设置 - 安全检查并发配置 - 信息
func (sm *SystemManage) SecurityScanConcurrencyInfo(ctx context.Context, res *typespec.SecurityScanConcurrencyInfoRes) error {
	cfg := services.GetSecurityScanConcurrency(ctx)
	res.HostConcurrent = cfg.HostConcurrent
	res.AppConcurrent = cfg.AppConcurrent
	res.DataConcurrent = cfg.DataConcurrent
	return nil
}

// SecurityScanConcurrencySave 系统设置 - 安全检查并发配置 - 保存
func (sm *SystemManage) SecurityScanConcurrencySave(ctx context.Context, req *typespec.SecurityScanConcurrencySaveReq) error {
	return services.SaveSecurityScanConcurrency(ctx, services.SecurityScanConcurrencyConfig{
		HostConcurrent: req.HostConcurrent,
		AppConcurrent:  req.AppConcurrent,
		DataConcurrent: req.DataConcurrent,
	})
}

// UseScoreInfo 业务设置 - 可以利用评分 - 信息
func (sm *SystemManage) UseScoreInfo(ctx context.Context, res *typespec.UseScoreInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.UseScoreSwitchMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	if res.IsOpen, err = strconv.Atoi(objValueStr); err != nil {
		return err
	}
	return nil
}

// UseScoreSave 业务设置 - 可以利用评分 - 保存
func (sm *SystemManage) UseScoreSave(ctx context.Context, req *typespec.UseScoreSaveReq) error {
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.UseScoreSwitchMapSetObjKey, strconv.Itoa(req.IsOpen), enums.UseScoreSwitchMapSetContent)
}

// TestScopeInfo 业务设置 - 测试范围校验开关 - 信息
func (sm *SystemManage) TestScopeInfo(ctx context.Context, res *typespec.TestScopeInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.TestScopeSwitchMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	if res.IsOpen, err = strconv.Atoi(objValueStr); err != nil {
		return err
	}
	return nil
}

// TestScopeSave 业务设置 - 测试范围校验开关 - 保存
func (sm *SystemManage) TestScopeSave(ctx context.Context, req *typespec.TestScopeSaveReq) error {
	var mapSetService services.MapSet
	return mapSetService.Create(ctx, enums.TestScopeSwitchMapSetObjKey, strconv.Itoa(req.IsOpen), enums.TestScopeSwitchMapSetContent)
}

// CpuInfo 系统管理 - 系统监控 - cpu
func (sm *SystemManage) CpuInfo(ctx context.Context, res *typespec.CpuInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.MonitorCpuMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var cpuMapSet services.MonitorCpuMemoryMapSet
	if err = json.Unmarshal([]byte(objValueStr), &cpuMapSet); err != nil {
		return err
	}
	res.List = cpuMapSet.List
	return nil
}

// MemoryInfo 系统管理 - 系统监控 - 内存
func (sm *SystemManage) MemoryInfo(ctx context.Context, res *typespec.MemoryInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.MonitorMemoryMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var memoryMapSet services.MonitorCpuMemoryMapSet
	if err = json.Unmarshal([]byte(objValueStr), &memoryMapSet); err != nil {
		return err
	}
	res.List = memoryMapSet.List
	return nil
}

// DiskInfo 系统管理 - 系统监控 - 磁盘
func (sm *SystemManage) DiskInfo(ctx context.Context, res *typespec.DiskInfoRes) error {
	var mapSetService services.MapSet
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.MonitorDiskMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		return nil
	}
	var diskMapSet services.MonitorDiskMapSet
	if err = json.Unmarshal([]byte(objValueStr), &diskMapSet); err != nil {
		return err
	}
	res.Free = diskMapSet.Free
	res.Used = diskMapSet.Used
	res.Total = diskMapSet.Total
	res.FreePercent = diskMapSet.FreePercent
	res.UsedPercent = diskMapSet.UsedPercent
	return nil
}

// RouteList 路由列表
func (sm *SystemManage) RouteList(res *typespec.SystemRouteListRes) error {
	var systemRoute network.SystemRoute
	routeList, err := systemRoute.RouteQuery()
	if err != nil {
		return err
	}
	res.List = make([]typespec.SystemRouteListItemRes, 0, len(routeList.List))
	for _, item := range routeList.List {
		res.List = append(res.List, typespec.SystemRouteListItemRes{
			Ip:      item.Ip,
			Netmask: item.Netmask,
			Gateway: item.Gateway,
		})
	}
	return nil
}

// RouteAdd 增加路由
func (sm *SystemManage) RouteAdd(req *typespec.SystemRouteAddReq) error {
	var systemRoute network.SystemRoute
	return systemRoute.RouteAdd(req.Ip, req.Netmask, req.Gateway)
}

// RouteDelete 删除路由
func (sm *SystemManage) RouteDelete(req *typespec.SystemRouteDeleteReq) error {
	var systemRoute network.SystemRoute
	return systemRoute.RouteDelete(req.Ip, req.Netmask, req.Gateway)
}

// SystemVersion 系统管理 - 升级还原 - 系统版本信息
func (sm *SystemManage) SystemVersion(ctx context.Context, res *typespec.SystemVersionRes) error {
	var (
		mapSetService       services.MapSet
		systemVersionMapSet services.SystemVersionMapSet
	)
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if err != nil {
		return err
	}
	if objValueStr == "" {
		systemVersionMapSet.CurrentVersion = enums.SystemBaseVersion
		objValueByte, err := json.Marshal(systemVersionMapSet)
		if err != nil {
			return err
		}
		if err = mapSetService.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent); err != nil {
			return err
		}
	}
	if objValueStr != "" {
		if err = json.Unmarshal([]byte(objValueStr), &systemVersionMapSet); err != nil {
			return err
		}
	}
	res.CurrentVersion = systemVersionMapSet.CurrentVersion
	res.UpdateTime = systemVersionMapSet.UpdateTime
	res.VulUpdateTime = systemVersionMapSet.VulUpdateTime
	res.VulVersion = systemVersionMapSet.VulVersion
	res.LastSystemVersion = systemVersionMapSet.LastSystemVersion
	res.LastVulVersion = systemVersionMapSet.LastVulVersion
	return nil
}

// GeneratePlatformToken 系统管理 - api秘钥 - 生成秘钥
func (sm *SystemManage) GeneratePlatformToken(ctx context.Context, req *typespec.GenerateTokenReq, resp *typespec.GenerateTokenResp) error {
	var (
		srv services.User
		err error
	)
	platformToken, err := srv.GeneratePlatformToken(ctx, req.UserName)
	if err != nil {
		return err
	}
	resp.Token = platformToken
	return nil
}

// TokenList 系统管理 - api秘钥 - 秘钥列表
func (sm *SystemManage) TokenList(ctx context.Context, req *typespec.TokenListReq, resp *typespec.TokenListResp) error {
	var (
		srv services.User
		err error
	)
	userList, total, err := srv.GetUserListByToken(ctx, "", req.Page, req.Size)
	if err != nil {
		return err
	}
	for _, user := range userList {
		resp.List = append(resp.List, typespec.TokenListItemRes{
			Username:   user.Username,
			Token:      user.Token,
			CreateTime: user.TokenCreateTime.Format(enums.TimeLayout),
		})
	}
	resp.Total = total
	return nil
}

//
//// 节点管理  - 设置是否开启分布式
//func (sm *SystemManage) SystemNodeSetDistribute(ctx context.Context, req *typespec.SystemNodeSetDistributeReq) error {
//	var param httpclients.OpenYakNodeSetDistributeReq
//	param.Status = req.Status
//
//	remoteRes, err := httpclients.EditDecisionYakNodeSetDistribute(ctx, param)
//	if err != nil {
//		return err
//	}
//	if remoteRes.Code != 200 {
//		return errors.New(remoteRes.Msg)
//	}
//	return nil
//}
//
//// 节点管理  - 获取是否开启分布式
//func (sm *SystemManage) SystemNodeGetDistribute(ctx context.Context, res *typespec.SystemNodeIsDistributeRes) error {
//	var param httpclients.OpenYakNodeGetDistributeReq
//	remoteRes, err := httpclients.GetDecisionYakNodeGetDistribute(ctx, param)
//	if err != nil {
//		return err
//	}
//	if remoteRes.Code != 200 {
//		return errors.New(remoteRes.Msg)
//	}
//
//	res.Status = remoteRes.Data.Status
//	return nil
//}
//
//// 节点管理 - 新增节点
//func (sm *SystemManage) SystemNodeAdd(ctx context.Context, req *typespec.SystemNodeAddReq) error {
//	var param httpclients.OpenYakNodeAddReq
//	param.Ip = req.Ip
//	param.Port = req.Port
//	param.Name = req.Name
//
//	remoteRes, err := httpclients.AddDecisionYakNode(ctx, param)
//	if err != nil {
//		return err
//	}
//	if remoteRes.Code != 200 {
//		return errors.New(remoteRes.Msg)
//	}
//	return nil
//}
//
//// 节点管理 - 节点列表
//func (sm *SystemManage) SystemNodeList(ctx context.Context, req *typespec.SystemNodeListReq, res *typespec.SystemNodeListRes) error {
//	var param httpclients.OpenYakNodeListReq
//	param.Page = req.Page
//	param.Size = req.Size
//	remoteList, err := httpclients.GetDecisionYakNodeList(ctx, param)
//	if err != nil {
//		return err
//	}
//	if remoteList.Code != 200 {
//		return errors.New(remoteList.Msg)
//	}
//
//	res.Total = remoteList.Data.Total
//	for _, item := range remoteList.Data.List {
//		res.List = append(res.List, typespec.SystemNodeListItem{
//			Id:            item.Id,
//			Name:          item.Name,
//			Ip:            item.Ip,
//			Port:          item.Port,
//			RunningNum:    item.RunningNum,
//			Status:        item.Status,
//			StatusEnum:    item.StatusEnum,
//			IsDisable:     item.IsDisable,
//			IsDisableEnum: item.IsDisableEnum,
//			CreateTime:    item.CreateTime,
//			UpdateTime:    item.CreateTime,
//		})
//	}
//	return nil
//}
//
//// 节点管理 - 删除节点
//func (sm *SystemManage) SystemNodeDel(ctx context.Context, req *typespec.SystemNodeDelReq, res *typespec.SystemNodeDelRes) error {
//	var param httpclients.OpenYakNodeDelReq
//	param.Id = req.Id
//	remoteRes, err := httpclients.DelDecisionYakNode(ctx, param)
//	if err != nil {
//		return err
//	}
//	if remoteRes.Code != 200 {
//		return errors.New(remoteRes.Msg)
//	}
//	return nil
////}
//
//// 节点管理 - 禁用｜启用节点
//func (sm *SystemManage) SystemNodeDisOrEnable(ctx context.Context, req *typespec.SystemNodeDisOrEnableReq, res *typespec.SystemNodeDisOrEnableRes) error {
//	var param httpclients.OpenYakNodeSetIsDisableReq
//	param.Id = req.Id
//	param.IsDisable = req.IsDisable
//	remoteRes, err := httpclients.EditDecisionYakNodeIsDisable(ctx, param)
//	if err != nil {
//		return err
//	}
//	if remoteRes.Code != 200 {
//		return errors.New(remoteRes.Msg)
//	}
//	return nil
//}
//
//// 节点管理 - 禁用｜启用节点
//func (sm *SystemManage) SystemNodeAllEnable(ctx context.Context, res *typespec.SystemNodeAllEnableRes) error {
//	var param httpclients.OpenYakNodeAllEnableReq
//	remoteRes, err := httpclients.GetDecisionYakNodeAllEnable(ctx, param)
//	if err != nil {
//		return err
//	}
//	if remoteRes.Code != 200 {
//		return errors.New(remoteRes.Msg)
//	}
//	for _, item := range remoteRes.Data.List {
//		res.List = append(res.List, typespec.SystemNodeAllEnableItem{
//			Id:   item.Id,
//			Name: item.Name,
//		})
//	}
//	return nil
//}

// SystemMessageAdd 消息中心 - 添加消息
func (sm *SystemManage) SystemMessageAdd(ctx context.Context) {

}

// runtimeTablesToClean 与当前库中实际存在的运行时表一致（不含已废弃/未建表的模型表）
var runtimeTablesToClean = []string{
	// 安全报告
	"security_report", "report_record",
	"report_verify_port", "report_verify_target", "report_verify_task", "report_verify_vul",
	// 主机 / 数据安全检查结果（表名以库为准：host_vul_*）
	"baseline_check_result", "host_vul_scan", "host_vul_finding",
	"host_malware_scan", "malware_check_result",
	"db_check_result", "sensitive_data_result",
	// 远程会话
	"remote_session",
	// 流程任务运行数据（功能已停用，表仍在库中）
	"flow_log", "flow_risk", "flow_target", "flow_task", "flow_base",
	// 资产扫描产生的运行时数据（保留 asset / asset_group / asset_connection）
	"asset_log", "asset_task_result", "asset_vul", "asset_port", "asset_risk_trend",
	// 渗透/扫描任务（子表优先）
	"task_log_info", "task_log", "task_result", "task_task_result",
	"task_vul", "task_evidence", "task_target_result", "task_task_info",
	"task_target", "task_configuration", "task_task",
}

func isMissingTableErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "1146") ||
		strings.Contains(msg, "doesn't exist") ||
		strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "no such table")
}

// clearRuntimeTableByName 按实际表名清空单表
func clearRuntimeTableByName(db *gorm.DB, tableName string) error {
	truncateSQL := "TRUNCATE TABLE `" + tableName + "`"
	if err := db.Exec(truncateSQL).Error; err != nil {
		if isMissingTableErr(err) {
			log.Warnf("[CleanupRuntimeData] 表 %s 不存在，跳过", tableName)
			return nil
		}
		log.Warnf("[CleanupRuntimeData] TRUNCATE %s 失败，改用 DELETE: %v", tableName, err)
		deleteSQL := "DELETE FROM `" + tableName + "`"
		if err := db.Exec(deleteSQL).Error; err != nil {
			if isMissingTableErr(err) {
				log.Warnf("[CleanupRuntimeData] 表 %s 不存在，跳过", tableName)
				return nil
			}
			return err
		}
	}
	return nil
}

// CleanupRuntimeData 清理运行数据（任务、结果、报告等），保留规则、漏洞库、策略
func (sm *SystemManage) CleanupRuntimeData(ctx context.Context) error {
	log.Warn("[CleanupRuntimeData] 开始清理运行数据...")

	db := mysql.FromContext(ctx)

	// 禁用外键检查
	_ = db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error

	var failedTables []string
	for _, tableName := range runtimeTablesToClean {
		log.Infof("[CleanupRuntimeData] 清理表: %s ...", tableName)
		if err := clearRuntimeTableByName(db, tableName); err != nil {
			log.Errorf("[CleanupRuntimeData] 清理表 %s 失败: %v", tableName, err)
			failedTables = append(failedTables, tableName)
		}
	}

	// 恢复外键检查
	_ = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error

	if len(failedTables) > 0 {
		return fmt.Errorf("以下表清理失败: %v", strings.Join(failedTables, ", "))
	}

	log.Warn("[CleanupRuntimeData] 运行数据清理完成")
	return nil
}

// TokenDelete 系统管理 - api秘钥 - 秘钥删除
func (sm *SystemManage) TokenDelete(ctx context.Context, req *typespec.TokenDeleteReq, resp *typespec.TokenDeleteResp) error {
	var srv services.User
	err := srv.TokenDelete(ctx, req.Username, req.Token)
	if err != nil {
		return err
	}
	return nil
}

// SystemManualRollback 手动回滚
func (sm *SystemManage) SystemManualRollback(ctx context.Context, req *typespec.SystemManualRollbackReq) error {
	// 1. 检查状态
	if services.IsUpgrading() {
		return errors.New("当前正在升级中，无法执行回滚操作")
	}

	status := services.GetUpgradeStatus()

	// Override type if provided in request
	if req != nil && req.Type != "" {
		if req.Type == "SYSTEM" {
			status.Type = services.UpgradeTypeSystem
		} else if req.Type == "VULN" {
			status.Type = services.UpgradeTypeVuln
		}
		// Update global status type so frontend sees correct type
		services.SetUpgradeType(status.Type)
	}

	if status.BackupDir == "" {
		backupRoot := filepath.Join(enums.SystemUpgradeProjectDir, "smart/backup")
		if entries, err := ioutil.ReadDir(backupRoot); err == nil {
			var prefixes []string
			switch status.Type {
			case services.UpgradeTypeSystem:
				prefixes = []string{"sys_"}
			case services.UpgradeTypeVuln:
				prefixes = []string{"vuln_"}
			default:
				prefixes = []string{"sys_", "vuln_"}
			}
			type candidate struct {
				path string
				ts   int64
			}
			var cands []candidate
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				name := e.Name()
				match := false
				for _, p := range prefixes {
					if strings.HasPrefix(name, p) {
						match = true
						break
					}
				}
				if !match {
					continue
				}
				idx := strings.LastIndex(name, "_")
				if idx < 0 || idx+1 >= len(name) {
					continue
				}
				ts, err2 := strconv.ParseInt(name[idx+1:], 10, 64)
				if err2 != nil {
					continue
				}
				cands = append(cands, candidate{path: filepath.Join(backupRoot, name), ts: ts})
			}
			if len(cands) > 0 {
				sort.Slice(cands, func(i, j int) bool { return cands[i].ts < cands[j].ts })
				status.BackupDir = cands[len(cands)-1].path
				services.SetUpgradeBackupDir(status.BackupDir)
				log.Infof("自动定位到最近的备份: %s", status.BackupDir)
			}
		}
	}

	if status.BackupDir == "" {
		return errors.New("未找到备份目录，无法回滚")
	}

	if !file.CheckPathExist(status.BackupDir) {
		return fmt.Errorf("备份目录不存在: %s", status.BackupDir)
	}

	// 2. 更新状态
	services.UpdateUpgradeStatus(enums.UpgradeStateRollbacking, "正在执行手动回滚...", 0)

	// 3. 异步执行回滚
	go func() {
		defer func() {
			if r := recover(); r != nil {
				services.SetRollbackError(fmt.Errorf("回滚过程发生 Panic: %v", r))
			}
		}()

		var strategy services.UpgradeStrategy
		// Determine strategy based on type
		switch status.Type {
		case services.UpgradeTypeSystem:
			strategy = &services.SystemUpgradeStrategy{}
		case services.UpgradeTypeVuln:
			strategy = &services.VulnUpgradeStrategy{}
		default:
			if strings.Contains(status.BackupDir, "vuln_") {
				strategy = &services.VulnUpgradeStrategy{}
			} else if strings.Contains(status.BackupDir, "sys_") {
				strategy = &services.SystemUpgradeStrategy{}
			} else {
				services.SetRollbackError(errors.New("未知升级类型，无法回滚"))
				return
			}
		}

		if err := strategy.Rollback(context.Background(), status.BackupDir); err != nil {
			services.SetRollbackError(fmt.Errorf("回滚失败: %v", err))
			return
		}

		// Restore success
		services.UpdateUpgradeStatus(enums.UpgradeStateRollbackSuccess, "手动回滚成功", 100)
	}()

	return nil
}
