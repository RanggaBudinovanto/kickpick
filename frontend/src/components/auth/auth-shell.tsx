import Image from "next/image";
import { Card } from "@/components/ui/card";

export function AuthShell({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
}) {
  return (
    <section className="mx-auto flex min-h-[75vh] max-w-[1200px] items-center justify-center px-4 py-8 md:py-12">
      <Card className="grid w-full grid-cols-1 overflow-hidden border-border/80 shadow-2xl md:grid-cols-2">
        {/* Left Column: Form Panel */}
        <div className="flex flex-col justify-center p-8 sm:p-12 md:p-14 bg-background">
          <div className="mb-6">
            <h1 className="font-display text-3xl font-bold tracking-tight text-foreground md:text-4xl">
              {title}
            </h1>
            {subtitle && (
              <p className="mt-2 text-sm text-muted">
                {subtitle}
              </p>
            )}
          </div>
          {children}
        </div>

        {/* Right Column: Editorial AI Sneaker Artwork Panel */}
        <div className="relative hidden md:flex flex-col justify-end overflow-hidden bg-neutral-950 p-12 text-white">
          <Image
            src="/hero/auth-banner.png"
            alt="KickPick Auth Banner"
            fill
            priority
            sizes="50vw"
            className="object-cover object-center opacity-85 transition-transform duration-1000 hover:scale-105"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-neutral-950 via-neutral-950/40 to-transparent" />

          {/* Bottom Editorial Quote */}
          <div className="relative z-10 space-y-2">
            <h2 className="font-display text-2xl font-bold tracking-tight text-white md:text-3xl">
              Compare Prices Across Top Stores
            </h2>
            <p className="text-xs text-neutral-300 leading-relaxed">
              Track restocks, verify authentic sneakers, and unlock real-time price drops from 40+ official brand stores.
            </p>
          </div>
        </div>
      </Card>
    </section>
  );
}
