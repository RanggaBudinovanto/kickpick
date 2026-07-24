import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api-client";

interface CreateReviewPayload {
  product_id: string;
  rating: number;
  comment: string;
  fit_feedback: string;
}

export function useCreateReview(productSlug: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateReviewPayload) =>
      apiFetch("/api/reviews", { method: "POST", auth: true, body: JSON.stringify(payload) }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["product", productSlug] }),
  });
}

export function useReportReview() {
  return useMutation({
    mutationFn: (reviewId: string) =>
      apiFetch(`/api/reviews/${reviewId}/report`, { method: "POST", auth: true }),
  });
}
