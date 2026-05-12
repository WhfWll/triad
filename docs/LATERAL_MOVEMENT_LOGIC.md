
## 6. 配置增强设计 (Configuration Enhancement Design)

根据最新的 UI 设计，横向移动配置模块将支持更灵活的策略控制。

### 6.1 配置结构定义

前端传递的 `LateralMove` 配置对象应扩展如下字段（与 UI 选项一一对应）：

```json
{
  "lateral_move": {
    "isOpen": true,                // [UI: 开关] 是否开启横向移动
    "strategy": "custom_range",    // [UI: 横向策略] 目标选取策略
                                   // 可选值: 
                                   // "same_subnet" (同网段探测 - 默认)
                                   // "neighbor" (邻居发现 ARP/NetBIOS)
                                   // "exclude_current" (排除同网段 - 多网卡)
                                   // "custom_range" (自定义范围)
    "range": "192.168.1.0/24",     // [UI: 目标范围] 当 strategy="custom_range" 时生效
    "ports": "22,445,3389",        // [UI: 端口范围] 自定义扫描端口，支持逗号分隔
    "timeout": 600                 // [UI: 超时时间] 整个横向任务或单个阶段的超时时间(秒)
  }
}
```

### 6.2 后端逻辑适配

后端 `YakPostExploitationService.Handle` 需做如下调整以适配新配置：

1.  **扫描目标确定 (Target Determination)**:
    *   **Same Subnet (同网段)**: 若 `strategy` 为 `same_subnet`（或为空），自动计算当前节点的 C 段 (`TargetIP/24`)。
    *   **Neighbor Discovery (邻居发现)**: 若 `strategy` 为 `neighbor`，利用 ARP/NetBIOS 协议（执行 `arp -a` 等命令）发现活跃邻居，作为扫描目标。
    *   **Exclude Current (排除同网段)**: 若 `strategy` 为 `exclude_current`，获取本机所有网卡信息，排除当前入侵入口所在的网段，重点探测其他网段（如内网办公段、服务器段），适合双网卡/多网卡主机跳板攻击。
    *   **Custom Range (自定义范围)**: 若 `strategy` 为 `custom_range`，强制使用 `range` 字段的值。若 `range` 为空则回退到同网段并记录警告。

2.  **端口范围 (Port Scope)**:
    *   读取 `ports` 字段。
    *   若有值，直接传递给 Yak 的 `PortScanRequest.Ports`。
    *   若为空，使用系统默认的高危端口组 (`22,80,443,445,3389,8080,9919`)。

3.  **资源控制 (Resource Control)**:
    *   **Timeout**: 在调用 Yak gRPC 接口（如 `PortScan`, `Exec`）时，设置 Context 的超时时间为 `timeout` 秒。

### 6.3 交互逻辑说明

*   **联动**: 当“横向策略”选择“自定义范围”时，“目标范围”输入框变为必填/可用状态；否则置灰或隐藏。
*   **默认值**:
    *   策略: 默认 "自动C段"
    *   端口: 默认空 (使用后端内置列表)
    *   超时: 默认 600秒

### 6.4 字段与后端映射

- isOpen → [taskConfiguration.go:L298-L303](file:///d:/goproject/smart/tools/enums/taskConfiguration.go#L298-L303)，在入口处生效 [yak.go:L54-L59](file:///d:/goproject/smart/services/post_exploitation/yak.go#L54-L59)
- strategy → [taskConfiguration.go:L298-L303](file:///d:/goproject/smart/tools/enums/taskConfiguration.go#L298-L303) 新增；在范围选择中使用 [yak.go:L172-L183](file:///d:/goproject/smart/services/post_exploitation/yak.go#L172-L183)
- range → [taskConfiguration.go:L298-L303](file:///d:/goproject/smart/tools/enums/taskConfiguration.go#L298-L303)；用于设置扫描目标 [yak.go:L172-L183](file:///d:/goproject/smart/services/post_exploitation/yak.go#L172-L183)
- ports → [taskConfiguration.go:L298-L303](file:///d:/goproject/smart/tools/enums/taskConfiguration.go#L298-L303) 新增；用于端口设置 [yak.go:L184-L193](file:///d:/goproject/smart/services/post_exploitation/yak.go#L184-L193)
- timeout → [taskConfiguration.go:L298-L303](file:///d:/goproject/smart/tools/enums/taskConfiguration.go#L298-L303) 新增；用于设置上下文超时 [yak.go:L184-L193](file:///d:/goproject/smart/services/post_exploitation/yak.go#L184-L193) 与脚本执行
