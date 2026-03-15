# gServ Docker 部署指南

## 快速开始

### 1. 准备环境
```bash
# 复制环境变量模板
cp docker.env .env

# 编辑环境变量（修改密码和密钥）
vim .env
vim docker.gserv.conf.yaml
```

### 2. 启动服务
```bash
# 使用 Docker Compose 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f gserv
```

### 3. 停止服务
```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

## 服务说明

### 容器服务
1. **mysql** (MySQL 8.0.36)
   - 端口: 3306 (可配置)
   - 数据卷: `mysql_data`
   - 自动初始化数据库结构

2. **gserv** (游戏服务端)
   - HTTP API: 8080 端口
   - TCP 长连接: 9090 端口
   - 依赖: MySQL 数据库

## 环境变量配置

### 必需配置
- `.env`
```bash
# MySQL 配置
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_DATABASE=gserv
MYSQL_USER=gserv_user
MYSQL_PASSWORD=gserv_password
```
- `docker.gserv.conf.yaml`
```yaml
# 邮箱配置
server:
  email:
    # 邮件服务主机
    host: "smtp.qq.com"
    # 邮件服务端口
    port: 465
    # 邮件服务邮箱
    email: "xxxxxxxxxxxx@qq.com"
    # 邮件服务密码
    password: "xxxxxxxxxxxxxxxxxxxx"
# MySQL 配置
database:
  # 数据库驱动
  driver: "mysql"
  # 数据库连接字符串
  dsn: "gserv_user:gserv_password@tcp(127.0.0.1:3306)/gserv?charset=utf8mb4&parseTime=True&loc=Local"
```

## 数据库管理

### 初始化数据库
首次启动时，gServ容器会自动进行数据迁移，创建：
- gserv 数据库
- 玩家表、房间表、游戏表等
- 默认游戏配置
- 定期清理事件

### 数据持久化
- MySQL数据存储在 `./app/mysql` 卷中
- 应用日志和配置文件存储在 `./app` 目录

### 备份与恢复
```bash
# 备份数据库
docker exec gserv-mysql mysqldump -u root -p gserv > backup.sql

# 恢复数据库
docker exec -i gserv-mysql mysql -u root -p gserv < backup.sql
```

## 健康检查

### 服务健康状态
```bash
# 检查所有服务健康状态
docker-compose ps

# 检查单个服务日志
docker-compose logs gserv

# 进入容器调试
docker exec -it gserv-app sh
```

### 应用健康端点
- HTTP健康检查: `http://localhost:8080/health`
- TCP连接测试: `telnet localhost 9090`

## 开发与调试

### 本地开发
```bash
# 1. 启动依赖服务（仅MySQL）
docker-compose up -d mysql

# 2. 本地运行Go应用
go run main.go

# 3. 构建Docker镜像
docker-compose build gserv
```

### 配置文件
应用配置文件位于 `docker.gserv.conf.yaml`，Docker构建前请及时修改。

### 日志查看
```bash
# 查看实时日志
docker-compose logs -f gserv

# 查看历史日志
docker-compose logs --tail=100 gserv

# 查看MySQL日志
docker-compose logs mysql
```

## 生产部署建议

### 1. 安全配置
- 修改所有默认密码和密钥
- 启用HTTPS（通过Nginx）
- 配置防火墙规则
- 定期更新镜像版本

### 2. 性能优化
- 根据负载调整容器资源限制
- 添加Redis缓存扩展
- 添加数据库连接池配置
- 启用Gzip压缩

### 3. 监控与告警
- 配置容器监控（Prometheus + Grafana）
- 设置日志聚合（ELK Stack）
- 配置健康检查告警
- 定期备份数据库

## 故障排除

### 常见问题

1. **端口冲突**
```yaml
# 修改 docker.gserv.conf.yaml 中的端口配置
server:
  # HTTP 服务端口
  http_port: 8080
  # TCP 服务端口
  tcp_port: 9090
```

2. **数据库连接失败**
```bash
# 检查MySQL服务状态
docker-compose logs mysql

# 检查网络连接
docker network inspect gserv_gserv-network
```

3. **应用启动失败**
```bash
# 查看应用日志
docker-compose logs gserv

# 检查环境变量
docker exec gserv-app env
```

### 调试命令
```bash
# 进入MySQL容器
docker exec -it gserv-mysql mysql -u gserv_user -p gserv

# 查看容器资源使用
docker stats

# 重启服务
docker-compose restart gserv

# 重建镜像
docker-compose build --no-cache gserv
```

## 更新与维护

### 更新应用
```bash
# 1. 拉取最新代码
git pull

# 2. 重建镜像
docker-compose build gserv

# 3. 重启服务
docker-compose up -d
```

### 清理资源
```bash
# 清理未使用的镜像
docker image prune

# 清理未使用的卷
docker volume prune

# 清理所有未使用资源
docker system prune -a
```

---

**注意**: 生产环境部署前，请务必修改所有默认密码和密钥，并配置适当的安全策略。