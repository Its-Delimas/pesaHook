import Link from "next/link";
import { ArrowLeft,CheckCircle2 } from "lucide-react";
import { getDeadLettersForEndpoint } from "@/lib/pesahook";
import { ReplayButton } from "./replay-button";

export default async function DeadLettersPage({params}:{params:Promise<{id:string}>;}){
    const {id} = await params
    const deadletters = await getDeadLettersForEndpoint(id)

    return (
        <main className="px-8 py-10 w-full">
            <Link href={`/endpoints/${id}`} className="inline-flex items-center gap-1.5 text-sm text-muted hover:text-ink transition-colors mb-6">
                <ArrowLeft size={14}/> Events
            </Link>

            <div className="mb-8">
                <h1 className="text-xl font-semibold tracking-tight">Dead letters</h1>
                <p className="text-sm text-muted mt-1">
                    Events that failed delivery after every retry. Replay once your destination is back up.
                </p>
            </div>

            {deadletters.length=== 0 ? (
                <div className="border border-dashed border-border rounded-xl px-6 py-16 flex flex-col items-center text-center">
                    <div className="w-10 h-10 rounded-full bg-accent-soft flex items-center justify-center mb-4">
                        <CheckCircle2 size={18} className="text-accent" strokeWidth={1.75} />
                    </div>
                    <p className="text-sm text-ink font-medium mb-1">All clear</p>
                    <p className="text-sm text-muted max-w-xs">
                        Every captured event has been delivered successfully
                    </p>
                </div>
            ):(
                <div className="border border-border rounded-xl overflow-hidden">
                    {deadletters.map((dl,i)=>(
                        <div key={dl.event.event_id} className={`flex items-center justify-between gap-4 px-5 py-4 ${
                            i !== 0 ? "border-t border-border":""
                        }`}>
                            <div className="flex flex-col gap-1 min-w-0">
                                <span className="font-mono text-sm">{dl.event.transaction_id}</span>
                                <span className="text-xs text-danger">{dl.last_error}</span>
                                <span className="text-xs text-muted">
                                    {dl.attempts} attempts · failed {new Date(dl.failed_at).toLocaleString()}
                                </span>
                            </div>

                            <ReplayButton eventID={dl.event.event_id} />
                        </div>
                    ))}
                </div>
            )}
        </main>
    );
}