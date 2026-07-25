"use client";

import { useEffect, useState } from "react";
import { useTheme } from "next-themes";
import { IconMoon, IconSun } from "@tabler/icons-react";
import { useTranslations } from "next-intl";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const t = useTranslations("Theme");

  // eslint-disable-next-line react-hooks/set-state-in-effect
  useEffect(() => setMounted(true), []);

  const isDark = mounted && resolvedTheme === "dark";

  return (
    <button
      type="button"
      aria-label={isDark ? t("light") : t("dark")}
      onClick={() => setTheme(isDark ? "light" : "dark")}
      className="flex h-11 w-11 items-center justify-center rounded-[var(--radius-control)] text-foreground hover:bg-surface"
    >
      {mounted && (isDark ? <IconSun size={20} /> : <IconMoon size={20} />)}
    </button>
  );
}
