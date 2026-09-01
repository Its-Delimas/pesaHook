import Link from "next/link";
import { ArrowLeft,CheckCircle2,XCircle } from "lucide-react";
import { getEventById } from "@/lib/pesahook";
import { notFound } from "next/navigation";

export default async function EventDetailPage ({params}:{params:Promise<{id:string}>;}){
    const {id} = await params;
    const event = await getEventById(id);

    if (!event) notFound();

    return (
        <main className="px-8 py-10 w-full">
            <Link href={`/endpoints/${event.endpoint_id}`} className="inline-flex items-center gap-1.5 text-sm text-muted hover:text-ink transition-colors mb-6">
                <ArrowLeft size={14} /> Events
            </Link>

            <div className="flex items-center gap-3 mb-1">
                <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${
                    event.status === "success" ? "bg-accent-soft" :"bg-danger-soft"
                }`}>
                    {event.status === "success" ? (
                        <CheckCircle2 size={15} className="text-accent" strokeWidth={2} />
                    ):(
                        <XCircle size={15} className="text-danger" strokeWidth={2} />
                    )}
                </div>
                <h1 className="text-xl font-semibold tracking-tight">
                    {event.transaction_id || "Unmatched event"}
                </h1>
            </div>
            <p className="font-mono text-xs text-muted mb-8 pl-11">{event.event_id}</p>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <section>
                    <h2 className="text-sm font-medium mb-3">Normalized</h2>
                    <dl className="border border-border rounded-xl divide-y divide-border font-mono text-sm overflow-hidden">
                        {Object.entries({
                            type:event.event_type,
                            status:event.status,
                            amount: `KES ${event.amount.toLocaleString()}`,
                            phone:event.phone_number,
                            received: new Date(event.received_at).toLocaleString(),
                        }).map(([Key,value])=>(
                            <div key={Key} className="flex justify-between px-4 py-2.5">
                                <dt className="text-muted">{Key}</dt>
                                <dd>{value}</dd>
                            </div>
                        ))}
                    </dl>
                </section>

                <section>
                    <h2 className="text-sm font-medium mb-3">Raw payload</h2>
                    <pre className="border border-border rounded-xl bg-surface p-4 font-mono text-xs overflow-x-auto whitespace-pre-wrap">
                        {JSON.stringify(event.raw,null,2)}
                    </pre>
                </section>
            </div>
        </main>
    );
}