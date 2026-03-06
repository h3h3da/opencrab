# OpenCrab 安全策略

OpenCrab 是安全优先的个人 AI 助手网关。本文档描述安全设计、威胁模型和部署建议。

## 安全设计原则

1. **默认安全**：未显式配置时，仅绑定 loopback，不暴露到网络
2. **显式暴露**：非 loopback 绑定必须配置认证（Token 或 Password）
3. **纵深防御**：速率限制、消息大小限制、安全响应头、Origin 校验
4. **最小依赖**：Go 标准库为主，依赖经 `govulncheck` 扫描

## 威胁模型

- **可信操作者**：与 OpenClaw 一致，单用户可信模型
- **不可信输入**：所有外部输入（WebSocket、HTTP）视为不可信
- **部署边界**：默认不暴露到公网；暴露时需认证 + 可选 TLS

## 安全措施

| 措施 | 说明 |
|------|------|
| Loopback 默认 | `gateway.bind` 默认 `127.0.0.1` |
| 强制认证 | 非 loopback 时 `AuthToken` 或 `AuthPassword` 必填 |
| 速率限制 | 按 IP 限流，默认 60 req/min |
| WebSocket 消息限制 | 最大 1MB，防内存耗尽 |
| 安全响应头 | X-Content-Type-Options, X-Frame-Options, CSP 等 |
| Origin 校验 | WebSocket 仅接受 loopback 来源 |
| 常数时间比较 | Token 比较防时序攻击 |

## 部署建议

### 本地开发
```bash
opencrab gateway
# 默认 127.0.0.1:18789，无需认证
```

### 远程访问（推荐：SSH 隧道）
```bash
# 本地运行 gateway，通过 SSH 隧道暴露
ssh -L 18789:127.0.0.1:18789 user@remote-host
```

### 直接暴露（需认证 + 建议 TLS）
```bash
export OPENCRAB_AUTH_TOKEN=$(openssl rand -hex 32)
export OPENCRAB_BIND_ALL=1
opencrab gateway
```

生产环境建议配置 TLS：
```go
cfg.TLSCertFile = "/path/to/cert.pem"
cfg.TLSKeyFile  = "/path/to/key.pem"
```

## 漏洞报告

如发现安全问题，请通过 GitHub Security Advisories 或邮件私下报告，勿在公开 issue 中披露。

## 依赖扫描

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

建议在 CI 中集成 `govulncheck`。
