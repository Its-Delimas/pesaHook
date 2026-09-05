import { pesahookFetch } from "./api";

export type Endpoint = {
  id:string;
  provider:string;
  shortcode:string;
  event_types:string[];
  destination_url: string;
  ingest_url: string;
  created_at:string;
};

export async function getEndpoints(): Promise<Endpoint[]> {
  return pesahookFetch("/endpoints");
}

export type NormalizedEvent = {
  event_id: string;
  endpoint_id:string;
  event_type:string;
  status: "success" | "failed";
  amount: number;
  transaction_id: string;
  phone_number : string;
  received_at: string;
}

export async function getEventsForEndpoints(endpointID: string): Promise<NormalizedEvent[]> {
  return pesahookFetch(`/events?endpoint_id=${endpointID}`);
}

export type DeadLetter = {
  event: NormalizedEvent;
  endpoint_id: string;
  last_error:string;
  failed_at:string;
  attempts:number;
};

export async function getDeadLettersForEndpoint(endpointId:string): Promise<DeadLetter[]> {
  return pesahookFetch(`/endpoints/${endpointId}/dead-letters`);
}

export async function getEventById (eventId:string): Promise<(NormalizedEvent & {raw:object}) | null> {
  try {
    return await pesahookFetch(`/events/${eventId}`);
  }catch {
    return null;
  }
}

export type APIKeySummary = {
  id:string;
  created_at:string;
}

export async function getAPIKeys(): Promise<APIKeySummary[]> {
    return pesahookFetch("/api-keys");
}