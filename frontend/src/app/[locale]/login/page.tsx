"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { AuthShell } from "@/components/auth/auth-shell";
import { FormField } from "@/components/auth/form-field";
import { Button } from "@/components/ui/button";
import { Link, useRouter } from "@/i18n/navigation";
import { useLogin } from "@/hooks/use-auth";
import { ApiError } from "@/lib/api-client";

const schema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
});

type FormValues = z.infer<typeof schema>;

export default function LoginPage() {
  const t = useTranslations("Auth");
  const router = useRouter();
  const login = useLogin();
  const [shake, setShake] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit = (values: FormValues) => {
    login.mutate(values, {
      onSuccess: () => {
        router.push("/");
      },
      onError: (err) => {
        setShake(true);
        setTimeout(() => setShake(false), 300);
        toast.error(err instanceof ApiError ? err.message : "Login failed");
      },
    });
  };

  return (
    <AuthShell title={t("loginTitle")}>
      <form onSubmit={handleSubmit(onSubmit)} noValidate className={shake ? "animate-shake" : ""}>
        <FormField
          label={t("email")}
          type="email"
          {...register("email")}
          error={errors.email ? t("emailInvalid") : undefined}
        />
        <FormField
          label={t("password")}
          type="password"
          {...register("password")}
          error={errors.password ? t("passwordMin") : undefined}
        />

        <Button type="submit" className="w-full" disabled={login.isPending}>
          {t("loginCta")}
        </Button>

        <div className="mt-4 flex items-center justify-between text-sm">
          <span className="text-muted">
            {t("noAccount")}{" "}
            <Link href="/registrasi" className="font-medium text-foreground underline">
              {t("registerLink")}
            </Link>
          </span>
          <Link href="/lupa-password" className="text-muted underline">
            {t("forgotLink")}
          </Link>
        </div>
      </form>
    </AuthShell>
  );
}
