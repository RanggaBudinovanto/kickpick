import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api-client";

interface ExchangeRateResponse {
  base_currency: string;
  target_currency: string;
  rate: number;
  recorded_date: string;
}

export function useExchangeRate() {
  return useQuery({
    queryKey: ["exchange-rate", "IDR", "USD"],
    queryFn: () => apiFetch<ExchangeRateResponse>("/api/exchange-rate"),
    staleTime: 60 * 60 * 1000, // kurs diperbarui harian, tidak perlu refetch tiap menit
    retry: false,
  });
}
