# User Center System

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

A comprehensive full-stack user center system with user registration/login, RBAC permission management, personal center, admin dashboard and enterprise-grade security features.

## 🌍 Documentation

- **[English](README.md)** (Current)
- **[中文文档](README_CN.md)** 
- **[Quick Start Guide](QUICK_START.md)**
- **[Project Status](PROJECT_STATUS.md)**

## 🚀 Quick Start

```bash
# One-command start with Docker
git clone https://github.com/namezzy/usercenter.git
cd usercenter && ./start.sh

# Access: Frontend (http://localhost:3000) | Backend (http://localhost:8080)
# Login: admin/admin123456 or user/user123456
```

## Project Structure

```
usercenter/
├── backend/                    # Go backend service
│   ├── cmd/
│   ├── configs/               # Configuration files
│   ├── internal/
│   │   ├── config/           # Configuration module
│   │   ├── database/         # Database connections
│   │   ├── handler/          # HTTP handlers
│   │   ├── middleware/       # Middleware
│   │   ├── models/           # Data models
│   │   ├── router/           # Route configuration
│   │   └── service/          # Business logic
│   ├── pkg/                  # Common packages
│   │   ├── captcha/          # Captcha
│   │   ├── crypto/           # Encryption
│   │   ├── email/            # Email service
│   │   ├── jwt/              # JWT authentication
│   │   └── sms/              # SMS service
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── Dockerfile
├── frontend/                  # React frontend application
│   ├── public/
│   ├── src/
│   │   ├── components/       # Common components
│   │   ├── hooks/            # Custom hooks
│   │   ├── pages/            # Page components
│   │   ├── services/         # API services
│   │   ├── store/            # State management
│   │   ├── types/            # TypeScript types
│   │   └── utils/            # Utility functions
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── index.html
├── docker-compose.yml         # Container orchestration
└── README.md
```

## Tech Stack

### Backend
- **Go 1.21** - Main programming language
- **Gin** - Web framework
- **GORM** - ORM framework
- **PostgreSQL** - Primary database
- **Redis** - Cache and session storage
- **JWT** - Authentication
- **Casbin** - Access control
- **Viper** - Configuration management
- **Zap** - Logging
- **Air** - Hot reload development tool
- **Wire** - Dependency injection
- **Swagger** - API documentation

### Frontend
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Ant Design** - UI component library
- **Redux Toolkit** - State management
- **React Router 6** - Routing
- **Axios** - HTTP client
- **React Hook Form** - Form handling
- **Styled Components** - CSS-in-JS
- **ECharts/AntV** - Data visualization
- **dayjs** - Date handling
- **ahooks** - React Hooks utilities

### Deployment
- **Docker** - Containerization
- **Docker Compose** - Service orchestration
- **Nginx** - Reverse proxy

## Core Features

### 1. User Registration & Login
- ✅ Multiple registration methods (email, phone)
- ✅ Image captcha protection
- ✅ Email/SMS verification codes
- ✅ Password strength validation
- ✅ JWT/Session dual mode
- ✅ Password reset functionality
- ✅ Account lockout mechanism
- ✅ Login rate limiting

### 2. RBAC Permission Management
- ✅ User, Role, Permission three-tier model
- ✅ Permission group management
- ✅ Dynamic permission assignment
- ✅ Permission inheritance
- ✅ Fine-grained access control
- ✅ Permission caching optimization

### 3. Personal Center
- ✅ Personal information management
- ✅ Avatar upload
- ✅ Password modification
- ✅ Email/phone binding
- ✅ Login history
- ✅ Security settings
- ✅ Device management
- ✅ Data export

### 4. Admin Dashboard
- ✅ User management (CRUD, status toggle, password reset)
- ✅ Role management (permission assignment)
- ✅ Permission management
- ✅ System monitoring
- ✅ Operation logs
- ✅ Data statistics
- ✅ User import/export

### 5. Security Features
- ✅ XSS protection
- ✅ SQL injection protection
- ✅ CSRF protection
- ✅ Password encryption storage
- ✅ Access rate limiting
- ✅ Sensitive operation encryption
- ✅ Security headers
- ✅ HTTPS support

### 6. System Features
- ✅ Notification message system
- ✅ Message push
- ✅ System logging
- ✅ Data backup & recovery
- ✅ System monitoring
- ✅ Error handling
- ✅ Health checks

## Quick Start

### Requirements
- Go 1.21+
- Node.js 18+
- PostgreSQL 13+
- Redis 6+
- Docker & Docker Compose

### Using Docker Compose

```bash
# Clone the project
git clone https://github.com/namezzy/usercenter.git
cd usercenter

# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

### Manual Setup

#### Backend Service

```bash
cd backend

# Install dependencies
go mod download

# Configure database connection (modify configs/config.yaml)
cp configs/config.yaml.example configs/config.yaml

# Run database migrations
go run main.go migrate

# Start service
go run main.go

# Or use Air for hot reload
air
```

#### Frontend Application

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

### Access the Application

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Documentation: http://localhost:8080/swagger/index.html

## Default Accounts

### Administrator
- Username: admin
- Email: admin@example.com
- Password: admin123456

### Regular User
- Username: user
- Email: user@example.com
- Password: user123456

## API Documentation

The backend provides complete RESTful APIs with Swagger documentation:

- Online docs: http://localhost:8080/swagger/index.html
- JSON format: http://localhost:8080/swagger/doc.json

Main API endpoints:

```
POST   /api/auth/register      # User registration
POST   /api/auth/login         # User login
POST   /api/auth/logout        # User logout
GET    /api/auth/refresh       # Refresh token
GET    /api/auth/captcha       # Get captcha

GET    /api/user/profile       # Get personal info
PUT    /api/user/profile       # Update personal info
POST   /api/user/avatar        # Upload avatar
PUT    /api/user/password      # Change password

GET    /api/admin/users        # Get user list
POST   /api/admin/users        # Create user
PUT    /api/admin/users/:id    # Update user
DELETE /api/admin/users/:id    # Delete user

GET    /api/admin/roles        # Get role list
POST   /api/admin/roles        # Create role
PUT    /api/admin/roles/:id    # Update role
DELETE /api/admin/roles/:id    # Delete role
```

## Configuration

### Backend Configuration (configs/config.yaml)

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

### Frontend Configuration (frontend/src/config.ts)

```typescript
export const config = {
  apiBaseUrl: 'http://localhost:8080/api',
  uploadUrl: 'http://localhost:8080/api/upload',
  websocketUrl: 'ws://localhost:8080/ws',
  
  // Pagination config
  pageSize: 10,
  
  // File upload config
  maxFileSize: 2 * 1024 * 1024, // 2MB
  allowedFileTypes: ['image/jpeg', 'image/png', 'image/gif'],
  
  // Password strength config
  passwordMinLength: 8,
  passwordRequireNumbers: true,
  passwordRequireSymbols: true,
};
```

## Deployment Guide

### Docker Deployment

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Scale services
docker-compose up -d --scale backend=3

# Update services
docker-compose pull
docker-compose up -d
```

### Traditional Deployment

#### Backend Deployment

```bash
# Build backend
cd backend
go build -o usercenter main.go

# Create systemd service
sudo cp usercenter.service /etc/systemd/system/
sudo systemctl enable usercenter
sudo systemctl start usercenter
```

#### Frontend Deployment

```bash
# Build frontend
cd frontend
npm run build

# Deploy to Nginx
sudo cp -r dist/* /var/www/html/
sudo systemctl reload nginx
```

## Development Guide

### Backend Development

```bash
# Hot reload development
air

# Run tests
go test ./...

# Code formatting
go fmt ./...

# Generate API docs
swag init

# Database migration
go run main.go migrate
```

### Frontend Development

```bash
# Start development server
npm run dev

# Type checking
npm run type-check

# Code formatting
npm run format

# Run tests
npm run test

# Build production
npm run build
```

### Code Standards

- Backend follows Go official code standards
- Frontend uses ESLint + Prettier
- Commit messages follow Conventional Commits specification
- Git hooks managed with Husky

## Advanced Features

### Third-party Login
- WeChat Login
- QQ Login
- GitHub Login
- Google Login

### Internationalization
- Simplified Chinese
- Traditional Chinese
- English
- Japanese

### Theme System
- Light theme
- Dark theme
- Custom themes

### Data Dashboard
- User statistics
- System monitoring
- Real-time data
- Chart displays

## Monitoring and Logging

### Application Monitoring
- Prometheus + Grafana
- Performance metrics collection
- Alert rule configuration
- Dashboard displays

### Log Management
- Structured logging
- Log level control
- Log rotation
- Centralized log collection

### Distributed Tracing
- OpenTelemetry
- Jaeger integration
- Distributed tracing
- Performance analysis

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check if database service is running
   - Verify connection configuration
   - Check network connectivity

2. **Redis Connection Failed**
   - Check Redis service status
   - Verify configuration parameters
   - Check firewall settings

3. **Frontend API Call Failed**
   - Check backend service status
   - Verify API address configuration
   - Check CORS settings

4. **Email Sending Failed**
   - Check SMTP configuration
   - Verify email authorization code
   - Check network connectivity

### Performance Optimization

1. **Database Optimization**
   - Add appropriate indexes
   - Optimize query statements
   - Configure connection pools

2. **Cache Optimization**
   - Redis cache for hot data
   - Set reasonable expiration times
   - Use cache preheating

3. **Frontend Optimization**
   - Code splitting
   - Lazy loading
   - CDN acceleration

## Contributing

1. Fork the project
2. Create a feature branch
3. Commit your changes
4. Create a Pull Request

### Commit Convention

```
feat: add user management feature
fix: fix login captcha issue
docs: update API documentation
style: code formatting
refactor: refactor permission module
test: add unit tests
chore: update dependencies
```

## License

MIT License

## Contact

- Project URL: https://github.com/namezzy/usercenter
- Issues: https://github.com/namezzy/usercenter/issues
- Email: contact@example.com

---

**Note**: This is a complete user center system implementation that includes various core features of modern web applications. The system uses a frontend-backend separation architecture, supports Docker deployment, and has good scalability and maintainability.

## 🙏 Contributors

Thanks to all contributors who have helped make this project better!

<a href="https://github.com/namezzy/usercenter/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=namezzy/usercenter" />
</a>

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=namezzy/usercenter&type=Date)](https://star-history.com/#namezzy/usercenter&Date)

## Language / 语言

- [English](README.md) - Current
- [中文](README_CN.md) - Chinese Documentation
