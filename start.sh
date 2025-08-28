#!/bin/bash

# 用户中心系统快速启动脚本
# Quick Start Script for User Center System

set -e

echo "🚀 用户中心系统快速启动脚本"
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
check_dependencies() {
    print_info "检查系统依赖..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    print_success "系统依赖检查完成"
}

# 停止现有服务
stop_services() {
    print_info "停止现有服务..."
    docker-compose down 2>/dev/null || true
    print_success "现有服务已停止"
}

# 启动服务
start_services() {
    print_info "启动用户中心系统服务..."
    
    # 启动数据库和Redis
    print_info "启动数据库和缓存服务..."
    docker-compose up -d postgres redis
    
    # 等待数据库启动
    print_info "等待数据库启动..."
    sleep 10
    
    # 启动后端服务
    print_info "启动后端服务..."
    docker-compose up -d backend
    
    # 等待后端启动
    print_info "等待后端服务启动..."
    sleep 15
    
    # 启动前端服务
    print_info "启动前端服务..."
    docker-compose up -d frontend
    
    print_success "所有服务启动完成"
}

# 检查服务状态
check_services() {
    print_info "检查服务状态..."
    
    # 检查服务状态
    docker-compose ps
    
    # 等待服务完全启动
    print_info "等待服务完全启动..."
    sleep 20
    
    # 检查后端健康状态
    print_info "检查后端服务..."
    if curl -f http://localhost:8080/api/health >/dev/null 2>&1; then
        print_success "后端服务运行正常"
    else
        print_warning "后端服务可能还在启动中，请稍等片刻"
    fi
    
    # 检查前端服务
    print_info "检查前端服务..."
    if curl -f http://localhost:3000 >/dev/null 2>&1; then
        print_success "前端服务运行正常"
    else
        print_warning "前端服务可能还在启动中，请稍等片刻"
    fi
}

# 显示访问信息
show_access_info() {
    echo ""
    echo "🎉 用户中心系统启动成功！"
    echo "=========================="
    echo ""
    echo "📱 访问地址:"
    echo "   前端应用: http://localhost:3000"
    echo "   后端API:  http://localhost:8080"
    echo "   API文档:  http://localhost:8080/swagger/index.html"
    echo ""
    echo "👤 默认账户:"
    echo "   管理员 - 用户名: admin, 密码: admin123456"
    echo "   普通用户 - 用户名: user, 密码: user123456"
    echo ""
    echo "🔧 常用命令:"
    echo "   查看日志: docker-compose logs -f"
    echo "   停止服务: docker-compose down"
    echo "   重启服务: docker-compose restart"
    echo "   查看状态: docker-compose ps"
    echo ""
    echo "📖 更多信息请查看 README.md"
}

# 主函数
main() {
    echo ""
    print_info "开始启动用户中心系统..."
    echo ""
    
    # 检查依赖
    check_dependencies
    
    # 停止现有服务
    stop_services
    
    # 启动服务
    start_services
    
    # 检查服务状态
    check_services
    
    # 显示访问信息
    show_access_info
}

# 参数处理
case "${1:-start}" in
    "start")
        main
        ;;
    "stop")
        print_info "停止用户中心系统..."
        docker-compose down
        print_success "系统已停止"
        ;;
    "restart")
        print_info "重启用户中心系统..."
        docker-compose down
        sleep 5
        main
        ;;
    "status")
        print_info "检查系统状态..."
        docker-compose ps
        echo ""
        print_info "服务日志:"
        docker-compose logs --tail=10
        ;;
    "logs")
        print_info "查看系统日志..."
        docker-compose logs -f
        ;;
    "clean")
        print_warning "清理系统数据和镜像..."
        read -p "确定要清理所有数据吗？这将删除数据库数据 (y/N): " confirm
        if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
            docker-compose down -v
            docker-compose rm -f
            docker system prune -f
            print_success "系统清理完成"
        else
            print_info "取消清理操作"
        fi
        ;;
    "help"|"-h"|"--help")
        echo "用户中心系统管理脚本"
        echo ""
        echo "用法: $0 [command]"
        echo ""
        echo "命令:"
        echo "  start    启动系统 (默认)"
        echo "  stop     停止系统"
        echo "  restart  重启系统"
        echo "  status   查看状态"
        echo "  logs     查看日志"
        echo "  clean    清理数据"
        echo "  help     显示帮助"
        echo ""
        ;;
    *)
        print_error "未知命令: $1"
        print_info "使用 '$0 help' 查看可用命令"
        exit 1
        ;;
esac
