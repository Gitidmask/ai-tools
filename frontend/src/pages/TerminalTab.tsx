import React, { useEffect, useRef } from 'react';

export default function TerminalTab() {
  const terminalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Terminal initialization would go here
    // Using xterm.js or similar in production
  }, []);

  return (
    <div className="terminal-tab" ref={terminalRef}>
      <div className="terminal-header">终端</div>
      <div className="terminal-body">
        <div className="terminal-placeholder">
          <p>终端标签页 - 用于执行命令和调试</p>
          <p className="terminal-hint">
            连接 WebSocket: ws://127.0.0.1:9800/terminal
          </p>
        </div>
      </div>
    </div>
  );
}
