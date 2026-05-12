# Neo4j 界面免密跳转设计方案

## 1. 需求背景
在 Smart 系统前端（Vue）中，增加一个入口，点击后可直接跳转到 Neo4j Browser 界面，并实现免密登录（SSO）。

## 2. 方案选择

已选择 **方案 B: ConnectURL 参数** 并完成实现。

| 特性 | 方案 B: ConnectURL 参数 (已选) |
| :--- | :--- |
| **安全性** | **中** (密码暴露在 URL/前端) |
| **SSO 机制** | URL 携带凭证 |
| **部署要求** | 无需 Nginx 修改，需配置环境变量 |
| **Neo4j 配置** | 保持开启 Auth |
| **适用场景** | 内网开发、测试、无法修改 Nginx |

---

## 3. 实现方案 (方案 B)

利用 Neo4j Browser 的 `connectURL` 参数自动填充密码。

### 3.1 架构
1.  **前端**: 调用后端接口 `/smart/system/neo4j/login`。
2.  **后端**: 拼接 `bolt://user:pass@host:port`，生成带 `connectURL` 的跳转链接。
3.  **前端**: 接收 URL 并 `window.open`。

### 3.2 后端接口设计 (已完成)
*   **路径**: `/smart/system/neo4j/login`
*   **方法**: `GET`
*   **配置**: (已硬编码)
    *   Host: `192.168.0.74`
    *   Browser Port: `7474`
    *   Bolt Port: `7687`
    *   User: `neo4j`
    *   Password: `4dogs.cn`

### 3.3 响应示例
```json
{
    "code": 200,
    "data": {
        "redirectUrl": "http://192.168.0.74:7474/browser/?connectURL=bolt%3A%2F%2Fneo4j%3A4dogs.cn%40192.168.0.74%3A7687"
    },
    "msg": "success"
}
```

### 3.4 前端调用示例 (Vue)
```javascript
// 获取跳转链接并打开
async function openNeo4j() {
  const { data } = await axios.get('/smart/system/neo4j/login');
  if (data.code === 200) {
     window.open(data.data.redirectUrl, '_blank');
  }
}
```

---
