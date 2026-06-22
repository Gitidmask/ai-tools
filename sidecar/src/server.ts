/**
 * DevBridge HTTP 服务器
 * 
 * 提供主进程 (Go Wails app) 需要的 API 端点：
 * - /api/sidecar-config   配置信息
 * - /api/health           健康检查
 * - /api/session          会话管理
 * - /api/terminal/ws      WebSocket 终端
 */

import type { Server } from 'bun';

export interface ServerOptions {
  host: string;
  port: number;
}

export interface SidecarServer {
  port: number;
  stop: () => Promise<void>;
}

/**
 * 启动 DevBridge HTTP 服务器
 */
export async function startServer(options: ServerOptions): Promise<SidecarServer> {
  const { host, port } = options;

  const bunServer = Bun.serve({
    hostname: host,
    port,
    websocket: {
      open(ws) {
        console.log('[ws] client connected');
      },
      message(ws, message) {
        console.log('[ws] received:', message);
        ws.send(`echo: ${message}`);
      },
      close(ws, code, reason) {
        console.log(`[ws] client disconnected: ${code} ${reason}`);
      },
    },
    async fetch(req) {
      const url = new URL(req.url);
      const path = url.pathname;

      // CORS headers
      const corsHeaders = {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization, New-Api-User',
      };

      // Handle CORS preflight
      if (req.method === 'OPTIONS') {
        return new Response(null, { headers: corsHeaders });
      }

      try {
        // Route handling
        if (path === '/api/health' || path === '/health') {
          return handleHealth(req, corsHeaders);
        }
        if (path === '/api/sidecar-config') {
          return handleSidecarConfig(req, corsHeaders);
        }
        if (path.startsWith('/api/session')) {
          return handleSession(req, corsHeaders);
        }
        if (path.startsWith('/terminal') || path.startsWith('/api/terminal')) {
          return handleTerminal(req, url);
        }

        // Claude Code internal routes are proxied
        if (path.startsWith('/api/')) {
          return await proxyToClaudeCode(req, path, corsHeaders);
        }

        // Default: 404
        return new Response(JSON.stringify({ error: 'not found', path }), {
          status: 404,
          headers: { 'Content-Type': 'application/json', ...corsHeaders },
        });
      } catch (err) {
        console.error('[devbridge] error:', err);
        return new Response(JSON.stringify({ error: String(err) }), {
          status: 500,
          headers: { 'Content-Type': 'application/json', ...corsHeaders },
        });
      }
    },
  });

  console.log(`[devbridge] listening on ${host}:${bunServer.port}`);

  return {
    port: bunServer.port,
    stop: async () => {
      bunServer.stop();
    },
  };
}

/**
 * GET /api/health
 * 健康检查端点 - 主进程定期调用确认 sidecar 存活
 */
function handleHealth(req: Request, cors: Record<string, string>): Response {
  return Response.json(
    {
      status: 'ok',
      timestamp: Date.now(),
      version: '1.0.0',
      pid: process.pid,
    },
    { headers: { 'Content-Type': 'application/json', ...cors } }
  );
}

/**
 * GET /api/sidecar-config
 * 返回 sidecar 配置信息
 */
function handleSidecarConfig(req: Request, cors: Record<string, string>): Response {
  return Response.json(
    {
      version: '1.0.0',
      platform: process.platform,
      arch: process.arch,
      pid: process.pid,
      uptime: process.uptime(),
      env: {
        NODE_ENV: process.env.NODE_ENV || 'production',
      },
    },
    { headers: { 'Content-Type': 'application/json', ...cors } }
  );
}

/**
 * POST /api/session
 * 会话管理 - 创建/查询/删除 AI 会话
 */
async function handleSession(req: Request, cors: Record<string, string>): Promise<Response> {
  const url = new URL(req.url);
  const method = req.method;

  if (method === 'GET') {
    // List sessions
    return Response.json({ sessions: [] }, { headers: { 'Content-Type': 'application/json', ...cors } });
  }

  if (method === 'POST') {
    // Create new session
    const body = await req.json().catch(() => ({}));
    return Response.json(
      {
        id: crypto.randomUUID(),
        baseURL: `http://127.0.0.1:${process.env.SIDECAR_PORT || '9800'}`,
        createdAt: Date.now(),
        ...body,
      },
      { headers: { 'Content-Type': 'application/json', ...cors }, status: 201 }
    );
  }

  return new Response('Method not allowed', { status: 405, headers: cors });
}

/**
 * WebSocket 终端处理
 * 返回 terminal.html 或升级为 WebSocket 连接
 */
function handleTerminal(req: Request, url: URL): Response {
  if (req.headers.get('Upgrade') === 'websocket') {
    // The Bun.serve websocket handler will take over
    return new Response(null, { status: 101 });
  }

  // Return terminal HTML page
  return new Response(
    `<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="UTF-8"/><title>AI 工具箱 · 终端</title>
<style>
html,body,#root{height:100%;margin:0;padding:0;background:#0d1017;color:#bfbdb6;font-family:monospace;font-size:13px;}
</style></head>
<body><div id="root"><h2>终端</h2><p>WebSocket: ${url.origin}/api/terminal/ws</p></div></body>
</html>`,
    {
      headers: { 'Content-Type': 'text/html; charset=utf-8' },
    }
  );
}

/**
 * 代理请求到 Claude Code 内部服务
 */
async function proxyToClaudeCode(req: Request, path: string, cors: Record<string, string>): Promise<Response> {
  // Claude Code runs as the sidecar process itself
  // API calls that need Claude Code's capabilities are handled here
  // For now, return a placeholder indicating the service is available
  return Response.json(
    {
      service: 'claude-code',
      path,
      available: true,
      message: 'Claude Code engine loaded and ready',
    },
    { headers: { 'Content-Type': 'application/json', ...cors } }
  );
}
