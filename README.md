# TopMQ

中文 [README-CN.md](README-CN.md)

TopMQ is a high-performance MQTT 5.0 Broker built on [Mochi MQTT Server](https://github.com/mochi-mqtt/server), providing a complete Web UI management console with support for client authentication, authorization, blacklist management, real-time monitoring, and more.

## ✨ Features

### Core Features
- **MQTT 5.0 Protocol Support** - Full support for MQTT 5.0 protocol specification
- **Web UI Management Console** - Modern web interface for easy management and monitoring
- **Client Authentication** - Username/password authentication with configurable authentication rules
- **Client Authorization (ACL)** - Fine-grained topic subscription and publish permission control
- **Blacklist Management** - Support for IP and client ID blacklists to enhance security
- **Real-time Monitoring** - Real-time viewing of client connection status, subscription information, message statistics, etc.
- **Metrics Monitoring** - Support for storing metrics data to InfluxDB or OpenGemini, and Prometheus metrics export for integration with various monitoring systems
- **Client Management** - View and manage all connected clients
- **Subscription Management** - View and manage all active subscriptions

### Technical Features
- **SQLite Database** - Built-in SQLite, no additional database installation required, ready to use out of the box
- **Redis Caching and Persistence** - Use Redis for data caching and message persistence
- **WebSocket Support** - Support for WebSocket connections, convenient for web application integration
- **TLS/SSL Support** - Support for encrypted connections to ensure data transmission security
- **Multi-protocol Support** - Simultaneous support for TCP and WebSocket protocols
- **High Performance** - Built with Go language for excellent performance

## 🚀 Quick Start

### Requirements

- Go 1.25.0 or higher
- Redis (for caching and persistence)
- InfluxDB or OpenGemini (for metrics data storage, optional)
- Node.js 18+ and pnpm (for building frontend UI, optional)

### Installation Steps

1. **Clone the repository**
```bash
git clone <repository-url>
cd topmq
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure Redis**
Ensure Redis service is running, default connection address is `127.0.0.1:6379`

4. **Configure the application**
Copy and edit the configuration file:
```bash
cp dist/config.yml ./config.yml
```

Main configuration items:
- `mqttServer.port`: MQTT TCP port (default: 1883)
- `mqttServer.wsPort`: WebSocket port (default: 1882)
- `web.port`: Web UI port (default: 8080)
- `redis.addr`: Redis address (default: 127.0.0.1:6379)
- `db.databaseType`: Database type (default: sqlite)

5. **Run the application**
```bash
go run main.go
```

Or build and run:
```bash
go build -o topmq
./topmq
```

6. **Access Web UI**
Open your browser and visit `http://localhost:8080`. Default username: `admin`, password: `admin`.

### Frontend Development

If you need to modify the frontend UI:

```bash
cd ui
pnpm install
pnpm dev
```

Build frontend:
```bash
pnpm build
```

## 📖 Usage

### Default Account

The default administrator account will be automatically created on first startup (if `db.autoCreateAdmin` is `true`). Please check the configuration or log output for the default username and password.

### Client Authentication

Client authentication is managed through the database, supporting dynamic addition, modification, and deletion of user accounts through the Web UI. All user information is stored in the database for unified management and maintenance.

### Client Authorization

Support for configuring ACL rules to control client subscription and publish permissions for topics.

### Blacklist

Support for IP address and client ID blacklists. Clients added to the blacklist will be unable to connect.

### Monitoring Features

- **Real-time Monitoring** - View server running status, connection count, message statistics, etc.
- **Client List** - View information of all connected clients
- **Subscription Management** - View and manage all active subscriptions
- **Metrics Export** - Prometheus format metrics export (if `exporterPort` is configured)

## 🔧 Configuration

Main configuration items:

```yaml
mqttServer:
  port: 1883          # MQTT TCP port
  wsPort: 1882        # WebSocket port
  tls: false          # Enable TLS
  persist: true       # Enable message persistence

web:
  port: 8080          # Web UI port
  debug: false        # Debug mode
  sessionTimeout: 60  # Session timeout (minutes)

db:
  databaseName: topmq
  databasePath: ./db  # SQLite database file path

redis:
  addr: 127.0.0.1:6379  # Redis address
  keyPrefix: "topmq:"   # Redis key prefix

metrics:                # Metrics monitoring configuration (optional)
  enabled: true        # Enable metrics monitoring
  type: InfluxDB       # Metrics database type: InfluxDB or OpenGemini
  address: localhost:8086  # Database address
  username: ""         # Username (optional)
  password: ""         # Password (supports encrypted storage)
  database: topmq      # Database name
  retention: 3        # Metrics data retention days
```

For more configuration options, please refer to the `dist/config.yml` file.

## 🏗️ Project Structure

```
topmq/
├── accounts/          # Account management module
├── acls/             # ACL authorization module
├── auth/             # Authentication module
├── black_list/       # Blacklist module
├── cache/            # Redis cache module
├── config/           # Configuration management
├── db/               # Database abstraction layer
├── metrics/          # Metrics monitoring module
├── mqttserver/       # MQTT server core
├── web/              # Web server and routes
├── ui/               # Frontend UI code
│   ├── src/
│   │   ├── api/      # API interface definitions
│   │   ├── pages/    # Page components
│   │   └── ...
└── main.go           # Program entry point
```

## 📦 Dependencies

### Backend Main Dependencies
- `github.com/mochi-mqtt/server/v2` - MQTT server core
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM framework
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/glebarez/sqlite` - SQLite driver

### Frontend Main Dependencies
- `vue` - Vue 3 framework
- `element-plus` - UI component library
- `vite` - Build tool
- `vue-router` - Router management
- `pinia` - State management

## 🔒 Security Recommendations

1. **Change default password** - Change the default administrator password immediately after first login
2. **Enable TLS** - Enable TLS encryption for production environments
3. **Configure firewall** - Restrict access sources for MQTT ports
4. **Regular updates** - Keep software versions updated
5. **Use strong passwords** - Set strong passwords for all accounts

## 📊 Monitoring and Metrics

TopMQ provides complete metrics monitoring functionality, supporting storage of metrics data to InfluxDB or OpenGemini time-series databases.

### Metrics Monitoring Configuration

The built-in metrics monitoring feature depends on the `metrics` configuration. Metrics data storage requires InfluxDB or OpenGemini, which can be configured through the configuration file:

```yaml
metrics:
  enabled: true        # Enable metrics monitoring
  type: InfluxDB       # Metrics database type: InfluxDB or OpenGemini
  address: localhost:8086  # Database address
  username: ""         # Username (optional)
  password: ""         # Password (supports encrypted storage)
  database: topmq      # Database name
  retention: 3        # Metrics data retention days
```

### Prometheus Metrics Export

TopMQ also supports Prometheus metrics export. You can enable the metrics export service by configuring `mqttServer.exporterPort`.

Main metrics include:
- Client connection count
- Message send/receive rate
- Subscription count
- Retained message count
- Server uptime

## 🤝 Contributing

Welcome to submit Issues and Pull Requests!

## 📄 License

This project is licensed under the MIT License. For details, please see the [LICENSE.md](LICENSE.md) file.

## 🔗 Related Links

- [Mochi MQTT Server](https://github.com/mochi-mqtt/server)
- [MQTT 5.0 Specification](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html)

## 📞 Support

For questions or suggestions, please contact us through:
- [Submit Issue](https://github.com/iotopo/topmq/issues)
- [IoTopo](https://www.iotopo.com)

---

**TopMQ** - Simple, efficient, and feature-rich MQTT Broker solution

