import React from 'react';

interface Message {
  id: string;
  role: 'user' | 'assistant' | 'tool';
  content: string;
  toolCalls?: ToolCall[];
  timestamp: number;
}

interface ToolCall {
  id: string;
  toolName: string;
  input: Record<string, unknown>;
  result?: string;
  isError?: boolean;
  status?: 'pending' | 'running' | 'success' | 'error';
}

interface MessageListProps {
  messages: Message[];
}

export default function MessageList({ messages }: MessageListProps) {
  const listRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [messages]);

  return (
    <div className="message-list" ref={listRef}>
      {messages.map((msg) => (
        <div key={msg.id} className={`message message-${msg.role}`}>
          <div className="message-avatar">
            {msg.role === 'user' ? '👤' : msg.role === 'assistant' ? '🤖' : '🔧'}
          </div>
          <div className="message-content">
            <div className="message-text">{msg.content}</div>
            {msg.toolCalls?.map((tc) => (
              <ToolCallBlock key={tc.id} call={tc} />
            ))}
          </div>
        </div>
      ))}
      {messages.length === 0 && (
        <div className="message-empty">
          <p>开始一段新的对话</p>
        </div>
      )}
    </div>
  );
}

function ToolCallBlock({ call }: { call: ToolCall }) {
  const [expanded, setExpanded] = React.useState(call.isError || false);

  return (
    <div className={`tool-call-block ${call.status || 'pending'}`}>
      <div className="tool-call-header" onClick={() => setExpanded(!expanded)}>
        <span className="tool-call-icon">🔌</span>
        <span className="tool-call-name">{call.toolName}</span>
        <span className="tool-call-status">
          {call.status === 'running' && '⏳'}
          {call.status === 'success' && '✅'}
          {call.status === 'error' && '❌'}
          {call.status === 'pending' && '⏸️'}
        </span>
      </div>
      {expanded && (
        <div className="tool-call-details">
          {call.result && (
            <div className={`tool-call-result ${call.isError ? 'error' : ''}`}>
              {call.result}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
