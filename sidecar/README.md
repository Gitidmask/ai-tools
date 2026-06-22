# claude-sidecar 分析

## 本质
**claude-sidecar-x86_64-pc-windows-msvc.exe** = Bun 运行时 + `@anthropic-ai/claude-code` npm 包，通过 `bun build --compile` 打包为单文件 exe。

## 为什么不直接还原源码
这个二进制是**第三方闭源软件 (Claude Code CLI)** 的编译产物，不是 AI 工具箱项目自己开发的。还原它等于逆向 Claude Code 本身，这不合理也无必要。

## 替代方案
在 `code/sidecar/` 中创建一个**自定义 sidecar**，提供相同的接口，但移除 Claude Code 闭源代码：

| 原始功能 | 替代实现 |
|----------|----------|
| Bun 运行时 | 保留 Bun，但用我们自己写的 TypeScript |
| Claude Code CLI | ❌ 移除 |
| DevBridge API | ✅ 自定义实现 |
| MCP 服务器管理 | ✅ 自定义实现 |
| WebSocket 终端 | ✅ 自定义实现 |
| 健康检查 | ✅ 自定义实现 |
