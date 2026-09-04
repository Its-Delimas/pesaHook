"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";
import { Menu, X, Webhook, KeyRound } from "lucide-react";
// import { Logo } from "./logo";

const links = [
  { href: "/", label: "Endpoints", icon: Webhook, match: (p: string) => p === "/" || p.startsWith("/endpoints") },
  { href: "/api-keys", label: "API keys", icon: KeyRound, match: (p: string) => p.startsWith("/api-keys") },
];

export function Nav() {
  const pathname = usePathname();
  const [open, setOpen] = useState(false);

  return (
    <header className="border-b border-border bg-white sticky top-0 z-10">
      <div className="flex items-center justify-between px-6 h-14">
        <div className="flex items-center gap-8">
          <Link href="/" className="flex items-center gap-2">
            {/* <Logo size={26} /> */}
            <span className="text-sm font-semibold tracking-tight">PesaHook</span>
          </Link>

          <nav className="hidden md:flex items-center gap-1">
            {links.map(({ href, label, icon: Icon, match }) => (
              <Link
                key={href}
                href={href}
                className={`flex items-center gap-1.5 text-sm px-3 py-1.5 rounded-md transition-colors ${
                  match(pathname)
                    ? "bg-accent-soft text-accent font-medium"
                    : "text-muted hover:text-ink hover:bg-surface"
                }`}
              >
                <Icon size={15} strokeWidth={2} />
                {label}
              </Link>
            ))}
          </nav>
        </div>

        <button
          onClick={() => setOpen(!open)}
          className="md:hidden p-1.5 text-muted hover:text-ink"
          aria-label={open ? "Close menu" : "Open menu"}
        >
          {open ? <X size={20} /> : <Menu size={20} />}
        </button>
      </div>

      {open && (
        <nav className="md:hidden border-t border-border px-4 py-2 flex flex-col gap-1">
          {links.map(({ href, label, icon: Icon, match }) => (
            <Link
              key={href}
              href={href}
              onClick={() => setOpen(false)}
              className={`flex items-center gap-2 text-sm px-3 py-2 rounded-md transition-colors ${
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
      )}
    </header>
  );
}