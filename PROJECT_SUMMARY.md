# AI 工具箱项目重建完成

## 重建内容

### Go 后端 (12 个业务包 + 1 个数据库层)

```
code/
├── cmd/main.go                          # 入口: 初始化 DB, 注册 12 个 Service, 启动窗口
├── internal/client/
│   ├── db/db.go                         # SQLite 初始化 + 12 张表的自动迁移
│   ├── newapi/service.go                # 多 API 账户管理 (ListAccounts, CRUD)
│   ├── claude/service.go                # Claude Preset 管理
│   ├── codex/service.go                 # Codex 配置 (Profiles, ProviderConfig)
│   ├── sidecar/service.go               # Sidecar 进程生命周期管理
│   ├── gateway/service.go               # AI 网关反向代理
│   ├── skills/service.go                # MCP 技能管理
│   ├── env/service.go                   # 环境信息
│   ├── npm/service.go                   # npm/Node.js 安装管理
│   ├── project/service.go               # 项目管理
│   ├── filesystem/service.go            # 文件系统操作
│   ├── clihistory/service.go            # CLI 会话历史
│   ├── heartbeat/service.go             # 健康检查
│   └── download/service.go              # 文件下载
├── go.mod                               # Go 模块定义
├── wails.json                           # Wails 项目配置
└── README.md                            # 构建说明
```

### React 前端 (7 个组件 + 3 个页面 + 2 个 Store)

```
frontend/
├── index.html, package.json, vite.config.ts, tsconfig.json
├── dist/assets/                         # 复用的前端构建产物
├── src/
│   ├── main.tsx                         # 入口
│   ├── components/
│   │   ├── App.tsx                      # 主应用壳 (含 Tab 导航)
│   │   ├── ChatInput.tsx                # 聊天输入框
│   │   ├── MessageList.tsx              # 消息列表 + ToolCallBlock
│   │   ├── ToolCallBlock.tsx            # 工具调用 UI (状态机渲染)
│   │   ├── ToolResultBlock.tsx          # 工具结果 UI
│   │   ├── AutoFix.tsx                  # 自动修复提示
│   │   └── ExtractedComponent.tsx       # 组件扩展占位
│   ├── pages/
│   │   ├── NativeChatTab.tsx            # 原生对话页
│   │   ├── TerminalTab.tsx              # 终端页 (占位)
│   │   └── WorkspacePanel.tsx           # 工作区面板
│   ├── stores/
│   │   ├── chatStore.ts                 # 聊天状态 (Zustand)
│   │   └── settingsStore.ts             # 设置状态
│   ├── styles/main.css                  # 全局样式 (暗色主题)
│   ├── utils/nativeSession.ts           # Wails 原生会话绑定
│   └── lib/wails.d.ts                   # Wails 运行时类型声明
```

### 复用的二进制文件

```
├── sidecars/claude-sidecar-x86_64-pc-windows-msvc.exe
├── uninstall.exe
└── frontend/dist/assets/ (7 个前端构建产物: main.js, react.js, dist.js 等)
```

## 数据库 Schema (12 张表)

| 表名 | 用途 |
|------|------|
| codex_profiles | AI 提供商配置 (API Key, 模型, Base URL) |
| cli_sessions | CLI 会话 (Claude/Codex 双源) |
| cli_session_tabs | 会话标签页 |
| accounts | 账号池 (配额, 冷却, 分级限流) |
| codex_tokens | OAuth 令牌 |
| codex_extra_providers | 自定义提供商 |
| codex_provider_config | 当前激活的提供商配置 |
| claude_presets | Claude 预设 (模型/颜色/图标/排序) |
| gateway_config | AI 网关反向代理配置 |
| gateway_keys | 网关 API 密钥 |
| projects | 工作区项目 |
| app_settings | 应用设置 (K/V) |

## 不包含的文件（逆向工具，非软件本身）

- `*.py` — 逆向分析脚本（pe_analyzer, focused_analysis, extract_* 等）
- `*.txt` — 分析报告（pe_analysis_report, focused_analysis, pe_summary）
- `*.zip` — 压缩包（ai_toolbox_source_code, 逆向分析报告）
- `ai_tools逆向分析报告.zip`
