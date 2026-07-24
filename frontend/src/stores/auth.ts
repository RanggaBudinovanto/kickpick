import { create } from "zustand";

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  onboarding_focus: string;
  preferred_language: string;
  preferred_currency: string;
  email_verified: boolean;
}

interface AuthState {
  accessToken: string | null;
  user: AuthUser | null;
  isBootstrapping: boolean;
  setSession: (accessToken: string, user?: AuthUser | null) => void;
  setUser: (user: AuthUser | null) => void;
  clearSession: () => void;
  finishBootstrap: () => void;
}

// Access token disimpan hanya di memory (bukan localStorage) sesuai Section 14 PRD,
// supaya tidak bisa dicuri lewat XSS. Hilang saat refresh halaman, dipulihkan via /api/auth/refresh.
export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  user: null,
  isBootstrapping: true,
  setSession: (accessToken, user) => set((state) => ({ accessToken, user: user ?? state.user })),
  setUser: (user) => set({ user }),
  clearSession: () => set({ accessToken: null, user: null }),
  finishBootstrap: () => set({ isBootstrapping: false }),
}));
