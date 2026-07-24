import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Privacy" });
  return { title: `${t("title")} | KickPick` };
}

export default async function PrivacyPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Privacy" });

  const sections = [1, 2, 3, 4].map((n) => ({
    title: t(`section${n}Title` as "section1Title"),
    body: t(`section${n}Body` as "section1Body"),
  }));

  return (
    <div className="mx-auto max-w-2xl px-4 py-10 md:px-6">
      <h1 className="mb-4 font-display text-3xl font-bold tracking-[-0.01em]">{t("title")}</h1>
      <p className="mb-10 text-sm text-muted">{t("intro")}</p>

      <div className="flex flex-col gap-6">
        {sections.map((s, i) => (
          <div key={i}>
            <h2 className="mb-2 font-display text-lg font-semibold">{s.title}</h2>
            <p className="text-sm text-muted">{s.body}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
