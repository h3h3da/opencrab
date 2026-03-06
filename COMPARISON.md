# OpenCrab vs OpenClaw 对比

OpenCrab 是 OpenClaw 的 Go 语言安全重构版本，专注于提升部署后的安全性和抗外部攻击能力。

## 项目概览

| 维度 | OpenClaw | OpenCrab |
|------|----------|----------|
| **语言** | TypeScript (Node.js ≥22) | Go 1.22+ |
| **定位** | 个人 AI 助手，多通道集成 | 同定位，安全优先 |
| **来源** | [openclaw/openclaw](https://github.com/openclaw/openclaw) | 本仓库 (opencrab) |

## 架构对比

### OpenClaw 架构
```
WhatsApp / Telegram / Slack / Discord / ... (20+ 通道)
                    │
                    ▼
┌─────────────────────────────────────┐
│            Gateway (Node.js)         │
│         ws://127.0.0.1:18789        │
└──────────────────┬──────────────────┘
                   │
    ├─ Pi agent (RPC, TypeScript)
    ├─ CLI (openclaw …)
    ├─ WebChat UI
    └─ macOS / iOS / Android 节点
```

### OpenCrab 架构
```
通道 (逐步迁移) → Gateway (Go)
                    │
                    ▼
┌─────────────────────────────────────┐
│         Gateway (Go, 安全加固)        │
│         ws://127.0.0.1:18789         │
│  - 默认 loopback 绑定                 │
│  - 强制认证 (非 loopback 时)           │
│  - 速率限制 / 防 DoS                  │
└──────────────────┬──────────────────┘
                   │
    ├─ CLI (opencrab …)
    └─ WebSocket 协议 (兼容设计)
```

## 安全性对比

### OpenClaw 安全现状
- **DM 配对策略**：未知发送者需配对码
- **沙箱**：`agents.defaults.sandbox.mode` 默认 `off`，exec 在主机运行
- **Web 界面**：建议仅本地使用，`gateway.bind="loopback"` 为默认
- **信任模型**：单用户可信操作者，非多租户
- **依赖**：大量 npm 依赖，需持续关注漏洞

### OpenCrab 安全增强

| 安全措施 | OpenClaw | OpenCrab |
|----------|----------|----------|
| **默认绑定** | loopback | loopback（强制，非 loopback 需显式配置） |
| **非 loopback 认证** | 可选 | **强制**：必须配置 AuthToken 或 AuthPassword |
| **速率限制** | 无内置 | **内置**：按 IP 限流，防 DoS |
| **WebSocket 消息大小** | 可配置 | **硬限制** 1MB，防内存耗尽 |
| **安全响应头** | 部分 | **完整**：X-Content-Type-Options, X-Frame-Options, CSP 等 |
| **Origin 校验** | 可配置 | **仅允许 loopback** 来源 |
| **认证比较** | 普通字符串比较 | **常数时间比较**，防时序攻击 |
| **TLS** | 可选 | 支持，推荐生产环境 |
| **依赖** | npm 生态 | Go 标准库 + 少量依赖，`govulncheck` 扫描 |

## 功能对比

### 当前 OpenCrab 实现状态

| 功能 | OpenClaw | OpenCrab |
|------|----------|----------|
| Gateway HTTP/WS | ✅ | ✅ 基础实现 |
| 多通道 (WhatsApp/Telegram/...) | ✅ 20+ | 🔲 计划中 |
| Pi Agent RPC | ✅ | 🔲 计划中 |
| CLI (gateway, agent, send) | ✅ | ✅ gateway, version |
| WebChat UI | ✅ | 🔲 计划中 |
| 会话管理 | ✅ | 🔲 计划中 |
| 沙箱 / Docker 执行 | ✅ | 🔲 计划中 |
| 技能平台 | ✅ | 🔲 计划中 |
| macOS/iOS/Android 应用 | ✅ | 🔲 计划中 |

## 部署安全建议

### OpenClaw
- 保持 `gateway.bind="loopback"`
- 远程访问用 SSH 隧道或 Tailscale
- 运行 `openclaw doctor` 检查配置

### OpenCrab
- **默认即安全**：loopback + 无认证时仅本地可访问
- 暴露到网络时**必须**设置 `OPENCRAB_AUTH_TOKEN`
- 生产环境建议启用 TLS
- 使用 `govulncheck` 定期扫描依赖

## 为什么用 Go 重构？

1. **内存安全**：Go 无 GC 外的内存管理错误，减少缓冲区溢出等风险
2. **单一二进制**：无 npm 依赖树，减少供应链攻击面
3. **强类型**：编译期捕获更多错误
4. **并发模型**：goroutine 比 Node 回调更易推理，减少竞态
5. **部署简单**：`go build` 产出静态二进制，易于容器化

## 路线图

1. **Phase 1**（当前）：Gateway 核心 + 安全中间件 ✅
2. **Phase 2**：WebSocket 协议兼容、会话管理
3. **Phase 3**：通道适配器（Telegram、Slack 等）
4. **Phase 4**：Agent 集成、工具调用
5. **Phase 5**：完整功能对等

## 参考

- [OpenClaw GitHub](https://github.com/openclaw/openclaw)
- [OpenClaw Security](https://github.com/openclaw/openclaw/blob/main/SECURITY.md)
- [Go Security Best Practices](https://go.dev/doc/security/best-practices)
- [OWASP Go Secure Coding](https://owasp.org/www-project-go-secure-coding-practices-guide)
