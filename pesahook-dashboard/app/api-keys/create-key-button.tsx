"use client";

import { useState } from "react";
import { Plus, Copy, Check, X } from "lucide-react";
import { createAPIKey } from "../../lib/actions";
import { useRouter } from "next/navigation";

export function CreateKeyButton() {
  const [newKey, setNewKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  async function handleCreate() {
    setLoading(true);
    const result = await createAPIKey();
    setNewKey(result.api_key);
    setLoading(false);
  }

  function handleCopy() {
    if (!newKey) return;
    navigator.clipboard.writeText(newKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  function handleDismiss() {
    setNewKey(null);
    router.refresh();
  }

  if (newKey) {
    return (
      <div className="fixed inset-0 bg-ink/40 flex items-center justify-center z-40 px-4">
        <div className="bg-white rounded-2xl shadow-lg max-w-md w-full p-6">
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-sm font-semibold">Your new API key</h2>
            <button onClick={handleDismiss} className="text-muted hover:text-ink">
              <X size={18} />
            </button>
          </div>
          <p className="text-sm text-muted mb-4">
            Copy this now — you won&apos;t be able to see it again.
          </p>
          <div className="flex items-center gap-2 bg-surface rounded-lg px-3 py-2.5 mb-4">
            <code className="font-mono text-xs flex-1 truncate">{newKey}</code>
            <button
              onClick={handleCopy}
              className="text-muted hover:text-ink shrink-0"
              aria-label="Copy key"
            >
              {copied ? <Check size={16} className="text-accent" /> : <Copy size={16} />}
            </button>
          </div>
          <button
            onClick={handleDismiss}
            className="w-full text-sm font-medium bg-ink text-white py-2 rounded-lg hover:bg-ink/90 transition-colors"
          >
            Done
          </button>
        </div>
      </div>
    );
  }

  return (
    <button
      onClick={handleCreate}
      disabled={loading}
      className="inline-flex items-center gap-1.5 text-sm font-medium bg-ink text-white pl-3 pr-4 py-2 rounded-lg shadow-sm hover:bg-ink/90 hover:shadow-md active:scale-[0.98] transition-all disabled:opacity-50"
    >
      <Plus size={15} strokeWidth={2.5} />
      {loading ? "Creating…" : "New API key"}
    </button>
  );
}