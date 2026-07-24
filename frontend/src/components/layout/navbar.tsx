"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { IconBell, IconHeart, IconMenu2, IconSearch, IconUser, IconX } from "@tabler/icons-react";
import { Link } from "@/i18n/navigation";
import { Input } from "@/components/ui/input";
import { ThemeToggle } from "@/components/layout/theme-toggle";
import { PreferenceSwitcher } from "@/components/layout/preference-switcher";
import { useAuthStore } from "@/stores/auth";
import { useUnreadCount } from "@/hooks/use-notifications";

export function Navbar() {
  const t = useTranslations("Nav");
  const [mobileOpen, setMobileOpen] = useState(false);
  const accessToken = useAuthStore((s) => s.accessToken);
  const { data: unread } = useUnreadCount();
  const unreadCount = unread?.count ?? 0;

  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur">
      <div className="mx-auto flex h-16 max-w-[1400px] items-center gap-4 px-4 md:px-6">
        <Link href="/" className="font-display text-xl font-bold tracking-[-0.01em]">
          KickPick
        </Link>

        <div className="hidden flex-1 max-w-md md:block">
          <div className="relative">
            <IconSearch
              size={18}
              className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted"
            />
            <Input placeholder={t("searchPlaceholder")} className="pl-10" />
          </div>
        </div>

        <nav className="hidden items-center gap-6 text-sm font-medium md:flex">
          <Link href="/cari" className="hover:text-muted">
            {t("search")}
          </Link>
          <Link href="/brand" className="hover:text-muted">
            {t("brand")}
          </Link>
          <Link href="/cari?filter=rare" className="hover:text-muted">
            {t("drops")}
          </Link>
        </nav>

        <div className="ml-auto flex items-center gap-1">
          <Link
            href="/notifikasi"
            aria-label={t("notifications")}
            className="relative hidden h-11 w-11 items-center justify-center rounded-[var(--radius-control)] hover:bg-surface md:flex"
          >
            <IconBell size={20} />
            {unreadCount > 0 && (
              <span className="absolute right-1.5 top-1.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-off-black px-1 text-[10px] font-bold text-pure-white dark:bg-pure-white dark:text-off-black">
                {unreadCount > 9 ? "9+" : unreadCount}
              </span>
            )}
          </Link>
          <Link
            href="/wishlist"
            aria-label={t("wishlist")}
            className="hidden h-11 w-11 items-center justify-center rounded-[var(--radius-control)] hover:bg-surface md:flex"
          >
            <IconHeart size={20} />
          </Link>
          <ThemeToggle />
          <div className="hidden md:block">
            <PreferenceSwitcher />
          </div>
          {accessToken ? (
            <Link
              href="/profil"
              aria-label="Profil"
              className="hidden h-11 w-11 items-center justify-center rounded-[var(--radius-control)] hover:bg-surface md:flex"
            >
              <IconUser size={20} />
            </Link>
          ) : (
            <Link
              href="/login"
              className="hidden h-11 items-center rounded-[var(--radius-control)] bg-off-black px-4 text-sm font-medium text-pure-white hover:opacity-90 dark:bg-pure-white dark:text-off-black md:flex"
            >
              {t("login")}
            </Link>
          )}

          <button
            type="button"
            aria-label="Menu"
            className="flex h-11 w-11 items-center justify-center rounded-[var(--radius-control)] hover:bg-surface md:hidden"
            onClick={() => setMobileOpen((v) => !v)}
          >
            {mobileOpen ? <IconX size={22} /> : <IconMenu2 size={22} />}
          </button>
        </div>
      </div>

      {mobileOpen && (
        <div className="border-t border-border px-4 py-4 md:hidden">
          <div className="relative mb-4">
            <IconSearch
              size={18}
              className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted"
            />
            <Input placeholder={t("searchPlaceholder")} className="pl-10" />
          </div>
          <nav className="flex flex-col gap-1 text-sm font-medium">
            <Link href="/cari" className="rounded-[var(--radius-control)] px-2 py-2 hover:bg-surface">
              {t("search")}
            </Link>
            <Link href="/brand" className="rounded-[var(--radius-control)] px-2 py-2 hover:bg-surface">
              {t("brand")}
            </Link>
            <Link href="/cari?filter=rare" className="rounded-[var(--radius-control)] px-2 py-2 hover:bg-surface">
              {t("drops")}
            </Link>
            <Link href="/wishlist" className="rounded-[var(--radius-control)] px-2 py-2 hover:bg-surface">
              {t("wishlist")}
            </Link>
            <Link href="/notifikasi" className="rounded-[var(--radius-control)] px-2 py-2 hover:bg-surface">
              {t("notifications")}
            </Link>
          </nav>
          <div className="mt-4 flex items-center justify-between">
            <PreferenceSwitcher />
            <Link
              href={accessToken ? "/profil" : "/login"}
              className="flex h-11 items-center rounded-[var(--radius-control)] bg-off-black px-4 text-sm font-medium text-pure-white dark:bg-pure-white dark:text-off-black"
            >
              {accessToken ? "Profil" : t("login")}
            </Link>
          </div>
        </div>
      )}
    </header>
  );
}
