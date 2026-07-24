"use client";

import { useTranslations } from "next-intl";
import { IconHeart, IconHeartFilled } from "@tabler/icons-react";
import { Button } from "@/components/ui/button";
import { Link } from "@/i18n/navigation";
import { useAuthStore } from "@/stores/auth";
import { useAddWishlist, useRemoveWishlist, useWishlist } from "@/hooks/use-wishlist";

export function WishlistButton({ productId }: { productId: string }) {
  const t = useTranslations("Product");
  const accessToken = useAuthStore((s) => s.accessToken);
  const { data } = useWishlist();
  const addWishlist = useAddWishlist();
  const removeWishlist = useRemoveWishlist();

  if (!accessToken) {
    return (
      <Link href="/login">
        <Button variant="secondary">
          <IconHeart size={18} />
          {t("addToWishlist")}
        </Button>
      </Link>
    );
  }

  const existing = data?.data.find((w) => w.product_id === productId);

  if (existing) {
    return (
      <Button variant="secondary" onClick={() => removeWishlist.mutate(existing.id)}>
        <IconHeartFilled size={18} />
        {t("removeFromWishlist")}
      </Button>
    );
  }

  return (
    <Button variant="secondary" onClick={() => addWishlist.mutate(productId)} disabled={addWishlist.isPending}>
      <IconHeart size={18} />
      {t("addToWishlist")}
    </Button>
  );
}
