import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api-client";
import { useAuthStore } from "@/stores/auth";
import type { UserProfile } from "@/types/api";

interface LoginPayload {
  email: string;
  password: string;
}

interface RegisterPayload {
  email: string;
  password: string;
  confirm_password: string;
  name: string;
}

export function useLogin() {
  const setSession = useAuthStore((s) => s.setSession);

  return useMutation({
    mutationFn: (payload: LoginPayload) =>
      apiFetch<{ access_token: string }>("/api/auth/login", {
        method: "POST",
        body: JSON.stringify(payload),
      }),
    onSuccess: (data) => setSession(data.access_token),
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: (payload: RegisterPayload) =>
      apiFetch<{ message: string }>("/api/auth/register", {
        method: "POST",
        body: JSON.stringify(payload),
      }),
  });
}

export function useForgotPassword() {
  return useMutation({
    mutationFn: (email: string) =>
      apiFetch<{ message: string }>("/api/auth/forgot-password", {
        method: "POST",
        body: JSON.stringify({ email }),
      }),
  });
}

export function useLogout() {
  const clearSession = useAuthStore((s) => s.clearSession);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => apiFetch("/api/auth/logout", { method: "POST" }),
    onSuccess: () => {
      clearSession();
      queryClient.clear();
    },
  });
}

export function useProfile() {
  const accessToken = useAuthStore((s) => s.accessToken);

  return useQuery({
    queryKey: ["profile"],
    queryFn: () => apiFetch<UserProfile>("/api/profile", { auth: true }),
    enabled: !!accessToken,
    retry: false,
  });
}

interface UpdateProfilePayload {
  name: string;
  onboarding_focus: string;
  preferred_language: string;
  preferred_currency: string;
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: UpdateProfilePayload) =>
      apiFetch<UserProfile>("/api/profile", {
        method: "PATCH",
        auth: true,
        body: JSON.stringify(payload),
      }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["profile"] }),
  });
}

export function useDeleteAccount() {
  const clearSession = useAuthStore((s) => s.clearSession);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (password: string) =>
      apiFetch("/api/profile", {
        method: "DELETE",
        auth: true,
        body: JSON.stringify({ password }),
      }),
    onSuccess: () => {
      clearSession();
      queryClient.clear();
    },
  });
}
