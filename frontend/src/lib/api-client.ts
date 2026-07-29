import { useAuthStore } from "@/stores/auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

async function refreshAccessToken(): Promise<string | null> {
  try {
    const res = await fetch(`${API_URL}/api/auth/refresh`, {
      method: "POST",
      credentials: "include",
      headers: { "X-Requested-With": "kickpick" },
    });
    if (!res.ok) return null;
    const data = await res.json();
    useAuthStore.getState().setSession(data.access_token);
    return data.access_token as string;
  } catch {
    // Network error (backend down, CORS preflight failure, etc.) —
    // treat as "no session" and continue silently.
    return null;
  }
}

interface ApiFetchOptions extends RequestInit {
  auth?: boolean;
}

export async function apiFetch<T>(path: string, options: ApiFetchOptions = {}): Promise<T> {
  const { auth = false, headers, ...rest } = options;

  const doFetch = async (token: string | null) => {
    const finalHeaders: Record<string, string> = {
      "Content-Type": "application/json",
      "X-Requested-With": "kickpick",
      ...(headers as Record<string, string>),
    };
    if (auth && token) {
      finalHeaders.Authorization = `Bearer ${token}`;
    }
    return fetch(`${API_URL}${path}`, {
      ...rest,
      headers: finalHeaders,
      credentials: "include",
    });
  };

  let token = useAuthStore.getState().accessToken;
  let res = await doFetch(token);

  if (auth && res.status === 401) {
    token = await refreshAccessToken();
    if (token) {
      res = await doFetch(token);
    }
  }

  if (!res.ok) {
    let message = "Terjadi kesalahan";
    try {
      const body = await res.json();
      message = body.error ?? message;
    } catch {
      // ignore body parse error
    }
    throw new ApiError(res.status, message);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json() as Promise<T>;
}

export { refreshAccessToken, API_URL };
