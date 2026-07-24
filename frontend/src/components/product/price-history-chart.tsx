"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Button } from "@/components/ui/button";
import { formatPrice } from "@/lib/currency";
import { usePriceHistory } from "@/hooks/use-products";

export function PriceHistoryChart({ slug }: { slug: string }) {
  const t = useTranslations("Product");
  const [days, setDays] = useState<30 | 90>(30);
  const { data, isLoading } = usePriceHistory(slug, days);

  const points = data?.data ?? [];

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="font-display text-xl font-semibold">{t("priceHistory")}</h2>
        <div className="flex gap-1">
          <Button variant={days === 30 ? "primary" : "secondary"} size="sm" onClick={() => setDays(30)}>
            {t("days30")}
          </Button>
          <Button variant={days === 90 ? "primary" : "secondary"} size="sm" onClick={() => setDays(90)}>
            {t("days90")}
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="h-64 animate-pulse rounded-[var(--radius-card)] bg-surface" />
      ) : points.length < 5 ? (
        <p className="text-sm text-muted">{t("notEnoughHistory")}</p>
      ) : (
        <div className="h-64 rounded-[var(--radius-card)] border border-border p-4">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={points}>
              <XAxis
                dataKey="date"
                tick={{ fontSize: 11 }}
                tickFormatter={(v) => v.slice(5)}
                stroke="currentColor"
                opacity={0.4}
              />
              <YAxis
                tick={{ fontSize: 11 }}
                tickFormatter={(v) => `${Math.round(v / 1000)}k`}
                stroke="currentColor"
                opacity={0.4}
                width={48}
              />
              <Tooltip
                formatter={(value) => formatPrice(Number(value))}
                contentStyle={{ fontSize: 12, fontFamily: "var(--font-mono)" }}
              />
              <Line
                type="monotone"
                dataKey="price"
                stroke="currentColor"
                strokeWidth={2}
                dot={false}
                isAnimationActive
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  );
}
