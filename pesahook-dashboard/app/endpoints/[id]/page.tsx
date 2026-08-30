import Link from "next/link";
import { getEventsForEndpoints } from "@/lib/pesahook";

export default async function EndpointEventsPage ({params,}:{params: Promise<{id: string}>}){
    const {id} = await params
    const events = await getEventsForEndpoints(id);

    return(
        <main className="min-h-screen bg-paper px-8 py-12 font-body text-ink">
            <Link href="/" className="text-xs text-ink/50 hover:text-ink mb-6 inline-block">
                 ← Endpoints
            </Link>

            <h1 className="font-display text-2xl font-bold tracking-tight mb-1">
                Events
            </h1>
            <p className="font-mono text-xs text-ink/50 mb-8">{id}</p>

            <div className="border border-ink/15 divide-y divide-ink/10">
                {events.length===0 && (
                    <div className="px-5 py-8 text-sm text-ink/50">
                        No events captured yet. Once Daraja sends a callback to this endpoint, it will appear here.
                    </div>
                )}

                {events.map((ev)=>(
                    <div key={ev.event_id} className="flex items-center justify-between px-5 py-4 hover:bg-ink/[0.02]">
                        <div className="flex items-center gap-4">
                            <span className={`w-2 h-2 rounded-full ${ev.status==="success"?"bg-ledger":"bg-brick"}`} aria-label={ev.status} />
                            <div className="flex flex-col gap-1">
                                <span className="font-mono text-sm">
                                    {ev.transaction_id || "-"}
                                </span>
                                <span className="text-xs uppercase tracking-wide text-ink/50">
                                {ev.event_type}
                                </span>
                            </div>
                        </div>

                        <div className="flex items-center gap-6">
                            <span className="fnt-mono text-sm">
                                {ev.amount > 0 ? `KES ${ev.amount.toLocaleString()}`:"-"}
                            </span>
                            <span className="font-mono text-xs text-ink/40">
                                {new Date(ev.received_at).toLocaleDateString()}
                            </span>
                        </div>
                    </div>
                ))}
            </div>
        </main>
    )
}