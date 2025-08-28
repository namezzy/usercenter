# 用户中心系统 Makefile
# User Center System Makefile

.PHONY: help start stop restart status logs clean dev dev-stop install test build docker-build docker-push

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
DOCKER_COMPOSE := docker-compose
PROJECT_NAME := usercenter
BACKEND_DIR := backend
FRONTEND_DIR := frontend

# 颜色定义
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

## 帮助信息
help:
	@echo "$(CYAN)用户中心系统 - 开发命令$(RESET)"
	@echo "================================="
	@echo ""
	@echo "$(GREEN)生产环境命令:$(RESET)"
	@echo "  make start      - 启动生产环境 (Docker)"
	@echo "  make stop       - 停止生产环境"
	@echo "  make restart    - 重启生产环境"
	@echo "  make status     - 查看服务状态"
	@echo "  make logs       - 查看服务日志"
	@echo "  make clean      - 清理环境和数据"
	@echo ""
	@echo "$(GREEN)开发环境命令:$(RESET)"
	@echo "  make dev        - 启动开发环境"
	@echo "  make dev-stop   - 停止开发环境"
	@echo "  make install    - 安装项目依赖"
	@echo ""
	@echo "$(GREEN)构建和测试:$(RESET)"
	@echo "  make test       - 运行测试"
	@echo "  make build      - 构建项目"
	@echo "  make docker-build - 构建Docker镜像"
	@echo ""
	@echo "$(GREEN)其他命令:$(RESET)"
	@echo "  make format     - 格式化代码"
	@echo "  make lint       - 代码检查"
	@echo "  make docs       - 生成文档"

## 生产环境 - 启动所有服务
start:
	@echo "$(GREEN)🚀 启动生产环境...$(RESET)"
	@./start.sh start

## 生产环境 - 停止所有服务
stop:
	@echo "$(YELLOW)⏹️  停止生产环境...$(RESET)"
	@./start.sh stop

## 生产环境 - 重启所有服务
restart:
	@echo "$(YELLOW)🔄 重启生产环境...$(RESET)"
	@./start.sh restart

## 查看服务状态
status:
	@echo "$(CYAN)📊 查看服务状态...$(RESET)"
	@./start.sh status

## 查看服务日志
logs:
	@echo "$(CYAN)📋 查看服务日志...$(RESET)"
	@./start.sh logs

## 清理环境和数据
clean:
	@echo "$(RED)🧹 清理环境和数据...$(RESET)"
	@./start.sh clean

## 开发环境 - 启动开发服务
dev:
	@echo "$(GREEN)🔧 启动开发环境...$(RESET)"
	@./dev-start.sh start

## 开发环境 - 停止开发服务
dev-stop:
	@echo "$(YELLOW)⏹️  停止开发环境...$(RESET)"
	@./dev-start.sh stop

## 安装项目依赖
install:
	@echo "$(CYAN)📦 安装项目依赖...$(RESET)"
	@echo "安装后端依赖..."
	@cd $(BACKEND_DIR) && go mod download
	@echo "安装前端依赖..."
	@cd $(FRONTEND_DIR) && npm install
	@echo "$(GREEN)✅ 依赖安装完成$(RESET)"

## 运行测试
test:
	@echo "$(CYAN)🧪 运行测试...$(RESET)"
	@echo "运行后端测试..."
	@cd $(BACKEND_DIR) && go test ./...
	@echo "运行前端测试..."
	@cd $(FRONTEND_DIR) && npm test
	@echo "$(GREEN)✅ 测试完成$(RESET)"

## 构建项目
build:
	@echo "$(CYAN)🔨 构建项目...$(RESET)"
	@echo "构建后端..."
	@cd $(BACKEND_DIR) && go build -o bin/$(PROJECT_NAME) main.go
	@echo "构建前端..."
	@cd $(FRONTEND_DIR) && npm run build
	@echo "$(GREEN)✅ 构建完成$(RESET)"

## 构建Docker镜像
docker-build:
	@echo "$(CYAN)🐳 构建Docker镜像...$(RESET)"
	@$(DOCKER_COMPOSE) build
	@echo "$(GREEN)✅ Docker镜像构建完成$(RESET)"

## 推送Docker镜像
docker-push:
	@echo "$(CYAN)📤 推送Docker镜像...$(RESET)"
	@$(DOCKER_COMPOSE) push
	@echo "$(GREEN)✅ Docker镜像推送完成$(RESET)"

## 格式化代码
format:
	@echo "$(CYAN)🎨 格式化代码...$(RESET)"
	@echo "格式化后端代码..."
	@cd $(BACKEND_DIR) && go fmt ./...
	@echo "格式化前端代码..."
	@cd $(FRONTEND_DIR) && npm run format
	@echo "$(GREEN)✅ 代码格式化完成$(RESET)"

## 代码检查
lint:
	@echo "$(CYAN)🔍 代码检查...$(RESET)"
	@echo "检查后端代码..."
	@cd $(BACKEND_DIR) && go vet ./...
	@echo "检查前端代码..."
	@cd $(FRONTEND_DIR) && npm run lint
	@echo "$(GREEN)✅ 代码检查完成$(RESET)"

## 生成文档
docs:
	@echo "$(CYAN)📚 生成文档...$(RESET)"
	@echo "生成后端API文档..."
	@cd $(BACKEND_DIR) && swag init
	@echo "生成前端文档..."
	@cd $(FRONTEND_DIR) && npm run docs || true
	@echo "$(GREEN)✅ 文档生成完成$(RESET)"

## 数据库迁移
migrate:
	@echo "$(CYAN)🗄️  运行数据库迁移...$(RESET)"
	@cd $(BACKEND_DIR) && go run main.go migrate
	@echo "$(GREEN)✅ 数据库迁移完成$(RESET)"

## 重置数据库
db-reset:
	@echo "$(RED)🗄️  重置数据库...$(RESET)"
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d postgres redis
	@sleep 10
	@$(MAKE) migrate
	@echo "$(GREEN)✅ 数据库重置完成$(RESET)"

## 备份数据库
db-backup:
	@echo "$(CYAN)💾 备份数据库...$(RESET)"
	@mkdir -p backups
	@docker exec $$(docker-compose ps -q postgres) pg_dump -U postgres $(PROJECT_NAME) > backups/$(PROJECT_NAME)_$$(date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✅ 数据库备份完成$(RESET)"

## 安装开发工具
dev-tools:
	@echo "$(CYAN)🛠️  安装开发工具...$(RESET)"
	@echo "安装Air热重载工具..."
	@go install github.com/cosmtrek/air@latest
	@echo "安装Swag文档生成工具..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "安装Wire依赖注入工具..."
	@go install github.com/google/wire/cmd/wire@latest
	@echo "$(GREEN)✅ 开发工具安装完成$(RESET)"

## 快速设置 - 首次使用
setup: dev-tools install migrate
	@echo "$(GREEN)🎉 项目设置完成！$(RESET)"
	@echo ""
	@echo "下一步："
	@echo "  1. 启动开发环境: make dev"
	@echo "  2. 或启动生产环境: make start"
	@echo "  3. 访问应用: http://localhost:3000"
	@echo ""

## 完整重建
rebuild: clean docker-build start
	@echo "$(GREEN)🔄 完整重建完成$(RESET)"

## 健康检查
health:
	@echo "$(CYAN)🏥 健康检查...$(RESET)"
	@echo "检查后端服务..."
	@curl -f http://localhost:8080/api/health || echo "后端服务未响应"
	@echo "检查前端服务..."
	@curl -f http://localhost:3000 || echo "前端服务未响应"
	@echo "检查数据库连接..."
	@$(DOCKER_COMPOSE) exec postgres pg_isready -U postgres || echo "数据库连接失败"
	@echo "检查Redis连接..."
	@$(DOCKER_COMPOSE) exec redis redis-cli ping || echo "Redis连接失败"

## 性能测试
perf:
	@echo "$(CYAN)⚡ 性能测试...$(RESET)"
	@echo "后端API性能测试..."
	@ab -n 1000 -c 10 http://localhost:8080/api/health || echo "请安装 apache2-utils"
	@echo "$(GREEN)✅ 性能测试完成$(RESET)"

## 安全扫描
security:
	@echo "$(CYAN)🔒 安全扫描...$(RESET)"
	@echo "扫描后端依赖..."
	@cd $(BACKEND_DIR) && go list -json -m all | nancy sleuth || echo "请安装 nancy"
	@echo "扫描前端依赖..."
	@cd $(FRONTEND_DIR) && npm audit
	@echo "$(GREEN)✅ 安全扫描完成$(RESET)"

## 显示项目信息
info:
	@echo "$(CYAN)📋 项目信息$(RESET)"
	@echo "===================="
	@echo "项目名称: $(PROJECT_NAME)"
	@echo "项目路径: $(PWD)"
	@echo "后端目录: $(BACKEND_DIR)"
	@echo "前端目录: $(FRONTEND_DIR)"
	@echo ""
	@echo "$(CYAN)服务地址:$(RESET)"
	@echo "  前端: http://localhost:3000"
	@echo "  后端: http://localhost:8080"
	@echo "  文档: http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "$(CYAN)数据库:$(RESET)"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis: localhost:6379"
