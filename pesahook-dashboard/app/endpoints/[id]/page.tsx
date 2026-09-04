import Link from "next/link";
import { ArrowLeft, CheckCircle2, XCircle, ArrowUpRight, Inbox, TrendingUp } from "lucide-react";
import { getEventsForEndpoints } from "../../../lib/pesahook";

export default async function EndpointEventsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const events = await getEventsForEndpoints(id);
  const successCount = events.filter((e) => e.status === "success").length;
  const totalAmount = events
    .filter((e) => e.status === "success")
    .reduce((sum, e) => sum + e.amount, 0);

  return (
    <main className="px-6 py-10 max-w-5xl mx-auto">
      <Link
        href="/"
        className="inline-flex items-center gap-1.5 text-sm text-muted hover:text-ink transition-colors mb-6"
      >
        <ArrowLeft size={14} />
        Endpoints
      </Link>

      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Events</h1>
          <p className="font-mono text-xs text-muted mt-1">{id}</p>
        </div>
        <Link
          href={`/endpoints/${id}/dead-letters`}
          className="text-sm font-medium text-ink bg-white border border-border px-4 py-2 rounded-lg shadow-sm hover:shadow-md transition-shadow"
        >
          Dead letters
        </Link>
      </div>

      {events.length > 0 && (
        <div className="grid grid-cols-3 gap-4 mb-8">
          <div className="rounded-2xl px-5 py-5 bg-gradient-to-br from-accent to-accent-dark text-white shadow-sm">
            <div className="flex items-center justify-between mb-3">
              <p className="text-xs text-white/80">Total received</p>
              <TrendingUp size={14} className="text-white/70" />
            </div>
            <p className="text-3xl font-semibold tracking-tight">
              KES {totalAmount.toLocaleString()}
            </p>
          </div>
          <div className="bg-white border border-border rounded-2xl px-5 py-5 shadow-sm hover:shadow-md transition-shadow">
            <p className="text-xs text-muted mb-3">Events captured</p>
            <p className="text-3xl font-semibold tracking-tight">{events.length}</p>
          </div>
          <div className="bg-white border border-border rounded-2xl px-5 py-5 shadow-sm hover:shadow-md transition-shadow">
            <p className="text-xs text-muted mb-3">Successful</p>
            <p className="text-3xl font-semibold tracking-tight">
              {successCount}/{events.length}
            </p>
          </div>
        </div>
      )}

      {events.length === 0 ? (
        <div className="bg-white border border-dashed border-border rounded-2xl px-6 py-16 flex flex-col items-center text-center shadow-sm">
          <div className="w-10 h-10 rounded-full bg-surface border border-border flex items-center justify-center mb-4">
            <Inbox size={18} className="text-muted" strokeWidth={1.75} />
          </div>
          <p className="text-sm text-ink font-medium mb-1">No events captured yet</p>
          <p className="text-sm text-muted max-w-xs">
            Once Daraja sends a callback to this endpoint, it will appear here.
          </p>
        </div>
      ) : (
        <div className="bg-white border border-border rounded-2xl overflow-hidden shadow-sm">
          {events.map((ev, i) => (
            <Link
              href={`/events/${ev.event_id}`}
              key={ev.event_id}
              className={`group flex flex-col gap-3 px-5 py-4 hover:bg-surface/50 transition-colors ${
                i !== 0 ? "border-t border-border" : ""
              }`}
            >
              <div className="flex items-center justify-between gap-4">
                <div className="flex items-center gap-3 min-w-0">
                  <div
                    className={`w-9 h-9 rounded-xl flex items-center justify-center shrink-0 ${
                      ev.status === "success" ? "bg-accent-soft" : "bg-danger-soft"
                    }`}
                  >
                    {ev.status === "success" ? (
                      <CheckCircle2 size={16} className="text-accent" strokeWidth={2} />
                    ) : (
                      <XCircle size={16} className="text-danger" strokeWidth={2} />
                    )}
                  </div>
                  <div className="flex flex-col gap-0.5 min-w-0">
                    <span className="font-mono text-sm font-medium">
                      {ev.transaction_id || "Unmatched"}
                    </span>
                    <span className="text-xs text-muted">{ev.phone_number}</span>
                  </div>
                </div>

                <div className="flex items-center gap-4 shrink-0">
                  <span className="font-mono text-sm font-medium">
                    {ev.amount > 0 ? `KES ${ev.amount.toLocaleString()}` : "—"}
                  </span>
                  <span className="text-xs text-muted hidden sm:inline">
                    {new Date(ev.received_at).toLocaleString()}
                  </span>
                  <ArrowUpRight
                    size={15}
                    className="text-muted opacity-0 group-hover:opacity-100 transition-opacity"
                  />
                </div>
              </div>

              <div className="flex gap-1.5 pl-12">
                <span className="text-xs text-muted bg-surface px-2 py-0.5 rounded-full">
                  {ev.event_type}
                </span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </main>
  );
}