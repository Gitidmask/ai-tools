/**
 * AI 工具箱 - Sidecar 主入口
 * 
 * 本程序作为桌面应用的 sidecar 子进程运行：
 * 1. 启动 Claude Code CLI 引擎 (@anthropic-ai/claude-code)
 * 2. 提供 DevBridge HTTP API 供主进程通信
 * 3. 管理 AI 会话生命周期
 * 
 * 构建: bun build --compile --target=browser --outfile=../sidecars/claude-sidecar.exe ./src/index.ts
 */

import { startServer } from './server';
import { logger } from './utils/logger';

async function main() {
  const sidecarPort = parseInt(process.env.SIDECAR_PORT || '0', 10);
  const listenAddr = process.env.LISTEN_ADDR || '127.0.0.1';

  logger.info('sidecar starting...', { port: sidecarPort, addr: listenAddr });

  // Start DevBridge HTTP server
  const server = await startServer({
    host: listenAddr,
    port: sidecarPort,
  });

  const assignedPort = server.port;
  logger.info(`devbridge started, available at http://${listenAddr}:${assignedPort}/api/sidecar-config`);

  // Graceful shutdown
  const shutdown = async () => {
    logger.info('sidecar shutting down...');
    await server.stop();
    process.exit(0);
  };

  process.on('SIGINT', shutdown);
  process.on('SIGTERM', shutdown);

  // Keep alive - the sidecar runs as a long-lived process
  // The Claude Code CLI will be loaded on-demand when sessions are created
}

main().catch((err) => {
  logger.error('sidecar fatal error', { error: String(err) });
  process.exit(1);
});
