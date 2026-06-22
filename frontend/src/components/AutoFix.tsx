import React from 'react';

interface AutoFixProps {
  issue: string;
  onFix: () => void;
  onDismiss: () => void;
}

export default function AutoFix({ issue, onFix, onDismiss }: AutoFixProps) {
  return (
    <div className="autofix-banner">
      <span className="autofix-icon">🔧</span>
      <span className="autofix-text">{issue}</span>
      <button className="autofix-btn" onClick={onFix}>
        自动修复
      </button>
      <button className="autofix-dismiss" onClick={onDismiss}>
        ✕
      </button>
    </div>
  );
}
