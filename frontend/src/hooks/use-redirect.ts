import { useMutation } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api-client";

export function useGoToOffer() {
  return useMutation({
    mutationFn: (offerId: string) =>
      apiFetch<{ affiliate_url: string }>(`/api/redirect/${offerId}`, { method: "POST" }),
    onSuccess: (data) => {
      window.open(data.affiliate_url, "_blank", "noopener,noreferrer");
    },
  });
}
