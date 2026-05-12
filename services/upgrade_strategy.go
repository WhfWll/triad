package services

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/file"
	"smart/tools/utils"
	"sort"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ==========================================
// 1. 数据结构定义 (Manifest)
// ==========================================

type UpgradeType string

const (
	UpgradeTypeSystem UpgradeType = "SYSTEM" // 系统升级
	UpgradeTypeVuln   UpgradeType = "VULN"   // 漏洞/脚本库升级
)

// VulnScope 定义漏洞升级范围
type VulnScope string

const (
	VulnScopePrinciple VulnScope = "PRINCIPLE" // 原理验证数据 (高频)
	VulnScopeFull      VulnScope = "FULL"      // 全量数据 (低频)
	VulnScopeScript    VulnScope = "SCRIPT"    // 仅脚本库 (SQLite)
	VulnScopeMixed     VulnScope = "MIXED"     // 混合模式
)

// UpgradeManifest 定义升级包内的 manifest.json 结构
type UpgradeManifest struct {
	Type        UpgradeType       `json:"type"`               // 升级包类型
	VulnScope   VulnScope         `json:"vuln_scope"`         // 漏洞升级范围
	Version     string            `json:"version"`            // 目标版本
	BuildTime   int64             `json:"build_time"`         // 构建时间
	Description string            `json:"description"`        // 描述
	MinVersion  string            `json:"min_system_version"` // 依赖的最低系统版本
	FileHash    map[string]string `json:"file_hash"`          // 文件指纹校验
	NeedRestart bool              `json:"need_restart"`       // 是否需要重启服务
}

// UpgradeResult 返回给前端的升级结果
type UpgradeResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// ==========================================
// 2. 接口定义 (Strategy Interface)
// ==========================================

type UpgradeStrategy interface {
	Verify(ctx context.Context, manifest *UpgradeManifest, unzipDir string) error
	Backup(ctx context.Context) (string, error)
	Apply(ctx context.Context, unzipDir string) error
	Rollback(ctx context.Context, backupDir string) error
	PostAction(ctx context.Context) error
}

// ==========================================
// Upgrade Status Management
// ==========================================

type UpgradeStatus struct {
	State     enums.UpgradeState `json:"state"`
	Message   string             `json:"message"`
	Progress  int                `json:"progress"` // 0-100
	Error     string             `json:"error,omitempty"`
	Version   string             `json:"version,omitempty"` // 目标版本
	StartTime time.Time          `json:"start_time"`
	EndTime   time.Time          `json:"end_time"`
	Type      UpgradeType        `json:"type,omitempty"`
	BackupDir string             `json:"backup_dir,omitempty"`
}

var (
	globalStatus      UpgradeStatus
	globalStatusMutex sync.RWMutex
)

// InitUpgradeStatus 初始化状态
func InitUpgradeStatus() {
	globalStatusMutex.Lock()
	defer globalStatusMutex.Unlock()
	globalStatus = UpgradeStatus{
		State:     enums.UpgradeStateVerifying, // 初始状态
		StartTime: time.Now(),
		Progress:  0,
		Message:   "开始处理...",
	}
}

// UpdateUpgradeStatus 更新状态
func UpdateUpgradeStatus(state enums.UpgradeState, message string, progress int) {
	globalStatusMutex.Lock()
	defer globalStatusMutex.Unlock()
	globalStatus.State = state
	globalStatus.Message = message
	if progress >= 0 {
		globalStatus.Progress = progress
	}
	if state == enums.UpgradeStateSuccess || state == enums.UpgradeStateFailed {
		globalStatus.EndTime = time.Now()
	}
}

// SetUpgradeVersion 设置目标版本
func SetUpgradeVersion(version string) {
	globalStatusMutex.Lock()
	defer globalStatusMutex.Unlock()
	globalStatus.Version = version
}

// SetUpgradeType 设置升级类型
func SetUpgradeType(t UpgradeType) {
	globalStatusMutex.Lock()
	defer globalStatusMutex.Unlock()
	globalStatus.Type = t
}

// SetUpgradeBackupDir 设置备份目录
func SetUpgradeBackupDir(dir string) {
	globalStatusMutex.Lock()
	defer globalStatusMutex.Unlock()
	globalStatus.BackupDir = dir
}

// SetUpgradeError 设置错误状态
func SetUpgradeError(err error) {
	globalStatusMutex.Lock()
	defer globalStatusMutex.Unlock()
	globalStatus.State = enums.UpgradeStateFailed
	globalStatus.Message = "升级失败"
	globalStatus.Error = err.Error()
	globalStatus.EndTime = time.Now()
}

// SetRollbackError 设置回滚错误状态
func SetRollbackError(err error) {
	globalStatusMutex.Lock()
	defer globalStatusMutex.Unlock()
	globalStatus.State = enums.UpgradeStateRollbackFailed
	globalStatus.Message = "回滚失败"
	globalStatus.Error = err.Error()
	globalStatus.EndTime = time.Now()
}

// GetUpgradeStatus 获取当前状态
func GetUpgradeStatus() UpgradeStatus {
	globalStatusMutex.RLock()
	defer globalStatusMutex.RUnlock()
	return globalStatus
}

// IsUpgrading 检查是否正在升级中
func IsUpgrading() bool {
	globalStatusMutex.RLock()
	defer globalStatusMutex.RUnlock()
	return globalStatus.State == enums.UpgradeStateVerifying ||
		globalStatus.State == enums.UpgradeStateBackingUp ||
		globalStatus.State == enums.UpgradeStateUpgrading
}

// ==========================================
// 3. 升级管理器 (Manager)
// ==========================================

// UpgradeManager 升级管理器
type UpgradeManager struct {
}

func NewUpgradeManager() *UpgradeManager {
	return &UpgradeManager{}
}

// ProcessUpgradeAsync 异步处理升级
func (m *UpgradeManager) ProcessUpgradeAsync(ctx context.Context, zipPath string) error {
	if IsUpgrading() {
		return fmt.Errorf("当前已有升级任务正在进行中")
	}

	InitUpgradeStatus()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("升级过程发生 Panic: %v", r)
				SetUpgradeError(fmt.Errorf("内部错误: %v", r))
			}
		}()

		// 1. 解压 Zip 到临时目录
		UpdateUpgradeStatus(enums.UpgradeStateVerifying, "正在解压升级包...", 5)
		unzipDir := strings.TrimSuffix(zipPath, filepath.Ext(zipPath)) + "_unzip"
		defer os.RemoveAll(unzipDir) // 确保清理

		if err := file.Unzip(zipPath, unzipDir, ""); err != nil {
			SetUpgradeError(fmt.Errorf("解压升级包失败: %v", err))
			return
		}

		// 2. 读取并解析 Manifest
		manifest, err := m.loadManifest(unzipDir)
		if err != nil {
			SetUpgradeError(fmt.Errorf("读取元数据失败: %v", err))
			return
		}
		SetUpgradeVersion(manifest.Version)
		SetUpgradeType(manifest.Type)

		// 3. 校验签名
		if err := m.verifySignature(unzipDir); err != nil {
			log.Warnf("Signature verify failed: %v", err)
			// return nil, fmt.Errorf("签名校验失败: %v", err) // 生产环境应开启
		}

		// 4. 选择策略
		var strategy UpgradeStrategy
		switch manifest.Type {
		case UpgradeTypeSystem:
			strategy = &SystemUpgradeStrategy{Manifest: manifest}
		case UpgradeTypeVuln:
			strategy = &VulnUpgradeStrategy{Manifest: manifest}
		default:
			SetUpgradeError(fmt.Errorf("未知的升级包类型: %s", manifest.Type))
			return
		}

		// 5. 执行升级原子流程
		if err := m.runAtomicUpgrade(context.Background(), strategy, manifest, unzipDir); err != nil {
			SetUpgradeError(err)
			return
		}

		UpdateUpgradeStatus(enums.UpgradeStateSuccess, "升级成功", 100)
	}()

	return nil
}

// ProcessUpgrade 处理升级请求的核心入口 (保留用于兼容，建议使用 Prepare + Execute)
func (m *UpgradeManager) ProcessUpgrade(ctx context.Context, zipPath string) (*UpgradeResult, error) {
	manifest, err := m.PrepareUpgrade(ctx, zipPath)
	if err != nil {
		return nil, err
	}

	if err := m.ExecuteUpgrade(ctx, zipPath); err != nil {
		return nil, err
	}

	return &UpgradeResult{
		Success: true,
		Message: "升级任务已提交",
		Version: manifest.Version,
	}, nil
}

// PrepareUpgrade 准备升级：仅解压元数据文件、校验、返回元数据
func (m *UpgradeManager) PrepareUpgrade(ctx context.Context, zipPath string) (*UpgradeManifest, error) {
	// 1. 准备解压目录
	unzipDir := filepath.Join(filepath.Dir(zipPath), "extract")
	// 注意：这里不再清空目录，因为可能是增量解压，或者Execute阶段会全量解压覆盖
	// 但为了避免旧文件干扰，最好还是清空
	os.RemoveAll(unzipDir)
	if err := file.CreateDir(unzipDir); err != nil {
		return nil, fmt.Errorf("创建解压目录失败: %v", err)
	}

	// 2. 仅解压 manifest.json 和 signature.bin (性能优化)
	// 使用 archive/zip 标准库手动解压特定文件
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("打开升级包失败: %v", err)
	}
	defer r.Close()

	var manifestData, signatureData []byte
	foundManifest := false

	for _, f := range r.File {
		if f.Name == "manifest.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			manifestData, err = ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
			// 写入磁盘供 loadManifest 使用
			if err := ioutil.WriteFile(filepath.Join(unzipDir, "manifest.json"), manifestData, 0644); err != nil {
				return nil, err
			}
			foundManifest = true
		} else if f.Name == "signature.bin" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			signatureData, err = ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
			// 写入磁盘供 verifySignature 使用
			if err := ioutil.WriteFile(filepath.Join(unzipDir, "signature.bin"), signatureData, 0644); err != nil {
				return nil, err
			}
		}
	}

	if !foundManifest {
		return nil, errors.New("升级包中缺少 manifest.json")
	}

	// 3. 读取并解析 Manifest
	manifest, err := m.loadManifest(unzipDir)
	if err != nil {
		return nil, fmt.Errorf("读取元数据失败: %v", err)
	}

	// 4. 校验签名 (只校验 manifest.json 的签名)
	// 注意：这里只校验了 manifest.json 的完整性，没有校验整个包
	// 完整的包校验推迟到 ExecuteUpgrade 阶段
	if err := m.verifySignature(unzipDir); err != nil {
		log.Warnf("Signature verify failed: %v", err)
		return nil, fmt.Errorf("签名校验失败: %v", err) // 生产环境应开启
	}

	return manifest, nil
}

// ExecuteUpgrade 执行已准备好的升级
func (m *UpgradeManager) ExecuteUpgrade(ctx context.Context, zipPath string) error {
	unzipDir := filepath.Join(filepath.Dir(zipPath), "extract")

	// 1. 再次读取 Manifest (确保一致性，且 Prepare 阶段已解压)
	manifest, err := m.loadManifest(unzipDir)
	if err != nil {
		// 如果 Prepare 阶段没解压成功，这里会失败，或者可以尝试重新解压 manifest
		// 理论上 Prepare 必须先调用成功
		return fmt.Errorf("读取元数据失败，请重新上传升级包: %v", err)
	}

	// 2. 选择策略
	var strategy UpgradeStrategy
	switch manifest.Type {
	case UpgradeTypeSystem:
		strategy = &SystemUpgradeStrategy{Manifest: manifest}
	case UpgradeTypeVuln:
		strategy = &VulnUpgradeStrategy{Manifest: manifest}
	default:
		return fmt.Errorf("未知的升级包类型: %s", manifest.Type)
	}

	// 3. 异步执行
	go func() {
		defer func() {
			if r := recover(); r != nil {
				SetUpgradeError(fmt.Errorf("Panic during upgrade: %v", r))
			}
		}()

		// 初始化升级状态
		InitUpgradeStatus()

		// 3.1 异步全量解压 (性能优化点：移到这里执行)
		UpdateUpgradeStatus(enums.UpgradeStateVerifying, "正在解压升级包...", 5)
		// 清理之前只解压了部分文件的目录，或者直接覆盖
		// 建议直接覆盖解压
		if err := file.Unzip(zipPath, unzipDir, ""); err != nil {
			SetUpgradeError(fmt.Errorf("解压升级包失败: %v", err))
			return
		}

		// 3.2 执行原子流程
		if err := m.runAtomicUpgrade(context.Background(), strategy, manifest, unzipDir); err != nil {
			log.Errorf("Async upgrade failed: %v", err)
			// 显式更新失败状态，避免卡在 50%
			SetUpgradeError(err)
		} else {
			UpdateUpgradeStatus(enums.UpgradeStateSuccess, "升级成功", 100)
		}
	}()

	return nil
}

func (m *UpgradeManager) runAtomicUpgrade(ctx context.Context, strategy UpgradeStrategy, manifest *UpgradeManifest, unzipDir string) error {
	// Step A: 业务校验
	UpdateUpgradeStatus(enums.UpgradeStateVerifying, "正在校验升级包内容...", 10)
	if err := strategy.Verify(ctx, manifest, unzipDir); err != nil {
		return fmt.Errorf("校验失败: %v", err)
	}

	// Step B: 备份
	UpdateUpgradeStatus(enums.UpgradeStateBackingUp, "正在备份数据...", 30)
	backupDir, err := strategy.Backup(ctx)
	if err != nil {
		return fmt.Errorf("备份失败: %v", err)
	}
	SetUpgradeBackupDir(backupDir)

	// Step C: 应用更新
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在应用更新...", 50)
	if err := strategy.Apply(ctx, unzipDir); err != nil {
		// 升级失败，记录日志并返回错误，但不自动回滚
		log.Errorf("升级失败: %v", err)
		return fmt.Errorf("升级失败，请手动恢复: %v", err)
	}

	// Step D: 后置操作
	go func() {
		// 异步执行后置操作 (如重启)
		// 增加延时，确保前端有足够时间(例如 10s)轮询到 100% 成功状态
		// 否则服务重启会导致内存中的状态丢失，前端会看到进度归零
		time.Sleep(10 * time.Second)
		if err := strategy.PostAction(context.Background()); err != nil {
			log.Errorf("后置操作失败: %v", err)
		}
	}()

	return nil
}

func (m *UpgradeManager) loadManifest(dir string) (*UpgradeManifest, error) {
	data, err := ioutil.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return nil, err
	}
	var manifest UpgradeManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (m *UpgradeManager) verifySignature(dir string) error {
	sigPath := filepath.Join(dir, "signature.bin")
	signature, err := ioutil.ReadFile(sigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("缺少签名文件 signature.bin")
		}
		return err
	}

	manifestPath := filepath.Join(dir, "manifest.json")
	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(manifestData)
	block, _ := pem.Decode([]byte(enums.SWPubKey))
	if block == nil {
		return errors.New("无法解析公钥 PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("非 RSA 公钥")
	}

	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
}

// ==========================================
// 4. 具体策略: 漏洞库升级 (VulnUpgradeStrategy)
// ==========================================

func getScannerDBPath() string {
	return filepath.Join(enums.SystemUpgradeProjectDir, "smart/scanner.db")
}

func getUserCustomDBPath() string {
	return filepath.Join(enums.SystemUpgradeProjectDir, "smart/user_custom.db")
}

type VulnUpgradeStrategy struct {
	Manifest *UpgradeManifest
}

func (s *VulnUpgradeStrategy) Verify(ctx context.Context, m *UpgradeManifest, dir string) error {
	jsonPath := filepath.Join(dir, "assets", "db", "vul_libraries.json")
	dbPath := filepath.Join(dir, "assets", "db", "scanner.db")

	hasJson := file.CheckPathExist(jsonPath)
	hasDb := file.CheckPathExist(dbPath)

	if !hasJson && !hasDb {
		return errors.New("升级包结构错误: 缺少 assets/db/vul_libraries.json 或 assets/db/scanner.db")
	}

	if m.VulnScope == VulnScopePrinciple && !hasJson {
		return errors.New("PRINCIPLE 模式需要 vul_libraries.json")
	}

	return nil
}

func (s *VulnUpgradeStrategy) Backup(ctx context.Context) (string, error) {
	backupDir := filepath.Join(enums.SystemUpgradeProjectDir, "smart/backup", fmt.Sprintf("vuln_%d", time.Now().Unix()))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	// 备份 manifest 信息，便于回滚时判断 Scope
	if s.Manifest != nil {
		if data, err := json.Marshal(s.Manifest); err == nil {
			ioutil.WriteFile(filepath.Join(backupDir, "manifest_backup.json"), data, 0644)
		}
	}

	if file.CheckPathExist(getScannerDBPath()) {
		if err := file.CopyFile(getScannerDBPath(), filepath.Join(backupDir, "scanner.db")); err != nil {
			return "", err
		}
	}
	if file.CheckPathExist(getUserCustomDBPath()) {
		if err := file.CopyFile(getUserCustomDBPath(), filepath.Join(backupDir, "user_custom.db")); err != nil {
			return "", err
		}
	}

	// 备份版本信息
	var mapSetService MapSet
	if val, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey); err == nil {
		ioutil.WriteFile(filepath.Join(backupDir, "version.json"), []byte(val), 0644)
	}

	db := mysql.FromContext(ctx)

	// 分批备份 vul_libraries 表
	vulnBackupDir := filepath.Join(backupDir, "vul_libraries")
	if err := os.MkdirAll(vulnBackupDir, 0755); err != nil {
		return "", err
	}

	batchSize := 10000
	var totalCount int64
	query := db.Model(&mysqls.VulLibraries{})
	if s.Manifest.VulnScope == VulnScopePrinciple {
		query = query.Where("verify_type = ?", enums.VulVerifyTypePrincipleVerification)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return "", err
	}

	for offset := int64(0); offset < totalCount; offset += int64(batchSize) {
		var libs []mysqls.VulLibraries
		if err := query.Limit(batchSize).Offset(int(offset)).Find(&libs).Error; err != nil {
			return "", err
		}

		data, err := json.Marshal(libs)
		if err != nil {
			return "", err
		}

		fileName := fmt.Sprintf("part_%d.json", offset/int64(batchSize))
		if err := ioutil.WriteFile(filepath.Join(vulnBackupDir, fileName), data, 0644); err != nil {
			return "", err
		}
	}

	// 备份 dictionary (数据量较小，保持单文件)
	var dicts []mysqls.Dictionary
	if err := db.Model(&mysqls.Dictionary{}).Find(&dicts).Error; err == nil {
		if data, err := json.Marshal(dicts); err == nil {
			if err := ioutil.WriteFile(filepath.Join(backupDir, "dictionary.json"), data, 0644); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	var fingers []mysqls.Finger
	if err := db.Model(&mysqls.Finger{}).Find(&fingers).Error; err == nil {
		if data, err := json.Marshal(fingers); err == nil {
			if err := ioutil.WriteFile(filepath.Join(backupDir, "finger.json"), data, 0644); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return backupDir, nil
}

func (s *VulnUpgradeStrategy) Apply(ctx context.Context, dir string) error {
	jsonPath := filepath.Join(dir, "assets", "db", "vul_libraries.json")
	dictPath := filepath.Join(dir, "assets", "db", "dictionary.json")
	fingerPath := filepath.Join(dir, "assets", "db", "finger.json")
	dbPath := filepath.Join(dir, "assets", "db", "scanner.db")

	// 1. 应用 MySQL 漏洞数据更新
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在更新漏洞库数据...", 50)
	if file.CheckPathExist(jsonPath) {
		if s.Manifest.VulnScope == VulnScopePrinciple {
			if err := s.applyPrincipleUpdate(ctx, jsonPath); err != nil {
				return err
			}
		} else {
			if err := s.applyFullUpdate(ctx, jsonPath); err != nil {
				return err
			}
		}
	}
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "漏洞库数据更新完成", 65)

	// 2. 应用 字典表 更新
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在更新字典表...", 65)
	if file.CheckPathExist(dictPath) {
		if err := s.applyDictionaryUpdate(ctx, dictPath); err != nil {
			return err
		}
	}
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "字典表更新完成", 75)

	// 3. 应用 指纹表 更新
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在更新指纹表...", 75)
	if file.CheckPathExist(fingerPath) {
		if err := s.applyFingerUpdate(ctx, fingerPath); err != nil {
			return err
		}
	}
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "指纹表更新完成", 85)

	// 4. 应用 SQLite 脚本库更新
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在更新脚本库...", 85)
	if file.CheckPathExist(dbPath) {
		if err := s.applyScriptDB(ctx, dbPath); err != nil {
			return err
		}
	}
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "脚本库更新完成", 95)

	// 5. 更新系统版本信息 (VulVersion)
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在更新版本信息...", 95)
	var mapSetService MapSet
	var systemVersionMapSet SystemVersionMapSet

	// 获取现有版本信息
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if err == nil && objValueStr != "" {
		json.Unmarshal([]byte(objValueStr), &systemVersionMapSet)
	} else {
		// defaults if not found
		systemVersionMapSet.CurrentVersion = enums.SystemBaseVersion
	}

	// 更新字段
	// 先记录上次工具库版本
	systemVersionMapSet.LastVulVersion = systemVersionMapSet.VulVersion
	systemVersionMapSet.VulVersion = s.Manifest.Version
	systemVersionMapSet.VulUpdateTime = time.Now().Format(enums.TimeLayout)

	// 序列化
	objValueByte, err := json.Marshal(systemVersionMapSet)
	if err != nil {
		log.Errorf("序列化版本信息失败: %v", err)
		return fmt.Errorf("序列化版本信息失败: %v", err)
	}

	// 保存
	if err = mapSetService.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent); err != nil {
		log.Errorf("更新版本信息失败: %v", err)
		return fmt.Errorf("更新版本信息失败: %v", err)
	}

	log.Infof("漏洞库版本更新成功: %s", s.Manifest.Version)

	return nil
}

// applyPrincipleUpdate 增量更新原理验证数据 (Memory Diff + Batch)
func (s *VulnUpgradeStrategy) applyPrincipleUpdate(ctx context.Context, jsonPath string) error {
	log.Info("开始执行漏洞库 PRINCIPLE 增量更新...")

	// A. 读取新数据
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	// 辅助函数: 解析时间字符串 (支持 MySQL 默认格式和 RFC3339)
	parseTime := func(tStr string) (time.Time, error) {
		if tStr == "" {
			return time.Time{}, nil
		}
		layouts := []string{
			"2006-01-02 15:04:05", // MySQL default
			time.RFC3339,          // JSON default
		}
		for _, layout := range layouts {
			if t, err := time.ParseInLocation(layout, tStr, time.Local); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("无法解析时间: %s", tStr)
	}

	// 定义导入结构体，匹配下划线风格 JSON
	// 针对 vul_libraries 表的字段映射
	type VulLibrariesImport struct {
		ID              int    `json:"id"`
		DataType        int    `json:"data_type"`
		VulID           string `json:"vul_id"`
		Name            string `json:"name"`
		Cve             string `json:"cve"`
		Risk            int    `json:"risk"`
		Type            int    `json:"type"`
		Class           int    `json:"class"`
		PublishedTime   string `json:"published_time"`
		Description     string `json:"description"`
		AffectRange     string `json:"affect_range"`
		ExploitImpact   int    `json:"exploit_impact"`
		FixSuggest      string `json:"fix_suggest"`
		Cnvd            string `json:"cnvd"`
		Cnnvd           string `json:"cnnvd"`
		Cncve           string `json:"cncve"`
		Bugtraq         string `json:"bugtraq"`
		Component       string `json:"component"`
		Status          int    `json:"status"`
		StatusMsg       string `json:"status_msg"`
		Priority        string `json:"priority"`
		OperatingSystem int    `json:"operating_system"`
		CreateTime      string `json:"create_time"`
		UpdateTime      string `json:"update_time"`
		Pocname         string `json:"pocname"`
		VerifyType      int    `json:"verify_type"`
		PatchUrl        string `json:"patch_url"`
		CvssVersion     string `json:"cvss_version"`
		CvssScore       string `json:"cvss_score"`
		PocOrExp        string `json:"poc_or_exp"`
		ScriptType      string `json:"script_type"`
	}

	var rawVulns []VulLibrariesImport
	if err := json.Unmarshal(data, &rawVulns); err != nil {
		return fmt.Errorf("解析漏洞数据失败: %v", err)
	}

	// 转换为模型对象
	var newVulns []mysqls.VulLibraries
	for _, raw := range rawVulns {
		ct, _ := parseTime(raw.CreateTime)
		ut, _ := parseTime(raw.UpdateTime)

		// 默认值处理
		dataType := raw.DataType
		if dataType == 0 {
			dataType = 1 // 默认为真实数据
		}

		newVulns = append(newVulns, mysqls.VulLibraries{
			ID:              raw.ID,
			DataType:        dataType,
			VulID:           raw.VulID,
			Name:            raw.Name,
			Cve:             raw.Cve,
			Risk:            raw.Risk,
			Type:            raw.Type,
			Class:           raw.Class,
			PublishedTime:   raw.PublishedTime,
			Description:     raw.Description,
			AffectRange:     raw.AffectRange,
			ExploitImpact:   raw.ExploitImpact,
			FixSuggest:      raw.FixSuggest,
			Cnvd:            raw.Cnvd,
			Cnnvd:           raw.Cnnvd,
			Cncve:           raw.Cncve,
			Bugtraq:         raw.Bugtraq,
			Component:       raw.Component,
			Status:          raw.Status,
			StatusMsg:       raw.StatusMsg,
			Priority:        raw.Priority,
			OperatingSystem: raw.OperatingSystem,
			CreateTime:      ct,
			UpdateTime:      ut,
			Pocname:         raw.Pocname,
			VerifyType:      raw.VerifyType,
			PatchUrl:        raw.PatchUrl,
			CvssVersion:     raw.CvssVersion,
			CvssScore:       raw.CvssScore,
			PocOrExp:        raw.PocOrExp,
			ScriptType:      raw.ScriptType,
		})
	}

	// B. 加载现有数据 (仅 ID 和 UpdateTime, VerifyType=1)
	db := mysql.FromContext(ctx)
	var existingVulns []mysqls.VulLibraries

	if err := db.Model(&mysqls.VulLibraries{}).
		Where("verify_type = ?", enums.VulVerifyTypePrincipleVerification).
		Select("id, vul_id, pocname, update_time"). // 增加 pocname
		Find(&existingVulns).Error; err != nil {
		return fmt.Errorf("查询现有漏洞数据失败: %v", err)
	}

	// C. 内存 Diff (改为使用 Pocname 作为 Key)
	existingMap := make(map[string]mysqls.VulLibraries)
	for _, v := range existingVulns {
		// 只有 Pocname 不为空时才建立映射，防止空值覆盖
		if v.Pocname != "" {
			existingMap[v.Pocname] = v
		}
	}

	var toInsert []mysqls.VulLibraries
	var toUpdate []mysqls.VulLibraries

	for _, v := range newVulns {
		v.VerifyType = enums.VulVerifyTypePrincipleVerification // 确保设置 VerifyType

		// 检查 Pocname 是否为空
		if v.Pocname == "" {
			// 如果 Pocname 为空，尝试用 VulID 匹配 (作为兜底)
			// 注意: 如果数据库里已有空 Pocname 记录，且有唯一索引，插入空 Pocname 会报错
			// 策略: 跳过 Pocname 为空的数据，记录警告日志
			log.Warnf("跳过 Pocname 为空的漏洞数据: %s", v.VulID)
			continue
		}

		// 使用 Pocname 进行匹配
		if old, exists := existingMap[v.Pocname]; exists {
			// 如果存在，复用现有 ID 进行更新
			v.ID = old.ID
			toUpdate = append(toUpdate, v)
		} else {
			// 如果不存在，清除 ID 以避免主键冲突 (让数据库自增)
			v.ID = 0
			toInsert = append(toInsert, v)
		}
	}

	// D. 批量操作
	if len(toInsert) > 0 {
		log.Infof("正在插入 %d 条新漏洞...", len(toInsert))
		if err := db.CreateInBatches(toInsert, 100).Error; err != nil {
			return fmt.Errorf("批量插入失败: %v", err)
		}
	}

	if len(toUpdate) > 0 {
		log.Infof("正在更新 %d 条现有漏洞...", len(toUpdate))
		// GORM 没有批量更新不同值的原生支持，只能循环更新或使用 CASE WHEN
		// 为保证稳定性，这里使用事务循环更新
		tx := db.Begin()
		for _, v := range toUpdate {
			// 只更新变更字段
			if err := tx.Model(&mysqls.VulLibraries{}).Where("id = ?", v.ID).Updates(v).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("更新漏洞 %s 失败: %v", v.VulID, err)
			}
		}
		tx.Commit()
	}

	log.Infof("PRINCIPLE 更新完成: 新增 %d, 更新 %d", len(toInsert), len(toUpdate))
	return nil
}

// applyFullUpdate 全量更新 (或低频更新)
func (s *VulnUpgradeStrategy) applyFullUpdate(ctx context.Context, jsonPath string) error {
	log.Info("开始执行漏洞库 FULL 更新...")
	// 复用之前的逻辑，逐条检查插入
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	// 定义导入结构体，确保所有字段的 JSON tag 与源数据（下划线风格）匹配
	type VulLibrariesImport struct {
		ID              int    `json:"id"`
		DataType        int    `json:"data_type"`
		VulID           string `json:"vul_id"`
		Name            string `json:"name"`
		Cve             string `json:"cve"`
		Risk            int    `json:"risk"`
		Type            int    `json:"type"`
		Class           int    `json:"class"`
		PublishedTime   string `json:"published_time"`
		Description     string `json:"description"`
		AffectRange     string `json:"affect_range"`
		ExploitImpact   int    `json:"exploit_impact"`
		FixSuggest      string `json:"fix_suggest"`
		Cnvd            string `json:"cnvd"`
		Cnnvd           string `json:"cnnvd"`
		Cncve           string `json:"cncve"`
		Bugtraq         string `json:"bugtraq"`
		Component       string `json:"component"`
		Status          int    `json:"status"`
		StatusMsg       string `json:"status_msg"`
		Priority        string `json:"priority"`
		OperatingSystem int    `json:"operating_system"`
		CreateTimeStr   string `json:"create_time"`
		UpdateTimeStr   string `json:"update_time"`
		Pocname         string `json:"pocname"`
		VerifyType      int    `json:"verify_type"`
		PatchUrl        string `json:"patch_url"`
		CvssVersion     string `json:"cvss_version"`
		CvssScore       string `json:"cvss_score"`
		PocOrExp        string `json:"poc_or_exp"`
		ScriptType      string `json:"script_type"`
		VulSource       int    `json:"vul_source"`
	}

	var rawVulns []VulLibrariesImport
	if err := json.Unmarshal(data, &rawVulns); err != nil {
		return fmt.Errorf("解析漏洞数据失败: %v", err)
	}

	var vulns []mysqls.VulLibraries
	for _, raw := range rawVulns {
		v := mysqls.VulLibraries{
			ID:              raw.ID,
			DataType:        raw.DataType,
			VulID:           raw.VulID,
			Name:            raw.Name,
			Cve:             raw.Cve,
			Risk:            raw.Risk,
			Type:            raw.Type,
			Class:           raw.Class,
			PublishedTime:   raw.PublishedTime,
			Description:     raw.Description,
			AffectRange:     raw.AffectRange,
			ExploitImpact:   raw.ExploitImpact,
			FixSuggest:      raw.FixSuggest,
			Cnvd:            raw.Cnvd,
			Cnnvd:           raw.Cnnvd,
			Cncve:           raw.Cncve,
			Bugtraq:         raw.Bugtraq,
			Component:       raw.Component,
			Status:          raw.Status,
			StatusMsg:       raw.StatusMsg,
			Priority:        raw.Priority,
			OperatingSystem: raw.OperatingSystem,
			Pocname:         raw.Pocname,
			VerifyType:      raw.VerifyType,
			PatchUrl:        raw.PatchUrl,
			CvssVersion:     raw.CvssVersion,
			CvssScore:       raw.CvssScore,
			PocOrExp:        raw.PocOrExp,
			ScriptType:      raw.ScriptType,
			VulSource:       raw.VulSource,
		}

		ct, err := parseTime(raw.CreateTimeStr)
		if err == nil && !ct.IsZero() {
			v.CreateTime = ct
		} else {
			v.CreateTime = time.Now() // 默认当前时间
		}
		ut, err := parseTime(raw.UpdateTimeStr)
		if err == nil && !ut.IsZero() {
			v.UpdateTime = ut
		} else {
			v.UpdateTime = time.Now() // 默认当前时间
		}
		vulns = append(vulns, v)
	}

	db := mysql.FromContext(ctx)
	successCount := 0

	// 使用 Upsert 语义 (MySQL: ON DUPLICATE KEY UPDATE)
	// GORM 的 Save 或 Clauses(clause.OnConflict{...})
	// 这里为了简单兼容，使用循环 Check-Insert/Update

	// 也可以分批处理
	batchSize := 100
	for i := 0; i < len(vulns); i += batchSize {
		end := i + batchSize
		if end > len(vulns) {
			end = len(vulns)
		}
		batch := vulns[i:end]

		// 简单起见，这里还是逐条，实际建议优化
		for _, v := range batch {
			var count int64
			// 修改匹配逻辑：使用 pocname 进行匹配
			db.Model(&mysqls.VulLibraries{}).Where("pocname = ?", v.Pocname).Count(&count)
			if count == 0 {
				if err := db.Omit("id").Create(&v).Error; err == nil {
					successCount++
				}
			} else {
				// 更新
				db.Model(&mysqls.VulLibraries{}).Where("pocname = ?", v.Pocname).Omit("id").Updates(v)
			}
		}
	}

	log.Infof("FULL 更新完成: 处理了 %d 条数据", len(vulns))
	return nil
}

// 辅助函数: 解析时间字符串 (支持 MySQL 默认格式和 RFC3339)
func parseTime(tStr string) (time.Time, error) {
	if tStr == "" {
		return time.Time{}, nil
	}
	layouts := []string{
		"2006-01-02 15:04:05", // MySQL default
		time.RFC3339,          // JSON default
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, tStr, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析时间: %s", tStr)
}

// applyDictionaryUpdate 增量更新字典表 (Name 匹配)
func (s *VulnUpgradeStrategy) applyDictionaryUpdate(ctx context.Context, jsonPath string) error {
	log.Info("开始执行字典表增量更新...")

	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	// 定义导入结构体，匹配下划线风格 JSON 和字符串时间
	type DictionaryImport struct {
		ID         int    `json:"id"`
		Sources    int    `json:"sources"`
		Service    int    `json:"service"`
		Name       string `json:"name"`
		Types      int    `json:"types"`
		IsDefault  int    `json:"is_default"`
		Content    string `json:"content"`
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
	}

	var rawDicts []DictionaryImport
	if err := json.Unmarshal(data, &rawDicts); err != nil {
		return fmt.Errorf("解析字典数据失败: %v", err)
	}

	// 转换为模型对象
	var newDicts []mysqls.Dictionary
	for _, raw := range rawDicts {
		ct, err := parseTime(raw.CreateTime)
		if err != nil || ct.IsZero() {
			ct = time.Now()
		}
		ut, err := parseTime(raw.UpdateTime)
		if err != nil || ut.IsZero() {
			ut = time.Now()
		}

		newDicts = append(newDicts, mysqls.Dictionary{
			ID:         raw.ID,
			Sources:    raw.Sources,
			Service:    raw.Service,
			Name:       raw.Name,
			Types:      raw.Types,
			IsDefault:  raw.IsDefault,
			Content:    raw.Content,
			CreateTime: ct,
			UpdateTime: ut,
		})
	}

	db := mysql.FromContext(ctx)
	var existingDicts []mysqls.Dictionary
	// 查询现有数据，只需要 ID 和 Name
	if err := db.Model(&mysqls.Dictionary{}).Select("id, name").Find(&existingDicts).Error; err != nil {
		return fmt.Errorf("查询现有字典数据失败: %v", err)
	}

	// 内存 Diff
	existingMap := make(map[string]mysqls.Dictionary)
	for _, v := range existingDicts {
		if v.Name != "" {
			existingMap[v.Name] = v
		}
	}

	var toInsert []mysqls.Dictionary
	var toUpdate []mysqls.Dictionary

	for _, v := range newDicts {
		if v.Name == "" {
			continue
		}
		if old, exists := existingMap[v.Name]; exists {
			v.ID = old.ID
			toUpdate = append(toUpdate, v)
		} else {
			v.ID = 0
			toInsert = append(toInsert, v)
		}
	}

	// 批量操作
	if len(toInsert) > 0 {
		log.Infof("正在插入 %d 条新字典...", len(toInsert))
		if err := db.CreateInBatches(toInsert, 100).Error; err != nil {
			return fmt.Errorf("批量插入字典失败: %v", err)
		}
	}

	if len(toUpdate) > 0 {
		log.Infof("正在更新 %d 条现有字典...", len(toUpdate))
		tx := db.Begin()
		for _, v := range toUpdate {
			if err := tx.Model(&mysqls.Dictionary{}).Where("id = ?", v.ID).Updates(v).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("更新字典 %s 失败: %v", v.Name, err)
			}
		}
		tx.Commit()
	}

	log.Infof("字典表更新完成: 新增 %d, 更新 %d", len(toInsert), len(toUpdate))
	return nil
}

// applyFingerUpdate 增量更新指纹表 (AppName 匹配)
func (s *VulnUpgradeStrategy) applyFingerUpdate(ctx context.Context, jsonPath string) error {
	log.Info("开始执行指纹表增量更新...")

	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	// 定义导入结构体，匹配下划线风格 JSON 和字符串时间
	type FingerImport struct {
		ID         int    `json:"id"`
		AppClass   int    `json:"app_class"`
		AppVersion string `json:"app_version"`
		CnName     string `json:"cn_name"`
		AppName    string `json:"app_name"`
		Flag       string `json:"flag"`
		Source     int    `json:"source"`
		FingerType int    `json:"finger_type"`
		Desc       string `json:"desc"`
		Level      string `json:"level"`
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
	}

	var rawFingers []FingerImport
	if err := json.Unmarshal(data, &rawFingers); err != nil {
		return fmt.Errorf("解析指纹数据失败: %v", err)
	}

	// 转换为模型对象
	var newFingers []mysqls.Finger
	for _, raw := range rawFingers {
		ct, err := parseTime(raw.CreateTime)
		if err != nil || ct.IsZero() {
			ct = time.Now()
		}
		ut, err := parseTime(raw.UpdateTime)
		if err != nil || ut.IsZero() {
			ut = time.Now()
		}

		newFingers = append(newFingers, mysqls.Finger{
			ID:         raw.ID,
			AppClass:   raw.AppClass,
			AppVersion: raw.AppVersion,
			CnName:     raw.CnName,
			AppName:    raw.AppName,
			Flag:       raw.Flag,
			Source:     raw.Source,
			FingerType: raw.FingerType,
			Desc:       raw.Desc,
			Level:      raw.Level,
			CreateTime: ct,
			UpdateTime: &ut,
		})
	}

	db := mysql.FromContext(ctx)
	var existingFingers []mysqls.Finger
	// 查询现有数据，只需要 ID 和 AppName
	if err := db.Model(&mysqls.Finger{}).Select("id, app_name").Find(&existingFingers).Error; err != nil {
		return fmt.Errorf("查询现有指纹数据失败: %v", err)
	}

	// 内存 Diff
	existingMap := make(map[string]mysqls.Finger)
	for _, v := range existingFingers {
		if v.AppName != "" {
			existingMap[v.AppName] = v
		}
	}

	var toInsert []mysqls.Finger
	var toUpdate []mysqls.Finger

	for _, v := range newFingers {
		if v.AppName == "" {
			continue
		}
		if old, exists := existingMap[v.AppName]; exists {
			v.ID = old.ID
			toUpdate = append(toUpdate, v)
		} else {
			v.ID = 0
			toInsert = append(toInsert, v)
		}
	}

	// 批量操作
	if len(toInsert) > 0 {
		log.Infof("正在插入 %d 条新指纹...", len(toInsert))
		if err := db.CreateInBatches(toInsert, 100).Error; err != nil {
			return fmt.Errorf("批量插入指纹失败: %v", err)
		}
	}

	if len(toUpdate) > 0 {
		log.Infof("正在更新 %d 条现有指纹...", len(toUpdate))
		tx := db.Begin()
		for _, v := range toUpdate {
			if err := tx.Model(&mysqls.Finger{}).Where("id = ?", v.ID).Updates(v).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("更新指纹 %s 失败: %v", v.AppName, err)
			}
		}
		tx.Commit()
	}

	log.Infof("指纹表更新完成: 新增 %d, 更新 %d", len(toInsert), len(toUpdate))
	return nil
}

// applyScriptDB 更新攻击脚本库 (SQLite)
func (s *VulnUpgradeStrategy) applyScriptDB(ctx context.Context, srcPath string) error {
	log.Info("开始更新攻击脚本库 (scanner.db)...")

	scannerDBPath := getScannerDBPath()

	// 1. 确保目标目录存在
	targetDir := filepath.Dir(scannerDBPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// 2. 替换 scanner.db (System Library)
	// 用户自定义数据应存储在 user_custom.db，因此替换 scanner.db 是安全的
	// 为了原子性，先 Copy 到 .tmp 然后 Rename (在 Linux 上)
	tmpPath := scannerDBPath + ".tmp"
	if err := file.CopyFile(srcPath, tmpPath); err != nil {
		return fmt.Errorf("复制 scanner.db 失败: %v", err)
	}

	// Rename 覆盖
	if err := os.Rename(tmpPath, scannerDBPath); err != nil {
		// Windows 上如果文件被占用可能会失败，尝试直接 Copy 覆盖
		if err := file.CopyFile(srcPath, scannerDBPath); err != nil {
			return fmt.Errorf("覆盖 scanner.db 失败: %v", err)
		}
		os.Remove(tmpPath)
	}

	log.Infof("scanner.db 更新成功")
	return nil
}

func (s *VulnUpgradeStrategy) Rollback(ctx context.Context, backupDir string) error {
	if backupDir == "" {
		return nil
	}

	// 1. 回滚 scanner.db
	backupDB := filepath.Join(backupDir, "scanner.db")
	if file.CheckPathExist(backupDB) {
		if err := file.CopyFile(backupDB, getScannerDBPath()); err != nil {
			return err
		}
	}

	// 2. 回滚 user_custom.db
	backupUserDB := filepath.Join(backupDir, "user_custom.db")
	if file.CheckPathExist(backupUserDB) {
		if err := file.CopyFile(backupUserDB, getUserCustomDBPath()); err != nil {
			return err
		}
	}

	// 3. 回滚 MySQL 数据
	// 读取 manifest_backup.json 确定 Scope
	var scope VulnScope = VulnScopeFull // 默认为 Full
	if data, err := ioutil.ReadFile(filepath.Join(backupDir, "manifest_backup.json")); err == nil {
		var m UpgradeManifest
		if err := json.Unmarshal(data, &m); err == nil {
			scope = m.VulnScope
		}
	}

	db := mysql.FromContext(ctx)

	// A. 恢复 vul_libraries
	vulnBackupDir := filepath.Join(backupDir, "vul_libraries")
	if file.CheckPathExist(vulnBackupDir) {
		log.Info("正在回滚 vul_libraries 表...")
		tx := db.Begin()

		// 清理数据
		if scope == VulnScopePrinciple {
			if err := tx.Where("verify_type = ?", enums.VulVerifyTypePrincipleVerification).Delete(&mysqls.VulLibraries{}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("清理 vul_libraries 失败: %v", err)
			}
		} else {
			if err := tx.Exec("TRUNCATE TABLE vul_libraries").Error; err != nil {
				// 如果 TRUNCATE 失败 (例如权限问题)，尝试 DELETE
				if err := tx.Where("1=1").Delete(&mysqls.VulLibraries{}).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("清空 vul_libraries 失败: %v", err)
				}
			}
		}

		// 读取分片并插入
		files, err := ioutil.ReadDir(vulnBackupDir)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("读取备份目录失败: %v", err)
		}

		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".json") {
				continue
			}
			data, err := ioutil.ReadFile(filepath.Join(vulnBackupDir, f.Name()))
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("读取备份文件 %s 失败: %v", f.Name(), err)
			}

			var libs []mysqls.VulLibraries
			if err := json.Unmarshal(data, &libs); err != nil {
				tx.Rollback()
				return fmt.Errorf("解析备份文件 %s 失败: %v", f.Name(), err)
			}

			if len(libs) > 0 {
				if err := tx.CreateInBatches(libs, 100).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("恢复数据失败: %v", err)
				}
			}
		}
		tx.Commit()
	}

	// B. 恢复 dictionary (全量覆盖)
	dictPath := filepath.Join(backupDir, "dictionary.json")
	if file.CheckPathExist(dictPath) {
		log.Info("正在回滚 dictionary 表...")
		data, err := ioutil.ReadFile(dictPath)
		if err == nil {
			var dicts []mysqls.Dictionary
			if err := json.Unmarshal(data, &dicts); err == nil {
				tx := db.Begin()
				tx.Exec("TRUNCATE TABLE dictionary")
				if len(dicts) > 0 {
					tx.CreateInBatches(dicts, 100)
				}
				tx.Commit()
			}
		}
	}

	// C. 恢复 finger (全量覆盖)
	fingerPath := filepath.Join(backupDir, "finger.json")
	if file.CheckPathExist(fingerPath) {
		log.Info("正在回滚 finger 表...")
		data, err := ioutil.ReadFile(fingerPath)
		if err == nil {
			var fingers []mysqls.Finger
			if err := json.Unmarshal(data, &fingers); err == nil {
				tx := db.Begin()
				tx.Exec("TRUNCATE TABLE finger")
				if len(fingers) > 0 {
					tx.CreateInBatches(fingers, 100)
				}
				tx.Commit()
			}
		}
	}

	var mapSetService MapSet
	ctx = mysql.NewContext(ctx, mysql.GetDB())
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if err == nil && objValueStr != "" {
		var systemVersionMapSet SystemVersionMapSet
		if json.Unmarshal([]byte(objValueStr), &systemVersionMapSet) == nil {
			if systemVersionMapSet.LastVulVersion != "" {
				systemVersionMapSet.VulVersion = systemVersionMapSet.LastVulVersion
				systemVersionMapSet.LastVulVersion = ""
				systemVersionMapSet.VulUpdateTime = time.Now().Format(enums.TimeLayout)
				if objValueByte, err := json.Marshal(systemVersionMapSet); err == nil {
					if err = mapSetService.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent); err != nil {
						log.Errorf("回滚漏洞库版本信息失败: %v", err)
					} else {
						log.Info("回滚漏洞库版本信息成功")
					}
				}
			}
		}
	}

	return nil
}

func (s *VulnUpgradeStrategy) PostAction(ctx context.Context) error {
	// 注入全局 DB 连接
	ctx = mysql.NewContext(ctx, mysql.GetDB())

	// 通知引擎重新加载规则
	// 实际项目中可能需要通过 RPC 通知 scanner 服务
	log.Info("通知扫描引擎重新加载规则库...")

	// 更新版本信息
	// 注意：版本信息更新已在 Apply 阶段完成 (VulnUpgradeStrategy.Apply)，
	// 这里不需要再次更新，否则会导致 LastVulVersion 被错误覆盖为当前版本。
	/*
		if s.Manifest.Version != "" {
			var mapSetService MapSet
			var systemVersionMapSet SystemVersionMapSet

			// 尝试获取现有记录
			objValueStr, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
			if err == nil && objValueStr != "" {
				json.Unmarshal([]byte(objValueStr), &systemVersionMapSet)
			} else {
				// 初始化默认值
				systemVersionMapSet.CurrentVersion = enums.SystemBaseVersion
			}

			// 记录上次版本
			systemVersionMapSet.LastVulVersion = systemVersionMapSet.VulVersion

			// 更新当前版本和时间
			systemVersionMapSet.VulVersion = s.Manifest.Version
			systemVersionMapSet.VulUpdateTime = time.Now().Format(enums.TimeLayout)

			// 保存
			objValueByte, err := json.Marshal(systemVersionMapSet)
			if err == nil {
				if err = mapSetService.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent); err != nil {
					log.Errorf("更新漏洞库版本信息失败: %v", err)
				} else {
					log.Infof("更新漏洞库版本信息成功: Version=%s, Time=%s", systemVersionMapSet.VulVersion, systemVersionMapSet.VulUpdateTime)
				}
			} else {
				log.Errorf("序列化漏洞库版本信息失败: %v", err)
			}
		}
	*/

	return nil
}

// ==========================================
// 5. 具体策略: 系统升级 (SystemUpgradeStrategy)
// ==========================================

// SystemUpgradeItem 定义系统升级项配置
type SystemUpgradeItem struct {
	Name                 string
	SourceRel            string      // 升级包内的相对路径 (e.g. assets/bin/smart)
	DestPath             string      // 系统绝对路径
	BackupName           string      // 备份目录下的名称
	IsDir                bool        // 是否为目录
	Chmod                os.FileMode // 需设置的权限 (如 0755)
	AlternativeSourceRel string      // 备用源路径 (兼容不同打包结构)
}

// getSystemUpgradeItems 获取升级配置列表
func getSystemUpgradeItems() []SystemUpgradeItem {
	base := enums.SystemUpgradeProjectDir
	smartDir := filepath.Join(base, enums.SystemUpgradeSmartDir)
	binDir := filepath.Join(base, enums.SystemUpgradeBinDir)
	nginxVueDir := filepath.Join(base, enums.SystemUpgradeNginxDir, enums.SystemUpgradeVueDir)

	return []SystemUpgradeItem{
		{
			Name:       "Smart组件",
			SourceRel:  filepath.Join(enums.SystemUpgradeAssetsDir, enums.SystemUpgradeBinDir, enums.SystemUpgradeSmartFile),
			DestPath:   filepath.Join(smartDir, enums.SystemUpgradeSmartFile),
			BackupName: enums.SystemUpgradeSmartFile,
			IsDir:      false,
			Chmod:      0755,
		},
		{
			Name:       "Scanner组件",
			SourceRel:  filepath.Join(enums.SystemUpgradeAssetsDir, enums.SystemUpgradeBinDir, enums.SystemUpgradeScannerFile),
			DestPath:   filepath.Join(smartDir, enums.SystemUpgradeScannerFile),
			BackupName: enums.SystemUpgradeScannerFile,
			IsDir:      false,
			Chmod:      0755,
		},
		{
			Name:       "DbWeb组件",
			SourceRel:  filepath.Join(enums.SystemUpgradeAssetsDir, enums.SystemUpgradeBinDir, enums.SystemUpgradeDbWebFile),
			DestPath:   filepath.Join(binDir, enums.SystemUpgradeDbWebFile),
			BackupName: enums.SystemUpgradeDbWebFile,
			IsDir:      false,
			Chmod:      0755,
		},
		{
			Name:                 "前端资源",
			SourceRel:            filepath.Join(enums.SystemUpgradeAssetsDir, enums.SystemUpgradeWebDir, enums.SystemUpgradeVueDir),
			AlternativeSourceRel: filepath.Join(enums.SystemUpgradeAssetsDir, enums.SystemUpgradeWebDir), // 兼容旧结构
			DestPath:             nginxVueDir,
			BackupName:           enums.SystemUpgradeVueDir,
			IsDir:                true,
			Chmod:                0755,
		},
	}
}

type SystemUpgradeStrategy struct {
	Manifest *UpgradeManifest
}

func (s *SystemUpgradeStrategy) Verify(ctx context.Context, m *UpgradeManifest, dir string) error {
	// 校验逻辑：检查是否有任意一个有效更新文件
	// 如果所有配置项在升级包中都不存在，则认为是无效包（或空包）
	items := getSystemUpgradeItems()
	foundAny := false
	for _, item := range items {
		srcPath := filepath.Join(dir, item.SourceRel)
		if file.CheckPathExist(srcPath) {
			foundAny = true
			continue
		}
		if item.AlternativeSourceRel != "" {
			if file.CheckPathExist(filepath.Join(dir, item.AlternativeSourceRel)) {
				foundAny = true
				continue
			}
		}
	}

	if !foundAny {
		if sqlFiles, err := findUpgradeSQLFiles(dir); err == nil && len(sqlFiles) > 0 {
			foundAny = true
		}
	}

	if !foundAny {
		return errors.New("升级包中未找到任何有效的系统更新文件")
	}
	return nil
}

func (s *SystemUpgradeStrategy) Backup(ctx context.Context) (string, error) {
	backupDir := filepath.Join(enums.SystemUpgradeProjectDir, "smart/backup", fmt.Sprintf("sys_%d", time.Now().Unix()))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	items := getSystemUpgradeItems()
	for _, item := range items {
		// 只有当目标文件存在时才备份
		if file.CheckPathExist(item.DestPath) {
			backupPath := filepath.Join(backupDir, item.BackupName)
			if item.IsDir {
				if err := copyDir(item.DestPath, backupPath); err != nil {
					return "", fmt.Errorf("备份 %s 失败: %v", item.Name, err)
				}
			} else {
				if err := file.CopyFile(item.DestPath, backupPath); err != nil {
					return "", fmt.Errorf("备份 %s 失败: %v", item.Name, err)
				}
			}
		}
	}

	// Backup configs (Always try to backup configs if they exist)
	smartDir := filepath.Join(enums.SystemUpgradeProjectDir, enums.SystemUpgradeSmartDir)
	if file.CheckPathExist(filepath.Join(smartDir, enums.SystemUpgradeSmartConfigFile)) {
		file.CopyFile(filepath.Join(smartDir, enums.SystemUpgradeSmartConfigFile), filepath.Join(backupDir, "smart_config.json"))
	}
	if file.CheckPathExist(enums.DockerComposeConfigPath) {
		file.CopyFile(enums.DockerComposeConfigPath, filepath.Join(backupDir, "docker-compose.yml"))
	}

	// 备份版本信息
	var mapSetService MapSet
	if val, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey); err == nil {
		ioutil.WriteFile(filepath.Join(backupDir, "version.json"), []byte(val), 0644)
	}

	return backupDir, nil
}

func (s *SystemUpgradeStrategy) Apply(ctx context.Context, dir string) error {
	log.Info("正在应用系统升级...")
	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在停止服务...", 50)

	// 1. 停止服务 (仅 Linux)
	if runtime.GOOS == "linux" {
		log.Info("正在停止服务...")
	}

	items := getSystemUpgradeItems()
	for _, item := range items {
		srcPath := filepath.Join(dir, item.SourceRel)
		realSrcPath := ""

		if file.CheckPathExist(srcPath) {
			realSrcPath = srcPath
		} else if item.AlternativeSourceRel != "" {
			altPath := filepath.Join(dir, item.AlternativeSourceRel)
			if file.CheckPathExist(altPath) {
				realSrcPath = altPath
			}
		}

		if realSrcPath != "" {
			UpdateUpgradeStatus(enums.UpgradeStateUpgrading, fmt.Sprintf("正在更新 %s...", item.Name), 60)
			log.Infof("应用更新: %s -> %s", realSrcPath, item.DestPath)

			// Ensure dest dir exists
			destDir := filepath.Dir(item.DestPath)
			os.MkdirAll(destDir, 0755)

			if item.IsDir {
				os.RemoveAll(item.DestPath) // Clean old dir
				if err := copyDir(realSrcPath, item.DestPath); err != nil {
					return fmt.Errorf("更新 %s 失败: %v", item.Name, err)
				}
			} else {
				// 使用原子替换策略：先复制到临时文件，再重命名覆盖
				// 这样可以避免 "text file busy" 错误 (当更新正在运行的二进制文件时)
				tmpDestPath := item.DestPath + ".tmp"
				if err := file.CopyFile(realSrcPath, tmpDestPath); err != nil {
					return fmt.Errorf("复制临时文件失败 %s: %v", item.Name, err)
				}

				// 恢复权限
				if item.Chmod != 0 {
					os.Chmod(tmpDestPath, item.Chmod)
				}

				// 原子替换
				if err := os.Rename(tmpDestPath, item.DestPath); err != nil {
					// 如果 Rename 失败 (例如跨设备)，尝试传统的删除再移动
					os.Remove(tmpDestPath) // 清理临时文件
					log.Warnf("Rename 失败，尝试强制覆盖: %v", err)

					// 尝试删除原文件 (对于正在运行的程序，Remove 是允许的)
					os.Remove(item.DestPath)

					if err := file.CopyFile(realSrcPath, item.DestPath); err != nil {
						return fmt.Errorf("更新 %s 失败: %v", item.Name, err)
					}
					if item.Chmod != 0 {
						os.Chmod(item.DestPath, item.Chmod)
					}
				}
			}
		} else {
			log.Infof("跳过更新 %s: 升级包中不存在", item.Name)
		}
	}

	if err := s.applyUpgradeSQL(ctx, dir); err != nil {
		return err
	}

	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "系统文件更新完成", 95)
	return nil
}

func (s *SystemUpgradeStrategy) applyUpgradeSQL(ctx context.Context, dir string) error {
	sqlFiles, err := findUpgradeSQLFiles(dir)
	if err != nil {
		return err
	}
	if len(sqlFiles) == 0 {
		return nil
	}

	UpdateUpgradeStatus(enums.UpgradeStateUpgrading, "正在执行数据库脚本...", 85)

	// Group files by target database
	dbFilesMap := make(map[string][]string)
	for _, p := range sqlFiles {
		targetDB := getUpgradeSQLTargetDB(p)
		dbFilesMap[targetDB] = append(dbFilesMap[targetDB], p)
	}

	// Sort DB names for deterministic execution order
	var dbNames []string
	for k := range dbFilesMap {
		dbNames = append(dbNames, k)
	}
	sort.Strings(dbNames)

	totalFiles := len(sqlFiles)
	doneFiles := 0

	for _, dbName := range dbNames {
		files := dbFilesMap[dbName]
		// Sort files within the same DB
		sort.Strings(files)

		dsn, err := getDSNForDatabase(dbName)
		if err != nil {
			return fmt.Errorf("获取数据库 %s 连接失败: %v", dbName, err)
		}

		gormDB, err := openUpgradeGormDB(dsn)
		if err != nil {
			return fmt.Errorf("连接数据库 %s 失败: %v", dbName, err)
		}

		for _, p := range files {
			doneFiles++
			UpdateUpgradeStatus(enums.UpgradeStateUpgrading, fmt.Sprintf("正在执行数据库脚本 (%d/%d)...", doneFiles, totalFiles), 85)
			if err := executeSQLFileWithGorm(ctx, gormDB, dbName, p); err != nil {
				return err
			}
		}
	}

	return nil
}

func findUpgradeSQLFiles(unzipDir string) ([]string, error) {
	candidates := []string{
		filepath.Join(unzipDir, "assets", "sql"),
		// Fallback for compatibility
		filepath.Join(unzipDir, "sql"),
		filepath.Join(unzipDir, "payload", "sql"),
	}

	var files []string
	for _, dir := range candidates {
		if !file.CheckPathExist(dir) {
			continue
		}
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			nameLower := strings.ToLower(info.Name())
			if strings.HasSuffix(nameLower, ".sql") {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	// Remove duplicates if any (though candidates are distinct dirs usually)
	uniqueFiles := make(map[string]bool)
	var result []string
	for _, f := range files {
		if !uniqueFiles[f] {
			uniqueFiles[f] = true
			result = append(result, f)
		}
	}
	sort.Strings(result)
	return result, nil
}

func getUpgradeSQLTargetDB(sqlFilePath string) string {
	base := filepath.Base(sqlFilePath)
	parent := filepath.Base(filepath.Dir(sqlFilePath))

	// 1. Priority: Parent directory name (if not a generic name)
	// This allows structure like: assets/sql/smart/01.sql -> DB: smart
	parentLower := strings.ToLower(parent)
	if parentLower != "sql" && parentLower != "assets" && parentLower != "db" && parentLower != "extract" && parentLower != "payload" {
		return parent
	}

	// 2. Fallback: Filename without extension
	// This supports: assets/sql/smart.sql -> DB: smart
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func getDSNForDatabase(targetDB string) (string, error) {
	dbSettings := make([]mysql.DatabaseSetting, 1)
	if err := config.Load("mysql", &dbSettings); err != nil {
		return "", err
	}

	var defaultDBSetting *mysql.DatabaseSetting
	for i := range dbSettings {
		if dbSettings[i].Name == "default" {
			defaultDBSetting = &dbSettings[i]
			break
		}
	}
	if defaultDBSetting == nil {
		if len(dbSettings) > 0 {
			defaultDBSetting = &dbSettings[0]
		} else {
			return "", errors.New("no mysql config found")
		}
	}

	smartDSN := defaultDBSetting.Master
	// If target is "smart" (convention for default DB), return default DSN directly
	if strings.ToLower(targetDB) == "smart" {
		return smartDSN, nil
	}

	return replaceDatabaseInDSN(smartDSN, targetDB)
}

func replaceDatabaseInDSN(dsn string, database string) (string, error) {
	slash := strings.LastIndex(dsn, "/")
	if slash < 0 || slash == len(dsn)-1 {
		return "", fmt.Errorf("无效的 MySQL DSN: %s", dsn)
	}
	rest := dsn[slash+1:]
	q := strings.Index(rest, "?")
	if q < 0 {
		return dsn[:slash+1] + database, nil
	}
	return dsn[:slash+1] + database + rest[q:], nil
}

func openUpgradeGormDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(gormmysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
	})
}

func executeSQLFileWithGorm(ctx context.Context, db *gorm.DB, dbName string, sqlFilePath string) error {
	stmts, err := readSQLStatements(sqlFilePath)
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		if err := db.WithContext(ctx).Exec(stmt).Error; err != nil {
			if isIgnorableSQLError(err) {
				log.Warnf("SQL 执行遇到幂等性错误 (自动忽略): db=%s file=%s err=%v", dbName, filepath.Base(sqlFilePath), err)
				continue
			}
			return fmt.Errorf("执行 SQL 失败: db=%s file=%s err=%v", dbName, sqlFilePath, err)
		}
	}
	return nil
}

func isIgnorableSQLError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// 常见幂等性错误码:
	// 1050: Table already exists
	// 1051: Unknown table (DROP TABLE 时)
	// 1060: Duplicate column name
	// 1061: Duplicate key name
	// 1091: Can't DROP 'x'; check that column/key exists
	// 1068: Multiple primary key defined
	ignorableCodes := []string{
		"Error 1050",
		"Error 1051",
		"Error 1060",
		"Error 1061",
		"Error 1091",
		"Error 1068",
	}

	for _, code := range ignorableCodes {
		if strings.Contains(msg, code) {
			return true
		}
	}
	return false
}

func readSQLStatements(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if bytes.HasPrefix(b, []byte{0xEF, 0xBB, 0xBF}) {
		b = b[3:]
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)

	delimiter := ";"
	inBlockComment := false
	var stmts []string
	var sb strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		upperTrimmed := strings.ToUpper(trimmed)

		if !inBlockComment && strings.HasPrefix(upperTrimmed, "DELIMITER ") {
			delimiter = strings.TrimSpace(trimmed[len("DELIMITER "):])
			continue
		}

		if inBlockComment {
			if idx := strings.Index(line, "*/"); idx >= 0 {
				line = line[idx+2:]
				inBlockComment = false
			} else {
				continue
			}
		}

		for {
			start := strings.Index(line, "/*")
			if start < 0 {
				break
			}
			end := strings.Index(line[start+2:], "*/")
			if end < 0 {
				line = line[:start]
				inBlockComment = true
				break
			}
			endIdx := start + 2 + end
			line = line[:start] + line[endIdx+2:]
		}

		trimmedNoSpace := strings.TrimSpace(line)
		if trimmedNoSpace == "" {
			continue
		}
		if strings.HasPrefix(trimmedNoSpace, "--") || strings.HasPrefix(trimmedNoSpace, "#") {
			continue
		}

		sb.WriteString(line)
		sb.WriteString("\n")

		if delimiter != "" {
			pending := sb.String()
			for {
				idx := strings.Index(pending, delimiter)
				if idx < 0 {
					break
				}
				stmt := strings.TrimSpace(pending[:idx])
				if stmt != "" {
					stmts = append(stmts, stmt)
				}
				pending = pending[idx+len(delimiter):]
			}
			sb.Reset()
			sb.WriteString(pending)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	last := strings.TrimSpace(sb.String())
	if last != "" {
		stmts = append(stmts, last)
	}
	return stmts, nil
}

func (s *SystemUpgradeStrategy) Rollback(ctx context.Context, backupDir string) error {
	items := getSystemUpgradeItems()
	for _, item := range items {
		backupPath := filepath.Join(backupDir, item.BackupName)
		if file.CheckPathExist(backupPath) {
			log.Infof("正在回滚 %s...", item.Name)
			if item.IsDir {
				os.RemoveAll(item.DestPath)
				copyDir(backupPath, item.DestPath)
			} else {
				file.CopyFile(backupPath, item.DestPath)
				if item.Chmod != 0 {
					os.Chmod(item.DestPath, item.Chmod)
				}
			}
		}
	}

	// Restore configs
	smartDir := filepath.Join(enums.SystemUpgradeProjectDir, enums.SystemUpgradeSmartDir)
	if file.CheckPathExist(filepath.Join(backupDir, "smart_config.json")) {
		file.CopyFile(filepath.Join(backupDir, "smart_config.json"), filepath.Join(smartDir, enums.SystemUpgradeSmartConfigFile))
	}

	// 恢复版本信息
	var mapSetService MapSet
	ctx = mysql.NewContext(ctx, mysql.GetDB())
	objValueStr, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
	if err == nil && objValueStr != "" {
		var systemVersionMapSet SystemVersionMapSet
		if json.Unmarshal([]byte(objValueStr), &systemVersionMapSet) == nil {
			if systemVersionMapSet.LastSystemVersion != "" {
				systemVersionMapSet.CurrentVersion = systemVersionMapSet.LastSystemVersion
				systemVersionMapSet.LastSystemVersion = ""
				systemVersionMapSet.UpdateTime = time.Now().Format(enums.TimeLayout)
				if objValueByte, err := json.Marshal(systemVersionMapSet); err == nil {
					mapSetService.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent)
				}
			}
		}
	}

	return nil
}

func (s *SystemUpgradeStrategy) PostAction(ctx context.Context) error {
	// 注入全局 DB 连接，防止传入的 context 缺失 DB
	ctx = mysql.NewContext(ctx, mysql.GetDB())

	// 1. 先更新系统版本信息 (确保服务启动时能读取到最新版本)
	// status := GetUpgradeStatus() // 弃用：改用 s.Manifest 直接获取，避免全局状态丢失
	log.Infof("PostAction 检查状态: Type=%s, Version=%s", s.Manifest.Type, s.Manifest.Version)
	if s.Manifest.Type == UpgradeTypeSystem && s.Manifest.Version != "" {
		var mapSetService MapSet
		var systemVersionMapSet SystemVersionMapSet
		objValueStr, err := mapSetService.GetMapValue(ctx, enums.SystemVersionMapSetObjKey)
		log.Infof("获取现有版本信息: objValueStr=%s, err=%v", objValueStr, err)
		if err == nil && objValueStr != "" {
			json.Unmarshal([]byte(objValueStr), &systemVersionMapSet)
		} else {
			systemVersionMapSet.CurrentVersion = enums.SystemBaseVersion
		}
		systemVersionMapSet.LastSystemVersion = systemVersionMapSet.CurrentVersion
		systemVersionMapSet.CurrentVersion = s.Manifest.Version
		systemVersionMapSet.UpdateTime = time.Now().Format(enums.TimeLayout)
		objValueByte, err := json.Marshal(systemVersionMapSet)
		if err == nil {
			log.Infof("准备更新版本信息: %s", string(objValueByte))
			if err = mapSetService.Create(ctx, enums.SystemVersionMapSetObjKey, string(objValueByte), enums.SystemVersionMapSetContent); err != nil {
				log.Errorf("更新系统版本信息失败: %v", err)
			} else {
				log.Info("更新系统版本信息成功")
			}
		} else {
			log.Errorf("序列化系统版本信息失败: %v", err)
		}
	} else {
		log.Warn("PostAction 跳过版本更新: 状态不符合要求")
	}

	log.Info("正在重启系统服务...")
	// 2. 启动服务 (仅 Linux)
	if runtime.GOOS == "linux" {
		var restartErr error
		if utils.SystemHasCommand(ctx, "supervisorctl") {
			restartErr = exec.Command("supervisorctl", "restart", enums.ServiceNameSmart).Run()
			if restartErr == nil {
				return nil
			}
			log.Errorf("supervisorctl restart 失败: %v", restartErr)
		}
		if utils.SystemHasCommand(ctx, "supervisord") {
			restartErr = exec.Command("bash", "-c", "supervisord ctl restart "+enums.ServiceNameSmart).Run()
			if restartErr == nil {
				return nil
			}
			log.Errorf("supervisord ctl restart 失败: %v", restartErr)
		}
		restartErr = exec.Command("pkill", "-f", enums.ServiceNameSmart).Run()
		if restartErr == nil {
			return nil
		}
		log.Errorf("pkill -f smart 失败: %v", restartErr)
		restartErr = exec.Command("bash", "-c", "ps aux | grep "+enums.ServiceNameSmart+" | grep -v grep | awk '{print $2}' | xargs kill -9").Run()
		if restartErr != nil {
			log.Errorf("手动 kill -9 失败: %v", restartErr)
			return restartErr
		}
	}
	return nil
}

func copyDir(src, dst string) error {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	for _, e := range entries {
		sp := filepath.Join(src, e.Name())
		dp := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(sp, dp); err != nil {
				return err
			}
		} else {
			// Skip non-regular files (like sockets, pipes, devices)
			// But allow symlinks (mode&os.ModeSymlink != 0) if needed,
			// though CopyFile implementation checks IsRegular which excludes symlinks too.
			// Let's explicitly check mode here to avoid calling CopyFile for sockets.
			if e.Mode()&os.ModeSocket != 0 {
				log.Warnf("Skipping socket file during backup: %s", sp)
				continue
			}
			// CopyFile internals also check IsRegular, but we want to avoid the error return.
			// So we check here: if it's not a regular file and not a symlink (if we wanted to support them), skip it.
			// For simplicity and safety, we only backup regular files.
			if !e.Mode().IsRegular() {
				log.Warnf("Skipping non-regular file during backup: %s (Mode: %s)", sp, e.Mode())
				continue
			}

			if err := file.CopyFile(sp, dp); err != nil {
				return err
			}
		}
	}
	return nil
}
