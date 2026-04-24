# go-base

一个基于 Go + Gin + Gorm + MySQL 的后端基础项目，内置用户注册、登录和 JWT 鉴权示例，适合作为 API 服务脚手架或学习项目。

## 技术栈

- Go
- Gin
- Gorm
- MySQL
- Viper（配置管理）
- JWT（登录态鉴权）

## 功能说明

- 用户注册：`POST /api/register`
- 用户登录：`POST /api/login`
- 获取当前用户信息（需要 JWT）：`GET /api/userInfo`

## 环境要求

- Go `1.25+`（以 `go.mod` 为准）
- MySQL `8.x`（推荐）

## 配置文件

项目默认读取 `config/config.yaml`：

- `server.port`：服务端口
- `database.*`：数据库连接参数
- `jwt.secret`：JWT 签名密钥（生产环境务必修改）
- `jwt.expire_hour`：Token 过期时间（小时）

示例数据库配置（已在项目内提供）：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: 123456
  dbname: test_db
```

## 快速开始

1. 安装依赖

```bash
go mod tidy
```

2. 准备 MySQL 数据库并修改配置（`config/config.yaml`）

3. 启动服务

```bash
go run main.go
```

服务默认运行在：`http://localhost:8888`

## 接口示例

### 1) 注册

```bash
curl -X POST "http://localhost:8888/api/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demo",
    "password": "123456",
    "email": "demo@example.com"
  }'
```

### 2) 登录

```bash
curl -X POST "http://localhost:8888/api/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demo",
    "password": "123456"
  }'
```

登录成功后会返回 `token`。

### 3) 获取当前用户信息

```bash
curl -X GET "http://localhost:8888/api/userInfo" \
  -H "Authorization: Bearer <your_token>"
```

## 项目结构

```text
go-base/
├── config/      # 配置加载与配置文件
├── database/    # 数据库初始化
├── handler/     # HTTP 处理层
├── middleware/  # 中间件（JWT）
├── model/       # 数据模型
├── response/    # 统一响应结构
├── router/      # 路由注册
├── service/     # 业务逻辑层
├── utils/       # 日志与工具函数
└── main.go      # 程序入口
```

## 说明

- 程序启动时会初始化日志文件 `gin.log`。
- 数据库连接失败或配置错误会直接导致服务启动失败，请优先检查配置与数据库可用性。
