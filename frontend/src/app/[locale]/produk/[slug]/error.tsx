"use client";

import { IconAlertTriangle } from "@tabler/icons-react";
import { Button } from "@/components/ui/button";

export default function ProductDetailError({ reset }: { reset: () => void }) {
  return (
    <div className="mx-auto flex max-w-[1400px] flex-col items-center gap-3 px-4 py-24 text-center">
      <IconAlertTriangle size={32} />
      <p className="text-sm font-medium">Gagal memuat data. Coba lagi.</p>
      <Button size="sm" onClick={reset}>
        Coba lagi
      </Button>
    </div>
  );
}
