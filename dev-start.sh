#!/bin/bash

# 开发环境启动脚本
# Development Environment Start Script

set -e

echo "🔧 用户中心系统 - 开发环境启动"
echo "================================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函数定义
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查依赖
check_dev_dependencies() {
    print_info "检查开发环境依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go 1.21+"
        exit 1
    fi
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        print_error "Node.js 未安装，请先安装 Node.js 18+"
        exit 1
    fi
    
    # 检查npm
    if ! command -v npm &> /dev/null; then
        print_error "npm 未安装，请先安装 npm"
        exit 1
    fi
    
    print_success "开发环境依赖检查完成"
}

# 安装后端依赖
install_backend_deps() {
    print_info "安装后端依赖..."
    cd backend
    go mod download
    cd ..
    print_success "后端依赖安装完成"
}

# 安装前端依赖
install_frontend_deps() {
    print_info "安装前端依赖..."
    cd frontend
    if [ ! -d "node_modules" ]; then
        npm install
    else
        print_info "前端依赖已存在，跳过安装"
    fi
    cd ..
    print_success "前端依赖安装完成"
}

# 启动数据库服务
start_databases() {
    print_info "启动数据库和缓存服务..."
    docker-compose up -d postgres redis
    sleep 10
    print_success "数据库服务启动完成"
}

# 启动后端服务
start_backend() {
    print_info "启动后端开发服务..."
    cd backend
    
    # 检查是否安装了Air
    if command -v air &> /dev/null; then
        print_info "使用 Air 热重载启动后端..."
        air &
    else
        print_warning "Air 未安装，使用 go run 启动..."
        go run main.go &
    fi
    
    BACKEND_PID=$!
    cd ..
    
    print_success "后端服务启动完成 (PID: $BACKEND_PID)"
    echo $BACKEND_PID > .backend.pid
}

# 启动前端服务
start_frontend() {
    print_info "启动前端开发服务..."
    cd frontend
    npm run dev &
    FRONTEND_PID=$!
    cd ..
    
    print_success "前端服务启动完成 (PID: $FRONTEND_PID)"
    echo $FRONTEND_PID > .frontend.pid
}

# 停止开发服务
stop_dev_services() {
    print_info "停止开发服务..."
    
    # 停止后端
    if [ -f ".backend.pid" ]; then
        BACKEND_PID=$(cat .backend.pid)
        if kill -0 $BACKEND_PID 2>/dev/null; then
            kill $BACKEND_PID
            print_success "后端服务已停止"
        fi
        rm -f .backend.pid
    fi
    
    # 停止前端
    if [ -f ".frontend.pid" ]; then
        FRONTEND_PID=$(cat .frontend.pid)
        if kill -0 $FRONTEND_PID 2>/dev/null; then
            kill $FRONTEND_PID
            print_success "前端服务已停止"
        fi
        rm -f .frontend.pid
    fi
    
    # 停止数据库服务
    docker-compose down
    print_success "数据库服务已停止"
}

# 检查服务状态
check_dev_status() {
    print_info "检查开发服务状态..."
    
    # 检查数据库
    if docker-compose ps | grep postgres | grep Up >/dev/null 2>&1; then
        print_success "PostgreSQL 运行中"
    else
        print_warning "PostgreSQL 未运行"
    fi
    
    if docker-compose ps | grep redis | grep Up >/dev/null 2>&1; then
        print_success "Redis 运行中"
    else
        print_warning "Redis 未运行"
    fi
    
    # 检查后端
    if [ -f ".backend.pid" ]; then
        BACKEND_PID=$(cat .backend.pid)
        if kill -0 $BACKEND_PID 2>/dev/null; then
            print_success "后端服务运行中 (PID: $BACKEND_PID)"
        else
            print_warning "后端服务未运行"
        fi
    else
        print_warning "后端服务未运行"
    fi
    
    # 检查前端
    if [ -f ".frontend.pid" ]; then
        FRONTEND_PID=$(cat .frontend.pid)
        if kill -0 $FRONTEND_PID 2>/dev/null; then
            print_success "前端服务运行中 (PID: $FRONTEND_PID)"
        else
            print_warning "前端服务未运行"
        fi
    else
        print_warning "前端服务未运行"
    fi
}

# 显示开发信息
show_dev_info() {
    echo ""
    echo "🎉 开发环境启动成功！"
    echo "===================="
    echo ""
    echo "📱 访问地址:"
    echo "   前端开发服务: http://localhost:5173"
    echo "   后端API服务:  http://localhost:8080"
    echo "   API文档:      http://localhost:8080/swagger/index.html"
    echo ""
    echo "🗄️  数据库连接:"
    echo "   PostgreSQL: localhost:5432 (用户名: postgres)"
    echo "   Redis:      localhost:6379"
    echo ""
    echo "🔧 开发工具:"
    echo "   热重载: 前端和后端都支持热重载"
    echo "   调试:   可在IDE中设置断点调试"
    echo "   测试:   npm test (前端) | go test ./... (后端)"
    echo ""
    echo "📖 常用命令:"
    echo "   停止服务: $0 stop"
    echo "   查看状态: $0 status"
    echo "   重启服务: $0 restart"
    echo "   查看日志: $0 logs"
}

# 显示日志
show_logs() {
    print_info "显示服务日志..."
    echo ""
    echo "=== 数据库日志 ==="
    docker-compose logs --tail=20 postgres redis
    echo ""
    echo "=== 后端日志 ==="
    if [ -f ".backend.pid" ]; then
        print_info "后端服务正在运行，请查看终端输出"
    else
        print_warning "后端服务未运行"
    fi
    echo ""
    echo "=== 前端日志 ==="
    if [ -f ".frontend.pid" ]; then
        print_info "前端服务正在运行，请查看终端输出"
    else
        print_warning "前端服务未运行"
    fi
}

# 主函数
main() {
    echo ""
    print_info "开始启动开发环境..."
    echo ""
    
    # 检查依赖
    check_dev_dependencies
    
    # 安装依赖
    install_backend_deps
    install_frontend_deps
    
    # 启动数据库
    start_databases
    
    # 启动后端
    start_backend
    
    # 等待后端启动
    sleep 10
    
    # 启动前端
    start_frontend
    
    # 等待前端启动
    sleep 5
    
    # 显示开发信息
    show_dev_info
    
    # 等待用户输入
    echo ""
    print_info "按 Ctrl+C 停止所有服务"
    
    # 捕获中断信号
    trap 'echo ""; print_info "正在停止服务..."; stop_dev_services; exit 0' INT
    
    # 保持脚本运行
    while true; do
        sleep 1
    done
}

# 参数处理
case "${1:-start}" in
    "start")
        main
        ;;
    "stop")
        stop_dev_services
        ;;
    "restart")
        print_info "重启开发环境..."
        stop_dev_services
        sleep 3
        main
        ;;
    "status")
        check_dev_status
        ;;
    "logs")
        show_logs
        ;;
    "install")
        check_dev_dependencies
        install_backend_deps
        install_frontend_deps
        print_success "依赖安装完成"
        ;;
    "clean")
        print_warning "清理开发环境..."
        stop_dev_services
        rm -f .backend.pid .frontend.pid
        docker-compose down -v
        print_success "开发环境清理完成"
        ;;
    "help"|"-h"|"--help")
        echo "开发环境管理脚本"
        echo ""
        echo "用法: $0 [command]"
        echo ""
        echo "命令:"
        echo "  start    启动开发环境 (默认)"
        echo "  stop     停止开发环境"
        echo "  restart  重启开发环境"
        echo "  status   查看状态"
        echo "  logs     查看日志"
        echo "  install  安装依赖"
        echo "  clean    清理环境"
        echo "  help     显示帮助"
        echo ""
        ;;
    *)
        print_error "未知命令: $1"
        print_info "使用 '$0 help' 查看可用命令"
        exit 1
        ;;
esac
