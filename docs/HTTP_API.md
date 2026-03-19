# gServ HTTP API 文档

## 概述

gServ（游戏服务端）是一个轻量、高效、安全的通用游戏服务端框架，提供玩家管理、房间管理、数据存储和消息转发等核心功能。本文档详细描述了gServ的HTTP RESTful API接口。

### 基础信息
- **项目名称**: gServ（全称：游戏服务端，英文名：gameServer）
- **服务地址**: `http://localhost:8080`（默认HTTP端口）
- **API前缀**: `/api`
- **认证方式**: JWT Token（玩家接口）、Auth-Code（管理接口）

### 响应格式
所有API响应均使用JSON格式，成功响应状态码为`200`，错误响应包含`error`字段。

### 认证方式
1. **玩家认证**: 使用JWT Token，通过`Authorization`请求头传递
```
Authorization: Bearer <token>
```
2. **管理员认证**: 使用Auth-Code，通过`Auth-Code`请求头传递
```
Auth-Code: Bearer <auth_code>
```

---

## 目录
- [健康检查](#健康检查)
- [游戏管理](#游戏管理)
- [房间管理](#房间管理)
- [验证码系统](#验证码系统)
- [玩家管理](#玩家管理)
- [数据存储](#数据存储)
- [管理员接口](#管理员接口)

---

## 健康检查

### 健康检查接口
检查服务是否正常运行。

**请求**
```
GET /api/health
```

**响应**
```json
null
```

**状态码**
- `200`: 服务正常运行

---

## 游戏管理

### 获取游戏列表
获取所有游戏的分页列表。

**请求**
```
GET /api/games
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| index | int | 是 | 起始索引（从1开始） |
| limit | int | 是 | 每页数量 |

**响应**
```json
[
  {
    "id": 1,
    "name": "游戏名称",
    "room_count": 5,
    "created_at": "2024-01-01T12:00:00Z"
  }
]
```

**状态码**
- `200`: 成功
- `400`: 请求参数错误
- `500`: 服务器内部错误

---

## 房间管理

### 获取房间列表
获取指定游戏的所有房间列表。

**请求**
```
GET /api/rooms
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |

**响应**
```json
[
  {
    "room_id": 1234567890,
    "name": "房间名称",
    "homeowner_id": 1,
    "max_player": 8,
    "player_count": 3,
    "created_at": "2024-01-01T12:00:00Z"
  }
]
```

**认证**: 需要玩家认证

### 创建房间
创建新的游戏房间。

**请求**
```
POST /api/room
```

**请求体**
```json
{
  "game_id": 1,
  "name": "房间名称",
  "max_player": 8
}
```

**响应**
```json
{
  "room_id": 1234567890
}
```

**认证**: 需要玩家认证

### 获取房间信息
获取指定房间的详细信息。

**请求**
```
GET /api/room/:game_id/:room_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |
| room_id | string | 是 | 房间ID |

**响应**
```json
{
  "room_id": 1234567890,
  "name": "房间名称",
  "homeowner_id": 1,
  "max_player": 8,
  "player_count": 3,
  "player_ids": [1, 2, 3],
  "created_at": "2024-01-01T12:00:00Z"
}
```

**认证**: 需要玩家认证

### 放逐玩家
将玩家从房间中移除。

**请求**
```
PUT /api/room/:game_id/:room_id/:player_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |
| room_id | string | 是 | 房间ID |
| player_id | string | 是 | 玩家ID |

**认证**: 需要玩家认证（房主权限）

### 锁定房间
锁定房间，禁止新玩家加入。

**请求**
```
PUT /api/room/:game_id/:room_id/lock
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |
| room_id | string | 是 | 房间ID |

**认证**: 需要玩家认证（房主权限）

### 解锁房间
解锁房间，允许新玩家加入。

**请求**
```
PUT /api/room/:game_id/:room_id/unlock
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |
| room_id | string | 是 | 房间ID |

**认证**: 需要玩家认证（房主权限）

### 删除房间
删除指定的房间。

**请求**
```
DELETE /api/room/:game_id/:room_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |
| room_id | string | 是 | 房间ID |

**认证**: 需要玩家认证（房主权限）

---

## 验证码系统

### 发送邮箱验证码
向指定邮箱发送验证码。

**请求**
```
POST /api/captcha/email
```

**请求体**
```json
{
  "email": "user@example.com",
  "email_type": 1
}
```

**参数说明**
- `email_type`: 验证码类型，可选值：
  1. `register`: 注册验证码
  2. `reset_password`: 修改密码验证码
  3. `change_email`: 修改邮箱验证码

**状态码**
- `200`: 成功
- `400`: 请求参数错误
- `500`: 服务器内部错误

---

## 玩家管理

### 玩家注册
注册新玩家账号。

**请求**
```
POST /api/player/register
```

**请求体**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "nickname": "玩家昵称",
  "captcha": "123456"
}
```

**响应**
```json
{
  "player_id": 1
}
```

**状态码**
- `200`: 注册成功
- `400`: 请求参数错误或验证码错误
- `500`: 服务器内部错误

### 玩家登录
玩家登录获取Token。

**请求**
```
POST /api/player/login
```

**请求体**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "nickname": "玩家昵称",
  "tcp_port": 9090
}
```

**状态码**
- `200`: 登录成功
- `400`: 请求参数错误
- `401`: 邮箱或密码错误
- `500`: 服务器内部错误

### 获取玩家信息
获取指定玩家的信息。

**请求**
```
GET /api/player/:player_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| player_id | string | 是 | 玩家ID |

**响应**
```json
{
  "id": 1,
  "email": "user@example.com",
  "nickname": "玩家昵称",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**认证**: 需要玩家认证（只能获取自己的信息）

### 更新玩家信息
更新玩家信息（目前仅支持修改昵称和邮箱）。

**请求**
```
PUT /api/player
```

**请求体**
```json
{
  "nickname": "新昵称"
}
```

**认证**: 需要玩家认证

**注意**: 修改邮箱需要附带发送到旧邮箱的验证码

### 修改密码
修改玩家密码。

**请求**
```
PUT /api/player/password
```

**请求体**
```json
{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

**认证**: 需要玩家认证

**注意**: 可能需要附带发送到邮箱的验证码

### 删除玩家
删除玩家账号。

**请求**
```
DELETE /api/player/:player_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| player_id | string | 是 | 玩家ID |

**认证**: 需要玩家认证（只能删除自己的账号）

---

## 数据存储

### 检查存档是否存在
检查指定游戏和玩家的存档是否存在。

**请求**
```
GET /api/data/:game_id/exists
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |

**响应**
```json
{
  "exists": true
}
```

**认证**: 需要玩家认证

### 创建JSON数据
创建新的JSON数据存档。

**请求**
```
POST /api/data
```

**请求体**
```json
{
  "game_id": 1,
  "data": {
    "level": 1,
    "score": 100,
    "items": ["sword", "shield"]
    // 更多数据
  }
}
```

**响应**
```json
{
  "id": 1
}
```

**认证**: 需要玩家认证

### 获取JSON数据
获取指定游戏和玩家的JSON数据存档。

**请求**
```
GET /api/data
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | uint | 是 | 游戏ID |

**响应**
```json
{
  "id": 1,
  "game_id": 1,
  "player_id": 1,
  "data": {
    "level": 1,
    "score": 100,
    "items": ["sword", "shield"]
    // 更多数据
  },
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

**认证**: 需要玩家认证

### 更新JSON数据
更新指定游戏的JSON数据存档。

**请求**
```
PUT /api/data/:game_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |

**请求体**
```json
{
  "data": {
    "level": 2,
    "score": 200,
    "items": ["sword", "shield", "potion"]
    // 更多数据
  }
}
```

**认证**: 需要玩家认证

### 删除JSON数据
删除指定游戏和玩家的JSON数据存档。

**请求**
```
DELETE /api/data
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | uint | 是 | 游戏ID |

**认证**: 需要玩家认证

---

## 管理员接口

所有管理员接口都需要在请求头中添加`Auth-Code`字段进行认证。

### 创建游戏
创建新的游戏实例。

**请求**
```
POST /admin/game
```

**请求头**
```
Auth-Code: Bearer <auth_code>
```

**请求体**
```json
{
  "name": "游戏名称"
}
```

**响应**
```json
{
  "game_id": 1
}
```

### 删除游戏
删除指定的游戏实例。

**请求**
```
DELETE /admin/game/:game_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 游戏ID |

**请求头**
```
Auth-Code: Bearer <auth_code>
```

### 强制删除房间
强制删除指定的房间。

**请求**
```
DELETE /admin/room/:game_id/:room_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| game_id | string | 是 | 房间ID |
| room_id | string | 是 | 房间ID |

**请求头**
```
Auth-Code: Bearer <auth_code>
```

### 封禁玩家
封禁指定玩家。

**请求**
```
DELETE /admin/ban/player/:player_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| player_id | string | 是 | 玩家ID |

**请求头**
```
Auth-Code: Bearer <auth_code>
```

### 解封玩家
解封指定玩家。

**请求**
```
PUT /admin/ban/player/:player_id
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| player_id | string | 是 | 玩家ID |

**请求头**
```
Auth-Code: Bearer <auth_code>
```

---

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 注意事项

1. **玩家认证**: 所有需要玩家认证的接口都需要在请求头中添加`Authorization: Bearer <token>`
2. **管理员认证**: 所有管理员接口都需要在请求头中添加`Auth-Code: Bearer <auth_code>`
3. **数据验证**: 所有请求参数都会进行验证，不符合要求的参数会返回400错误
4. **权限控制**: 玩家只能操作自己的数据，房主只能操作自己的房间
5. **数据隔离**: 不同游戏的数据完全隔离，确保数据安全

## 版本信息

- **API版本**: v1.0
- **最后更新**: 2024年1月
- **维护者**: gServ开发团队

---

*本文档根据gServ项目实际代码生成，如有更新请参考最新代码。*