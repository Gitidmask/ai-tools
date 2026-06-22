// Wails runtime bindings for native session management
// These functions are exposed by the Go backend via Wails v3

export interface NativeSession {
  id: string;
  baseURL: string;
  lessonCount: number;
  currentLessonID: string | null;
}

export async function createNativeSession(): Promise<NativeSession> {
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.CreateSession') as Promise<NativeSession>;
  }
  throw new Error('Wails runtime not available');
}

export async function deleteNativeSession(sessionId: string): Promise<void> {
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.DeleteSession', sessionId) as Promise<void>;
  }
}

export async function attachExistingNativeSession(
  sessionId: string
): Promise<NativeSession> {
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.AttachSession', sessionId) as Promise<NativeSession>;
  }
  throw new Error('Wails runtime not available');
}

export async function listNativeSessions(): Promise<NativeSession[]> {
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.ListSessions') as Promise<NativeSession[]>;
  }
  return [];
}
