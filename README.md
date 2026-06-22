# 项目构建说明

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.23+ |
| 桌面框架 | Wails v3 (v3.0.0-alpha.84) + WebView2 |
| 前端 | React 18 + TypeScript + Vite + Zustand |
| 数据库 | SQLite (modernc.org/sqlite, 纯 Go 实现) |
| Sidecar | Bun 运行时 (自定义 TypeScript DevBridge + @anthropic-ai/claude-code) |

## 目录结构

```
code/
├── cmd/main.go                 # 入口
├── internal/client/             # 主业务包
│   ├── db/              # SQLite 层 + 12 张表 schema
│   ├── newapi/          # 多 AI API 账户池管理
│   ├── claude/          # Claude Preset 管理
│   ├── codex/           # Codex 网关配置/Profiles
│   ├── sidecar/         # Sidecar 进程生命周期管理 (启动/停止/健康检查)
│   ├── devbridge/       # DevBridge HTTP API 服务
│   ├── gateway/         # AI 网关反向代理
│   ├── skills/          # MCP 技能管理
│   ├── env/             # 环境信息
│   ├── npm/             # npm 包管理
│   ├── project/         # 项目管理
│   ├── filesystem/      # 文件系统操作
│   ├── clihistory/      # CLI 会话历史
│   ├── heartbeat/       # 健康检查
│   └── download/        # 文件下载
├── frontend/             # React 前端 (TypeScript)
├── sidecar/              # Sidecar TypeScript 源码 (Bun 项目)
│   ├── src/index.ts     # 入口: 启动 DevBridge + Claude Code
│   ├── src/server.ts    # HTTP/WebSocket 服务器
│   ├── src/utils/       # 工具函数
│   └── package.json     # Bun 项目配置
├── sidecars/             # Sidecar 编译产物
├── scripts/              # 构建脚本
└── go.mod / wails.json
```

## 前置要求

- [Go](https://go.dev/) 1.23+
- [Bun](https://bun.sh/) 1.0+ (用于构建 sidecar)
- [Node.js](https://nodejs.org/) 18+ (用于前端构建)
- [Wails CLI](https://wails.io/) v3-alpha (用于打包桌面应用)
- WebView2 Runtime (Windows 10 1803+ 自带)
- Windows 10/11 (当前仅支持 Windows)

## 构建

### 一键构建 (推荐)

```bash
# Windows
scripts\build.bat

# Linux/macOS (仅用于交叉编译)
bash scripts/build.sh
```

### 分步构建

```bash
# 0. 确保用户已提供 Claude Code CLI
#    用户需要自行安装: npm install -g @anthropic-ai/claude-code
#    或让 sidecar 的 package.json 自动处理依赖

# 1. 构建 Sidecar (Bun → 单文件 exe)
cd sidecar
bun install                    # 安装依赖 (含 @anthropic-ai/claude-code)
bun build --compile --target=browser \
    --outfile=../sidecars/claude-sidecar-x86_64-pc-windows-msvc.exe \
    ./src/index.ts
cd ..

# 2. 安装前端依赖并构建
cd frontend
npm install
npm run build
cd ..

# 3. 安装 Wails CLI (首次)
go install github.com/wailsapp/wails/v3/cmd/wails@v3.0.0-alpha.84

# 4. 构建桌面应用
wails build -name ai_tools -platform windows/amd64

# 或开发模式
wails dev
```

### 开发模式

```bash
# 终端 1: 启动前端开发服务器
cd frontend
npm run dev

# 终端 2: 启动 sidecar 开发模式
cd sidecar
bun run dev

# 终端 3: 启动 Wails 开发模式
wails dev
```

## Sidecar 说明

sidecar (`code/sidecar/`) 是一个独立的 Bun TypeScript 项目，提供：

| 功能 | 端点 | 说明 |
|------|------|------|
| 健康检查 | `GET /api/health` | 主进程定期 ping 确认存活 |
| 配置信息 | `GET /api/sidecar-config` | 返回 sidecar 运行时信息 |
| 会话管理 | `POST/GET /api/session` | AI 会话创建/查询/删除 |
| WebSocket | `/api/terminal/ws` | 终端 WebSocket 连接 |
| Claude Code | (内部集成) | 通过 dependency 集成 |

构建为单文件 exe 后可独立运行，由主进程 (`ai_tools.exe`) 启动并管理生命周期。

## 环境要求

- Go 1.23+
- Node.js 18+
- Bun 1.0+
- WebView2 Runtime (Windows 10+ 自带)
- Windows 10/11 (当前仅支持 Windows)
