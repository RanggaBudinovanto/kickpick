"use client";

import { Button } from "@/components/ui/button";
import { useRouter } from "@/i18n/navigation";

export function ResetFilterLink() {
  const router = useRouter();
  return (
    <Button size="sm" onClick={() => router.push("/cari")}>
      Reset filter
    </Button>
  );
}
