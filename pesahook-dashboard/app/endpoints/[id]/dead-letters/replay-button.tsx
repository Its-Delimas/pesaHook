"use client"

import { useState } from "react"
import { replayEvent } from "@/lib/actions"

export function ReplayButton({eventID}:{eventID:string}){
    const [status, setStatus] = useState<"idle" | "loading" | "done" | "error">("idle");

    async function handleReplay(){
        setStatus('loading')
        const result = await replayEvent(eventID);
        setStatus(result.success ? "done":"error")
    }

    if (status === "done"){
        return <span className="font-mono text-xs text-ledger">Delivered ✓</span>;
    }

    return (
        <button
            onClick={handleReplay}
            disabled={status==="loading"}
            className="font-mono text-xs uppercase tracking-wide px-3 py-1.5 border border-ink/30 hover:bg-ink hover:text-paper transition-colors disabled:opacity-40"
        >
            {status === "loading" ? "Replaying..." : status === "error" ?"Failed - retry":"Replay"}
        </button>
    )
}