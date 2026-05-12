package global

//
//import (
//	"encoding/json"
//	log "github.com/sirupsen/logrus"
//	"gitlabee.4dogs.cn/common/config"
//	"sync"
//	"unicode"
//)
//
//var Global globalStruct
//
//type globalStruct struct {
//	// 渗透利用 通知
//	SeepUseMap sync.Map
//	// 关闭被动流量 通知
//	PassiveTrafficMitmCloseChan map[string]chan bool
//	// 关闭被动流量 数据监听 通知
//	PassiveTrafficListenDataCloseChan map[string]chan bool
//
//	// 反弹shell通信
//	// [checkResultId]chan
//	BounceShellChan sync.Map
//}
//
//var GlobalFingerType map[string]string
//
//func init() {
//	LoadGlobalConf()
//}
//
//func LoadGlobalConf() {
//	var fingerTypeFile string
//	err := config.Load("fingerTypeFile", &fingerTypeFile)
//	if err != nil {
//		log.Panic("filepathprefix config load err:" + err.Error())
//	}
//	//jsonStr := gfile.GetBytes(fingerTypeFile)
//	jsonStr := []byte{}
//	// 转换map
//	type dataStruct struct {
//		Name       string `json:"name"`
//		FingerType string `json:"finger_type"`
//	}
//	var data []dataStruct
//	if err = json.Unmarshal(jsonStr, &data); err != nil {
//		panic("fingerTypeFile data err")
//	}
//
//	// 转换map
//	GlobalFingerType = make(map[string]string)
//	for _, v := range data {
//		// 1 空格、中线、下滑线：去掉
//		// 2 统一转小写字母
//		GlobalFingerType[ConvLower(v.Name)] = v.FingerType
//	}
//}
//func ConvLower(s string) string {
//	//if gstr.ContainsI(s, " ") {
//	//	s = strings.Replace(s, " ", "", -1)
//	//}
//	//if gstr.ContainsI(s, "-") {
//	//	s = strings.Replace(s, "-", "", -1)
//	//}
//	//
//	//if gstr.ContainsI(s, "_") {
//	//	s = strings.Replace(s, "_", "", -1)
//	//}
//	returnStr := ""
//	for _, v := range s {
//		v = unicode.ToLower(v)
//		returnStr += string(v)
//	}
//	return returnStr
//}
