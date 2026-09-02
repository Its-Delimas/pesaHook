const API_URL = process.env.PESAHOOK_API_URL;
const API_KEY = process.env.PESAHOOK_API_KEY;

export async function pesahookFetch(path: string, options: RequestInit = {}) {
  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers: {
      Authorization: `Bearer ${API_KEY}`,
      "Content-Type": "application/json",
      ...options.headers,
    },
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`PesaHook API error: ${res.status} ${await res.text()}`);
  }

  return res.json();
}