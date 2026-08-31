import Link from "next/link";
import { getDeadLettersForEndpoint } from "@/lib/pesahook";
import { ReplayButton } from "./replay-button";


export default async function DeadLettersPage({params}:{params: Promise<{id: string}>;}) {
    const {id} = await params
    const deadLetters = await getDeadLettersForEndpoint(id);

    return (
        <main className="min-h-screen bg-paper px-8 py-12 font-body text-ink">
            <Link href={`/endpoints/${id}`} className="text-xs text-ink/50 hover:text-ink mb-6 inline-block">
                  ← Events
            </Link>

            <h1 className="font-display text-2xl font-bold tracking-tight mb-1">
                Dead Letters
            </h1>
            <p className="text-sm text-ink/60 mb-8">
                Events that failed delivery after every retry. Replay once your destination is back up
            </p>

            <div className="border border-brick/30 divide-y divide-brick/15">
                {deadLetters.length === 0 && (
                    <div className="px-5 py-8 text-sm text-ink/50">
                        No dead letters ~ evenry captured event has been delivered successfully.
                    </div>
                )}

                {deadLetters.map((dl)=>(
                    <div key={dl.event.event_id} className="flex items-center justify-between px-5 py-4">
                        <div className="flex flex-col gap-1">
                            <span className="font-mono text-sm">{dl.event.transaction_id}</span>
                            <span className="text-xs text-brick">{dl.last_error}</span>
                            <span className="font-mono text-[11px] text-ink/40">
                                {dl.attempts} attempts · failed {new Date(dl.failed_at).toLocaleString()}
                            </span>
                        </div>

                        <ReplayButton eventID={dl.event.event_id} />
                    </div>
                ))}
            </div>
        </main>
    )
    
}