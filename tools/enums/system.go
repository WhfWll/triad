package enums

const SystemNoAuth = "系统未授权 请前往授权"

var (
	SystemUpgradeProjectDir = "/opt/laozhi/" //项目目录
	DockerComposeConfigPath = "/opt/laozhi/docker-compose.yml"
)

// InitSystemPaths 初始化系统路径
func InitSystemPaths(baseDir string) {
	if baseDir != "" {
		SystemUpgradeProjectDir = baseDir
		lastChar := SystemUpgradeProjectDir[len(SystemUpgradeProjectDir)-1]
		if lastChar != '/' && lastChar != '\\' {
			SystemUpgradeProjectDir += "/"
		}
		DockerComposeConfigPath = SystemUpgradeProjectDir + "docker-compose.yml"
	}
}

const (
	TargetIpWhiteBlackIsOpenOn  = 1
	TargetIpWhiteBlackIsOpenOff = 2
	TargetIpWhiteBlackTypeWhite = 1
	TargetIpWhiteBlackTypeBlack = 2
)

// SystemConfigBackupDir 系统配置备份目录
const SystemConfigBackupDir = "backup/systems/"

// 升级还原路径部分
const (
	SystemBaseVersion               = "V0.4.0.230921R"   //系统基准版本号
	SystemUpgradeFileDir            = "upgrade/"         //系统升级zip文件所在目录
	SystemUpgradeZipFileDir         = "upgrade/zip"      //系统升级zip文件所在目录
	SystemUpgradeUnZipFileDir       = "upgrade/unzip/"   //系统升级unzip文件所在目录
	SystemUpgradeWorkDir            = "upgrade/work/"    //系统升级收集打包文件所在目录
	SystemUpgradePassword           = "admin123321"      //系统离线更新包密码
	SystemUpgradeScriptFilename     = "update_static.sh" //系统离线更新中的更新脚本
	SystemUpgradeDataDir            = "data"             //项目目录-mysql目录名
	SystemUpgradeNginxDir           = "nginx"            //项目目录-nginx目录名
	SystemUpgradeVueDir             = "smart_vue"        //项目目录-nginx目录名-vue目录名（unzip目录-vue目录名）
	SystemUpgradeSmartDir           = "smart"            //项目目录-smart目录名（unzip目录-smart目录名）
	SystemUpgradeDecisionDir        = "decision"         //项目目录-decision目录名（unzip目录-decision目录名）
	SystemUpgradeSmartFile          = "smart"            //项目目录-smart目录名-smart文件（unzip目录-smart目录名-smart文件）
	SystemUpgradeSmartConfigFile    = "config.json"      //项目目录-smart目录名-smart配置文件（unzip目录-smart目录名-smart配置文件）
	SystemUpgradeDecisionFile       = "decision"         //项目目录-smart目录名-decision文件（unzip目录-smart目录名-decision文件）
	SystemUpgradeDecisionConfigFile = "config.json"      //项目目录-smart目录名-decision配置文件（unzip目录-smart目录名-decision配置文件）
	SystemUpgradeBinDir             = "bin"              //项目目录-bin目录名
	SystemUpgradeWebDir             = "web"              //升级包-web目录名
	SystemUpgradeScannerFile        = "scanner"          //项目目录-smart目录名-scanner文件
	SystemUpgradeDbWebFile          = "dbweb"            //项目目录-bin目录名-dbweb文件
	SystemUpgradeYakFile            = "yak"              //项目目录-yak文件（unzip目录-yak文件）
	SystemUpgradeMd5Filename        = "data.md5"         //unzip目录-md5文件
	SystemUpgradeAssetsDir          = "assets"           //升级包-assets目录名
	SystemUpgradeBackupMainDir      = "smart/backup"     //升级备份主目录
	ServiceNameSmart                = "smart"            //Smart服务名
	ServiceNameDecision             = "decision"         //Decision服务名
)

const (
	UbuntuNetworkConfigDir = "/etc/netplan"
	CentosNetworkConfigDir = "/etc/sysconfig/network-scripts"
)

// 系统管理中的各配置开启/关闭按钮
const (
	ConfigOpen  = 1
	ConfigClose = 2
)

// syslog日志类型
const (
	SyslogTypeAudit = "1"
	SyslogTypeDebug = "2"
	SyslogTypeWarn  = "3"
	SyslogTypeError = "4"
)

// 升级还原状态
const (
	Uploading    = "上传中"
	Uploaded     = "已上传"
	UploadFail   = "上传失败"
	Downloading  = "下载中"
	Downloaded   = "已下载"
	DownloadFail = "下载失败"
)

// 消息中心
const (
	MessageTypeNotice = 1 //消息类型 - 通知
	MessageTypeWarn   = 2 //消息类型 - 警告
	MessageTypeError  = 3 //消息类型 - 异常
)

const (
	MessageStatusRead   = 1 //已读
	MessageStatusUnread = 2 //未读
)

const (
	MonitorCpu    = "cpu占用率"
	MonitorMemory = "内存占用率"
	MonitorDisk   = "磁盘占用率"
	MonitorFlow   = "网络输入"
)

// MessageTypeMap 系统消息类型map
var MessageTypeMap = map[int]string{
	MessageTypeNotice: "通知",
	MessageTypeWarn:   "警告",
	MessageTypeError:  "异常",
}

// MessageStatusMap 系统消息状态map
var MessageStatusMap = map[int]string{
	MessageStatusRead:   "已读",
	MessageStatusUnread: "未读",
}
