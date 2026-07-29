"use client";

import { useState, useEffect } from "react";
import { IconHeartOff } from "@tabler/icons-react";
import { AuthGuard } from "@/components/auth/auth-guard";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { useRemoveWishlist, useSetWishlistAlert, useWishlist } from "@/hooks/use-wishlist";

interface WishlistAlertCheckboxProps {
  item: {
    id: string;
    alert_active: boolean;
  };
  onMutate: (params: { id: string; alertActive: boolean; alertType: string }) => void;
}

function WishlistAlertCheckbox({ item, onMutate }: WishlistAlertCheckboxProps) {
  const [checked, setChecked] = useState(item.alert_active);

  useEffect(() => {
    setChecked(item.alert_active);
  }, [item.alert_active]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const nextChecked = e.target.checked;
    setChecked(nextChecked);
    onMutate({ id: item.id, alertActive: nextChecked, alertType: "restock" });
  };

  return (
    <input
      type="checkbox"
      checked={checked}
      onChange={handleChange}
    />
  );
}

function WishlistContent() {
  const { data, isLoading, isError, refetch } = useWishlist();
  const removeWishlist = useRemoveWishlist();
  const setAlert = useSetWishlistAlert();

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <h1 className="mb-6 font-display text-3xl font-bold tracking-[-0.01em]">Wishlist</h1>

      {isLoading && (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="aspect-[3/4] animate-pulse rounded-[var(--radius-card)] bg-surface" />
          ))}
        </div>
      )}

      {isError && (
        <div className="flex flex-col items-center gap-3 py-16 text-center">
          <p className="text-sm font-medium">Failed to load data. Please try again.</p>
          <Button size="sm" onClick={() => refetch()}>
            Try again
          </Button>
        </div>
      )}

      {!isLoading && !isError && data?.data.length === 0 && (
        <div className="flex flex-col items-center gap-3 py-16 text-center">
          <IconHeartOff size={32} />
          <p className="text-sm font-medium">No saved sneakers yet</p>
          <Link href="/cari">
            <Button size="sm">Start browsing sneakers</Button>
          </Link>
        </div>
      )}

      {!isLoading && !isError && data && data.data.length > 0 && (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          {data.data.map((item) => (
            <div
              key={item.id}
              className="flex items-center justify-between gap-4 rounded-[var(--radius-card)] border border-border p-4"
            >
              <Link href={`/produk/${item.product_slug}`} className="font-medium hover:underline">
                {item.product_name}
              </Link>
              <div className="flex items-center gap-3">
                <label className="flex items-center gap-2 text-sm">
                  <WishlistAlertCheckbox item={item} onMutate={setAlert.mutate} />
                  Alert
                </label>
                <Button variant="ghost" size="sm" onClick={() => removeWishlist.mutate(item.id)}>
                  Remove
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default function WishlistPage() {
  return (
    <AuthGuard>
      <WishlistContent />
    </AuthGuard>
  );
}
