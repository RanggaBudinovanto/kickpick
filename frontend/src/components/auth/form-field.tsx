import { useId } from "react";
import { IconAlertTriangle } from "@tabler/icons-react";
import { Input } from "@/components/ui/input";

export function FormField({
  label,
  error,
  id,
  ...props
}: React.InputHTMLAttributes<HTMLInputElement> & { label: string; error?: string }) {
  const generatedId = useId();
  const inputId = id ?? generatedId;

  return (
    <div className="mb-4">
      <label htmlFor={inputId} className="mb-1.5 block text-sm font-medium">
        {label}
      </label>
      <Input id={inputId} aria-invalid={!!error} {...props} />
      {error && (
        <p className="mt-1.5 flex items-center gap-1.5 text-sm font-medium text-foreground">
          <IconAlertTriangle size={14} />
          {error}
        </p>
      )}
    </div>
  );
}
