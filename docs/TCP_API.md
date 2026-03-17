# gServ TCP API 协议文档

## 概述

gServ（游戏服务端）是一个简单通用的游戏服务端，主要功能包括：
1. 转发游戏客户端数据
2. 在线存储客户端需要存放的数据
3. 管理玩家连接和房间通信

本文档面向游戏客户端开发者，描述了TCP通信协议格式和消息类型。

## 协议基础格式

所有TCP消息都遵循以下基础格式：

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 |
| ProtocolType | uint8 | 1字节 | 协议类型 |
| ProtocolData | 变长 | 变长 | 协议数据 |

### 协议版本号
| 值 | 描述 |
|----|------|
| 0 | 未知版本 |
| 1 | 第一版协议 |

### 协议类型
| 值 | 类型 | 描述 |
|----|------|------|
| 0 | UNKNOWN | 未知协议类型 |
| 1 | HEADER_ERROR | 协议头错误 |
| 2 | UNIVERSAL_RESPONSE_ERROR | 通用错误响应 |
| 3 | AUTH_PLAYER_LOGIN | 玩家登录认证 |
| 4 | AUTH_PLAYER_LOGOUT | 玩家登出 |
| 5 | JOIN_ROOM | 加入房间 |
| 6 | ROOM_BROADCAST_DATA | 房间广播数据 |
| 7 | ROOM_UNICAST_DATA | 房间单播数据 |
| 8 | ROOM_MULTICAST_DATA | 房间组播数据 |

## 协议详细说明

### 1. 协议头错误 (HEADER_ERROR)

当协议头验证失败时返回。

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (1) |
| ErrorLength | uint32 | 4字节 | 错误信息长度 |
| Error | []byte | 变长 | 错误信息内容 |

### 2. 通用错误响应 (UNIVERSAL_RESPONSE_ERROR)

通用错误响应格式。

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (2) |
| Status | uint8 | 1字节 | 状态码 |
| MessageType | uint8 | 1字节 | 消息类型 |
| MessageLength | uint32 | 4字节 | 消息长度 |
| Message | []byte | 变长 | 消息内容 |

#### 消息类型
| 值 | 类型 | 描述 |
|----|------|------|
| 0 | UNKNOWN | 未知消息类型 |
| 1 | TEXT | 文本消息 |
| 2 | ERROR | 错误消息 |

### 3. 玩家登录认证 (AUTH_PLAYER_LOGIN)

#### 3.1 登录请求

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (3) |
| TokenLength | uint16 | 2字节 | Token长度 |
| Token | []byte | 变长 | JWT Token |
| GameID | uint32 | 4字节 | 游戏ID |

#### 3.2 登录响应

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (3) |
| Status | uint8 | 1字节 | 登录状态 |
| MessageType | uint8 | 1字节 | 消息类型 |
| MessageLength | uint32 | 4字节 | 消息长度 |
| Message | []byte | 变长 | 消息内容 |

#### 登录状态码
| 值 | 状态 | 描述 |
|----|------|------|
| 0 | FAILURE | 登录失败 |
| 1 | SUCCESS | 登录成功 |
| 2 | ALREADY_LOGIN_OTHER_GAME | 已登录其他游戏（登录成功） |
| 3 | ALREADY_LOGIN_OTHER_PLACE | 异地登录（登录成功） |

### 4. 玩家登出 (AUTH_PLAYER_LOGOUT)

#### 4.1 登出请求

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (4) |
| TokenLength | uint16 | 2字节 | Token长度 |
| Token | []byte | 变长 | JWT Token |
| GameID | uint32 | 4字节 | 游戏ID |

#### 4.2 登出响应

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (4) |
| Status | uint8 | 1字节 | 登出状态 |
| MessageType | uint8 | 1字节 | 消息类型 |
| MessageLength | uint32 | 4字节 | 消息长度 |
| Message | []byte | 变长 | 消息内容 |

#### 登出状态码
| 值 | 状态 | 描述 |
|----|------|------|
| 0 | FAILURE | 登出失败 |
| 1 | SUCCESS | 登出成功 |
| 2 | NOT_LOGIN | 未登录 |

### 5. 加入房间 (JOIN_ROOM)

#### 5.1 加入房间请求

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (5) |
| TokenLength | uint16 | 2字节 | Token长度 |
| Token | []byte | 变长 | JWT Token |
| GameID | uint32 | 4字节 | 游戏ID |
| RoomID | uint64 | 8字节 | 房间ID |

#### 5.2 加入房间响应

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (5) |
| Status | uint8 | 1字节 | 加入状态 |
| MessageType | uint8 | 1字节 | 消息类型 |
| MessageLength | uint32 | 4字节 | 消息长度 |
| Message | []byte | 变长 | 消息内容 |

#### 加入房间状态码
| 值 | 状态 | 描述 |
|----|------|------|
| 0 | FAILURE | 加入失败 |
| 1 | SUCCESS | 加入成功 |

### 6. 房间广播数据 (ROOM_BROADCAST_DATA)

#### 6.1 广播数据请求

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (6) |
| TokenLength | uint16 | 2字节 | Token长度 |
| Token | []byte | 变长 | JWT Token |
| GameID | uint32 | 4字节 | 游戏ID |
| RoomID | uint64 | 8字节 | 房间ID |
| DataLength | uint32 | 4字节 | 数据长度 |
| Data | []byte | 变长 | 数据内容 |

#### 6.2 广播数据响应

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (6) |
| GameID | uint32 | 4字节 | 游戏ID |
| RoomID | uint64 | 8字节 | 房间ID |
| SourcePlayerID | uint32 | 4字节 | 发送者玩家ID |
| DataLength | uint32 | 4字节 | 数据长度 |
| Data | []byte | 变长 | 数据内容 |

#### 广播数据状态码
| 值 | 状态 | 描述 |
|----|------|------|
| 0 | FAILURE | 广播失败 |

### 7. 房间单播数据 (ROOM_UNICAST_DATA)

#### 7.1 单播数据请求

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (7) |
| TokenLength | uint16 | 2字节 | Token长度 |
| Token | []byte | 变长 | JWT Token |
| GameID | uint32 | 4字节 | 游戏ID |
| RoomID | uint64 | 8字节 | 房间ID |
| PlayerID | uint32 | 4字节 | 目标玩家ID |
| DataLength | uint32 | 4字节 | 数据长度 |
| Data | []byte | 变长 | 数据内容 |

#### 7.2 单播数据响应

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (7) |
| GameID | uint32 | 4字节 | 游戏ID |
| RoomID | uint64 | 8字节 | 房间ID |
| SourcePlayerID | uint32 | 4字节 | 发送者玩家ID |
| DataLength | uint32 | 4字节 | 数据长度 |
| Data | []byte | 变长 | 数据内容 |

### 8. 房间组播数据 (ROOM_MULTICAST_DATA)

#### 8.1 组播数据请求

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (8) |
| TokenLength | uint16 | 2字节 | Token长度 |
| Token | []byte | 变长 | JWT Token |
| GameID | uint32 | 4字节 | 游戏ID |
| RoomID | uint64 | 8字节 | 房间ID |
| PlayerIDCount | uint32 | 4字节 | 目标玩家数量 |
| PlayerIDs | []uint32 | 变长 | 目标玩家ID列表 |
| DataLength | uint32 | 4字节 | 数据长度 |
| Data | []byte | 变长 | 数据内容 |

#### 8.2 组播数据响应

| 字段名 | 类型 | 长度 | 描述 |
|--------|------|------|------|
| ProtocolVersion | uint8 | 1字节 | 协议版本号 (1) |
| ProtocolType | uint8 | 1字节 | 协议类型 (8) |
| GameID | uint32 | 4字节 | 游戏ID |
| RoomID | uint64 | 8字节 | 房间ID |
| SourcePlayerID | uint32 | 4字节 | 发送者玩家ID |
| DataLength | uint32 | 4字节 | 数据长度 |
| Data | []byte | 变长 | 数据内容 |

## 通信流程

### 1. 连接建立
1. 客户端连接到gServ TCP服务器
2. 服务器等待客户端发送登录请求

### 2. 玩家登录
1. 客户端发送AUTH_PLAYER_LOGIN请求
2. 服务器验证Token和游戏ID
3. 服务器返回登录响应
4. 登录成功后，玩家连接被管理

### 3. 房间操作
1. 玩家发送JOIN_ROOM请求加入房间
2. 服务器验证权限并加入房间
3. 玩家可以在房间内发送数据

### 4. 数据传输
1. 广播数据：发送给房间内所有玩家
2. 单播数据：发送给指定玩家
3. 组播数据：发送给多个指定玩家

### 5. 玩家登出
1. 客户端发送AUTH_PLAYER_LOGOUT请求
2. 服务器清理玩家状态
3. 服务器返回登出响应
4. 连接关闭

## 错误处理

### 协议头错误
- 当协议头验证失败时，返回HEADER_ERROR
- 包含具体的错误信息

### 通用错误
- 使用UNIVERSAL_RESPONSE_ERROR响应
- 包含状态码和错误消息

### 连接管理
- 玩家必须登录后才能进行其他操作
- Token过期或无效会导致连接断开
- 房间操作需要验证玩家权限

## 注意事项

1. **字节序**：所有多字节字段使用大端序（Big Endian）
2. **Token验证**：每次请求都需要携带有效的JWT Token
3. **连接保持**：登录成功后需要保持TCP连接
4. **数据格式**：Data字段可以是任意二进制数据，由游戏逻辑解析
5. **错误重试**：客户端应实现适当的错误重试机制

## 示例通信

### 登录流程
```
客户端 → 服务器: [ProtocolVersion=1, ProtocolType=3, Token, GameID]
服务器 → 客户端: [ProtocolVersion=1, ProtocolType=3, Status=1, Message="登录成功"]
```

### 广播数据
```
客户端 → 服务器: [ProtocolVersion=1, ProtocolType=6, Token, GameID, RoomID, Data]
服务器 → 所有玩家: [ProtocolVersion=1, ProtocolType=6, GameID, RoomID, SourcePlayerID, Data]
```

## 版本历史

| 版本 | 日期 | 描述 |
|------|------|------|
| 1.0 | 2026-03-17 | 初始版本，基于gServ项目实现 |

---

*本文档基于gServ项目代码生成，如有更新请参考实际代码实现。*