package services

import (
	log "github.com/sirupsen/logrus"
	"smart/tools/mylog/syslog"
	logrus_syslog "smart/tools/mysyslog"
	"strconv"

	"smart/tools/enums"
	"strings"
)

type SyslogServer struct {
}

type CustomHook struct {
	*logrus_syslog.SyslogHook
	LevelList []log.Level
}

func (h *CustomHook) Levels() []log.Level {
	return h.LevelList
}

// HandleSyslog 同步日志到syslog服务器
func (s *SyslogServer) HandleSyslog(ip, syslogTypes string, isOpen, port int) error {
	//直接替换钩子列表
	log.StandardLogger().Hooks = make(log.LevelHooks)
	if isOpen == enums.ConfigClose {
		return nil
	}
	//设置syslog日志级别
	syslogLevels := make([]log.Level, 0)
	for _, level := range strings.Split(syslogTypes, ",") {
		if level == enums.SyslogTypeAudit {
			syslogLevels = append(syslogLevels, log.InfoLevel)
		} else if level == enums.SyslogTypeDebug {
			syslogLevels = append(syslogLevels, log.DebugLevel)
		} else if level == enums.SyslogTypeWarn {
			syslogLevels = append(syslogLevels, log.WarnLevel)
		} else if level == enums.SyslogTypeError {
			syslogLevels = append(syslogLevels, log.ErrorLevel)
		} else {
			syslogLevels = append(syslogLevels, log.FatalLevel)
			log.StandardLogger().Hooks = make(log.LevelHooks)
			return nil
		}
	}
	//初始化syslog钩子
	hook, err := logrus_syslog.NewSyslogHook("udp", ip+":"+strconv.Itoa(port), syslog.LOG_INFO, "smartSyslogTag")
	if err != nil {
		log.Error("connect syslog fail, err: ", err)
		return err
	}
	customHook := &CustomHook{hook, syslogLevels}
	//添加syslog钩子
	log.AddHook(customHook)

	return nil
}
