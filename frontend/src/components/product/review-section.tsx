"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { IconFlag, IconStarFilled } from "@tabler/icons-react";
import { Button } from "@/components/ui/button";
import { Link, useRouter } from "@/i18n/navigation";
import { useAuthStore } from "@/stores/auth";
import { useCreateReview, useReportReview } from "@/hooks/use-reviews";
import { ApiError } from "@/lib/api-client";
import { revalidateProduct } from "@/lib/revalidate";
import type { Review } from "@/types/api";

const FIT_OPTIONS = ["kekecilan", "pas", "kebesaran"] as const;

export function ReviewSection({ productId, productSlug, reviews }: { productId: string; productSlug: string; reviews: Review[] }) {
  const t = useTranslations("Product");
  const router = useRouter();
  const accessToken = useAuthStore((s) => s.accessToken);
  const createReview = useCreateReview(productSlug);
  const reportReview = useReportReview();

  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState("");
  const [fitFeedback, setFitFeedback] = useState<(typeof FIT_OPTIONS)[number]>("pas");
  const [showForm, setShowForm] = useState(false);

  function submit() {
    createReview.mutate(
      { product_id: productId, rating, comment, fit_feedback: fitFeedback },
      {
        onSuccess: async () => {
          toast.success(t("reviewSubmitted"));
          setShowForm(false);
          setComment("");
          await revalidateProduct(productSlug);
          router.refresh();
        },
        onError: (err) => toast.error(err instanceof ApiError ? err.message : "Gagal mengirim review"),
      },
    );
  }

  function report(reviewId: string) {
    reportReview.mutate(reviewId, {
      onSuccess: () => toast.success(t("reportSubmitted")),
    });
  }

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="font-display text-xl font-semibold">{t("reviews")}</h2>
        {accessToken ? (
          !showForm && (
            <Button size="sm" variant="secondary" onClick={() => setShowForm(true)}>
              {t("writeReview")}
            </Button>
          )
        ) : (
          <Link href="/login" className="text-sm underline">
            {t("loginToReview")}
          </Link>
        )}
      </div>

      {showForm && (
        <div className="mb-6 rounded-[var(--radius-card)] border border-border p-4">
          <label className="mb-1.5 block text-sm font-medium">{t("rating")}</label>
          <div className="mb-3 flex gap-1">
            {[1, 2, 3, 4, 5].map((n) => (
              <button key={n} type="button" onClick={() => setRating(n)}>
                <IconStarFilled size={20} className={n <= rating ? "" : "opacity-20"} />
              </button>
            ))}
          </div>

          <label className="mb-1.5 block text-sm font-medium">{t("fitFeedback")}</label>
          <div className="mb-3 flex gap-2">
            {FIT_OPTIONS.map((f) => (
              <button
                key={f}
                type="button"
                onClick={() => setFitFeedback(f)}
                className={`rounded-[var(--radius-control)] border border-zinc-500/40 px-3 py-1.5 text-sm ${
                  fitFeedback === f ? "bg-off-black text-pure-white dark:bg-pure-white dark:text-off-black" : ""
                }`}
              >
                {f === "kekecilan" ? t("fitSmall") : f === "pas" ? t("fitTrue") : t("fitBig")}
              </button>
            ))}
          </div>

          <label htmlFor="review-comment" className="mb-1.5 block text-sm font-medium">
            {t("comment")}
          </label>
          <textarea
            id="review-comment"
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            rows={3}
            className="mb-3 w-full rounded-[var(--radius-control)] border border-zinc-500/40 bg-background p-3 text-sm outline-none focus:ring-2 focus:ring-off-black dark:focus:ring-pure-white"
          />

          <Button size="sm" onClick={submit} disabled={createReview.isPending}>
            {t("submitReview")}
          </Button>
        </div>
      )}

      {reviews.length === 0 ? (
        <p className="text-sm text-muted">{t("noReviews")}</p>
      ) : (
        <ul className="flex flex-col gap-4">
          {reviews.map((r) => (
            <li key={r.id} className="border-b border-border pb-4 last:border-b-0">
              <div className="mb-1 flex items-center justify-between">
                <p className="flex items-center gap-1 text-sm font-medium">
                  <IconStarFilled size={14} />
                  {r.rating}
                  <span className="text-muted">· {r.user_name}</span>
                  {r.fit_feedback && (
                    <span className="text-xs uppercase tracking-[0.01em] text-muted">· {r.fit_feedback}</span>
                  )}
                </p>
                <button
                  onClick={() => report(r.id)}
                  className="flex items-center gap-1 text-xs text-muted hover:text-foreground"
                >
                  <IconFlag size={12} />
                  {t("report")}
                </button>
              </div>
              <p className="text-sm">{r.comment}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
