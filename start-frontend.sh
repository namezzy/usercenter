#!/bin/bash

# 用户中心前端启动脚本

echo "正在启动用户中心前端服务..."

# 检查Node.js是否已安装
if ! command -v node &> /dev/null; then
    echo "错误: Node.js未安装，请先安装Node.js"
    exit 1
fi

# 检查npm是否已安装
if ! command -v npm &> /dev/null; then
    echo "错误: npm未安装，请先安装npm"
    exit 1
fi

# 进入前端目录
cd frontend

# 检查依赖是否已安装
if [ ! -d "node_modules" ]; then
    echo "正在安装依赖..."
    npm install
fi

# 启动开发服务器
echo "正在启动开发服务器..."
npm run dev
