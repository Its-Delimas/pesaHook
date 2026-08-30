
export type Endpoint = {
  id: string;
  provider: string;
  shortcode: string;
  event_types: string[];
  destination_url: string;
  created_at: string;
};

// TODO: replace mock data with real call to PesaHook once GET /endpoints exists on the backend
export async function getEndpoints(): Promise<Endpoint[]> {
  return [
    {
      id: "53146ddac1f3b124",
      provider: "daraja",
      shortcode: "600000",
      event_types: ["stk_push", "c2b_confirmation"],
      destination_url: "https://example.com/webhooks/mpesa",
      created_at: "2026-08-26T13:38:15Z",
    },
    {
      id: "7a2c9e1f4b8d3a56",
      provider: "daraja",
      shortcode: "600001",
      event_types: ["b2c_result"],
      destination_url: "https://payroll.example.com/webhooks/payout",
      created_at: "2026-08-20T09:12:00Z",
    },
  ];
}