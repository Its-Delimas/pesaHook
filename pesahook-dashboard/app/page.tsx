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

        <div className="grid grid-cols-3 gap-4 mb-8">
          <div className="border border-border rounded-xl px-5 py-4">
              <p className="text-xs text-muted mb-1">Endpoints</p>
              <p className="text-2xl font-semibold tracking-tight">{endpoints.length}</p>
          </div>
          <div className="border border-border rounded-xl px-5 py-4">
              <p className="text-xs text-muted mb-1">Providers</p>
              <p className="text-2xl font-semibold tracking-tight">{providerCount}</p>
          </div>
                    <div className="border border-border rounded-xl px-5 py-4">
              <p className="text-xs text-muted mb-1">Event types</p>
              <p className="text-2xl font-semibold tracking-tight">{totalEventTypes}</p>
          </div>
        </div>

        {endpoints.length===0 ? (
          <div className="border border-dashed border-border rounded-xl px-6 py-16 flex flex-col items-center text-center">
            <div className="w-10 h-10 rounded-full bg-surface border border-border flex items-center justify-center mb-8">
                <Webhook size={18} className="text-muted" strokeWidth={1.75} />
            </div>
            <p className="text-sm text-ink font-medium mb-1">No endpoints yet</p>
            <p className="text-sm text-muted mb-4 max-w-xs">
              Register your first shortcode to start receiving normalized callbacks.
            </p>
            <button className="inline-flex items-center gap-1.5 text-sm font-medium bg-ink text-white pl-3 pr-4 py-2 rounded-lg hover:bg-ink/90 transition-colors">
              <Plus size={15} strokeWidth={2.5} />
              New endpoint
            </button>
          </div>
        ):(
          <div className="border border-border rounded-xl overflow-hidden">
              {endpoints.map((ep,i)=>(
                <Link href={`/endpoints/${ep.id}`} key={ep.id} className={`group flex flex-col gap-3 px-5 py-4 hover:bg-surface transition-colors ${
                   i !== 0 ? "border-t border-border":""
                }`}>
                  
                </Link>
              ))}
          </div>
        )};
      </main>
    )
}

