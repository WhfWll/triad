package sqlite

import (
	"os"
	"path/filepath"
	"smart/tools/enums"
	"sync"

	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	scannerDB   *gorm.DB
	scannerOnce sync.Once
	initErr     error
)

// GetScannerDB 获取 scanner.db 数据库连接 (单例)
// 注意：该方法依赖 enums.SystemUpgradeProjectDir 变量，请确保在使用前该变量已正确初始化
func GetScannerDB() (*gorm.DB, error) {
	scannerOnce.Do(func() {
		// 1. Try default path: /opt/laozhi/smart/scanner.db
		// enums.SystemUpgradeProjectDir is typically "/opt/laozhi/"
		dbPath := filepath.Join(enums.SystemUpgradeProjectDir, "smart", "scanner.db")

		// 2. Check if file exists. If not, try current directory and other fallbacks.
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			cwd, _ := os.Getwd()

			// Priority fallback: Check current working directory
			localPath := filepath.Join(cwd, "scanner.db")
			if _, err := os.Stat(localPath); err == nil {
				log.Warnf("Default scanner.db not found at %s, using local file: %s", dbPath, localPath)
				dbPath = localPath
			} else {
				// Secondary fallback: Try one level up (if running in a subdir)
				parentPath := filepath.Join(filepath.Dir(cwd), "scanner.db")
				if _, err := os.Stat(parentPath); err == nil {
					log.Warnf("Default scanner.db not found, using parent dir file: %s", parentPath)
					dbPath = parentPath
				} else {
					// Last resort for local Windows dev: Check project root if we can guess it
					// e.g. D:\goproject\smart\scanner.db
					// This helps if we are running from a subdir in dev environment
					devPath := filepath.Join(cwd, "..", "scanner.db") // Simplified assumption
					if _, err := os.Stat(devPath); err == nil {
						dbPath = devPath
					}
				}
			}
		}

		log.Infof("Initializing scanner.db connection at: %s", dbPath)

		// 使用 glebarez/sqlite (纯Go实现，无需CGO)
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			initErr = err
			log.Errorf("连接 scanner.db 失败: %v", err)
			return
		}

		// 获取通用数据库对象 sql.DB 以设置连接池
		sqlDB, err := db.DB()
		if err != nil {
			initErr = err
			log.Errorf("获取 sql.DB 失败: %v", err)
			return
		}

		// SQLite 仅支持单写，建议限制连接数
		sqlDB.SetMaxOpenConns(1)

		scannerDB = db
		log.Infof("scanner.db 连接初始化成功")
	})

	if initErr != nil {
		return nil, initErr
	}
	return scannerDB, nil
}
