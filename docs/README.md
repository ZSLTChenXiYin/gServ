# gServ - 游戏服务端

## 概述

gServ（全称游戏服务端，英文名 gameServer）是一个简单通用的游戏服务端，主要功能包括：

1. **转发游戏客户端数据** - 支持房间内的广播、单播、组播通信
2. **在线存储客户端数据** - 提供玩家数据存储和管理
3. **玩家连接管理** - 处理玩家登录、登出和连接保持
4. **房间管理** - 创建、加入、清理游戏房间

## 参考文档
- [HTTP API 参考文档](HTTP_API.md)
- [TCP API 参考文档](TCP_API.md)

## 使用浊水楼台服务
1. 发送注册邮件到 `imjfoy@163.com`，邮件注明游戏名称和用途，一个工作日内会收到回复邮件
2. 获取服务地址和认证信息
3. 直接调用远程 API

## 客户端开发流程

### 第 1 步：游戏注册

游戏开发者需要先注册游戏才能使用服务：

#### 选项 A：使用浊水楼台免费服务
1. 向 `imjfoy@163.com` 发送游戏注册说明邮件
2. 邮件内容需包含游戏基本信息
3. 一个工作日内完成注册

#### 选项 B：本地部署
1. 调用管理员接口添加游戏
2. 使用 `auth_code` 进行管理员认证
3. 通过 HTTP API 注册游戏

### 第 2 步：玩家注册流程

```
│ 获取验证码 │────▶│ 用户注册 │────▶│ 用户登录 │────▶│ 获取TCP │
```

#### 2.1 [获取验证码](HTTP_API.md#发送邮箱验证码)
```http
POST /captcha
Content-Type: application/json

{
  "email": "player@example.com"
  // 其他字段
}
```

#### 2.2 [用户注册](HTTP_API.md#玩家注册)
```http
POST /register
Content-Type: application/json

{
  "email": "player@example.com",
  "password": "secure_password",
  "captcha": "123456"
  // 其他字段
}
```

#### 2.3 [用户登录](HTTP_API.md#玩家登录)
```http
POST /login
Content-Type: application/json

{
  "email": "player@example.com",
  "password": "secure_password"
}
```

### 第 3 步：TCP 长连接

#### 3.1 建立 TCP 连接
以下是Go语言的代码示例：
```go
conn, err := net.Dial("tcp", "localhost:9090")
```

#### 3.2 [TCP 登录认证](TCP_API.md#3-玩家登录认证-AUTH_PLAYER_LOGIN)
发送 `AUTH_PLAYER_LOGIN` 协议包：
- ProtocolVersion: 1
- ProtocolType: 3 (AUTH_PLAYER_LOGIN)
- Token: JWT Token（从登录接口获取）
- GameID: 游戏ID

#### 3.3 加入房间
1. 通过 HTTP 接口[创建房间](HTTP_API.md#创建房间)：
```http
POST /room/create
Authorization: Bearer {token}
Content-Type: application/json

{
  "game_id": 1,
  "room_name": "游戏房间1",
  "max_players": 4
}
```

2. 通过 TCP [加入房间](TCP_API.md#5-加入房间-JOIN_ROOM)：
发送 `JOIN_ROOM` 协议包：
- ProtocolVersion: 1
- ProtocolType: 5 (JOIN_ROOM)
- Token: JWT Token
- GameID: 游戏ID
- RoomID: 房间ID

### 第 4 步：游戏内通信

#### 4.1 [广播数据](TCP_API.md#6-房间广播数据-ROOM_BROADCAST_DATA) (ROOM_BROADCAST_DATA)
向房间内所有玩家发送数据：
- ProtocolType: 6
- 包含游戏ID、房间ID和数据内容

#### 4.2 [单播数据](TCP_API.md#7-房间单播数据-ROOM_UNICAST_DATA) (ROOM_UNICAST_DATA)
向指定玩家发送数据：
- ProtocolType: 7
- 包含目标玩家ID

#### 4.3 [组播数据](TCP_API.md#8-房间组播数据-ROOM_MULTICAST_DATA) (ROOM_MULTICAST_DATA)
向多个指定玩家发送数据：
- ProtocolType: 8
- 包含目标玩家ID列表

## 数据存储

### 存储类型
gServ 目前提供JSON数据在线存储方式。

### 数据格式
在线存储的数据将 JSON 格式结构体请求和响应，而非 JSON 编码字符串，由游戏客户端自行解析。

## 错误处理

### HTTP 错误码
- `200` - 成功
- `400` - 请求参数错误
- `401` - 未授权
- `403` - 禁止访问
- `404` - 资源不存在
- `500` - 服务器内部错误

### TCP 错误响应
- `HEADER_ERROR` - 协议头错误
- `UNIVERSAL_RESPONSE_ERROR` - 通用错误

## 自动清理机制

gServ 包含以下自动清理协程：

1. **房间自动清理** - 每 5 分钟清理空闲房间
2. **验证码自动清理** - 每 5 分钟清理过期验证码

## 开发注意事项

### 1. 字节序
所有多字节字段使用 **大端序 (Big Endian)**。

### 2. Token 管理
- JWT Token 有效期为 24 小时
- 每次 TCP 请求都需要携带有效 Token
- Token 过期需要重新登录获取

### 3. 连接保持
- TCP 连接需要保持活跃
- 建议实现心跳机制
- 断线后需要重新登录

### 4. 数据安全
- 敏感数据建议客户端加密
- 不要通过 gServ 传输未加密的敏感信息
- 使用 HTTPS 保护 HTTP 通信

## 示例客户端

### Go 客户端示例
```go
package main

import (
    "encoding/binary"
    "fmt"
    "net"
)

type GameClient struct {
    conn net.Conn
    token string
    gameID uint32
}

func (c *GameClient) Login() error {
    // 构建登录包
    packet := make([]byte, 1+1+2+len(c.token)+4)
    packet[0] = 1 // ProtocolVersion
    packet[1] = 3 // AUTH_PLAYER_LOGIN
    
    // Token 长度和内容
    binary.BigEndian.PutUint16(packet[2:4], uint16(len(c.token)))
    copy(packet[4:4+len(c.token)], c.token)
    
    // GameID
    binary.BigEndian.PutUint32(packet[4+len(c.token):], c.gameID)
    
    // 发送
    _, err := c.conn.Write(packet)
    return err
}

func (c *GameClient) JoinRoom(roomID uint64) error {
    // 构建加入房间包
    packet := make([]byte, 1+1+2+len(c.token)+4+8)
    packet[0] = 1 // ProtocolVersion
    packet[1] = 5 // JOIN_ROOM
    
    // Token 长度和内容
    binary.BigEndian.PutUint16(packet[2:4], uint16(len(c.token)))
    copy(packet[4:4+len(c.token)], c.token)
    
    // GameID
    offset := 4 + len(c.token)
    binary.BigEndian.PutUint32(packet[offset:offset+4], c.gameID)
    
    // RoomID
    binary.BigEndian.PutUint64(packet[offset+4:offset+12], roomID)
    
    // 发送
    _, err := c.conn.Write(packet)
    return err
}
```

## 常见问题

### Q1: 如何获取游戏ID？
A: 游戏开发者需要通过管理员接口添加游戏，或使用浊水楼台服务注册游戏。

### Q2: Token 过期怎么办？
A: 需要重新调用登录接口获取新的 Token。

### Q3: TCP 连接断开如何处理？
A: 实现重连机制，重新建立连接并登录。

### Q4: 如何保证消息顺序？
A: TCP 协议保证数据顺序，但客户端需要处理并发发送。

### Q5: 支持多少并发玩家？
A: 取决于服务器配置，默认配置支持数百并发连接。

## 版本历史

| 版本 | 日期 | 描述 |
|------|------|------|
| 1.0 | 2026-03-17 | 初始版本，支持基础游戏服务功能 |

## 技术支持

- 文档：本项目文档
- 邮箱：imjfoy@163.com（浊水楼台服务）
- 问题反馈：GitHub Issues

---

*本文档基于 gServ 项目代码和接口文档生成，具体实现以实际代码为准。*