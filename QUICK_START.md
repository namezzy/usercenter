# 快速入门指南

欢迎使用用户中心系统！这是一个现代化的全栈Web应用，提供完整的用户管理、权限控制和管理后台功能。

## 🚀 5分钟快速开始

### 方式一：使用Docker（推荐）

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd usercenter

# 2. 一键启动
./start.sh

# 3. 访问应用
# 前端: http://localhost:3000
# 后端: http://localhost:8080
# API文档: http://localhost:8080/swagger/index.html
```

### 方式二：使用Makefile

```bash
# 安装和启动
make setup    # 首次使用，安装工具和依赖
make start    # 启动生产环境

# 或启动开发环境
make dev      # 启动开发环境（支持热重载）
```

### 方式三：开发环境

```bash
# 启动开发环境
./dev-start.sh

# 或分别启动
make install  # 安装依赖
make dev      # 启动开发服务
```

## 🎯 默认账户

启动后可以使用以下账户登录：

**管理员账户**
- 用户名：`admin`
- 邮箱：`admin@example.com`
- 密码：`admin123456`

**普通用户**
- 用户名：`user`  
- 邮箱：`user@example.com`
- 密码：`user123456`

## 🌟 核心功能

### 🔐 用户认证
- 邮箱/用户名登录
- 用户注册
- 密码重置
- JWT认证
- 会话管理

### 👥 用户管理
- 个人信息管理
- 头像上传
- 密码修改
- 安全设置

### 🛡️ 权限控制
- RBAC角色权限
- 动态权限分配
- 页面级权限控制
- API级权限验证

### 🏢 管理后台
- 用户管理（增删改查）
- 角色管理
- 权限分配
- 系统监控

### 🔒 安全特性
- XSS防护
- CSRF防护
- SQL注入防护
- 密码加密
- 访问限流

## 📁 项目结构

```
usercenter/
├── backend/           # Go后端服务
│   ├── cmd/          # 命令行工具
│   ├── configs/      # 配置文件
│   ├── internal/     # 内部包
│   ├── pkg/          # 公共包
│   └── docs/         # API文档
├── frontend/         # React前端应用
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── pages/       # 页面
│   │   ├── services/    # API服务
│   │   ├── store/       # 状态管理
│   │   └── types/       # 类型定义
│   └── public/       # 静态资源
├── docker-compose.yml # 容器编排
├── start.sh          # 生产启动脚本
├── dev-start.sh      # 开发启动脚本
├── Makefile          # 开发命令
└── README.md         # 项目文档
```

## 🛠️ 技术栈

### 后端
- **Go 1.21** - 服务端语言
- **Gin** - Web框架
- **GORM** - ORM框架  
- **PostgreSQL** - 数据库
- **Redis** - 缓存
- **JWT** - 身份认证
- **Casbin** - 权限控制
- **Swagger** - API文档

### 前端
- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Ant Design** - UI组件
- **Redux Toolkit** - 状态管理
- **React Router** - 路由管理
- **Vite** - 构建工具
- **Axios** - HTTP客户端

## 📋 常用命令

### 使用Makefile
```bash
make help      # 查看所有命令
make start     # 启动生产环境
make dev       # 启动开发环境
make stop      # 停止服务
make status    # 查看状态
make logs      # 查看日志
make test      # 运行测试
make build     # 构建项目
```

### 使用脚本
```bash
./start.sh help        # 生产环境命令
./dev-start.sh help    # 开发环境命令
```

### Docker Compose
```bash
docker-compose up -d   # 启动所有服务
docker-compose down    # 停止所有服务
docker-compose ps      # 查看服务状态
docker-compose logs -f # 查看日志
```

## 🔧 开发指南

### 环境要求
- **Go 1.21+**
- **Node.js 18+**
- **Docker & Docker Compose**
- **PostgreSQL 13+**
- **Redis 6+**

### 开发环境设置
```bash
# 1. 安装开发工具
make dev-tools

# 2. 安装项目依赖
make install

# 3. 启动数据库
docker-compose up -d postgres redis

# 4. 运行数据库迁移
make migrate

# 5. 启动开发服务
make dev
```

### 代码规范
```bash
make format    # 格式化代码
make lint      # 代码检查
make test      # 运行测试
```

## 🐛 故障排除

### 常见问题

**1. 端口占用**
```bash
# 检查端口占用
lsof -i :3000  # 前端端口
lsof -i :8080  # 后端端口
lsof -i :5432  # 数据库端口
```

**2. 数据库连接失败**
```bash
# 重启数据库
docker-compose restart postgres redis

# 检查数据库状态
make health
```

**3. 前端编译错误**
```bash
# 清理依赖重新安装
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**4. 后端启动失败**
```bash
# 检查Go模块
cd backend
go mod tidy
go mod download
```

### 完全重置
```bash
# 清理所有数据和容器
make clean

# 重新构建和启动
make rebuild
```

## 📖 API文档

启动服务后访问API文档：
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON格式**: http://localhost:8080/swagger/doc.json

### 主要API端点

**认证相关**
```
POST /api/auth/register  # 用户注册
POST /api/auth/login     # 用户登录
POST /api/auth/logout    # 用户登出
GET  /api/auth/refresh   # 刷新Token
```

**用户相关**
```
GET  /api/user/profile   # 获取个人信息
PUT  /api/user/profile   # 更新个人信息
PUT  /api/user/password  # 修改密码
POST /api/user/avatar    # 上传头像
```

**管理相关**
```
GET    /api/admin/users     # 获取用户列表
POST   /api/admin/users     # 创建用户
PUT    /api/admin/users/:id # 更新用户
DELETE /api/admin/users/:id # 删除用户

GET    /api/admin/roles     # 获取角色列表
POST   /api/admin/roles     # 创建角色
PUT    /api/admin/roles/:id # 更新角色
DELETE /api/admin/roles/:id # 删除角色
```

## 🚀 部署指南

### Docker部署（推荐）
```bash
# 生产环境部署
docker-compose -f docker-compose.prod.yml up -d

# 扩容服务
docker-compose up -d --scale backend=3
```

### 手动部署
```bash
# 构建后端
cd backend
go build -o bin/usercenter main.go

# 构建前端
cd frontend
npm run build

# 部署到服务器
# ... 复制文件到服务器
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

### 提交规范
```
feat: 添加新功能
fix: 修复问题
docs: 更新文档
style: 代码格式化
refactor: 代码重构
test: 添加测试
chore: 维护任务
```

## 📝 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 获取详细更新记录。

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 💬 支持

如果你遇到问题或有建议，请：

1. 查看 [FAQ](docs/FAQ.md)
2. 搜索 [Issues](../../issues)
3. 创建新的 [Issue](../../issues/new)
4. 发送邮件到：support@example.com

---

**快速链接**
- [完整文档](README.md)
- [API文档](http://localhost:8080/swagger/index.html)
- [项目状态](PROJECT_STATUS.md)
- [常见问题](docs/FAQ.md)

🎉 **祝你使用愉快！**
