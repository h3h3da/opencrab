# 🦀 OpenCrab — Secure Personal AI Assistant

OpenCrab 是 [OpenClaw](https://github.com/openclaw/openclaw) 的 **Go 语言安全重构版本**，在保持个人 AI 助手核心能力的同时，强化部署后的安全性和抗外部攻击能力。

## 特性

- **安全优先**：默认 loopback 绑定、强制认证、速率限制、防 DoS
- **Go 实现**：单一二进制、少依赖、易部署
- **兼容设计**：Gateway 端口与协议与 OpenClaw 对齐，便于渐进迁移

## 快速开始

### 要求

- Go 1.22+

### 安装

```bash
git clone https://github.com/h3h3da/opencrab.git
cd opencrab
go build -o opencrab ./cmd/opencrab
```

### 运行 Gateway

```bash
./opencrab gateway
```

默认监听 `127.0.0.1:18789`，仅本地可访问。

### 暴露到网络时（需认证）

```bash
export OPENCRAB_AUTH_TOKEN=$(openssl rand -hex 32)
export OPENCRAB_BIND_ALL=1
./opencrab gateway
```

客户端需携带 `Authorization: Bearer <token>`。

## 安全

详见 [SECURITY.md](SECURITY.md)。

核心原则：

- 默认仅绑定 loopback
- 非 loopback 时必须配置认证
- 速率限制、消息大小限制、安全响应头
- 常数时间 token 比较，防时序攻击

## 与 OpenClaw 对比

详见 [COMPARISON.md](COMPARISON.md)。

## 开发

```bash
go build ./...
go test ./...
govulncheck ./...
```

## License

MIT — see [LICENSE](LICENSE)
