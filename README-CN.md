# TopMQ

[English](https://github.com/iotopo/topmq/blog/main/README.md) [简体中文](https://github.com/iotopo/topmq/blog/main/README-CN.md)

TopMQ 是一个基于 [Mochi MQTT Server](https://github.com/mochi-mqtt/server) 实现的高性能 MQTT 5.0 Broker，提供了完整的 Web UI 管理控制台，支持客户端认证、授权、黑名单管理、实时监控等功能。

## ✨ 特性

### 核心功能
- **MQTT 5.0 协议支持** - 完整支持 MQTT 5.0 协议规范
- **Web UI 管理控制台** - 现代化的 Web 界面，方便管理和监控
- **客户端认证** - 支持用户名密码认证，可配置认证规则
- **客户端授权（ACL）** - 细粒度的主题订阅和发布权限控制
- **黑名单管理** - 支持 IP 和客户端 ID 黑名单，增强安全性
- **实时监控** - 实时查看客户端连接状态、订阅信息、消息统计等
- **指标监控** - 支持将指标数据存储到 InfluxDB 或 OpenGemini，同时支持 Prometheus 指标导出，可集成到各种监控系统
- **客户端管理** - 查看和管理所有连接的客户端
- **订阅管理** - 查看和管理所有活跃的订阅关系

### 技术特性
- **SQLite 数据库** - 内置 SQLite，无需额外安装数据库，开箱即用
- **Redis 缓存和持久化** - 使用 Redis 进行数据缓存和消息持久化
- **WebSocket 支持** - 支持 WebSocket 连接，方便 Web 应用集成
- **TLS/SSL 支持** - 支持加密连接，保障数据传输安全
- **多协议支持** - 同时支持 TCP 和 WebSocket 协议
- **高性能** - 基于 Go 语言实现，性能优异

## 🚀 快速开始

### 环境要求

- Go 1.25.0 或更高版本
- Redis（用于缓存和持久化）
- InfluxDB 或 OpenGemini（用于指标数据存储，可选）
- Node.js 18+ 和 pnpm（用于构建前端 UI，可选）

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd topmq
```

2. **安装依赖**
```bash
go mod download
```

3. **配置 Redis**
确保 Redis 服务正在运行，默认连接地址为 `127.0.0.1:6379`

4. **配置应用**
复制并编辑配置文件：
```bash
cp dist/config.yml ./config.yml
```

主要配置项：
- `mqttServer.port`: MQTT TCP 端口（默认：1883）
- `mqttServer.wsPort`: WebSocket 端口（默认：1882）
- `web.port`: Web UI 端口（默认：8080）
- `redis.addr`: Redis 地址（默认：127.0.0.1:6379）
- `db.databaseType`: 数据库类型（默认：sqlite）

5. **运行应用**
```bash
go run main.go
```

或者编译后运行：
```bash
go build -o topmq
./topmq
```

6. **访问 Web UI**
打开浏览器访问 `http://localhost:8080`。默认用户名：`admin`，密码：`admin`。

### 前端开发

如果需要修改前端 UI：

```bash
cd ui
pnpm install
pnpm dev
```

构建前端：
```bash
pnpm build
```

## 📖 使用说明

### 默认账户

首次启动会自动创建默认管理员账户（如果 `db.autoCreateAdmin` 为 `true`），默认用户名和密码请查看配置或日志输出。

### 客户端认证

客户端认证通过数据库进行管理，支持通过 Web UI 动态添加、修改和删除用户账户。所有用户信息存储在数据库中，便于统一管理和维护。

### 客户端授权

支持配置 ACL 规则，控制客户端对主题的订阅和发布权限。

### 黑名单

支持 IP 地址和客户端 ID 黑名单，被加入黑名单的客户端将无法连接。

### 监控功能

- **实时监控** - 查看服务器运行状态、连接数、消息统计等
- **客户端列表** - 查看所有连接的客户端信息
- **订阅管理** - 查看和管理所有活跃的订阅
- **指标导出** - Prometheus 格式的指标导出（如果配置了 `exporterPort`）

## 🔧 配置说明

主要配置项说明：

```yaml
mqttServer:
  port: 1883          # MQTT TCP 端口
  wsPort: 1882        # WebSocket 端口
  tls: false          # 是否启用 TLS
  persist: true       # 是否启用消息持久化

web:
  port: 8080          # Web UI 端口
  debug: false        # 调试模式
  sessionTimeout: 60  # 会话超时时间（分钟）

db:
  databaseName: topmq
  databasePath: ./db  # SQLite 数据库文件路径

redis:
  addr: 127.0.0.1:6379  # Redis 地址
  keyPrefix: "topmq:"   # Redis key 前缀

metrics:                # 指标监测配置（可选）
  enabled: true        # 是否启用指标监测
  type: InfluxDB       # 指标数据库类型：InfluxDB 或 OpenGemini
  address: localhost:8086  # 数据库地址
  username: ""         # 用户名（可选）
  password: ""         # 密码（支持加密保存）
  database: topmq      # 数据库名称
  retention: 3        # 指标数据保存天数
```

更多配置选项请参考 `dist/config.yml` 文件。

## 🏗️ 项目结构

```
topmq/
├── accounts/          # 账户管理模块
├── acls/             # ACL 授权模块
├── auth/             # 认证模块
├── black_list/       # 黑名单模块
├── cache/            # Redis 缓存模块
├── config/           # 配置管理
├── db/               # 数据库抽象层
├── metrics/          # 指标监控模块
├── mqttserver/       # MQTT 服务器核心
├── web/              # Web 服务器和路由
├── ui/               # 前端 UI 代码
│   ├── src/
│   │   ├── api/      # API 接口定义
│   │   ├── pages/    # 页面组件
│   │   └── ...
└── main.go           # 程序入口
```

## 📦 依赖

### 后端主要依赖
- `github.com/mochi-mqtt/server/v2` - MQTT 服务器核心
- `github.com/gin-gonic/gin` - Web 框架
- `gorm.io/gorm` - ORM 框架
- `github.com/redis/go-redis/v9` - Redis 客户端
- `github.com/glebarez/sqlite` - SQLite 驱动

### 前端主要依赖
- `vue` - Vue 3 框架
- `element-plus` - UI 组件库
- `vite` - 构建工具
- `vue-router` - 路由管理
- `pinia` - 状态管理

## 🔒 安全建议

1. **修改默认密码** - 首次登录后立即修改默认管理员密码
2. **启用 TLS** - 生产环境建议启用 TLS 加密
3. **配置防火墙** - 限制 MQTT 端口的访问来源
4. **定期更新** - 保持软件版本更新
5. **使用强密码** - 为所有账户设置强密码

## 📊 监控和指标

TopMQ 提供了完整的指标监测功能，支持将指标数据存储到 InfluxDB 或 OpenGemini 时序数据库中。

### 指标监测配置

内置的指标监测功能依赖 `metrics` 配置。指标数据存储需要依赖 InfluxDB 或 OpenGemini，可通过配置文件进行配置：

```yaml
metrics:
  enabled: true        # 是否启用指标监测
  type: InfluxDB       # 指标数据库类型：InfluxDB 或 OpenGemini
  address: localhost:8086  # 数据库地址
  username: ""         # 用户名（可选）
  password: ""         # 密码（支持加密保存）
  database: topmq      # 数据库名称
  retention: 3        # 指标数据保存天数
```

### Prometheus 指标导出

TopMQ 还支持 Prometheus 指标导出，可以通过配置 `mqttServer.exporterPort` 启用指标导出服务。

主要指标包括：
- 客户端连接数
- 消息发送/接收速率
- 订阅数量
- 保留消息数量
- 服务器运行时间

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](https://github.com/iotopo/topmq/blob/main/LICENSE) 文件。

## 🔗 相关链接

- [Mochi MQTT Server](https://github.com/mochi-mqtt/server)
- [MQTT 5.0 规范](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html)

## 📞 支持

如有问题或建议，请通过以下方式联系：
- [提交 Issue](https://github.com/iotopo/topmq/issues)
- [图扑物联](https://www.iotopo.com)

---

**TopMQ** - 简单、高效、功能完善的 MQTT Broker 解决方案

