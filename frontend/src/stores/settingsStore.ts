import { create } from 'zustand';

interface SettingsState {
  theme: 'dark' | 'light';
  listenAddr: string;
  activeProfile: string | null;
  setTheme: (theme: 'dark' | 'light') => void;
  setListenAddr: (addr: string) => void;
  setActiveProfile: (name: string | null) => void;
}

export const useSettingsStore = create<SettingsState>((set) => ({
  theme: 'dark',
  listenAddr: '127.0.0.1:9800',
  activeProfile: null,

  setTheme: (theme) => set({ theme }),
  setListenAddr: (addr) => set({ listenAddr: addr }),
  setActiveProfile: (name) => set({ activeProfile: name }),
}));
