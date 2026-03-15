# gServ | Game Server Framework #
[**中文**](./README.md) | [**English**](./README_EN.md)

gServ is a lightweight, efficient, and secure universal game server framework designed for multiplayer online games, providing core functionalities such as player management, room management, data storage, and message forwarding.

## Table of Contents ##
* [Introduction](#introduction)
* [Features](#features)
* [Architecture](#architecture)
* [Quick Start](#quick-start)
* [Configuration](#configuration)
* [Deployment](#deployment)

## Introduction ##
#### What is gServ?
gServ is a universal game server framework designed to provide stable and reliable backend support for multiplayer online games. Developed in Go, it offers high performance and low latency, supporting rapid development of various types of multiplayer online games.

#### Core Capabilities
* **Player Management**: Player registration, login, online status management, data archiving
* **Room Management**: Create, join, leave rooms, room locking/unlocking
* **Data Storage**: Persistent player data storage with SQLite and MySQL support
* **Message Forwarding**: TCP long-connection support for real-time message forwarding
* **Verification System**: Email verification codes, JWT authentication

#### 浊水楼台 Free Services
* **中国-成都**：chengdu-gserv.zslt-official.com

## Features ##
#### Player System
- Player registration and login (email + password)
- Player data archiving and restoration
- Online player status management
- Player banning and unbanning functionality

#### Room System
- Create game rooms (supports setting maximum players)
- Join/leave rooms
- Room locking and unlocking
- Automatic room cleanup (idle rooms deleted after 5 minutes)
- Room owner permission management

#### Game Management
- Game instance creation and management
- Multi-game support (can run multiple game services simultaneously)
- Game data isolation storage

#### Network Communication
- HTTP RESTful API (player management, room operations)
- TCP long-connection service (real-time message forwarding)
- Custom communication protocol support

#### Security Features
- JWT authentication
- Email verification code system
- Password hash storage (bcrypt)
- CORS cross-origin support
- Administrator authentication mechanism

## Architecture ##
```
gServ/
├── core/           # Core modules
│   ├── config/     # Configuration management
│   ├── gameserv/   # Game service core
│   ├── httpserv/   # HTTP service
│   ├── log/        # Logging system
│   ├── repository/ # Data repository
│   ├── tcpserv/    # TCP service
│   └── validate/   # Data validation
├── pkg/            # Common packages
│   ├── gserv/      # Game service models
│   ├── hash/       # Hash utilities
│   ├── jwt/        # JWT utilities
│   ├── middleware/ # Middleware
│   └── model/      # Data models
└── main.go         # Program entry point
```

## Quick Start ##
#### Requirements
- Go 1.24.11 or higher
- SQLite or MySQL database

#### Installation Steps
1. Clone the project
```bash
git clone https://github.com/ZSLTChenXiYin/gServ.git
cd gServ
```

2. Install dependencies
```bash
go mod download
```

3. Configure the service
Copy the example configuration file and modify it:
```bash
cp example.gserv.conf.yaml gserv.conf.yaml
```

Edit `gserv.conf.yaml` to configure database and server parameters.

4. Start the service
```bash
go run main.go
```

After service starts:
- HTTP service listens on `http_port`
- TCP service listens on `tcp_port`

## Configuration ##
#### Server Configuration
```yaml
server:
  mode: "dev"                          # Run mode: "prod"/"dev"
  http_port: 8080                      # HTTP service port
  tcp_port: 9090                       # TCP service port
  log: "gserv.log"                     # Log file path
  jwt: "your-jwt-secret"               # JWT secret key
  auth_code: "admin-auth-code"         # Administrator authentication code
  email:                               # Email service configuration
    host: "smtp.qq.com"
    port: 465
    email: "your-email@qq.com"
    password: "your-email-password"
```

#### Database Configuration
```yaml
database:
  driver: "sqlite"      # Database driver: "sqlite"/"mysql"
  dsn: "gserv.db"       # Database connection string
  # MySQL example: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

## Deployment ##
#### Local Build
```bash
# 1. Configure environment
cp example.gserv.conf.yaml gserv.conf.yaml
# Edit configuration file

# 2. Build executable
go build -ldflags "-s -w"

# 3. Start service
./gServ
```

#### Docker Deployment
```bash
# 1. Configure environment
cp docker.env .env
vim .env
vim docker.gserv.conf.yaml
# Edit configuration file

# 2. Start service
docker compose up -d
```
- [Docker Deployment Reference Document](README_DOCKER.md)

## Development Guide ##
#### Code Standards
- Follow Go official coding standards
- Use gofmt for code formatting
- Use standard error type for error handling
- Log levels: Debug, Info, Warn, Error

#### Extending Functionality
1. Adding new API endpoints
   - Create new controller in `core/httpserv/`
   - Register routes in `core/httpserv/init.go`

2. Adding new game logic
   - Extend functionality in `core/gameserv/`
   - Define data structures in `pkg/gserv/`

3. Custom communication protocols
   - Modify `core/tcpserv/protocol.go`
   - Implement custom message handling logic

## License ##
This project is licensed under the MIT License. See LICENSE file for details.

## Contact ##
For questions or suggestions, please contact:
- Email: imjfoy@163.com
- GitHub Issues: [project-url](https://github.com/ZSLTChenXiYin/gServ/issues)

---
**gServ - Making Game Development Simpler**