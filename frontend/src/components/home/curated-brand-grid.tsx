import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { IconArrowUpRight } from "@tabler/icons-react";

export interface CuratedTile {
  brandName: string;
  brandSlug: string;
  /** Background photo (product shot). Used when no logoUrl is supplied. */
  imageUrl: string;
  /** Official brand logo (SVG/PNG). When provided the tile switches to
   *  logo-on-dark-background layout instead of full-bleed photo. */
  logoUrl?: string;
  /** Accent hex colour (e.g. "#f5a623") used for the subtle glow behind the logo. */
  accentColor?: string;
  /** Spans 2 grid columns on desktop for the bento asymmetry DESIGN.md § 5 requires. */
  wide?: boolean;
}

interface CuratedBrandGridProps {
  title: string;
  tiles: CuratedTile[];
}

export function CuratedBrandGrid({ title, tiles }: CuratedBrandGridProps) {
  if (tiles.length === 0) return null;

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <div className="mb-6 flex items-end justify-between">
        <h2 className="font-display text-2xl font-semibold tracking-[-0.005em] md:text-4xl">
          {title}
        </h2>
        <Link href="/brand" className="text-sm font-medium underline hover:text-foreground">
          Lihat semua brand
        </Link>
      </div>

      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        {tiles.map((tile) => {
          const hasLogo = Boolean(tile.logoUrl);

          return (
            <Link
              key={tile.brandSlug}
              href={`/brand/${tile.brandSlug}`}
              className={`group relative overflow-hidden rounded-2xl border border-border/50 shadow-sm transition-all duration-500 hover:shadow-xl ${
                tile.wide ? "col-span-2 aspect-[2/1]" : "col-span-1 aspect-square"
              } ${hasLogo ? "bg-neutral-950" : "bg-surface"}`}
            >
              {hasLogo ? (
                /* ── Logo-on-dark layout ── */
                <>
                  {/* Subtle accent glow */}
                  {tile.accentColor && (
                    <div
                      className="pointer-events-none absolute inset-0 opacity-20 transition-opacity duration-500 group-hover:opacity-35"
                      style={{
                        background: `radial-gradient(ellipse at center, ${tile.accentColor} 0%, transparent 70%)`,
                      }}
                    />
                  )}

                  {/* Official brand logo — centred, white-filtered */}
                  <div className="absolute inset-0 flex items-center justify-center p-8 md:p-10">
                    <Image
                      src={tile.logoUrl!}
                      alt={tile.brandName}
                      fill
                      sizes={tile.wide ? "(max-width: 768px) 100vw, 50vw" : "(max-width: 768px) 50vw, 25vw"}
                      className="object-contain p-8 md:p-12 brightness-0 invert opacity-80 transition-all duration-500 group-hover:opacity-100 group-hover:scale-105"
                    />
                  </div>

                  {/* Bottom label */}
                  <div className="absolute inset-x-4 bottom-4 flex items-center justify-between text-white">
                    <div>
                      <span className="block text-[10px] font-mono font-medium uppercase tracking-widest text-white/50">
                        Featured Brand
                      </span>
                      <p className="font-display text-lg font-bold uppercase tracking-tight text-white md:text-xl">
                        {tile.brandName}
                      </p>
                    </div>
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-white/15 text-white backdrop-blur-md transition-all duration-300 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 group-hover:bg-white group-hover:text-black">
                      <IconArrowUpRight className="h-4 w-4" />
                    </div>
                  </div>
                </>
              ) : (
                /* ── Full-bleed product photo layout (fallback) ── */
                <>
                  {tile.imageUrl && (
                    <Image
                      src={tile.imageUrl}
                      alt={tile.brandName}
                      fill
                      sizes={tile.wide ? "(max-width: 768px) 100vw, 50vw" : "(max-width: 768px) 50vw, 25vw"}
                      className="object-cover opacity-85 transition-all duration-700 group-hover:scale-105 group-hover:opacity-100"
                    />
                  )}

                  {/* Gradient overlay */}
                  <div className="absolute inset-0 bg-gradient-to-t from-black/85 via-black/25 to-black/5 transition-opacity duration-300 group-hover:opacity-90" />

                  {/* Brand title & icon */}
                  <div className="absolute inset-x-4 bottom-4 flex items-center justify-between text-white">
                    <div>
                      <span className="block text-[10px] font-mono font-medium uppercase tracking-widest text-white/70">
                        Featured Brand
                      </span>
                      <p className="font-display text-xl font-bold uppercase tracking-tight text-white md:text-2xl">
                        {tile.brandName}
                      </p>
                    </div>
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-white/20 text-white backdrop-blur-md transition-transform duration-300 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 group-hover:bg-white group-hover:text-black">
                      <IconArrowUpRight className="h-4 w-4" />
                    </div>
                  </div>
                </>
              )}
            </Link>
          );
        })}
      </div>
    </section>
  );
}
