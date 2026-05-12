# SSH 横向移动专用脚本设计方案 (Smart Probe)

## 1. 概述
本方案旨在设计一个轻量级、可独立运行的探测脚本（Smart Probe），专用于通过 SSH 通道下发至目标主机执行。该脚本承担“内网哨兵”的角色，负责在受控主机上进行本地环境感知、邻网资产发现及轻量级漏洞检测，并将结果以标准格式回传给 Smart Server。

**核心设计原则**：
*   **无依赖 (Dependency-Free)**: 采用静态编译的 Go 二进制文件或通用 Python/Shell 脚本，无需目标机安装额外库。
*   **低噪 (Stealthy)**: 仅在内存中运行或运行后立即自毁，不修改系统配置。
*   **标准化 (Standardized)**: 输入输出采用 JSON 格式，便于 Smart Server 自动化解析。

## 2. 脚本功能模块设计

该脚本 (`smart_probe`) 包含以下三个核心模块：

### 2.1 目标与存活探测 (Discovery Module)
*   **输入**: 扫描网段 (CIDR)，例如 `192.168.1.0/24`。
*   **功能**:
    *   **ARP 嗅探**: 读取 `/proc/net/arp` 或发送 ARP 包发现邻居。
    *   **ICMP Ping**: 对目标网段进行 Ping 扫描（需 Root 权限）或 UDP 探测。
    *   **TCP SYN/Connect**: 对常见端口（22, 80, 445）进行快速连接探测。
*   **输出**: 存活主机 IP 列表及对应的 MAC 地址。

### 2.2 端口与指纹识别 (Enumeration Module)
*   **输入**: 存活 IP 列表。
*   **功能**:
    *   **端口扫描**: 扫描 Top 100 常见高危端口。
    *   **Banner 抓取**: 获取服务 Banner 信息（如 SSH 版本, Web Server 类型）。
    *   **OS 识别**: 基于 TTL 或 Banner 推测操作系统类型 (Windows/Linux)。
*   **输出**: 开放端口及服务指纹信息。

### 2.3 漏洞检测 (Vulnerability Module)
*   **输入**: 开放端口及服务信息。
*   **功能**: 执行轻量级、低风险的 PoC 检测。
    *   **Weak Password**: 针对 SSH/MySQL/Redis 的 Top 5 弱口令尝试。
    *   **MS17-010**: SMB 协议版本检测（仅探测，不攻击）。
    *   **Web Logic**: 简单的 HTTP 请求探测。
    *   **Unauth Access**: Redis/MongoDB 未授权访问检测。
*   **输出**: 漏洞名称、风险等级、验证证据。

## 3. 运行流程 (Smart Server 视角)

1.  **上传 (Upload)**:
    *   Smart Server 通过 `sshutils` 的 SFTP 功能，将 `smart_probe` 上传至目标机临时目录（如 `/tmp/.X11-unix/sp`）。
2.  **授权 (Chmod)**:
    *   发送指令: `chmod +x /tmp/.X11-unix/sp`。
3.  **执行 (Execute)**:
    *   发送指令: `/tmp/.X11-unix/sp -target 192.168.1.0/24 -ports 22,80,445 -json`。
    *   *注*: 建议使用 `nohup` 或后台运行，防止 SSH 超时。
4.  **回传 (Retrieve)**:
    *   脚本执行完毕后将结果输出到标准输出 (Stdout) 或临时文件。
    *   Smart Server 读取输出内容并解析 JSON。
5.  **清理 (Cleanup)**:
    *   发送指令: `rm -f /tmp/.X11-unix/sp`。

## 4. 数据结构定义

### 4.1 输入参数 (CLI Args)
```bash
./smart_probe -target "192.168.1.0/24" -mode "all" -timeout 500
```

### 4.2 输出格式

#### 方式 A: 完整 JSON (Default)
适用于小规模扫描，脚本执行完毕后一次性输出所有结果。

```json
{
  "meta": {
    "start_time": "2023-10-27T10:00:00Z",
    "duration": "15s",
    "host_ip": "192.168.1.100"
  },
  "targets": [
    ... // 目标列表
  ]
}
```

#### 方式 B: 实时流式输出 (Stream/NDJSON)
适用于大规模扫描或需要实时反馈进度的场景。建议使用 **NDJSON (Newline Delimited JSON)** 格式，即每一行是一个独立的 JSON 对象。

Smart Server 端可以通过 SSH 的 `StdoutPipe` 实时读取每一行并解析，从而在前端动态展示扫描进度。

**输出流示例**:

```json
{"type": "progress", "data": {"stage": "discovery", "percent": 10, "message": "ARP scanning..."}}
{"type": "host_found", "data": {"ip": "192.168.1.101", "mac": "00:11:22:33:44:55", "os": "Windows"}}
{"type": "port_open", "data": {"ip": "192.168.1.101", "port": 445, "service": "microsoft-ds"}}
{"type": "vuln_detected", "data": {"ip": "192.168.1.101", "vuln_id": "MS17-010", "severity": "High"}}
{"type": "progress", "data": {"stage": "finished", "percent": 100, "message": "Scan complete"}}
```

**解析逻辑**:
1. Server 建立 SSH 连接并执行命令。
2. 开启一个 Goroutine 持续读取 `Stdout`。
3. 按换行符 `\n` 切分数据。
4. 解析 JSON 的 `type` 字段：
   - `host_found`: 立即在数据库创建/更新资产。
   - `vuln_detected`: 立即创建漏洞记录。
   - `progress`: 推送 WebSocket 消息给前端更新进度条。

## 5. 实现技术选型

### 方案 A: Golang (推荐)
*   **优点**: 编译为静态二进制文件，无依赖，跨平台（Linux/Windows/macOS），并发性能极强，适合快速扫描。
*   **缺点**: 文件体积稍大 (2MB-5MB)。
*   **实现**: 使用 `net` 库做扫描，`golang.org/x/net` 做协议交互。

### 方案 B: Python
*   **优点**: 脚本体积小，易于修改。
*   **缺点**: 依赖目标机 Python 环境（Python 2 vs 3 兼容性地狱），缺少第三方库（如 `requests` 可能不存在），只能用标准库实现。

### 方案 C: Shell/Bash
*   **优点**: 几乎所有 Linux 都有。
*   **缺点**: 功能极其有限，难以实现多线程和复杂协议解析，解析结果困难。

**结论**: 建议采用 **方案 A (Golang)** 开发 `smart_probe` 工具，Smart Server 根据目标机架构（x86/arm, linux/windows）自动选择对应的编译版本上传。

## 6. 与 Smart Scanner 及 Backend 的集成

### 6.1 架构流程
用户计划将此探测逻辑集成到 `scanner` 程序中，并由后端处理结果。完整数据流如下：

1.  **Scanner Service (执行层)**:
    *   通过 SSH 连接目标主机。
    *   上传并执行探测脚本（或直接运行命令）。
    *   实时读取脚本的 `Stdout` (NDJSON 流)。
    *   将 NDJSON 解析为内部结构体，并通过 gRPC Stream 回传给 Smart Backend。

2.  **Smart Backend (处理层)**:
    *   **入口**: `services/vulnerability_scanner.go` 中的结果处理逻辑。
    *   **分发**: 根据 JSON 中的 `type` 字段，调用不同的 Service 方法存储数据。

### 6.2 数据映射策略

后端需在 `processVulnerabilityResult` 或类似方法中新增对 Lateral Movement 结果的处理分支：

#### A. 资产发现 (type: host_found)
*   **动作**: 新增或更新 `TaskTarget` 表。
*   **代码参考**: `models.TaskTarget.Create`
*   **字段映射**:
    *   `data.ip` -> `task_target.target`
    *   `data.os` -> `task_target.os_type` (需转换为枚举值)
    *   `data.mac` -> `task_target.mac_address` (需新增字段或存入扩展字段)

#### B. 端口开放 (type: port_open)
*   **动作**: 更新 `TaskTarget` 的端口列表字段，或插入 `PortService` 表（如有）。
*   **代码参考**: `models.TaskTarget.Update`

#### C. 漏洞检出 (type: vuln_detected)
*   **动作**: 插入 `TaskVul` 表，并更新目标风险等级。
*   **代码参考**: `services.VulnerabilityScanner.handleVulnerabilityResult`
*   **逻辑复用**:
    *   可直接复用现有的 `handleVulnerabilityResult` 方法，将脚本返回的 `vuln_id` 映射到 `VulLibraries` 中的 POC Name。
    *   例如：脚本返回 `MS17-010` -> 查找系统漏洞库 -> 创建 `TaskVul` 记录。

#### D. 进度与日志 (type: progress)
*   **动作**: 写入任务日志，供前端展示。
*   **代码参考**: `services.TaskLogInfo.AddTaskLogInfo`

### 6.3 伪代码示例 (Backend)

```go
func HandleLateralStream(streamData []byte) {
    var event Event
    json.Unmarshal(streamData, &event)

    switch event.Type {
    case "host_found":
        // 发现新主机，添加到目标列表
        models.TaskTarget{
            Target: event.Data.IP,
            Source: "lateral_movement",
            ParentID: currentHostID, // 记录是谁发现了它（拓扑关系）
        }.Create()
        
    case "vuln_detected":
        // 发现漏洞，复用现有逻辑
        vulScanner.handleVulnerabilityResult(...)
        
    case "progress":
        // 记录日志
        taskLog.Add(event.Data.Message)
    }
}
```

## 7. 安全规避策略
*   **文件名混淆**: 随机化文件名，如 `kworker_u`。
*   **CPU 限制**: 脚本内部通过 `time.Sleep` 控制扫描速率，避免 CPU 飙升引起注意。
*   **内存执行**: 在 Linux 下尝试使用 `memfd_create` (需内核支持) 直接从内存加载 ELF 文件，避免落地。
