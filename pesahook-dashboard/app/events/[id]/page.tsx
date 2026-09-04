import Link from "next/link";
import { ArrowLeft, CheckCircle2, XCircle } from "lucide-react";
import { getEventById } from "../../../lib/pesahook";
import { notFound } from "next/navigation";

export default async function EventDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const event = await getEventById(id);

  if (!event) notFound();

  return (
    <main className="px-6 py-10 max-w-5xl mx-auto">
      <Link
        href={`/endpoints/${event.endpoint_id}`}
        className="inline-flex items-center gap-1.5 text-sm text-muted hover:text-ink transition-colors mb-6"
      >
        <ArrowLeft size={14} />
        Events
      </Link>

      <div className="flex items-center gap-3 mb-1">
        <div
          className={`w-9 h-9 rounded-xl flex items-center justify-center shadow-sm ${
            event.status === "success" ? "bg-accent-soft" : "bg-danger-soft"
          }`}
        >
          {event.status === "success" ? (
            <CheckCircle2 size={16} className="text-accent" strokeWidth={2} />
          ) : (
            <XCircle size={16} className="text-danger" strokeWidth={2} />
          )}
        </div>
        <h1 className="text-2xl font-semibold tracking-tight">
          {event.transaction_id || "Unmatched event"}
        </h1>
      </div>
      <p className="font-mono text-xs text-muted mb-8 pl-12">{event.event_id}</p>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <section>
          <h2 className="text-sm font-semibold mb-3">Normalized</h2>
          <dl className="bg-white border border-border rounded-2xl divide-y divide-border font-mono text-sm overflow-hidden shadow-sm">
            {Object.entries({
              type: event.event_type,
              status: event.status,
              amount: `KES ${event.amount.toLocaleString()}`,
              phone: event.phone_number,
              received: new Date(event.received_at).toLocaleString(),
            }).map(([key, value]) => (
              <div key={key} className="flex justify-between px-4 py-3">
                <dt className="text-muted">{key}</dt>
                <dd className="font-medium">{value}</dd>
              </div>
            ))}
          </dl>
        </section>

        <section>
          <h2 className="text-sm font-semibold mb-3">Raw payload</h2>
          <pre className="bg-ink text-white/90 rounded-2xl p-4 font-mono text-xs overflow-x-auto whitespace-pre-wrap shadow-sm">
            {JSON.stringify(event.raw, null, 2)}
          </pre>
        </section>
      </div>
    </main>
  );
}