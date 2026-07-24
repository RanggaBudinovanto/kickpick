"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslations } from "next-intl";
import { AuthShell } from "@/components/auth/auth-shell";
import { FormField } from "@/components/auth/form-field";
import { Button } from "@/components/ui/button";
import { useForgotPassword } from "@/hooks/use-auth";

const schema = z.object({
  email: z.string().email(),
});

type FormValues = z.infer<typeof schema>;

export default function ForgotPasswordPage() {
  const t = useTranslations("Auth");
  const forgotPassword = useForgotPassword();
  const [submitted, setSubmitted] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit = (values: FormValues) => {
    forgotPassword.mutate(values.email, {
      onSuccess: () => setSubmitted(true),
    });
  };

  if (submitted) {
    return (
      <AuthShell title={t("forgotTitle")}>
        <p className="text-sm">{t("forgotSuccess")}</p>
      </AuthShell>
    );
  }

  return (
    <AuthShell title={t("forgotTitle")}>
      <form onSubmit={handleSubmit(onSubmit)} noValidate>
        <FormField
          label={t("email")}
          type="email"
          {...register("email")}
          error={errors.email ? t("emailInvalid") : undefined}
        />

        <Button type="submit" className="w-full" disabled={forgotPassword.isPending}>
          {t("forgotCta")}
        </Button>
      </form>
    </AuthShell>
  );
}
