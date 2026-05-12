# 系统升级与漏洞库更新方案设计

## 1. 核心设计理念

本方案旨在解决当前系统升级逻辑中耦合度高、缺乏原子性保障以及无法灵活应对不同更新频率需求的问题。
核心思想概括为：**“统一入口，分策略执行，原子化回滚”**。

*   **分离关注点**：将“系统功能升级”（低频、高风险、需重启）与“漏洞库/规则库升级”（高频、低风险、热更新）彻底分离。
*   **安全优先**：引入强制签名校验，防止升级包被篡改。
*   **原子性保障**：升级过程要么全部成功，要么完全回滚，杜绝“中间状态”。

---

## 2. 升级包规范设计

所有的更新（无论是系统还是漏洞库）均使用统一的 `.zip` 格式，通过包内的元数据文件区分身份。

### 2.1 目录结构

```text
update.zip
├── manifest.json       # [核心] 升级包元数据
├── signature.bin       # [核心] 数字签名 (当前实现：仅对 manifest.json 验签)
└── assets/             # 资源目录
    ├── bin/            # (System) 二进制文件 (smart, decision)
    ├── web/            # (System) 前端静态资源
    ├── sql/            # (System) 数据库变更脚本
    ├── db/             # (Vuln) 漏洞库数据
    │   ├── {table}.json   # [推荐] 增量更新数据，每张表一个 JSON 文件 (e.g. vul_libraries.json)，紧凑格式无缩进
    │   └── scanner.db     # [遗留] 全量替换文件 (SQLite)
    └── scripts/        # 钩子脚本 (当前实现：未启用 pre_install/post_install 执行)
```

### 2.2 元数据定义 (`manifest.json`)

```json
{
  "type": "SYSTEM",             // 枚举: "SYSTEM" 或 "VULN"
  "vuln_scope": "PRINCIPLE",    // [新增] 漏洞范围: "PRINCIPLE", "FULL", "SCRIPT", "MIXED"
  "version": "1.2.0",           // 目标版本号
  "build_time": 1698123456,     // 构建时间戳
  "description": "修复了XX漏洞，增加了XX功能",
  "min_system_version": "1.0.0",// 兼容性检查：依赖的最低系统版本
  "need_restart": true,         // 是否需要重启服务
  "file_hash": {                // 关键文件指纹 (当前实现：未对 file_hash 做强制校验)
    "assets/bin/smart": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }
}
```

---

## 3. 核心流程架构 (Strategy Pattern)

后端采用**策略模式** (Strategy Pattern) 实现。定义统一的 `Upgrader` 接口，不同类型的升级包实现各自的逻辑。

### 3.1 接口定义

```go
type UpgradeStrategy interface {
    Verify(ctx) error   // 校验 (签名, 版本兼容性, 磁盘空间)
    Backup(ctx) path    // 备份 (返回备份路径)
    Apply(ctx) error    // 执行 (替换文件, 迁移数据)
    Rollback(ctx, path) // 回滚 (出错时恢复备份)
    PostAction(ctx)     // 后置 (重启服务, 通知重载)
}
```

### 3.2 通用执行流 (Manager)

1.  **Upload**: 上传 Zip 包至临时目录，并做基础校验（如扩展名）。
2.  **Prepare (预检)**: 解压并读取 `manifest.json`，对 `manifest.json` 执行验签，返回元信息供前端确认。
3.  **Execute (异步升级)**: 前端确认后触发异步升级任务，任务会全量解压，再按 `manifest.type` 选择策略执行 `Verify()` -> `Backup()` -> `Apply()`，并异步触发 `PostAction()`。
4.  **Rollback (手动)**: 当前实现不会在 `Apply()` 失败时自动回滚；需要调用手动回滚接口触发策略 `Rollback()`。

---

## 4. 场景详细策略

### 4.1 场景 A：系统升级 (System Upgrade)

*   **特点**: 涉及二进制替换，必须停止服务，风险较高。
*   **流程细节**:
    1.  **Verify**: 当前实现主要检查升级包中是否包含任一可更新的系统文件；暂未对 `min_system_version` 做强制校验。
    2.  **Backup**: 
        *   **策略**: **文件级全量备份**。
        *   **路径**: `/opt/laozhi/smart/backup/sys_{timestamp}` (例如 `sys_1704355200`)。
        *   **核心服务**: 
            *   `/opt/laozhi/smart/smart` (主程序)
            *   `/opt/laozhi/smart/config.yaml` (主配置)
            *   `/opt/laozhi/decision/decision` (决策引擎)
            *   `/opt/laozhi/decision/config.yaml` (决策配置)
        *   **前端资源**: 递归备份 `/opt/laozhi/nginx/html/` (Vue 静态资源)。
        *   **编排文件**: 备份 `/opt/laozhi/docker-compose.yml`。
        *   **数据目录**: (可选) `/opt/laozhi/data/`。
    3.  **Apply**:
        *   停止服务：当前实现未在 Apply 阶段显式停止 `smart/decision`（仅保留了 Linux 下停止服务的注释/占位逻辑）。
        *   替换二进制文件、配置文件和静态资源。
        *   执行 `assets/sql/` 下的数据库迁移脚本。
    4.  **Rollback**:
        *   触发方式：当前实现通过手动回滚接口触发回滚流程。
        *   将备份目录中的文件和文件夹递归覆盖回原路径。
        *   恢复文件执行权限 (chmod +x)。
        *   重启服务：当前实现仅在系统升级的 PostAction 阶段尝试重启 `smart` 服务进程。
    5.  **PostAction**:
        *   更新系统版本信息（写入 `map_set` 的 `systemVersion` 对象）。
        *   Linux 环境尝试重启 `smart` 服务（优先 supervisorctl，其次 supervisord，最后 pkill/kill -9 兜底）。
        *   当前实现未包含 `/api/health` 自检与失败自动回滚逻辑。

### 4.2 场景 B：漏洞库升级 (Vuln Upgrade)

*   **特点**: 更新 MySQL 漏洞库数据，无需停止服务，追求零停机。
*   **子策略 (`vuln_scope`)**:
    *   **PRINCIPLE**: 仅包含“原理验证”类高价值数据。
    *   **FULL**: 包含所有历史 CVE/CNVD 数据。
    *   **SCRIPT/MIXED**: 包含脚本库更新。
*   **流程细节**:
    1.  **Verify**: 检查 `assets/db/` 目录下数据文件的完整性。
    2.  **Backup**: 
        *   **策略**: **混合备份 (文件 + 逻辑导出)**。
        *   **路径**: `/opt/laozhi/smart/backup/vuln_{timestamp}`。
        *   **SQLite**: 复制 `/opt/laozhi/smart/scanner.db` 和 `user_custom.db` 到备份目录。
        *   **MySQL (逻辑备份)**: 
            *   查询当前数据库中即将被更新的关键表（`vul_libraries`, `finger`, `dictionary`）。
            *   **分片备份**: 针对 `vul_libraries` 这种大表，采用分页查询（Page Size 10000），将数据分片导出为多个 JSON 文件（如 `vul_libraries/part_0.json`, `vul_libraries/part_1.json`），避免一次性加载过多数据导致内存溢出 (OOM)。
            *   **按需备份**: 若升级包为 `PRINCIPLE` (原理验证) 类型，备份时仅导出 `verify_type=1` 的记录，大幅减少备份数据量；若为 `FULL` 类型，则全量导出。
            *   小表（如 `dictionary`）保持单文件导出。
            *   *理由*: 相比物理备份 MySQL 数据目录，逻辑备份更轻量且不影响数据库运行；分片处理解决了大表导出时的内存瓶颈。
    3.  **Apply**:
        *   **遍历加载**: 扫描 `assets/db/` 目录下的所有 `.json` 文件。
        *   **执行 (内存 Diff + 批量处理)**:
            *   读取 JSON 文件内容。
            *   加载数据库中现有数据的指纹（`VulID` 或 `Name`）到内存。
            *   比对差异，生成插入列表和更新列表。
            *   执行批量插入（Batch Insert）和批量更新。
    4.  **Rollback**:
        *   **触发方式**: 当前实现主要通过手动回滚接口触发 Rollback 流程。
        *   **SQLite**: 将备份的 `scanner.db` 和 `user_custom.db` 覆盖回原路径。
        *   **MySQL**: 读取备份的 JSON 文件（支持分片文件 `part_*.json`），将数据恢复到升级前的状态。
            *   *恢复策略*: 清空当前表数据，然后批量重新导入 JSON 备份数据。
        *   **日志**: 记录详细的回滚原因和状态，便于排查。
    5.  **PostAction**:
        *   无需重启主 Web 服务。
        *   通知扫描引擎重新加载规则库。

### 4.3 场景 C：攻击脚本库更新 (Script DB Upgrade)

*   **特点**: 涉及 `scanner.db` (SQLite) 文件的更新，该文件存储具体的 PoC/Exp 攻击脚本。
*   **注意**: 此场景在代码实现上归属于 `VULN` 类型，通过 `vuln_scope` 为 `SCRIPT` 或 `MIXED` 触发。
*   **数据分离策略**:
    *   **System DB (`scanner.db`)**: 官方维护的核心脚本库。更新策略为**全量替换**。
    *   **User DB (`user_custom.db`)**: 用户自定义添加的脚本库。更新策略为**完全保留** (Ignore)。
*   **流程细节**:
    1.  **Verify**: 校验升级包中 `assets/db/scanner.db` 的签名/Hash。
    2.  **Backup**: 同漏洞库升级策略，备份旧的 `scanner.db`。
    3.  **Apply**:
        *   直接使用新文件覆盖 `scanner.db`。
        *   *注意：系统启动时需配置扫描引擎同时加载 `scanner.db` 和 `user_custom.db`。*
    4.  **PostAction**:
        *   通知扫描引擎重载数据库连接 (Reload/Hot-Swap)。

---

## 5. 关键技术点与建议

1.  **数据更新策略**
    *   为避免覆盖用户自定义数据或已存在的记录，采用**增量插入**策略。
    *   升级包提供 JSON 格式的全量或增量数据，系统在应用时自动比对去重（基于 `VulID`），仅插入新数据。
    *   **JSON 优化**: 导出时采用 `json.Marshal` 紧凑格式，去除空格和换行，显著减小文件体积。

2.  **签名机制**
    *   开发环境：生成一对 RSA 密钥。
    *   打包工具 (Packager)：使用私钥对 Hash 签名，写入 `signature.bin`。
    *   生产环境：代码内置公钥 (`enums.SWPubKey`)；当前实现对 `manifest.json` 强制验签，未对 `assets/` 内容做签名/哈希强制校验。

3.  **回滚机制**
    *   这是系统稳定性的最后一道防线。必须确保 `Rollback` 函数本身是极其健壮的（例如简单的文件 Copy 还原），不依赖复杂的逻辑。

---

## 6. API 接口定义

后端提供 RESTful API 用于上传升级包和查询进度。

### 6.1 上传升级包

*   **URL**: `/system/uploadupgradefile`
*   **Method**: `POST`
*   **Content-Type**: `multipart/form-data`
*   **Params**:
    *   `file`: (File) .zip 升级包文件
*   **Response**:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": {
        "filename": "update_1.2.0_SYSTEM.zip",
        "version": "1.2.0",
        "type": "SYSTEM",
        "typeDesc": "系统升级包",
        "vulnScope": "PRINCIPLE",
        "scopeDesc": "快速更新包",
        "description": "修复了XX漏洞，增加了XX功能",
        "needRestart": true,
        "buildTime": 1698123456
      }
    }
    ```
*   **Logic**:
    1.  接收文件并保存到临时目录。
    2.  预检升级包：读取并解析 `manifest.json`，对 `manifest.json` 验签。
    3.  返回 `manifest` 信息供前端确认；上传阶段不会直接开始升级。

### 6.2 确认升级

*   **URL**: `/system/confirmupgrade`
*   **Method**: `POST`
*   **Request**:
    ```json
    {
      "filename": "update_1.2.0_SYSTEM.zip"
    }
    ```
*   **Response**:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": "升级任务已启动"
    }
    ```
*   **Logic**:
    1.  校验上传的 Zip 文件是否仍存在。
    2.  启动异步升级任务，不等待升级完成。

### 6.3 查询升级状态

*   **URL**: `/system/upgradestatus`
*   **Method**: `GET`
*   **Response**:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": {
        "state": "UPGRADING", // 枚举: IDLE, VERIFYING, BACKINGUP, UPGRADING, SUCCESS, FAILED
        "message": "正在应用更新...",
        "progress": 45,       // 进度百分比 (0-100)
        "error": "",          // 错误信息 (如果有)
        "version": "1.2.0",   // 目标版本
        "start_time": "2023-10-27T10:00:00Z",
        "end_time": "0001-01-01T00:00:00Z"
      }
    }
    ```

### 6.4 手动回滚

*   **URL**: `/system/manualrollback`
*   **Method**: `POST`
*   **Response**:
    ```json
    {
      "code": 200,
      "message": "回滚任务已开始，请查询进度",
      "data": null
    }
    ```
*   **Logic**:
    1.  检查是否存在有效的备份目录。
    2.  根据升级类型（SYSTEM/VULN）选择对应的回滚策略。
    3.  异步执行回滚操作，并更新全局状态。
