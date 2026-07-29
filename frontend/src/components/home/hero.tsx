import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { IconSearch, IconArrowRight } from "@tabler/icons-react";

interface HeroProps {
  headline: string;
  subtext: string;
  searchCtaText: string;
  browseBrandsText: string;
}

export function Hero({ headline, subtext, searchCtaText, browseBrandsText }: HeroProps) {
  return (
    <section className="relative mx-auto grid min-h-[calc(100vh-4rem)] max-w-[1400px] grid-cols-1 items-center gap-8 px-4 py-12 md:grid-cols-[60%_40%] md:px-6 md:py-16">
      <div className="flex flex-col justify-center space-y-6 md:space-y-8">
        <div className="space-y-4">
          <div className="inline-flex items-center gap-2 rounded-full border border-border bg-surface px-3 py-1 text-xs font-mono tracking-wider text-muted uppercase">
            <span className="h-1.5 w-1.5 rounded-full bg-foreground animate-pulse" />
            Live Price Comparison Engine
          </div>
          <h1 className="font-display text-4xl font-bold uppercase leading-[0.95] tracking-tight text-foreground sm:text-5xl md:text-6xl lg:text-7xl">
            {headline}
          </h1>
          <p className="max-w-xl text-base text-muted sm:text-lg">
            {subtext}
          </p>
        </div>

        {/* Quick Search Input / Link */}
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <Link href="/cari" className="flex-1">
            <div className="flex items-center gap-3 rounded-lg border border-border bg-surface px-4 py-3.5 text-muted transition-colors hover:border-foreground/40 hover:text-foreground">
              <IconSearch className="h-5 w-5 shrink-0" />
              <span className="text-sm font-medium">{searchCtaText}...</span>
            </div>
          </Link>
          <Link href="/brand">
            <Button variant="secondary" size="default" className="w-full sm:w-auto h-[50px] px-6 font-medium">
              {browseBrandsText}
            </Button>
          </Link>
        </div>
      </div>

      {/* 40% Visual Area */}
      <div className="relative aspect-[4/5] w-full overflow-hidden rounded-2xl border border-border bg-surface shadow-2xl">
        <Image
          src="/hero-sneakers.png"
          alt="KickPick Monochrome Sneakers Hero Showcase"
          fill
          priority
          sizes="(max-width: 768px) 100vw, 40vw"
          className="object-cover grayscale transition-transform duration-700 hover:scale-105"
        />
        {/* Grayscale overlay vignette */}
        <div className="absolute inset-0 bg-gradient-to-t from-background/80 via-transparent to-transparent" />
        
        {/* Floating Minimal Price Tag / Data Badge */}
        <div className="absolute bottom-4 left-4 right-4 flex items-center justify-between rounded-xl border border-border/80 bg-background/90 p-4 backdrop-blur-md">
          <div>
            <span className="block text-[10px] font-mono uppercase tracking-widest text-muted">Best Market Price</span>
            <span className="font-display text-sm font-bold tracking-tight text-foreground">IDR 1,299,000</span>
          </div>
          <Link href="/cari" className="inline-flex items-center gap-1 text-xs font-semibold text-foreground hover:underline">
            Compare <IconArrowRight className="h-3.5 w-3.5" />
          </Link>
        </div>
      </div>
    </section>
  );
}
