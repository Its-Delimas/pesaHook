import Link from "next/link";

export function Sidebar() {
  return (
    <aside className="w-56 shrink-0 border-r border-border bg-surface min-h-screen px-4 py-6 hidden md:flex md:flex-col gap-1">
      <div className="px-2 mb-6">
        <span className="text-sm font-semibold tracking-tight">PesaHook</span>
      </div>

      <Link
        href="/"
        className="text-sm px-2 py-1.5 rounded-md text-ink hover:bg-white transition-colors"
      >
        Endpoints
      </Link>
    </aside>
  );
}