"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import { AuthGuard } from "@/components/auth/auth-guard";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRouter } from "@/i18n/navigation";
import { useDeleteAccount, useLogout, useProfile, useUpdateProfile } from "@/hooks/use-auth";
import { ApiError } from "@/lib/api-client";

function ProfileContent() {
  const { data: profile, isLoading } = useProfile();
  const updateProfile = useUpdateProfile();
  const deleteAccount = useDeleteAccount();
  const logout = useLogout();
  const router = useRouter();

  const [name, setName] = useState("");
  const [confirmingDelete, setConfirmingDelete] = useState(false);
  const [password, setPassword] = useState("");

  useEffect(() => {
    if (profile) setName(profile.name);
  }, [profile]);

  function saveProfile() {
    if (!profile) return;
    updateProfile.mutate(
      {
        name,
        onboarding_focus: profile.onboarding_focus,
        preferred_language: profile.preferred_language,
        preferred_currency: profile.preferred_currency,
      },
      {
        onSuccess: () => toast.success("Profil berhasil disimpan"),
        onError: (err) => toast.error(err instanceof ApiError ? err.message : "Gagal menyimpan profil"),
      },
    );
  }

  function confirmDelete() {
    deleteAccount.mutate(password, {
      onSuccess: () => {
        toast.success("Akun berhasil dihapus");
        router.push("/");
      },
      onError: (err) => toast.error(err instanceof ApiError ? err.message : "Gagal menghapus akun"),
    });
  }

  if (isLoading || !profile) {
    return (
      <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
        <div className="h-64 max-w-md animate-pulse rounded-[var(--radius-card)] bg-surface" />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <h1 className="mb-6 font-display text-3xl font-bold tracking-[-0.01em]">Profil</h1>

      <div className="max-w-md">
        <div className="mb-4">
          <label htmlFor="profile-email" className="mb-1.5 block text-sm font-medium">
            Email
          </label>
          <Input id="profile-email" value={profile.email} disabled />
        </div>
        <div className="mb-6">
          <label htmlFor="profile-name" className="mb-1.5 block text-sm font-medium">
            Nama
          </label>
          <Input id="profile-name" value={name} onChange={(e) => setName(e.target.value)} />
        </div>

        <div className="flex gap-3">
          <Button onClick={saveProfile} disabled={updateProfile.isPending}>
            Simpan
          </Button>
          <Button variant="secondary" onClick={() => logout.mutate()}>
            Keluar
          </Button>
        </div>

        <div className="mt-10 border-t border-border pt-6">
          {!confirmingDelete ? (
            <Button variant="ghost" size="sm" onClick={() => setConfirmingDelete(true)}>
              Hapus akun
            </Button>
          ) : (
            <div>
              <p className="mb-3 text-sm font-medium">
                Masukkan password untuk konfirmasi penghapusan akun. Tindakan ini tidak bisa dibatalkan.
              </p>
              <label htmlFor="profile-delete-password" className="sr-only">
                Password
              </label>
              <Input
                id="profile-delete-password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Password"
                className="mb-3"
              />
              <div className="flex gap-3">
                <Button variant="secondary" size="sm" onClick={confirmDelete} disabled={deleteAccount.isPending}>
                  Konfirmasi hapus
                </Button>
                <Button variant="ghost" size="sm" onClick={() => setConfirmingDelete(false)}>
                  Batal
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default function ProfilePage() {
  return (
    <AuthGuard>
      <ProfileContent />
    </AuthGuard>
  );
}
