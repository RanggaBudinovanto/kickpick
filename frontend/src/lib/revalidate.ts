// Client-side helper to invalidate a specific product's ISR cache after a
// mutation (e.g. review submit) so other visitors don't wait up to an hour
// to see it. Best-effort: a failure here shouldn't block the user's own flow
// (router.refresh() already gives them an up-to-date view regardless).
export async function revalidateProduct(slug: string) {
  try {
    await fetch("/api/revalidate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ tag: `product:${slug}` }),
    });
  } catch {
    // best-effort, ignore
  }
}
