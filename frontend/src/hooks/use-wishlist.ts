import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api-client";
import { useAuthStore } from "@/stores/auth";
import type { WishlistItem } from "@/types/api";

export function useWishlist() {
  const accessToken = useAuthStore((s) => s.accessToken);

  return useQuery({
    queryKey: ["wishlist"],
    queryFn: () => apiFetch<{ data: WishlistItem[] }>("/api/wishlist", { auth: true }),
    enabled: !!accessToken,
  });
}

export function useAddWishlist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (productId: string) =>
      apiFetch("/api/wishlist", {
        method: "POST",
        auth: true,
        body: JSON.stringify({ product_id: productId }),
      }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["wishlist"] }),
  });
}

export function useRemoveWishlist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => apiFetch(`/api/wishlist/${id}`, { method: "DELETE", auth: true }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["wishlist"] }),
  });
}

export function useSetWishlistAlert() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, alertActive, alertType }: { id: string; alertActive: boolean; alertType?: string }) =>
      apiFetch(`/api/wishlist/${id}/alert`, {
        method: "PATCH",
        auth: true,
        body: JSON.stringify({ alert_active: alertActive, alert_type: alertType ?? "" }),
      }),
    // Optimistic update: the checkbox is a controlled component keyed off
    // item.alert_active, so without this it visually snaps back to its old
    // state for the duration of the round-trip before the refetch lands.
    onMutate: async ({ id, alertActive }) => {
      await queryClient.cancelQueries({ queryKey: ["wishlist"] });
      const previous = queryClient.getQueryData<{ data: WishlistItem[] }>(["wishlist"]);

      queryClient.setQueryData<{ data: WishlistItem[] }>(["wishlist"], (old) =>
        old
          ? {
              data: old.data.map((item) =>
                item.id === id ? { ...item, alert_active: alertActive } : item,
              ),
            }
          : old,
      );

      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) {
        queryClient.setQueryData(["wishlist"], context.previous);
      }
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey: ["wishlist"] }),
  });
}
