// file: pesahook-dashboard/components/sidebar.tsx
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";
import { Webhook, KeyRound, BookOpen, Settings, Menu, X } from "lucide-react";
import { Logo } from "./logo";

const links = [
  { href: "/", label: "Endpoints", icon: Webhook, match: (p: string) => p === "/" || p.startsWith("/endpoints") },
  { href: "/api-keys", label: "API keys", icon: KeyRound, match: (p: string) => p.startsWith("/api-keys") },
];

const secondaryLinks = [
  { href: "https://github.com/Its-Delimas/pesahook", label: "Documentation", icon: BookOpen, external: true },
  { href: "/settings", label: "Settings", icon: Settings, external: false },
];

function NavLinks({ pathname, onNavigate }: { pathname: string; onNavigate?: () => void }) {
  return (
    <>
      <nav className="flex flex-col gap-1">
        {links.map(({ href, label, icon: Icon, match }) => (
          <Link
            key={href}
            href={href}
            onClick={onNavigate}
            className={`flex items-center gap-2.5 text-sm px-3 py-2 rounded-lg transition-colors ${
              match(pathname)
                ? "bg-accent-soft text-accent font-medium"
                : "text-muted hover:text-ink hover:bg-surface"
            }`}
          >
            <Icon size={16} strokeWidth={2} />
            {label}
          </Link>
        ))}
      </nav>

      <div className="mt-6 pt-6 border-t border-border flex flex-col gap-1">
        {secondaryLinks.map(({ href, label, icon: Icon, external }) => (
          <Link
            key={href}
            href={href}
            onClick={onNavigate}
            target={external ? "_blank" : undefined}
            className="flex items-center gap-2.5 text-sm px-3 py-2 rounded-lg text-muted hover:text-ink hover:bg-surface transition-colors"
          >
            <Icon size={16} strokeWidth={2} />
            {label}
          </Link>
        ))}
      </div>
    </>
  );
}

export function Sidebar() {
  const pathname = usePathname();
  const [open, setOpen] = useState(false);

  return (
    <>
      {/* Desktop sidebar */}
      <aside className="hidden md:flex w-60 shrink-0 flex-col border-r border-border bg-white min-h-screen px-4 py-6">
        <Link href="/" className="flex items-center gap-2 px-2 mb-8">
          <Logo size={28} />
          <span className="text-sm font-semibold tracking-tight">PesaHook</span>
        </Link>

        <NavLinks pathname={pathname} />

        <div className="mt-auto px-3 pt-4 border-t border-border">
          <p className="text-xs text-muted">v1.0.0</p>
        </div>
      </aside>

      {/* Mobile top bar */}
      <header className="md:hidden flex items-center justify-between px-4 h-14 border-b border-border bg-white sticky top-0 z-20">
        <Link href="/" className="flex items-center gap-2">
          <Logo size={26} />
          <span className="text-sm font-semibold tracking-tight">PesaHook</span>
        </Link>
        <button
          onClick={() => setOpen(true)}
          className="p-1.5 text-muted hover:text-ink"
          aria-label="Open menu"
        >
          <Menu size={20} />
        </button>
      </header>

      {/* Mobile slide-over drawer */}
      {open && (
        <div className="md:hidden fixed inset-0 z-30 flex">
          <div
            className="absolute inset-0 bg-ink/40"
            onClick={() => setOpen(false)}
          />
          <div className="relative w-64 bg-white h-full px-4 py-6 flex flex-col shadow-lg">
            <div className="flex items-center justify-between mb-8 px-2">
              <div className="flex items-center gap-2">
                <Logo size={26} />
                <span className="text-sm font-semibold tracking-tight">PesaHook</span>
              </div>
              <button onClick={() => setOpen(false)} className="text-muted hover:text-ink" aria-label="Close menu">
                <X size={20} />
              </button>
            </div>
            <NavLinks pathname={pathname} onNavigate={() => setOpen(false)} />
          </div>
        </div>
      )}
    </>
  );
}