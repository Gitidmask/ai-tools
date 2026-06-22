// Wails runtime type declarations

export {};

interface WailsRuntime {
  call: (method: string, ...args: unknown[]) => Promise<unknown>;
  events: {
    on: (event: string, callback: (...args: unknown[]) => void) => void;
    off: (event: string, callback: (...args: unknown[]) => void) => void;
    emit: (event: string, ...args: unknown[]) => void;
  };
  window: {
    close: () => void;
    minimize: () => void;
    maximize: () => void;
    unmaximize: () => void;
    isMaximised: () => Promise<boolean>;
    setTitle: (title: string) => void;
  };
}

declare global {
  interface Window {
    runtime?: WailsRuntime;
    _wails?: {
      environment: { OS: string; Arch: string; Debug: boolean };
      flags: Record<string, unknown>;
      dispatchWailsEvent: (event: string) => void;
    };
    chrome?: {
      webview?: {
        postMessage: (message: string) => void;
      };
    };
  }
}
