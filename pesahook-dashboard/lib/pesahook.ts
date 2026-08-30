
export type Endpoint = {
  id: string;
  provider: string;
  shortcode: string;
  event_types: string[];
  destination_url: string;
  created_at: string;
};

export type NormalizedEvent = {
  event_id:string;
  endpoint_id:string;
  event_type: string;
  status: "success" | "failed";
  amount:number,
  transaction_id:string;
  phone_number:string;
  received_at:string;
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

// TODO: replace mock data with real GET /events?endpoint_id=
export async function getEventsForEndpoints( endpointId: string): Promise<NormalizedEvent[]> {
  const mockEvents : NormalizedEvent[] = [
    {
      event_id:"evt-1",
      endpoint_id:"53146ddac1f3b124",
      event_type:"stk_push",
      status:"success",
      amount:1500,
      transaction_id:"NLJ7RT61SV",
      phone_number:"254708374149",
      received_at:"2026-08-29T10:15:00Z"
    },
    {
      event_id: "evt-2",
      endpoint_id: "53146ddac1f3b124",
      event_type: "stk_push",
      status: "failed",
      amount: 0,
      transaction_id: "",
      phone_number: "254712345678",
      received_at: "2026-08-29T09:42:00Z",
    },
    {
      event_id: "evt-3",
      endpoint_id: "53146ddac1f3b124",
      event_type: "c2b_confirmation",
      status: "success",
      amount: 500,
      transaction_id: "OKJ8RT62SW",
      phone_number: "254798765432",
      received_at: "2026-08-28T16:03:00Z",
    },
  ];

  return mockEvents.filter((e)=>e.endpoint_id === endpointId);
}

