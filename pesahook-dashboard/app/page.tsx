// file: pesahook-dashboard/app/page.tsx
import Link from "next/link";
import { getEndpoints } from "@/lib/pesahook";

export default async function EndpointsPage() {
  const endpoints = await getEndpoints();

  return (
    <main className="px-8 py-10 max-w-4xl">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Endpoints</h1>
          <p className="text-sm text-muted mt-1">
            Destinations receiving your Daraja callbacks.
          </p>
        </div>
        <button className="text-sm font-medium bg-ink text-white px-4 py-2 rounded-lg hover:bg-ink/90 transition-colors">
          New endpoint
        </button>
      </div>

      {endpoints.length === 0 ? (
        <div className="border border-dashed border-border rounded-xl px-6 py-16 text-center">
          <p className="text-sm text-ink font-medium mb-1">No endpoints yet</p>
          <p className="text-sm text-muted mb-4">
            Register your first shortcode to start receiving normalized callbacks.
          </p>
          <button className="text-sm font-medium bg-ink text-white px-4 py-2 rounded-lg hover:bg-ink/90 transition-colors">
            New endpoint
          </button>
        </div>
      ) : (
        <div className="border border-border rounded-xl overflow-hidden">
          {endpoints.map((ep, i) => (
            <Link
              href={`/endpoints/${ep.id}`}
              key={ep.id}
              className={`flex items-center justify-between px-5 py-4 hover:bg-surface transition-colors ${
                i !== 0 ? "border-t border-border" : ""
              }`}
            >
              <div className="flex flex-col gap-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{ep.shortcode}</span>
                  <span className="text-xs text-accent bg-accent-soft px-1.5 py-0.5 rounded">
                    {ep.provider}
                  </span>
                </div>
                <span className="font-mono text-xs text-muted truncate">
                  {ep.destination_url}
                </span>
              </div>

              <div className="flex items-center gap-4 shrink-0">
                <div className="hidden sm:flex gap-1.5">
                  {ep.event_types.map((type) => (
                    <span
                      key={type}
                      className="text-xs text-muted border border-border px-2 py-0.5 rounded-full"
                    >
                      {type}
                    </span>
                  ))}
                </div>
                <span className="text-xs text-muted">
                  {new Date(ep.created_at).toLocaleDateString()}
                </span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </main>
  );
}