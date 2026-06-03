package crons

import (
	"github.com/robfig/cron"
)

func Start() {
	c := cron.New()
	//c.AddFunc("5 * * * * *", TargetCreate) //任务创建,每秒执行一次
	c.AddFunc("@every 10s", TargetCreate)             //任务发送,每10秒执行一次
	c.AddFunc("@every 10s", TaskCycle)                // 定时或周期任务执行,每10秒执行一次
	c.AddFunc("@every 1m", TaskRunningPeriod)         // 任务运行时间段，是否允许运行
	c.AddFunc("0 0 0 * * ?", CheckAccountValidPeriod) // 检查用户账号是否过期
	c.AddFunc("@every 30s", TaskInfoStat)             // 每30秒执行下任务统计
	c.AddFunc("@every 3m", TargetTimeout)             // 每10分钟执行下检测目标超时
	c.AddFunc("@every 10s", ReportGenerate)           // 每10s执行下报告生成
	//c.AddFunc("@every 10s", TripartiteXray)            // 每10s执行下三方工具 xray
	//c.AddFunc("@every 10s", TripartiteBurpsuite)       // 每10s执行下三方工具 burpsuite
	c.AddFunc("@every 1m", CheckSystemAuth)        // 每10m执行系统授权检查（非 Linux 会标记未授权）
	c.AddFunc("@every 1m", LogBackupCron)          // 每1分钟判断一次是否执行日志备份
	c.AddFunc("@every 1m", SystemConfigBackupCron) // 每1分钟判断一次是否执行系统配置备份
	c.AddFunc("@every 2m", SystemMonitor)          // 每2m执行一次系统信息统计
	//c.AddFunc("@every 10s", BasExec)                   // 每10s执行一次bas任务
	//c.AddFunc("@every 1m", BasNodeOnlineStatus)        // 每1m执行一次更新Bas节点的在线状态
	c.AddFunc("@every 10s", TaskReportVerify)          // 每10s执行一次报告验证任务
	c.AddFunc("0 0 2 * * ?", DelExpirationTimeLog)     // 每天2点删除过期的日志
	c.AddFunc("0 0 2 * * ?", CheckPasswordValidPeriod) // 每天2点进行密码过期检查
	// c.AddFunc("@every 12s", VulScanReportUpdate)       // 每12s执行下报告生成
	c.AddFunc("@every 10s", scannerNodeIsAlive) //检测分布式节点是否存活,每10秒执行一次

	c.AddFunc("@every 60m", SystemRunMonitor) // 每60m执行一次系统运行监控
	//c.AddFunc("@every 30s", MonitorDBConnection) // 每30s执行一次数据库连接监控

	c.AddFunc("@every 10m", RemoteSessionStatusCheck) // remote session 状态检测

	//c.AddFunc("@every 25s", TaskGroupInfoStat)         // 每天进行任务组信息统计
	//c.AddFunc("@every 1m", AssetSync)                  // 资产同步
	//c.AddFunc("@every 8s", LogicTaskExec) // 每隔8s进行一次逻辑漏洞任务检测

	c.Start()
	go CheckSystemAuth() // 启动时立即校验（非 Linux 会标记未授权）
}
