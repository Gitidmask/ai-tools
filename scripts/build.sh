#!/usr/bin/env bash
# AI 工具箱 - 完整构建脚本
# 构建 sidecar (Bun → exe) + 主应用 (Go + Wails v3 + React 前端)
set -euo pipefail

APP_NAME="ai_tools"
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BUILD_DIR="${ROOT_DIR}/build"
SIDECAR_DIR="${ROOT_DIR}/sidecar"
FRONTEND_DIR="${ROOT_DIR}/frontend"
SIDECAR_OUT="${ROOT_DIR}/sidecars/claude-sidecar-x86_64-pc-windows-msvc.exe"

echo "=== AI 工具箱 构建脚本 ==="
echo "Root: ${ROOT_DIR}"

# ---- Step 1: Sidecar (Bun -> exe) ----
echo ""
echo "--- Step 1/3: 构建 Sidecar ---"

if command -v bun &> /dev/null; then
    echo "Bun 已安装: $(bun --version)"

    if [ -f "${SIDECAR_DIR}/package.json" ]; then
        cd "${SIDECAR_DIR}"
        echo "安装 sidecar 依赖..."
        bun install

        echo "编译 sidecar..."
        mkdir -p "$(dirname "${SIDECAR_OUT}")"
        bun build --compile --target=browser \
            --outfile="${SIDECAR_OUT}" \
            ./src/index.ts
        echo "Sidecar 编译完成: ${SIDECAR_OUT}"
    else
        echo "警告: 未找到 sidecar/package.json，跳过 sidecar 构建"
    fi
else
    echo "警告: Bun 未安装，跳过 sidecar 构建"
    echo "  请安装 Bun: https://bun.sh/"
    echo "  或手动运行: cd sidecar && bun build --compile ./src/index.ts"
fi

# ---- Step 2: Frontend (Vite build) ----
echo ""
echo "--- Step 2/3: 构建前端 ---"

if [ -f "${FRONTEND_DIR}/package.json" ]; then
    cd "${FRONTEND_DIR}"
    if [ ! -d "node_modules" ]; then
        echo "安装前端依赖..."
        npm install
    fi
    echo "构建前端..."
    npm run build
    echo "前端构建完成"
else
    echo "警告: 未找到 frontend/package.json，使用已存在的前端构建产物"
fi

# ---- Step 3: Go 主应用 (Wails build) ----
echo ""
echo "--- Step 3/3: 构建主应用 (Wails) ---"

cd "${ROOT_DIR}"
echo "使用 Wails 构建桌面应用..."
wails build \
    -name "${APP_NAME}" \
    -platform windows/amd64 \
    -o "${BUILD_DIR}/${APP_NAME}.exe"

echo ""
echo "=== 构建完成! ==="
echo "主应用: ${BUILD_DIR}/${APP_NAME}.exe"
echo "Sidecar: ${SIDECAR_OUT}"
echo ""
echo "启动: ${BUILD_DIR}/${APP_NAME}.exe"
