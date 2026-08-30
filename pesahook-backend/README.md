# PesaHook

A reliable webhook gateway for Daraja (M-Pesa) callbacks — built for any developer, in any language.

Daraja's callback payloads are inconsistent, undocumented in practice, and painful to debug locally. PesaHook sits between Safaricom and your server: it normalizes STK Push, C2B, and B2C callbacks into one consistent JSON shape, guarantees delivery with automatic retries, and gives you a dead-letter log and manual replay when your server was down.

## Features

- **Payload normalization** — STK Push, C2B, and B2C callbacks are mapped into one consistent event shape, regardless of Daraja's inconsistent raw formats
- **Guaranteed delivery** — automatic retries with backoff if your server is down or slow
- **Dead-letter logging** — failed deliveries (after all retries) are recorded, not lost
- **Manual replay** — resend any captured event to your server on demand
- **HMAC-signed deliveries** — verify that incoming webhooks really came from PesaHook (see [Verifying webhook signatures](#verifying-webhook-signatures))
- **Language-agnostic** — plain REST/JSON API; use it from Node, Python, PHP, or any language that can make an HTTP request
- **Per-endpoint routing** — one registered endpoint per shortcode, so multiple Paybills/Tills stay cleanly separated

## Use cases

- **SaaS billing** — reliably capture C2B payment confirmations for subscription/invoice systems without building your own retry logic
- **E-commerce checkout** — STK Push payment confirmation for online stores, with dead-letter visibility if your webhook receiver goes down during a sale
- **Payroll/disbursement systems** — B2C result tracking for salary or refund payouts
- **Local dev/testing** — a stable, inspectable callback URL instead of juggling ngrok tunnels and reading raw Daraja JSON

## Quick start

1. Register an endpoint:
```bash
curl -X POST https://your-pesahook-instance/endpoints \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "daraja",
    "shortcode": "600000",
    "event_types": ["stk_push", "c2b_confirmation"],
    "destination_url": "https://yourapp.com/webhooks/mpesa"
  }'
```

2. Point your Daraja callback URL (in the Daraja portal or your STK Push request) at the returned `ingest_url`.

3. Your server receives normalized, signed events at `destination_url` — no more parsing Daraja's raw payload shapes yourself.

## Verifying webhook signatures

Every event PesaHook delivers to your `destination_url` includes an `X-PesaHook-Signature` header — an HMAC-SHA256 signature computed using the `secret` you received when registering your endpoint. Verify it before trusting the payload, to confirm the request genuinely came from PesaHook and wasn't spoofed.

**Go:**
```go
import "github.com/Its-Delimas/pesaHook/pkg/verify"

signature := r.Header.Get("X-PesaHook-Signature")
body, _ := io.ReadAll(r.Body)

if !verify.Verify(body, yourEndpointSecret, signature) {
    http.Error(w, "invalid signature", http.StatusUnauthorized)
    return
}
```

**Any other language:** compute an HMAC-SHA256 of the raw request body using your endpoint's `secret` as the key, hex-encode the result, and compare it (using a constant-time comparison, not `==`) against the `X-PesaHook-Signature` header.

## Status

v1.0.0 — core pipeline complete (ingest, normalize, deliver, dead-letter, replay). Currently backed by an in-memory store; Postgres persistence planned for v1.1.

## License

MIT [LICENSE](LICENSE)