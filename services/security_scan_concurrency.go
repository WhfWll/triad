package services

import (
	"context"
	"encoding/json"
	"sync"

	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
)

// SecurityScanConcurrencyConfig 三类安全检查单任务内目标并发上限
type SecurityScanConcurrencyConfig struct {
	HostConcurrent int `json:"hostConcurrent"`
	AppConcurrent  int `json:"appConcurrent"`
	DataConcurrent int `json:"dataConcurrent"`
}

const (
	DefaultHostScanConcurrent = 10
	DefaultAppScanConcurrent  = 5
	DefaultDataScanConcurrent = 5
	minSecurityScanConcurrent = 1
	maxSecurityScanConcurrent = 50
)

var (
	secScanConcMu    sync.RWMutex
	secScanConcCache SecurityScanConcurrencyConfig
	secScanConcReady bool
)

func defaultSecurityScanConcurrency() SecurityScanConcurrencyConfig {
	return SecurityScanConcurrencyConfig{
		HostConcurrent: DefaultHostScanConcurrent,
		AppConcurrent:  DefaultAppScanConcurrent,
		DataConcurrent: DefaultDataScanConcurrent,
	}
}

func normalizeSecurityScanConcurrency(cfg SecurityScanConcurrencyConfig) SecurityScanConcurrencyConfig {
	def := defaultSecurityScanConcurrency()
	if cfg.HostConcurrent <= 0 {
		cfg.HostConcurrent = def.HostConcurrent
	}
	if cfg.AppConcurrent <= 0 {
		cfg.AppConcurrent = def.AppConcurrent
	}
	if cfg.DataConcurrent <= 0 {
		cfg.DataConcurrent = def.DataConcurrent
	}
	if cfg.HostConcurrent < minSecurityScanConcurrent {
		cfg.HostConcurrent = minSecurityScanConcurrent
	}
	if cfg.AppConcurrent < minSecurityScanConcurrent {
		cfg.AppConcurrent = minSecurityScanConcurrent
	}
	if cfg.DataConcurrent < minSecurityScanConcurrent {
		cfg.DataConcurrent = minSecurityScanConcurrent
	}
	if cfg.HostConcurrent > maxSecurityScanConcurrent {
		cfg.HostConcurrent = maxSecurityScanConcurrent
	}
	if cfg.AppConcurrent > maxSecurityScanConcurrent {
		cfg.AppConcurrent = maxSecurityScanConcurrent
	}
	if cfg.DataConcurrent > maxSecurityScanConcurrent {
		cfg.DataConcurrent = maxSecurityScanConcurrent
	}
	return cfg
}

// GetSecurityScanConcurrency 读取安全检查并发配置（带内存缓存）
func GetSecurityScanConcurrency(ctx context.Context) SecurityScanConcurrencyConfig {
	secScanConcMu.RLock()
	if secScanConcReady {
		cfg := secScanConcCache
		secScanConcMu.RUnlock()
		return cfg
	}
	secScanConcMu.RUnlock()

	cfg := loadSecurityScanConcurrencyFromDB(ctx)
	secScanConcMu.Lock()
	secScanConcCache = cfg
	secScanConcReady = true
	secScanConcMu.Unlock()
	return cfg
}

func loadSecurityScanConcurrencyFromDB(ctx context.Context) SecurityScanConcurrencyConfig {
	var mapSet MapSet
	raw, err := mapSet.GetMapValue(ctx, enums.SecurityScanConcurrencyMapSetObjKey)
	if err != nil || raw == "" {
		return defaultSecurityScanConcurrency()
	}
	var cfg SecurityScanConcurrencyConfig
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		log.Warnf("security scan concurrency config invalid: %v", err)
		return defaultSecurityScanConcurrency()
	}
	return normalizeSecurityScanConcurrency(cfg)
}

// SaveSecurityScanConcurrency 保存并刷新缓存
func SaveSecurityScanConcurrency(ctx context.Context, cfg SecurityScanConcurrencyConfig) error {
	cfg = normalizeSecurityScanConcurrency(cfg)
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	var mapSet MapSet
	if err := mapSet.Create(ctx, enums.SecurityScanConcurrencyMapSetObjKey, string(b), enums.SecurityScanConcurrencyMapSetContent); err != nil {
		return err
	}
	secScanConcMu.Lock()
	secScanConcCache = cfg
	secScanConcReady = true
	secScanConcMu.Unlock()
	return nil
}

func GetHostScanConcurrent(ctx context.Context) int {
	return GetSecurityScanConcurrency(ctx).HostConcurrent
}

func GetAppScanConcurrent(ctx context.Context) int {
	return GetSecurityScanConcurrency(ctx).AppConcurrent
}

func GetDataScanConcurrent(ctx context.Context) int {
	return GetSecurityScanConcurrency(ctx).DataConcurrent
}
