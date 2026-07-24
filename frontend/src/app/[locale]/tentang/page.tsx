import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { AccordionItem } from "@/components/ui/accordion-item";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "About" });
  return { title: `${t("title")} | KickPick` };
}

export default async function AboutPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "About" });

  const faqs = [1, 2, 3, 4].map((n) => ({
    q: t(`faq${n}Q` as "faq1Q"),
    a: t(`faq${n}A` as "faq1A"),
  }));

  return (
    <div className="mx-auto max-w-2xl px-4 py-10 md:px-6">
      <h1 className="mb-4 font-display text-3xl font-bold tracking-[-0.01em]">{t("title")}</h1>
      <p className="mb-10 text-sm text-muted">{t("intro")}</p>

      <h2 className="mb-4 font-display text-xl font-semibold">{t("faqTitle")}</h2>
      <div>
        {faqs.map((faq, i) => (
          <AccordionItem key={i} question={faq.q} answer={faq.a} />
        ))}
      </div>
    </div>
  );
}
