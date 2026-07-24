"use client";

import { useState } from "react";
import { IconChevronDown } from "@tabler/icons-react";

export function AccordionItem({ question, answer }: { question: string; answer: string }) {
  const [open, setOpen] = useState(false);

  return (
    <div className="border-b border-border py-4">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex w-full items-center justify-between text-left"
      >
        <span className="font-medium">{question}</span>
        <IconChevronDown
          size={18}
          className={`shrink-0 transition-transform duration-200 ${open ? "rotate-180" : ""}`}
        />
      </button>
      {open && <p className="mt-3 text-sm text-muted">{answer}</p>}
    </div>
  );
}
