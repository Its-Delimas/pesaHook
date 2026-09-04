import { Webhook,Plus,ArrowUpRight } from "lucide-react";
import Link from "next/link";
import { getEndpoints } from "@/lib/pesahook";

export default async function EndpointsPage(){
    const endpoints = await getEndpoints();
    const providerCount = new Set(endpoints.map((e)=>e.provider)).size;
    const totalEventTypes = new Set(endpoints.flatMap((e)=>e.event_types)).size;

    return (
      <main className="px-6 py-10 max-w-5xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-xl font-semibold tracking-tight">Endpoints</h1>
            <p className="text-sm text-muted mt-1">
              Destinations receiving your Daraja callbacks.
            </p>
          </div>
          <button className="inline-flex items-center gap-1.5 text-sm font-medium bg-ink text-white pl-3 pr-4 py-2 rounded-lg hover:bg-ink/90 active:scale-[0.98] transition-all focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible::outline-ink">
            <Plus size={15} strokeWidth={2.5} />
            New endpoint
          </button>
        </div>
      </main>
    )
}

