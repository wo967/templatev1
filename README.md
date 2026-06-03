# TemplateV1 项目

## 项目概述

这是一个基于 go-zero 框架的微服务项目，包含：
- **API服务**：HTTP接口层，处理用户请求
- **RPC服务**：gRPC服务层，处理业务逻辑和数据访问
- **mDNS扫描模块**：网络资产发现功能

## 快速开始

### 1. 环境要求

- Go 1.23.7+
- MySQL 数据库
- PowerShell（Windows）或 Bash（Linux/Mac）

### 2. 数据库准备

创建MySQL数据库并执行以下SQL：

```sql
CREATE DATABASE IF NOT EXISTS lyx CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE lyx;

CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` VARCHAR(64) NOT NULL COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码',
  `register_time` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `last_login_time` DATETIME DEFAULT NULL COMMENT '最后登录时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 3. 配置修改

修改 `rpc/etc/template.yaml` 中的数据库配置：

```yaml
Mysql:
  datasource: "用户名:密码@tcp(主机:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local"
```

### 4. 构建项目

在项目根目录执行：

```powershell
.\build.ps1
```

或者手动构建：

```powershell
# 同步依赖
go work sync

# 构建RPC服务
cd rpc
go mod tidy
go build -o rpc.exe .

# 构建API服务
cd ..\api
go mod tidy
go build -o api.exe .
```

### 5. 启动服务

**第一步：启动RPC服务**

```powershell
cd rpc
.\rpc.exe
```

**第二步：启动API服务**（新终端）

```powershell
cd api
.\api.exe
```

### 6. 测试接口

**用户注册：**

```bash
curl -X POST http://localhost:8888/v1/user/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"test","password":"123456","email":"test@test.com"}'
```

**用户登录：**

```bash
curl -X POST http://localhost:8888/v1/user/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"test","password":"123456"}'
```

**mDNS扫描：**

```bash
curl -X POST http://localhost:8888/v1/scan \
  -H 'Content-Type: application/json' \
  -d '{"cidr":"192.168.1.0/24","portRanges":["5353"],"timeout":2,"concurrency":100}'
```

## 项目结构

```
templatev1/
├── api/                    # API服务
│   ├── etc/               # 配置文件
│   ├── internal/          # 内部代码
│   │   ├── config/        # 配置定义
│   │   ├── handler/       # HTTP处理器
│   │   ├── logic/         # 业务逻辑
│   │   ├── svc/           # 服务上下文
│   │   └── types/         # 类型定义
│   └── main.go            # 入口文件
├── rpc/                   # RPC服务
│   ├── etc/              # 配置文件
│   ├── internal/         # 内部代码
│   │   ├── config/       # 配置定义
│   │   ├── logic/        # 业务逻辑
│   │   ├── server/       # gRPC服务器
│   │   ├── svc/          # 服务上下文
│   │   └── infrastructrue/ # 基础设施层
│   ├── template/         # protobuf生成的代码
│   └── template.go       # 入口文件
├── pkg/                  # 公共包
│   └── mdns/            # mDNS扫描模块
│       ├── model/       # 数据模型
│       ├── protocol/    # 协议实现
│       └── scanner/     # 扫描引擎
└── go.work              # Go工作区配置
```

## 主要功能

### 1. 用户认证

- 用户注册：验证用户名唯一性，创建新用户
- 用户登录：验证用户凭据，返回JWT Token

### 2. mDNS扫描

- 支持自定义CIDR网段扫描
- 支持多端口范围配置
- 并发扫描控制
- 发现网络设备和服务

## 开发说明

### 添加新功能

1. **API层**：
   - 在 `api/internal/types/types.go` 定义请求/响应类型
   - 在 `api/internal/handler/` 创建处理器
   - 在 `api/internal/logic/` 实现业务逻辑
   - 在 `api/internal/handler/routes.go` 注册路由

2. **RPC层**：
   - 在 `rpc/template.proto` 定义protobuf消息和服务
   - 执行 `proto.sh` 生成代码
   - 在 `rpc/internal/logic/` 实现业务逻辑

### 注意事项

1. 密码应该使用bcrypt等算法加密存储（当前为演示目的未加密）
2. JWT Token应该使用真实的密钥生成和验证（当前返回固定值）
3. 生产环境应启用HTTPS和适当的认证机制

## 常见问题

### 1. 编译错误

确保已正确配置Go工作区：

```powershell
go work sync
```

### 2. 数据库连接失败

检查 `rpc/etc/template.yaml` 中的数据库配置是否正确。

### 3. RPC连接失败

确保RPC服务已启动，且API配置中的 `UserRpc.Target` 地址正确。

## 许可证

MIT License
