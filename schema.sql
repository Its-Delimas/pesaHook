CREATE TABLE endpoints (
    id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    shortcode TEXT NOT NULL,
    event_types TEXT[] NOT NULL,
    destination_url TEXT NOT NULL,
    secret TEXT NOT NULL,
    ingest_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE events(
    event_id TEXT PRIMARY KEY,
    endpoint_id TEXT NOT NULL REFERENCES endpoints(id),
    event_type TEXT NOT NULL,
    provider TEXT NOT NULL,
    shortcode TEXT,
    transaction_id TEXT,
    amount DOUBLE PRECISION,
    phone_number TEXT,
    account_references TEXT,
    status TEXT NOT NULL,
    result_code INT,
    status_reason TEXT,
    occured_at TIMESTAMPTZ,
    received_at TIMESTAMPTZ NOT NULL,
    provider_meta JSONB,
    raw JSONB
);

CREATE INDEX idx_events_endpoint_id ON events(endpoint_id)