'use client'

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Webhook } from "lucide-react";

export function Sidebar() {
    const pathname = usePathname();
    const isEndpoints = pathname === "/" || pathname.startsWith("/endpoints")
  return (
    <aside className="w-56 shrink-0 border-r border-border bg-surface min-h-screen px-4 py-6 hidden md:flex md:flex-col gap-1">
        <div className="flex items-center gap-2 px-2 mb-6">
            <div className="w-6 h-6 rounded-md bg-ink flex items-center justify-center ">
                <Webhook size={14}className="text-white" strokeWidth={2.5} />
            </div>
            <span className="text-lg font-semibold tracking-tight">pesahook</span>
        </div>


      <Link
        href="/"
        className={`flex items-center gap-2.5 text-sm px-2 py-1.5 rounded-md transition-colors ${
          isEndpoints
            ? "bg-white text-ink font-medium shadow-sm"
            : "text-muted hover:text-ink hover:bg-white/60"
        }`}
      >
        <Webhook size={15} strokeWidth={2} />
        Endpoints
      </Link>

            <div className="mt-auto px-2 pt-4 border-t border-border">
        <p className="text-xs text-muted">v1.0.0</p>
      </div>
    </aside>
  );
}