import { getEventById } from "@/lib/pesahook";
import Link from "next/link";
import { notFound } from "next/navigation";

export default async function   EventDetailPage({params}:{params:Promise<{id:string}>;}){
    const {id} = await params
    const event = await getEventById(id)
    
    if (!event) notFound();

    return (
        <main className="min-h-screen bg-paper px-8 py-12 font-body text-ink">
            <Link href={`/endpoints/${event.endpoint_id}`} className="text-xs text-ink/50 hover:text-ink mb-6 inline-block">
                   ← Events
            </Link>

            <div className="flex items-center gap-3 mb-1">
                <span className={`w-2 h-2 rounded-fill ${event.status === "success" ? "bg-ledger": "-brick"}`} />
                <h1 className="font-display text-2xl font-bold tracking-tight">
                    {event.transaction_id || "Unmatched event"}
                </h1>
            </div>
            <p className="font-mono text-xs text-ink/50 mb-8">{event.event_id  }</p>

            <div className="grid grid-cols-2 gap-8">
                <section>
                    <h2 className="text-xs uppercase tracking-wide text-ink/50 mb-3">Normalized</h2>
                    <dl className="border border-ink/15 divide-y divide-ink/10 font-mono text-xs">
                        {Object.entries({
                            type: event.event_type,
                            status: event.status,
                            amount: `KES ${event.amount.toLocaleString()}`,
                            phone:event.phone_number,
                            received: new Date(event.received_at).toLocaleString(),
                        }).map(([key,value])=>(
                            <div key={key} className="flex justify-between px-4 py-2.5">
                                <dt className="text-ink/50">{key}</dt>
                                <dd>{value}</dd>
                            </div>
                        ))}
                    </dl>
                </section>

                <section>
                    <h2 className="text-xs uppercase tracking-wide text-ink/50 mb-3">
                        Raw payload
                    </h2>
                    <pre className="border border-ink/15 bg-ink/[0.02] p-4 font-mono text-xs overflow-x-auto whitespace-pre-wrap">
                        {JSON.stringify(event.raw,null,2)}
                    </pre>
                </section>
            </div>
        </main>
    )

}