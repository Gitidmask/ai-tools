import React from 'react';

interface ToolResultBlockProps {
  content: string;
  isError?: boolean;
  toolName?: string;
}

export default function ToolResultBlock({
  content,
  isError,
  toolName,
}: ToolResultBlockProps) {
  const truncated = content.length > 200;
  const [expanded, setExpanded] = React.useState(isError || !truncated);

  return (
    <div className={`tool-result-block ${isError ? 'error' : ''}`}>
      {toolName && <div className="tool-result-header">{toolName}</div>}
      <div className="tool-result-content">
        {expanded ? content : content.slice(0, 200) + '...'}
      </div>
      {truncated && (
        <button className="tool-result-expand" onClick={() => setExpanded(!expanded)}>
          {expanded ? '收起' : '展开全部'}
        </button>
      )}
    </div>
  );
}
