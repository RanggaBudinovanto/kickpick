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
import { Link } from "@/i18n/navigation";
import { useRegister } from "@/hooks/use-auth";
import { ApiError } from "@/lib/api-client";

const schema = z
  .object({
    name: z.string().min(1),
    email: z.string().email(),
    password: z.string().min(8),
    confirm_password: z.string().min(8),
  })
  .refine((data) => data.password === data.confirm_password, {
    path: ["confirm_password"],
  });

type FormValues = z.infer<typeof schema>;

export default function RegisterPage() {
  const t = useTranslations("Auth");
  const register_ = useRegister();
  const [submitted, setSubmitted] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit = (values: FormValues) => {
    register_.mutate(values, {
      onSuccess: () => setSubmitted(true),
      onError: (err) => {
        toast.error(err instanceof ApiError ? err.message : "Registrasi gagal");
      },
    });
  };

  if (submitted) {
    return (
      <AuthShell title={t("registerTitle")}>
        <p className="text-sm">{t("registerSuccess")}</p>
      </AuthShell>
    );
  }

  return (
    <AuthShell title={t("registerTitle")}>
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <FormField
          label={t("name")}
          {...register("name")}
          error={errors.name ? t("nameRequired") : undefined}
        />
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
        <FormField
          label={t("confirmPassword")}
          type="password"
          {...register("confirm_password")}
          error={errors.confirm_password ? t("passwordMismatch") : undefined}
        />

        <Button type="submit" className="w-full" disabled={register_.isPending}>
          {t("registerCta")}
        </Button>

        <p className="mt-4 text-center text-sm text-muted">
          {t("haveAccount")}{" "}
          <Link href="/login" className="font-medium text-foreground underline">
            {t("loginLink")}
          </Link>
        </p>
      </form>
    </AuthShell>
  );
}
