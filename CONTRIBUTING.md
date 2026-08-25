# Contributing to PesaHook

Thanks for considering contributing — PesaHook is early-stage infra, and feedback/PRs from real Daraja integrators are especially valuable.

## Getting started

1. Fork and clone the repo
2. Requires Go 1.25+
3. Run the test suite before making changes, to confirm your environment is clean:
```bash
   go test ./...
```

## Making changes

- Keep provider-specific logic isolated inside its own package (e.g. `internal/daraja/`) — the normalization contract (`event.NormalizedEvent`) should stay provider-agnostic so other providers (MTN MoMo, etc.) can be added without touching shared code
- Add or update tests for any change to a mapper, handler, or delivery logic — mapper bugs silently corrupt data, so coverage matters more here than in most projects
- Match existing patterns: interfaces for anything that needs a swappable implementation (see `store.EventStore`), constructors (`NewX`) over exported structs with public fields where state needs to stay consistent

## Submitting a PR

- One logical change per PR — easier to review, easier to revert if something's wrong
- Include a clear description of what changed and why, especially for anything touching payload parsing (link to Daraja docs or a real payload example if relevant)
- Run `go vet ./...` and `gofmt -l .` before submitting — CI will check for both

## Reporting bugs

If you've hit a real Daraja payload shape that doesn't normalize correctly, please include:
- The raw payload (redact any real phone numbers, amounts, or transaction IDs)
- What you expected vs what PesaHook produced
- Which event type (STK Push, C2B, B2C)

## Questions / ideas

Open an issue — especially for provider requests or features you'd want as a consumer of the gateway.# Contributing to PesaHook

Thanks for considering contributing — PesaHook is early-stage infra, and feedback/PRs from real Daraja integrators are especially valuable.
