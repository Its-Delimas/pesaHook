import Link from "next/link";
import { Webhook, Plus, ArrowUpRight } from "lucide-react";
import { getEndpoints } from "@/lib/pesahook";

export default async function EndpointsPage() {
  const endpoints = await getEndpoints();

  return (
    <main className="px-8 py-10 w-full">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Endpoints</h1>
          <p className="text-sm text-muted mt-1">
            Destinations receiving your Daraja callbacks.
          </p>
        </div>
        <button className="inline-flex items-center gap-1.5 text-sm font-medium bg-ink text-white pl-3 pr-4 py-2 rounded-lg hover:bg-ink/90 active:scale-[0.98] transition-all focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ink">
          <Plus size={15} strokeWidth={2.5} />
          New endpoint
        </button>
      </div>

      {endpoints.length === 0 ? (
        <div className="border border-dashed border-border rounded-xl px-6 py-16 flex flex-col items-center text-center">
          <div className="w-10 h-10 rounded-full bg-surface border border-border flex items-center justify-center mb-4">
            <Webhook size={18} className="text-muted" strokeWidth={1.75} />
          </div>
          <p className="text-sm text-ink font-medium mb-1">No endpoints yet</p>
          <p className="text-sm text-muted mb-4 max-w-xs">
            Register your first shortcode to start receiving normalized callbacks.
          </p>
          <button className="inline-flex items-center gap-1.5 text-sm font-medium bg-ink text-white pl-3 pr-4 py-2 rounded-lg hover:bg-ink/90 transition-colors">
            <Plus size={15} strokeWidth={2.5} />
            New endpoint
          </button>
        </div>
      ) : (
        <div className="border border-border rounded-xl overflow-hidden">
          {endpoints.map((ep, i) => (
            <Link
              href={`/endpoints/${ep.id}`}
              key={ep.id}
              className={`group flex flex-col gap-3 px-5 py-4 hover:bg-surface transition-colors ${
                i !== 0 ? "border-t border-border" : ""
              }`}
            >
              <div className="flex items-center justify-between gap-4">
                <div className="flex items-center gap-3 min-w-0">
                  <div className="w-8 h-8 rounded-lg bg-accent-soft flex items-center justify-center shrink-0">
                    <Webhook size={15}  strokeWidth={2} />
                  </div>
                  <div className="flex flex-col gap-0.5 min-w-0">
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
                </div>

                <div className="flex items-center gap-4 shrink-0">
                  <span className="text-xs text-muted hidden sm:inline">
                    {new Date(ep.created_at).toLocaleDateString()}
                  </span>
                  <ArrowUpRight
                    size={15}
                    className="text-muted opacity-0 group-hover:opacity-100 transition-opacity"
                  />
                </div>
              </div>

              <div className="flex gap-1.5 pl-11">
                {ep.event_types.map((type) => (
                  <span
                    key={type}
                    className="text-xs text-muted border border-border px-2 py-0.5 rounded-full"
                  >
                    {type}
                  </span>
                ))}
              </div>
            </Link>
          ))}
        </div>
      )}
    </main>
  );
}