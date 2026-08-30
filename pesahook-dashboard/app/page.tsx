import { getEndpoints,type Endpoint } from "@/lib/pesahook";
import Link from "next/link";

export default async function EndpointsPage(){
  const endpoints = await getEndpoints();
  return (
   <main className="min-h-screen bg-paper px-8 py-12 font-body text-ink" >
    <h1 className="font-display text-2xl font-bold tracking-tight mb-1">
      Endpoints
    </h1>
    <p className="text-sm text-ink/60 mb-8">
      Registered destinations receiving your Daraja callbacks
    </p>

    <div className="border border-ink/15 divide-y divide-ink/10">
      {endpoints.map((ep)=>(
          <Link href={`/endpoints/${ep.id}`} key={ep.id} className="flex items-center justify-between px-5 py-4 hover:bg-ink/[0.02] transition-colors">
              <div key={ep.id} className="flex items-center justify-between px-5 py-4 hover:bg-ink/[0.02]">
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-1">
                  <span className="font-mono text-sm font-medium">{ep.shortcode}</span>
                  <span className="text-xs uppercase tracking-wide text-ledger">{ep.provider}</span> 
                </div>
                  <span className="font-mono text-xs text-ink/50">{ep.destination_url}</span>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex gap-1">
                {ep.event_types.map((type)=>(
                  <span key={type} className="font-mono text-[10px] uppercase px-2 py-0.5 border border-ink/20 test-ink/60">
                    {type}
                  </span>
                ))}
              </div>
              <span className="font-mono text-xs text-ink/40">
                {new Date(ep.created_at).toLocaleDateString()}
              </span>
            </div>
            </div>
          </Link>
      ))}
    </div>
   </main>
  )
}