# 升级包制作工具 (Packager) 设计方案

## 1. 设计目标

为了配合新的升级系统，我们需要一个标准化、自动化的工具（CLI）来制作升级包。该工具主要负责：
1.  **自动化构建**：根据配置自动收集文件。
2.  **元数据生成**：自动计算哈希、时间戳，生成 `manifest.json`。
3.  **安全签名**：使用私钥对升级包进行签名，确保不可抵赖和防篡改。
4.  **统一打包**：输出符合规范的 `.zip` 文件。

此工具通常集成在 CI/CD 流水线中，也可以由开发人员在本地手动执行。

---

## 2. 工具架构 (smart-packager)

建议开发一个独立的 CLI 工具，命名为 `smart-packager`。

### 2.1 命令结构

```bash
smart-packager [command] [flags]

Commands:
  build-system   构建系统升级包
  build-vuln     构建漏洞库升级包
  gen-keys       生成 RSA 公私钥对 (用于初始化环境)
  verify         校验已有的升级包 (测试用)
```

### 2.2 核心流程

1.  **Input (输入)**:
    *   **上传模式**: 接收用户上传的文件（如 SQL 脚本、二进制文件），直接作为资源文件。
    *   **DB 模式**: 连接远程 MySQL 数据库导出表数据（支持多表），或连接远程服务器通过 SFTP 下载文件（如 `scanner.db`）。
    *   **参数**: 版本号、最小系统版本、描述信息、私钥等。

2.  **Staging (组装)**:
    *   创建临时工作目录 (`packager_build_{timestamp}`)。
    *   **上传模式**: 将上传文件移动到临时目录，根据配置映射到 `assets/` 路径。
    *   **DB 模式**:
        *   导出表数据为紧凑格式 JSON 文件 (`{table}.json`)，存入 `assets/db/`。
        *   或者下载远程 `scanner.db` 到 `assets/db/`。
        *   或者生成 `update.sql` 到 `assets/sql/`。

3.  **Manifest (元数据)**:
    *   遍历 `assets` 目录，计算关键文件的 SHA256。
    *   填充 Version, BuildTime, MinSystemVersion, VulnScope 等字段。
    *   生成 `manifest.json`。

4.  **Signing (签名)**:
    *   读取 `manifest.json` 的内容。
    *   使用 **Private Key** 对内容进行 RSA 签名。
    *   生成 `signature.bin`。

5.  **Packaging (打包)**: 将工作目录下的所有内容（`manifest.json`, `signature.bin`, `assets/`）压缩为 `update_{version}_{type}.zip`。

6.  **Cleanup**: 打包完成后清理临时目录。

---

## 3. 配置文件示例 (`packager.yaml`)

为了简化重复操作，建议使用配置文件来描述打包规则。

```yaml
# 基础信息
app_name: "smart-security-system"
output_dir: "./dist"
private_key_path: "./keys/private.pem" # 严密保管

# 任务定义
tasks:
  # 任务A: 系统升级包
  - name: "system-release-v1.2"
    type: "SYSTEM"
    version: "1.2.0"
    min_system_version: "1.0.0"
    description: "季度功能迭代：增加AI分析模块"
    need_restart: true
    # 文件映射规则 (源路径 -> 包内路径)
    files:
      - src: "./bin/smart_linux_amd64"
        dst: "assets/bin/smart"
      - src: "./bin/decision_linux_amd64"
        dst: "assets/bin/decision"
      - src: "./web/dist/"
        dst: "assets/web/"
      - src: "./scripts/migrations/v1.2/"
        dst: "assets/sql/"

  # 任务B: 漏洞库升级包
  - name: "vuln-weekly-2023W42"
    type: "VULN"
    vuln_scope: "PRINCIPLE" # [新增] 漏洞范围: "PRINCIPLE", "FULL", "SCRIPT", "MIXED"
    version: "2023.10.24"
    min_system_version: "1.0.0"
    description: "更新 Log4j 检测规则"
    need_restart: false
    files:
      - src: "./data/scanner_core.db"
        dst: "assets/db/scanner.db"

### 3.1 元数据文件 (`manifest.json`)

`smart-packager` 会自动生成 `manifest.json`，其结构如下：

```json
{
  "type": "SYSTEM",             // 枚举: "SYSTEM" 或 "VULN"
  "vuln_scope": "PRINCIPLE",    // 漏洞范围: "PRINCIPLE", "FULL", "SCRIPT", "MIXED"
  "version": "1.2.0",           // 目标版本号
  "build_time": 1698123456,     // 构建时间戳
  "description": "修复了XX漏洞，增加了XX功能",
  "min_system_version": "1.0.0",// 兼容性检查：依赖的最低系统版本
  "need_restart": true,         // 是否需要重启服务
  "file_hash": {                // 关键文件指纹
    "payload/bin/smart": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }
}
```

```

---

## 4. 签名与验签机制详解

这是保障升级安全的核心。

### 4.1 密钥管理
*   **私钥 (Private Key)**: 仅存在于打包服务器（CI Server）或发布管理员的离线电脑中。**绝对不能打包进代码库或升级包中**。
*   **公钥 (Public Key)**: 硬编码在 `smart` 应用程序代码中，或者作为初始配置文件部署在客户服务器上。

### 4.2 签名生成逻辑 (Packager 端)

```go
// 伪代码
func SignManifest(manifestPath string, privateKeyPath string) error {
    // 1. 读取 Manifest 内容
    data, _ := ioutil.ReadFile(manifestPath)
    
    // 2. 计算 Hash (防止大文件处理慢，虽Manifest很小，但保持习惯)
    hashed := sha256.Sum256(data)
    
    // 3. RSA 签名
    privKey := LoadPrivateKey(privateKeyPath)
    signature, _ := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
    
    // 4. 写入文件
    return ioutil.WriteFile("signature.bin", signature, 0644)
}
```

### 4.3 验签逻辑 (Server 端 - 升级时执行)

```go
// 伪代码
func VerifyPackage(dir string) error {
    // 1. 读取包内的 Manifest 和 Signature
    manifestData, _ := ioutil.ReadFile(filepath.Join(dir, "manifest.json"))
    signature, _ := ioutil.ReadFile(filepath.Join(dir, "signature.bin"))
    
    // 2. 计算 Manifest Hash
    hashed := sha256.Sum256(manifestData)
    
    // 3. 使用内置公钥验签
    err := rsa.VerifyPKCS1v15(global.PublicKey, crypto.SHA256, hashed[:], signature)
    if err != nil {
        return errors.New("升级包签名无效，可能已被篡改！")
    }
    return nil
}
```

---

## 5. 开发建议

1.  **独立项目**: 建议在 `tools/packager` 下创建一个独立的 `main` 包。
2.  **CI 集成**: 在 GitLab CI / Jenkins 中，构建成功后直接调用 `smart-packager`，并将产物上传到制品库。
3.  **版本控制**: `packager.yaml` 应该纳入 Git 版本控制，确保每次发布的构建配置可追溯。
