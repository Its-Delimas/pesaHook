import { KeyRound } from "lucide-react";
import { getAPIKeys } from "../../lib/pesahook";
import { CreateKeyButton } from "./create-key-button"

export default async function APIKeysPage() {
  const keys = await getAPIKeys();

  return (
    <main className="px-6 py-10 max-w-5xl mx-auto">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">API keys</h1>
          <p className="text-sm text-muted mt-1">
            Used to authenticate requests to the PesaHook API.
          </p>
        </div>
        <CreateKeyButton />
      </div>

      {keys.length === 0 ? (
        <div className="bg-white border border-dashed border-border rounded-2xl px-6 py-16 flex flex-col items-center text-center shadow-sm">
          <div className="w-10 h-10 rounded-full bg-surface border border-border flex items-center justify-center mb-4">
            <KeyRound size={18} className="text-muted" strokeWidth={1.75} />
          </div>
          <p className="text-sm text-ink font-medium mb-1">No API keys yet</p>
          <p className="text-sm text-muted max-w-xs">
            Create one to start registering endpoints and calling the PesaHook API.
          </p>
        </div>
      ) : (
        <div className="bg-white border border-border rounded-2xl overflow-hidden shadow-sm">
          {keys.map((key, i) => (
            <div
              key={key.id}
              className={`flex items-center justify-between px-5 py-4 ${
                i !== 0 ? "border-t border-border" : ""
              }`}
            >
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 rounded-xl bg-accent-soft flex items-center justify-center">
                  <KeyRound size={16} className="text-accent" strokeWidth={2} />
                </div>
                <div className="flex flex-col gap-0.5">
                  <span className="font-mono text-sm font-medium">ph_••••••••••••••••</span>
                  <span className="text-xs text-muted">
                    Created {new Date(key.created_at).toLocaleDateString()}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </main>
  );
}