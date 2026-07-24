import { Card } from "@/components/ui/card";

export function AuthShell({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section className="mx-auto flex min-h-[70vh] max-w-[1400px] items-center justify-center px-4 py-12">
      <Card className="w-full max-w-md p-8">
        <h1 className="mb-6 font-display text-2xl font-bold tracking-[-0.01em]">{title}</h1>
        {children}
      </Card>
    </section>
  );
}
