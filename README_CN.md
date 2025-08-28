# 用户中心系统

一个完整的用户中心系统，包含用户注册登录、RBAC权限管理、个人中心、管理后台等功能。

## 项目结构

```
usercenter/
├── backend/                    # Go后端服务
│   ├── cmd/
│   ├── configs/               # 配置文件
│   ├── internal/
│   │   ├── config/           # 配置模块
│   │   ├── database/         # 数据库连接
│   │   ├── handler/          # HTTP处理器
│   │   ├── middleware/       # 中间件
│   │   ├── models/           # 数据模型
│   │   ├── router/           # 路由配置
│   │   └── service/          # 业务逻辑
│   ├── pkg/                  # 公共包
│   │   ├── captcha/          # 验证码
│   │   ├── crypto/           # 加密
│   │   ├── email/            # 邮件服务
│   │   ├── jwt/              # JWT认证
│   │   └── sms/              # 短信服务
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── Dockerfile
├── frontend/                  # React前端应用
│   ├── public/
│   ├── src/
│   │   ├── components/       # 公共组件
│   │   ├── hooks/            # 自定义hooks
│   │   ├── pages/            # 页面组件
│   │   ├── services/         # API服务
│   │   ├── store/            # 状态管理
│   │   ├── types/            # TypeScript类型
│   │   └── utils/            # 工具函数
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── index.html
├── docker-compose.yml         # 容器编排
└── README.md
```

## 技术栈

### 后端
- **Go 1.21** - 主要编程语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **PostgreSQL** - 主数据库
- **Redis** - 缓存和会话存储
- **JWT** - 身份认证
- **Casbin** - 权限控制
- **Viper** - 配置管理
- **Zap** - 日志管理
- **Air** - 热重载开发工具
- **Wire** - 依赖注入
- **Swagger** - API文档

### 前端
- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Ant Design** - UI组件库
- **Redux Toolkit** - 状态管理
- **React Router 6** - 路由管理
- **Axios** - HTTP客户端
- **React Hook Form** - 表单处理
- **Styled Components** - CSS-in-JS
- **ECharts/AntV** - 数据可视化
- **dayjs** - 日期处理
- **ahooks** - React Hooks工具集

### 部署
- **Docker** - 容器化
- **Docker Compose** - 服务编排
- **Nginx** - 反向代理

## 核心功能

### 1. 用户注册登录
- ✅ 多种注册方式（邮箱、手机号）
- ✅ 图片验证码防护
- ✅ 邮箱/短信验证码
- ✅ 密码强度检查
- ✅ JWT/Session双模式
- ✅ 密码重置功能
- ✅ 账户锁定机制
- ✅ 登录频率限制

### 2. RBAC权限管理
- ✅ 用户、角色、权限三层模型
- ✅ 权限分组管理
- ✅ 动态权限分配
- ✅ 权限继承
- ✅ 细粒度权限控制
- ✅ 权限缓存优化

### 3. 个人中心
- ✅ 个人信息管理
- ✅ 头像上传
- ✅ 密码修改
- ✅ 邮箱/手机绑定
- ✅ 登录日志查看
- ✅ 安全设置
- ✅ 设备管理
- ✅ 数据导出

### 4. 管理后台
- ✅ 用户管理（CRUD、状态切换、密码重置）
- ✅ 角色管理（权限分配）
- ✅ 权限管理
- ✅ 系统监控
- ✅ 操作日志
- ✅ 数据统计
- ✅ 用户导入导出

### 5. 安全特性
- ✅ XSS防护
- ✅ SQL注入防护
- ✅ CSRF防护
- ✅ 密码加密存储
- ✅ 访问频率限制
- ✅ 敏感操作加密
- ✅ 安全头设置
- ✅ HTTPS支持

### 6. 系统功能
- ✅ 通知消息系统
- ✅ 消息推送
- ✅ 系统日志
- ✅ 数据备份恢复
- ✅ 系统监控
- ✅ 错误处理
- ✅ 健康检查

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 13+
- Redis 6+
- Docker & Docker Compose

### 使用Docker Compose启动

```bash
# 克隆项目
git clone <repository-url>
cd usercenter

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 手动启动

#### 后端服务

```bash
cd backend

# 安装依赖
go mod download

# 配置数据库连接（修改configs/config.yaml）
cp configs/config.yaml.example configs/config.yaml

# 运行数据库迁移
go run main.go migrate

# 启动服务
go run main.go

# 或使用Air热重载
air
```

#### 前端应用

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

### 访问应用

- 前端应用: http://localhost:3000
- 后端API: http://localhost:8080
- API文档: http://localhost:8080/swagger/index.html

## 默认账户

### 管理员账户
- 用户名: admin
- 邮箱: admin@example.com
- 密码: admin123456

### 普通用户
- 用户名: user
- 邮箱: user@example.com
- 密码: user123456

## API文档

后端提供完整的RESTful API，支持Swagger文档：

- 在线文档: http://localhost:8080/swagger/index.html
- JSON格式: http://localhost:8080/swagger/doc.json

主要API端点：

```
POST   /api/auth/register      # 用户注册
POST   /api/auth/login         # 用户登录
POST   /api/auth/logout        # 用户登出
GET    /api/auth/refresh       # 刷新Token
GET    /api/auth/captcha       # 获取验证码

GET    /api/user/profile       # 获取个人信息
PUT    /api/user/profile       # 更新个人信息
POST   /api/user/avatar        # 上传头像
PUT    /api/user/password      # 修改密码

GET    /api/admin/users        # 获取用户列表
POST   /api/admin/users        # 创建用户
PUT    /api/admin/users/:id    # 更新用户
DELETE /api/admin/users/:id    # 删除用户

GET    /api/admin/roles        # 获取角色列表
POST   /api/admin/roles        # 创建角色
PUT    /api/admin/roles/:id    # 更新角色
DELETE /api/admin/roles/:id    # 删除角色
```

## 配置说明

### 后端配置 (configs/config.yaml)

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: usercenter
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key
  expire: 7200

email:
  smtp_host: smtp.example.com
  smtp_port: 587
  username: noreply@example.com
  password: your-email-password

sms:
  access_key: your-access-key
  secret_key: your-secret-key
  region: cn-hangzhou
```

### 前端配置 (frontend/src/config.ts)

```typescript
export const config = {
  apiBaseUrl: 'http://localhost:8080/api',
  uploadUrl: 'http://localhost:8080/api/upload',
  websocketUrl: 'ws://localhost:8080/ws',
  
  // 分页配置
  pageSize: 10,
  
  // 文件上传配置
  maxFileSize: 2 * 1024 * 1024, // 2MB
  allowedFileTypes: ['image/jpeg', 'image/png', 'image/gif'],
  
  // 密码强度配置
  passwordMinLength: 8,
  passwordRequireNumbers: true,
  passwordRequireSymbols: true,
};
```

## 部署指南

### Docker部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 扩容服务
docker-compose up -d --scale backend=3

# 更新服务
docker-compose pull
docker-compose up -d
```

### 传统部署

#### 后端部署

```bash
# 编译后端
cd backend
go build -o usercenter main.go

# 创建systemd服务
sudo cp usercenter.service /etc/systemd/system/
sudo systemctl enable usercenter
sudo systemctl start usercenter
```

#### 前端部署

```bash
# 构建前端
cd frontend
npm run build

# 部署到Nginx
sudo cp -r dist/* /var/www/html/
sudo systemctl reload nginx
```

## 开发指南

### 后端开发

```bash
# 热重载开发
air

# 运行测试
go test ./...

# 代码格式化
go fmt ./...

# 生成API文档
swag init

# 数据库迁移
go run main.go migrate
```

### 前端开发

```bash
# 启动开发服务器
npm run dev

# 类型检查
npm run type-check

# 代码格式化
npm run format

# 运行测试
npm run test

# 构建生产版本
npm run build
```

### 代码规范

- 后端遵循Go官方代码规范
- 前端使用ESLint + Prettier
- 提交信息遵循Conventional Commits规范
- 使用Husky进行Git Hooks管理

## 扩展功能

### 第三方登录
- 微信登录
- QQ登录
- GitHub登录
- Google登录

### 国际化支持
- 中文简体
- 中文繁体
- English
- 日本語

### 主题系统
- 亮色主题
- 暗色主题
- 自定义主题

### 数据大屏
- 用户统计
- 系统监控
- 实时数据
- 图表展示

## 监控和日志

### 应用监控
- Prometheus + Grafana
- 性能指标收集
- 告警规则配置
- 仪表盘展示

### 日志管理
- 结构化日志
- 日志级别控制
- 日志轮转
- 集中式日志收集

### 链路追踪
- OpenTelemetry
- Jaeger集成
- 分布式追踪
- 性能分析

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务是否启动
   - 确认连接配置正确
   - 检查网络连通性

2. **Redis连接失败**
   - 检查Redis服务状态
   - 确认配置参数
   - 检查防火墙设置

3. **前端API调用失败**
   - 检查后端服务状态
   - 确认API地址配置
   - 检查跨域设置

4. **邮件发送失败**
   - 检查SMTP配置
   - 确认邮箱授权码
   - 检查网络连通性

### 性能优化

1. **数据库优化**
   - 添加合适索引
   - 优化查询语句
   - 配置连接池

2. **缓存优化**
   - Redis缓存热点数据
   - 设置合理过期时间
   - 使用缓存预热

3. **前端优化**
   - 代码分割
   - 懒加载
   - CDN加速

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码
4. 发起Pull Request

### 提交规范

```
feat: 添加用户管理功能
fix: 修复登录验证码问题
docs: 更新API文档
style: 代码格式化
refactor: 重构权限模块
test: 添加单元测试
chore: 更新依赖包
```

## 许可证

MIT License

## 联系方式

- 项目地址: [https://github.com/your-username](https://github.com/namezzy)/usercenter
- 问题反馈: [https://github.com/your-username](https://github.com/namezzy)/usercenter/issues
- 邮箱: xmj@xx-xmj.com

---

**注意**: 这是一个完整的用户中心系统实现，包含了现代Web应用的各种核心功能。系统采用前后端分离架构，支持Docker部署，具有良好的扩展性和可维护性。
