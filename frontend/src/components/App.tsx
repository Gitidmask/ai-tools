import React, { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import NativeChatTab from '../pages/NativeChatTab';
import TerminalTab from '../pages/TerminalTab';
import WorkspacePanel from '../pages/WorkspacePanel';

const queryClient = new QueryClient();

type Tab = 'chat' | 'terminal' | 'workspace';

export default function App() {
  const [activeTab, setActiveTab] = useState<Tab>('chat');

  return (
    <QueryClientProvider client={queryClient}>
      <div className="app-container">
        <nav className="app-nav">
          <button
            className={`nav-btn ${activeTab === 'chat' ? 'active' : ''}`}
            onClick={() => setActiveTab('chat')}
          >
            💬 对话
          </button>
          <button
            className={`nav-btn ${activeTab === 'terminal' ? 'active' : ''}`}
            onClick={() => setActiveTab('terminal')}
          >
            ⌨️ 终端
          </button>
          <button
            className={`nav-btn ${activeTab === 'workspace' ? 'active' : ''}`}
            onClick={() => setActiveTab('workspace')}
          >
            📂 工作区
          </button>
        </nav>
        <main className="app-main">
          {activeTab === 'chat' && <NativeChatTab />}
          {activeTab === 'terminal' && <TerminalTab />}
          {activeTab === 'workspace' && <WorkspacePanel />}
        </main>
      </div>
    </QueryClientProvider>
  );
}
