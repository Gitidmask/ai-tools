import { create } from 'zustand';

interface Message {
  id: string;
  role: 'user' | 'assistant' | 'tool';
  content: string;
  toolCalls?: Array<{
    id: string;
    toolName: string;
    input: Record<string, unknown>;
    result?: string;
    isError?: boolean;
    status?: 'pending' | 'running' | 'success' | 'error';
  }>;
  timestamp: number;
}

interface ChatState {
  messages: Message[];
  isLoading: boolean;
  sessionId: string | null;
  addMessage: (msg: Message) => void;
  sendMessage: (text: string) => Promise<void>;
  clearMessages: () => void;
  setSessionId: (id: string) => void;
}

export const useChatStore = create<ChatState>((set, get) => ({
  messages: [],
  isLoading: false,
  sessionId: null,

  addMessage: (msg) => {
    set((state) => ({ messages: [...state.messages, msg] }));
  },

  sendMessage: async (text) => {
    set({ isLoading: true });
    const userMsg: Message = {
      id: `msg-${Date.now()}`,
      role: 'user',
      content: text,
      timestamp: Date.now(),
    };
    set((state) => ({ messages: [...state.messages, userMsg] }));

    try {
      // @ts-ignore - Wails runtime binding
      if (window.runtime?.call) {
        // TODO: implement actual AI provider call
        // const response = await window.runtime.call('newapi.SendMessage', text);
        const assistantMsg: Message = {
          id: `msg-${Date.now() + 1}`,
          role: 'assistant',
          content: `已收到: "${text}"\n\n(该消息为占位响应，实际 AI 集成将在此处处理)`,
          timestamp: Date.now(),
        };
        set((state) => ({ messages: [...state.messages, assistantMsg] }));
      }
    } catch (err) {
      const errorMsg: Message = {
        id: `msg-${Date.now() + 2}`,
        role: 'assistant',
        content: `错误: ${String(err)}`,
        timestamp: Date.now(),
      };
      set((state) => ({ messages: [...state.messages, errorMsg] }));
    } finally {
      set({ isLoading: false });
    }
  },

  clearMessages: () => set({ messages: [] }),

  setSessionId: (id) => set({ sessionId: id }),
}));
