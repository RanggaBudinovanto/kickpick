import Image from "next/image";
import { Link } from "@/i18n/navigation";

const CATEGORIES = [
  { id: "running", label: "Running", sub: "Speed & endurance", href: "/cari?kategori=running", imageUrl: "/categories/running.png" },
  { id: "lifestyle", label: "Lifestyle", sub: "Street & everyday", href: "/cari?kategori=lifestyle", imageUrl: "/hero/slide-1.png" },
  { id: "training", label: "Training", sub: "Gym & performance", href: "/cari?kategori=training", imageUrl: "/hero/slide-2.png" },
  { id: "rare", label: "Rare Finds", sub: "Limited & exclusive", href: "/cari?filter=rare", imageUrl: "/hero/slide-3.png" },
  { id: "formal", label: "Formal", sub: "Clean & refined", href: "/cari?kategori=formal", imageUrl: "/hero/trending-banner.png" },
];

export function ShopByCategory() {
  return (
    <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <div className="mb-6 flex items-end justify-between border-b border-border/60 pb-4">
        <div>
          <h2 className="font-display text-2xl font-bold tracking-tight text-foreground md:text-4xl">Shop by Category</h2>
          <p className="mt-1 text-xs text-muted md:text-sm">Find exactly what you&apos;re looking for across every sneaker style.</p>
        </div>
      </div>
      <div className="grid grid-cols-2 gap-3 md:grid-cols-4 md:gap-4">
        <Link href={CATEGORIES[0].href} className="group relative col-span-2 overflow-hidden rounded-[var(--radius-card)] bg-neutral-900 md:row-span-2" style={{ minHeight: "260px" }}>
          <Image src={CATEGORIES[0].imageUrl} alt={CATEGORIES[0].label} fill sizes="(max-width: 768px) 100vw, 50vw" className="object-cover transition-transform duration-700 group-hover:scale-105" />
          <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent" />
          <div className="absolute bottom-0 left-0 p-5">
            <p className="text-[10px] font-mono font-bold uppercase tracking-widest text-white/60">{CATEGORIES[0].sub}</p>
            <h3 className="font-display text-2xl font-bold text-white md:text-3xl">{CATEGORIES[0].label}</h3>
            <span className="mt-2 inline-block text-xs font-semibold text-white underline underline-offset-4 opacity-0 transition-opacity duration-300 group-hover:opacity-100">Shop now →</span>
          </div>
        </Link>
        {CATEGORIES.slice(1).map((cat) => (
          <Link key={cat.id} href={cat.href} className="group relative overflow-hidden rounded-[var(--radius-card)] bg-neutral-900" style={{ minHeight: "180px" }}>
            <Image src={cat.imageUrl} alt={cat.label} fill sizes="(max-width: 768px) 50vw, 25vw" className="object-cover transition-transform duration-700 group-hover:scale-105" />
            <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent" />
            <div className="absolute bottom-0 left-0 p-4">
              <p className="text-[9px] font-mono font-bold uppercase tracking-widest text-white/60">{cat.sub}</p>
              <h3 className="font-display text-base font-bold text-white md:text-lg">{cat.label}</h3>
              <span className="mt-1 inline-block text-xs font-semibold text-white underline underline-offset-4 opacity-0 transition-opacity duration-300 group-hover:opacity-100">Shop now →</span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
