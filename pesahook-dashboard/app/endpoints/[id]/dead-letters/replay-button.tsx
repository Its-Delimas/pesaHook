"use client"

import { useState } from "react"
import { RotateCw,CheckCircle2,AlertCircle } from "lucide-react"
import { replayEvent } from "@/lib/actions"

export function ReplayButton ({eventId}:{eventId:string}){
    const [status,setStatus] = useState<"idle" | "loading" | "done" | "error">("idle");

    async function handleReplay(){
        setStatus("loading");
        const result = await replayEvent(eventId);
        setStatus(result.success ? "done" : "error");
    }

    if (status==="done"){
        return (
            <span className="inline-flex items-center gap-1.5 text-sm text-accent font-medium">
                <CheckCircle2 size={15} />
                Delivered
            </span>
        );
    }

    return (
        <button onClick={handleReplay} disabled={status==="loading"}
            className={`inline-flex items-center gap-1.5 text-sm font-medium px-3 py-1.5 rounded-lg border transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 ${
                status === "error" 
                    ? "border-danger/30 text-danger hover:bg-danger-soft focus-visible:outline-danger" 
                    : "border-border text-ink hover:bg-surface focus-visible:outline-ink"
            }`}>
                {status === "loading" ? (
                    <RotateCw size={14} className="animate-spin" />
                ): status === "error" ? (
                    <AlertCircle size={14} />
                ):(
                    <RotateCw size={14} />
                )}
                {status === "loading" ? "Replaying" : status === "error" ? "Retry" : "Replay"}
            </button>
    );
}