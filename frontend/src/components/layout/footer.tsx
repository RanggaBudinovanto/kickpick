import Image from "next/image";
import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export function Footer() {
  const t = useTranslations("Footer");

  return (
    <footer className="border-t border-border bg-card/30 backdrop-blur-sm">
      <div className="mx-auto max-w-[1400px] px-4 py-12 md:px-6">
        <div className="grid grid-cols-1 gap-10 md:grid-cols-4">
          <div className="flex flex-col gap-3">
            <Link href="/" aria-label="KickPick" className="w-fit">
              <Image
                src="/logo-wordmark.png"
                alt="KickPick"
                width={140}
                height={49}
                className="h-7 w-auto dark:invert"
              />
            </Link>
            <p className="text-xs leading-relaxed text-muted max-w-xs">
              The most complete sneaker price comparison platform from official stores & marketplaces.
            </p>
          </div>

          <div>
            <p className="mb-3 text-xs font-semibold uppercase tracking-wider text-muted">
              {t("categories")}
            </p>
            <ul className="flex flex-col gap-2.5 text-sm">
              <li>
                <Link href="/cari?kategori=running" className="text-foreground/80 hover:text-foreground transition-colors">
                  Running
                </Link>
              </li>
              <li>
                <Link href="/cari?kategori=lifestyle" className="text-foreground/80 hover:text-foreground transition-colors">
                  Lifestyle
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <p className="mb-3 text-xs font-semibold uppercase tracking-wider text-muted">
              {t("brands")}
            </p>
            <ul className="flex flex-col gap-2.5 text-sm">
              <li>
                <Link href="/brand" className="text-foreground/80 hover:text-foreground transition-colors">
                  {t("brands")}
                </Link>
              </li>
              <li>
                <Link href="/tentang" className="text-foreground/80 hover:text-foreground transition-colors">
                  {t("about")}
                </Link>
              </li>
              <li>
                <Link href="/privasi" className="text-foreground/80 hover:text-foreground transition-colors">
                  {t("privacy")}
                </Link>
              </li>
              <li>
                <Link href="/disclosure" className="text-foreground/80 hover:text-foreground transition-colors">
                  {t("disclosure")}
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <p className="mb-3 text-xs font-semibold uppercase tracking-wider text-muted">
              {t("subscribeTitle")}
            </p>
            <form className="flex gap-2">
              <Input placeholder={t("subscribePlaceholder")} type="email" className="text-xs" />
              <Button type="submit" size="sm">
                {t("subscribeCta")}
              </Button>
            </form>
          </div>
        </div>

        <div className="mt-12 flex flex-col items-center justify-between gap-4 border-t border-border/50 pt-6 md:flex-row">
          <p className="text-xs text-muted">
            © {new Date().getFullYear()} KickPick. {t("rights")}
          </p>
        </div>
      </div>
    </footer>
  );
}
