import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export function Footer() {
  const t = useTranslations("Footer");

  return (
    <footer className="border-t border-border">
      <div className="mx-auto max-w-[1400px] px-4 py-12 md:px-6">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
          <div>
            <p className="font-display text-lg font-bold">KickPick</p>
          </div>

          <div>
            <p className="mb-3 text-xs font-medium uppercase tracking-[0.01em] text-muted">
              {t("categories")}
            </p>
            <ul className="flex flex-col gap-2 text-sm">
              <li>
                <Link href="/cari?kategori=running" className="hover:text-muted">
                  Running
                </Link>
              </li>
              <li>
                <Link href="/cari?kategori=lifestyle" className="hover:text-muted">
                  Lifestyle
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <p className="mb-3 text-xs font-medium uppercase tracking-[0.01em] text-muted">
              {t("brands")}
            </p>
            <ul className="flex flex-col gap-2 text-sm">
              <li>
                <Link href="/brand" className="hover:text-muted">
                  {t("brands")}
                </Link>
              </li>
              <li>
                <Link href="/tentang" className="hover:text-muted">
                  {t("about")}
                </Link>
              </li>
              <li>
                <Link href="/privasi" className="hover:text-muted">
                  {t("privacy")}
                </Link>
              </li>
              <li>
                <Link href="/disclosure" className="hover:text-muted">
                  {t("disclosure")}
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <p className="mb-3 text-xs font-medium uppercase tracking-[0.01em] text-muted">
              {t("subscribeTitle")}
            </p>
            <form className="flex gap-2">
              <Input placeholder={t("subscribePlaceholder")} type="email" />
              <Button type="submit" size="sm">
                {t("subscribeCta")}
              </Button>
            </form>
          </div>
        </div>

        <p className="mt-10 text-xs text-muted">
          © {new Date().getFullYear()} KickPick. {t("rights")}
        </p>
      </div>
    </footer>
  );
}
