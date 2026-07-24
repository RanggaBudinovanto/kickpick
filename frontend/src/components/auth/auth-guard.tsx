"use client";

import { useEffect } from "react";
import { useRouter } from "@/i18n/navigation";
import { useAuthStore } from "@/stores/auth";

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const accessToken = useAuthStore((s) => s.accessToken);
  const isBootstrapping = useAuthStore((s) => s.isBootstrapping);
  const router = useRouter();

  useEffect(() => {
    if (!isBootstrapping && !accessToken) {
      router.push("/login");
    }
  }, [isBootstrapping, accessToken, router]);

  if (isBootstrapping) {
    return (
      <div className="mx-auto max-w-[1400px] px-4 py-16">
        <div className="h-8 w-48 animate-pulse rounded bg-surface" />
      </div>
    );
  }

  if (!accessToken) {
    return null;
  }

  return <>{children}</>;
}
