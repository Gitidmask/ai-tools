import React from 'react';

interface ToolCallBlockProps {
  toolName: string;
  input: Record<string, unknown>;
  result?: string;
  isError?: boolean;
  status?: 'pending' | 'running' | 'success' | 'error';
}

export default function ToolCallBlock({
  toolName,
  input,
  result,
  isError,
  status = 'success',
}: ToolCallBlockProps) {
  const [expanded, setExpanded] = React.useState(isError || false);

  return (
    <div className={`tool-call-block ${status}`}>
      <div className="tool-call-header" onClick={() => setExpanded(!expanded)}>
        <span className="tool-call-icon">🔌</span>
        <span className="tool-call-name">{toolName}</span>
        <span className="tool-call-status">
          {status === 'running' && '⏳'}
          {status === 'success' && '✅'}
          {status === 'error' && '❌'}
        </span>
      </div>
      {expanded && (
        <div className="tool-call-details">
          <div className="tool-call-input">
            <pre>{JSON.stringify(input, null, 2)}</pre>
          </div>
          {result && (
            <div className={`tool-call-result ${isError ? 'error' : ''}`}>
              {result}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
