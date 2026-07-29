"use client";

import { IconBellOff } from "@tabler/icons-react";
import { AuthGuard } from "@/components/auth/auth-guard";
import { Link } from "@/i18n/navigation";
import { useMarkNotificationRead, useNotifications } from "@/hooks/use-notifications";

function NotificationsContent() {
  const { data, isLoading } = useNotifications();
  const markRead = useMarkNotificationRead();

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <h1 className="mb-6 font-display text-3xl font-bold tracking-[-0.01em]">Notifications</h1>

      {isLoading && (
        <div className="flex flex-col gap-3">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="h-16 animate-pulse rounded-[var(--radius-card)] bg-surface" />
          ))}
        </div>
      )}

      {!isLoading && data?.data.length === 0 && (
        <div className="flex flex-col items-center gap-3 py-16 text-center">
          <IconBellOff size={32} />
          <p className="text-sm font-medium">No notifications yet</p>
        </div>
      )}

      {!isLoading && data && data.data.length > 0 && (
        <ul className="flex flex-col gap-2">
          {data.data.map((n) => (
            <li key={n.id}>
              <Link
                href={n.action_url || "/"}
                onClick={() => !n.is_read && markRead.mutate(n.id)}
                className={`flex items-start gap-3 rounded-[var(--radius-card)] border border-border p-4 ${
                  n.is_read ? "" : "bg-surface"
                }`}
              >
                {!n.is_read && <span className="mt-1.5 h-2 w-2 shrink-0 rounded-full bg-off-black dark:bg-pure-white" />}
                <div>
                  <p className="text-sm font-medium">{n.title}</p>
                  <p className="text-sm text-muted">{n.body}</p>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default function NotificationsPage() {
  return (
    <AuthGuard>
      <NotificationsContent />
    </AuthGuard>
  );
}
