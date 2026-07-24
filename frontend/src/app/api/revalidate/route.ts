import { revalidateTag } from "next/cache";
import { NextResponse } from "next/server";

// Called after a mutation that should invalidate ISR cache for a specific
// product — either by the frontend itself (after a review submit) or by the
// Go scraper worker (after a price update).
//
// No secret gate here on purpose: this only marks a cache tag stale (forces
// the next request to refetch from the Go API), it doesn't expose or change
// any data, so the worst case of an open endpoint is a bit of extra origin
// load. What it must NOT allow is revalidating arbitrary Next.js cache tags,
// so the tag is restricted to the "product:<slug>" shape this app actually
// uses (a NEXT_PUBLIC_ secret would've been shipped to the browser anyway —
// no real protection — so this is the honest boundary instead).
const TAG_PATTERN = /^product:[a-z0-9-]+$/;

export async function POST(request: Request) {
  const { tag } = await request.json();

  if (typeof tag !== "string" || !TAG_PATTERN.test(tag)) {
    return NextResponse.json({ error: "invalid tag" }, { status: 400 });
  }

  // { expire: 0 } (not "max") because the submitting user's own
  // router.refresh() runs right after this call and needs the fresh data
  // immediately — "max" would give them stale-while-revalidate semantics,
  // serving the OLD cached page on that very next request.
  revalidateTag(tag, { expire: 0 });
  return NextResponse.json({ revalidated: true, tag });
}
