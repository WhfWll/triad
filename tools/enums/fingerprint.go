package enums

// 5 应用层 表示目标资产具体运行的业务应用，例如:
// 4 支撑层 表示目标资产的开发语言/框架
// 3 服务层 表示目标资产运行的服务
// 2 系统层 表示目标资产运行的系统
// 1 硬件层 表示目标资产运行的硬件设备或安防设备
// 9 默认为9 表示未分层指纹

const (
	FingerLevelApp      = 5
	FingerLevelSupport  = 4
	FingerLevelService  = 3
	FingerLevelSystem   = 2
	FingerLevelHardware = 1
	FingerLevelDefault  = 9
)
