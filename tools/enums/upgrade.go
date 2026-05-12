package enums

const (
	UpgradePackagePath = "upgrade/"
)

const UpgradeServerPlatformToken = "e99e994cc6b05d325af21896c35402dc"

type UpgradeState string

const (
	UpgradeStateIdle            UpgradeState = "IDLE"             // 空闲
	UpgradeStateVerifying       UpgradeState = "VERIFYING"        // 校验中
	UpgradeStateBackingUp       UpgradeState = "BACKINGUP"        // 备份中
	UpgradeStateUpgrading       UpgradeState = "UPGRADING"        // 升级中
	UpgradeStateSuccess         UpgradeState = "SUCCESS"          // 成功
	UpgradeStateFailed          UpgradeState = "FAILED"           // 失败
	UpgradeStateRollbacking     UpgradeState = "ROLLBACKING"      // 回滚中
	UpgradeStateRollbackSuccess UpgradeState = "ROLLBACK_SUCCESS" // 回滚成功
	UpgradeStateRollbackFailed  UpgradeState = "ROLLBACK_FAILED"  // 回滚失败
)
