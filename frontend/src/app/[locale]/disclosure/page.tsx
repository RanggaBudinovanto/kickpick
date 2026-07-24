import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Disclosure" });
  return { title: `${t("title")} | KickPick` };
}

export default async function DisclosurePage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Disclosure" });

  return (
    <div className="mx-auto max-w-2xl px-4 py-10 md:px-6">
      <h1 className="mb-6 font-display text-3xl font-bold tracking-[-0.01em]">{t("title")}</h1>
      <div className="flex flex-col gap-4 text-sm text-muted">
        <p>{t("body1")}</p>
        <p>{t("body2")}</p>
      </div>
    </div>
  );
}
