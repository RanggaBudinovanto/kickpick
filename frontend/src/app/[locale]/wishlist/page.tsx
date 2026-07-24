"use client";

import { IconHeartOff } from "@tabler/icons-react";
import { AuthGuard } from "@/components/auth/auth-guard";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { useRemoveWishlist, useSetWishlistAlert, useWishlist } from "@/hooks/use-wishlist";

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
          <p className="text-sm font-medium">Gagal memuat data. Coba lagi.</p>
          <Button size="sm" onClick={() => refetch()}>
            Coba lagi
          </Button>
        </div>
      )}

      {!isLoading && !isError && data?.data.length === 0 && (
        <div className="flex flex-col items-center gap-3 py-16 text-center">
          <IconHeartOff size={32} />
          <p className="text-sm font-medium">Belum ada sepatu yang disimpan</p>
          <Link href="/cari">
            <Button size="sm">Mulai cari sepatu</Button>
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
                  <input
                    type="checkbox"
                    checked={item.alert_active}
                    onChange={(e) =>
                      setAlert.mutate({ id: item.id, alertActive: e.target.checked, alertType: "restock" })
                    }
                  />
                  Alert
                </label>
                <Button variant="ghost" size="sm" onClick={() => removeWishlist.mutate(item.id)}>
                  Hapus
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
