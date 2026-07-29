import Image from "next/image";
import { Link } from "@/i18n/navigation";

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

      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {tiles.map((tile) => {
          const hasLogo = Boolean(tile.logoUrl);

          return (
            <Link
              key={tile.brandSlug}
              href={`/brand/${tile.brandSlug}`}
              className="group relative flex aspect-[16/9] items-center justify-center overflow-hidden rounded-2xl border border-border/80 bg-card p-6 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:border-border-hover hover:shadow-md"
            >
              {/* Official brand logo — only logo */}
              <div className="relative h-12 w-full">
                {hasLogo ? (
                  <Image
                    src={tile.logoUrl!}
                    alt={tile.brandName}
                    fill
                    sizes="(max-width: 768px) 50vw, 20vw"
                    className="object-contain dark:invert opacity-80 transition-all duration-300 group-hover:opacity-100 group-hover:scale-105"
                  />
                ) : tile.imageUrl ? (
                  <Image
                    src={tile.imageUrl}
                    alt={tile.brandName}
                    fill
                    sizes="(max-width: 768px) 50vw, 20vw"
                    className="object-cover opacity-85 transition-transform duration-300 group-hover:scale-105"
                  />
                ) : (
                  <span className="font-display text-base font-bold tracking-tight text-foreground/80">
                    {tile.brandName}
                  </span>
                )}
              </div>
            </Link>
          );
        })}
      </div>
    </section>
  );
}
