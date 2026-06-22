// Wails runtime bindings for native session management
// These functions are exposed by the Go backend via Wails v3

export interface NativeSession {
  id: string;
  baseURL: string;
  lessonCount: number;
  currentLessonID: string | null;
}

export async function createNativeSession(): Promise<NativeSession> {
  // @ts-ignore - Wails runtime binding
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.CreateSession');
  }
  throw new Error('Wails runtime not available');
}

export async function deleteNativeSession(sessionId: string): Promise<void> {
  // @ts-ignore
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.DeleteSession', sessionId);
  }
}

export async function attachExistingNativeSession(
  sessionId: string
): Promise<NativeSession> {
  // @ts-ignore
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.AttachSession', sessionId);
  }
  throw new Error('Wails runtime not available');
}

export async function listNativeSessions(): Promise<NativeSession[]> {
  // @ts-ignore
  if (window.runtime?.call) {
    return window.runtime.call('sidecar.ListSessions');
  }
  return [];
}
