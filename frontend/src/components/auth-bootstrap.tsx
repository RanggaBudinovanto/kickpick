"use client";

import { useEffect } from "react";
import { refreshAccessToken } from "@/lib/api-client";
import { useAuthStore } from "@/stores/auth";

// Access token lives only in memory, so on a fresh page load we try to silently
// restore the session from the httpOnly refresh_token cookie.
export function AuthBootstrap() {
  const finishBootstrap = useAuthStore((s) => s.finishBootstrap);

  useEffect(() => {
    refreshAccessToken().finally(finishBootstrap);
  }, [finishBootstrap]);

  return null;
}
