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