# gServ | 游戏服务端框架 #
[**中文**](./README.md) | [**English**](./README_EN.md)

gServ是一个轻量、高效、安全的通用游戏服务端框架，专为多人在线游戏设计，提供玩家管理、房间管理、数据存储和消息转发等核心功能。

## 内容导引 ##
* [介绍](#介绍)
* [功能特性](#功能特性)
* [架构设计](#架构设计)
* [快速开始](#快速开始)
* [配置说明](#配置说明)
* [部署指南](#部署指南)

## 介绍 ##
#### gServ是什么
gServ是一个通用的游戏服务端框架，旨在为多人在线游戏提供稳定可靠的后端支持。它采用Go语言开发，具有高性能、低延迟的特点，支持快速构建各种类型的多人在线游戏。

#### 核心功能
* **玩家管理**：玩家注册、登录、在线状态管理、数据存档
* **房间管理**：创建、加入、离开房间，房间锁定/解锁
* **数据存储**：玩家数据持久化存储，支持SQLite和MySQL
* **消息转发**：TCP长连接支持，实时消息转发
* **验证系统**：邮箱验证码验证，JWT身份认证

#### 浊水楼台免费服务
* **中国-成都**：chengdu-gserv.zslt-official.com

## 功能特性 ##
#### 玩家系统
- 玩家注册与登录（邮箱+密码）
- 玩家数据存档与恢复
- 在线玩家状态管理
- 玩家封禁与解封功能

#### 房间系统
- 创建游戏房间（支持设置最大人数）
- 加入/离开房间
- 房间锁定与解锁
- 房间自动清理（空闲5分钟自动删除）
- 房主权限管理

#### 游戏管理
- 游戏实例创建与管理
- 多游戏支持（可同时运行多个游戏服务）
- 游戏数据隔离存储

#### 网络通信
- HTTP RESTful API（玩家管理、房间操作）
- TCP长连接服务（实时消息转发）
- 自定义通信协议支持

#### 安全特性
- JWT身份验证
- 邮箱验证码系统
- 密码哈希存储（bcrypt）
- CORS跨域支持
- 管理员鉴权机制

## 架构设计 ##
```
gServ/
├── core/           # 核心模块
│   ├── config/     # 配置管理
│   ├── gameserv/   # 游戏服务核心
│   ├── httpserv/   # HTTP服务
│   ├── log/        # 日志系统
│   ├── repository/ # 数据仓库
│   ├── tcpserv/    # TCP服务
│   └── validate/   # 数据验证
├── pkg/            # 公共包
│   ├── gserv/      # 游戏服务模型
│   ├── hash/       # 哈希工具
│   ├── jwt/        # JWT工具
│   ├── middleware/ # 中间件
│   └── model/      # 数据模型
└── main.go         # 程序入口
```

## 快速开始 ##
#### 环境要求
- Go 1.24.11 或更高版本
- SQLite 或 MySQL 数据库

#### 安装步骤
1. 克隆项目
```bash
git clone https://github.com/ZSLTChenXiYin/gServ.git
cd gServ
```

2. 安装依赖
```bash
go mod download
```

3. 配置服务
复制示例配置文件并修改：
```bash
cp example.gserv.conf.yaml gserv.conf.yaml
```

编辑 `gserv.conf.yaml`，配置数据库和服务器参数。

4. 启动服务
```bash
go run main.go
```

服务启动后：
- HTTP服务监听在 `http_port`
- TCP服务监听在 `tcp_port`

## 配置说明 ##
#### 服务器配置
```yaml
server:
  mode: "dev"                          # 运行模式: "prod"/"dev"
  http_port: 8080                      # HTTP服务端口
  tcp_port: 9090                       # TCP服务端口
  log: "gserv.log"                     # 日志文件路径
  jwt: "your-jwt-secret"               # JWT密钥
  auth_code: "admin-auth-code"         # 管理员鉴权码
  email:                               # 邮件服务配置
     template: "res/captcha/email_captcha.html"
     host: "smtp.qq.com"
     port: 465
     email: "your-email@qq.com"
     password: "your-email-password"
```

#### 数据库配置
```yaml
database:
  driver: "sqlite"      # 数据库驱动: "sqlite"/"mysql"
  dsn: "gserv.db"       # 数据库连接字符串
  # MySQL示例: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

## 部署指南 ##
#### 本地构建
```bash
# 1. 配置环境
cp example.gserv.conf.yaml gserv.conf.yaml
# 编辑配置文件

# 2. 构建可执行文件
go build -ldflags "-s -w"

# 3. 启动服务
./gServ
```

#### Docker 部署
```bash
# 1. 配置环境
cp docker.env .env
vim .env
vim docker.gserv.conf.yaml
# 编辑配置文件

# 2. 启动服务
docker compose up -d
```
- [Docker 部署参考文档](README_DOCKER.md)

## 开发说明 ##
#### 代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 错误处理使用标准error类型
- 日志分级：Debug、Info、Warn、Error

#### 扩展开发
1. 添加新的API接口
   - 在 `core/httpserv/` 下创建新的controller
   - 在 `core/httpserv/init.go` 中注册路由

2. 添加新的游戏逻辑
   - 在 `core/gameserv/` 下扩展功能
   - 在 `pkg/gserv/` 下定义数据结构

3. 自定义通信协议
   - 修改 `core/tcpserv/protocol.go`
   - 实现自定义的消息处理逻辑

## 许可证 ##
本项目采用MIT许可证。详见LICENSE文件。

## 联系方式 ##
如有问题或建议，请通过以下方式联系：
- 邮箱：imjfoy@163.com
- GitHub Issues：[项目地址](https://github.com/ZSLTChenXiYin/gServ/issues)

---
**gServ - 让游戏开发更简单**