# 快速开始指南

## 🚀 5分钟启动项目

### 前置条件

1. ✅ Go 1.23.7+ 已安装
2. ✅ MySQL 数据库已安装并运行
3. ✅ PowerShell（Windows自带）

### 第一步：准备数据库

在MySQL中执行以下命令：

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

### 第二步：配置数据库连接

编辑文件 `rpc/etc/template.yaml`，修改数据库连接信息：

```yaml
Mysql:
  datasource: "root:你的密码@tcp(127.0.0.1:3306)/lyx?charset=utf8mb4&parseTime=True&loc=Local"
```

将 `你的密码` 替换为实际的MySQL密码。

### 第三步：构建项目

在项目根目录打开PowerShell，执行：

```powershell
.\build.ps1
```

等待构建完成，看到 "✓ API服务构建成功" 和 "✓ RPC服务构建成功" 即可。

### 第四步：启动服务

**方式1：使用启动脚本（推荐）**

```powershell
.\start.ps1
```

选择选项 `1` 启动所有服务。

**方式2：手动启动**

打开两个PowerShell窗口：

窗口1 - 启动RPC服务：
```powershell
cd rpc
.\rpc.exe
```

窗口2 - 启动API服务：
```powershell
cd api
.\api.exe
```

### 第五步：测试接口

**方式1：使用测试脚本**

```powershell
.\test.ps1
```

**方式2：使用curl命令**

测试注册：
```bash
curl -X POST http://localhost:8888/v1/user/register -H 'Content-Type: application/json' -d '{"username":"test","password":"123456","email":"test@test.com"}'
```

测试登录：
```bash
curl -X POST http://localhost:8888/v1/user/login -H 'Content-Type: application/json' -d '{"username":"test","password":"123456"}'
```

**方式3：使用Postman或浏览器插件**

- 注册接口：POST http://localhost:8888/v1/user/register
  ```json
  {
    "username": "test",
    "password": "123456",
    "email": "test@test.com"
  }
  ```

- 登录接口：POST http://localhost:8888/v1/user/login
  ```json
  {
    "username": "test",
    "password": "123456"
  }
  ```

## ✅ 预期结果

### 注册成功响应
```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "userId": 1,
    "username": "test",
    "message": "注册成功"
  }
}
```

### 登录成功响应
```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "token": "123"
  }
}
```

## 🔧 常见问题

### 1. 编译失败：无法解析依赖

**解决方案：**
```powershell
go work sync
go mod tidy
```

### 2. 数据库连接失败

**检查项：**
- MySQL服务是否运行
- 用户名密码是否正确
- 数据库是否已创建
- 端口是否正确（默认3306）

### 3. RPC连接失败

**检查项：**
- RPC服务是否已启动
- 端口8083是否被占用
- API配置文件 `api/etc/user-api.yaml` 中的 `UserRpc.Target` 是否为 `127.0.0.1:8083`

### 4. 端口被占用

**解决方案：**
- 修改配置文件中的端口号
- 或者关闭占用端口的程序

## 📝 下一步

- 查看 [README.md](README.md) 了解更多功能
- 查看项目结构了解代码组织
- 尝试mDNS扫描功能

## 🛠️ 开发提示

1. **日志位置**：
   - API日志：`api/logs/`
   - RPC日志：`rpc/logs/`

2. **停止服务**：
   - 在运行服务的终端按 `Ctrl+C`

3. **重新构建**：
   ```powershell
   .\build.ps1
   ```

祝你使用愉快！🎉
