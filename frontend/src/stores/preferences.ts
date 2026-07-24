import { create } from "zustand";
import { persist } from "zustand/middleware";

export type Currency = "IDR" | "USD";

interface PreferencesState {
  currency: Currency;
  setCurrency: (currency: Currency) => void;
}

export const usePreferencesStore = create<PreferencesState>()(
  persist(
    (set) => ({
      currency: "IDR",
      setCurrency: (currency) => set({ currency }),
    }),
    { name: "kickpick-preferences" },
  ),
);
