import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api-client";
import { useAuthStore } from "@/stores/auth";
import type { NotificationItem } from "@/types/api";

export function useNotifications() {
  const accessToken = useAuthStore((s) => s.accessToken);

  return useQuery({
    queryKey: ["notifications"],
    queryFn: () => apiFetch<{ data: NotificationItem[] }>("/api/notifications", { auth: true }),
    enabled: !!accessToken,
  });
}

export function useUnreadCount() {
  const accessToken = useAuthStore((s) => s.accessToken);

  return useQuery({
    queryKey: ["notifications", "unread-count"],
    queryFn: () => apiFetch<{ count: number }>("/api/notifications/unread-count", { auth: true }),
    enabled: !!accessToken,
    refetchInterval: 30_000,
  });
}

export function useMarkNotificationRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      apiFetch(`/api/notifications/${id}/read`, { method: "PATCH", auth: true }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
  });
}
