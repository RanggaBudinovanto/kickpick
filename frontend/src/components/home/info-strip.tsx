interface InfoStripProps {
  items: string[];
}

// Monochrome replacement for the colored promo strip pattern seen on sites
// like JD Sports (DESIGN.md forbids color entirely, so this is bold text on
// a flat black/white bar instead of a yellow banner).
export function InfoStrip({ items }: InfoStripProps) {
  if (items.length === 0) return null;

  return (
    <div className="border-y border-border bg-foreground text-background">
      <div className="mx-auto grid max-w-[1400px] grid-cols-1 divide-y divide-background/20 px-4 md:grid-cols-3 md:divide-x md:divide-y-0 md:px-6">
        {items.map((item) => (
          <p
            key={item}
            className="py-3 text-center text-xs font-semibold uppercase tracking-[0.01em] md:text-sm"
          >
            {item}
          </p>
        ))}
      </div>
    </div>
  );
}
