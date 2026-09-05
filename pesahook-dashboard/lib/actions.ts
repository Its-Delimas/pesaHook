"use server";

import { pesahookFetch } from "./api";

// TODO: wire to real POST /events/{id}/replay once backed route ownership check is confimerd agains dashboard API Key
export async function replayEvent(eventID: string): Promise<{success: boolean;attempts?:number}> {
    const res = await fetch(
        `${process.env.PESAHOOK_API_URL}/events/${eventID}/replay`,
        { 
            method:"POST",
            headers:{Authorization: `BEARER ${process.env.PESAHOOK_API_KEY}`},
        }        
    );
    if (!res.ok){
        return { success: false }
    }
    const data = await res.json();
    return { success:true, attempts:data.attempts};
}

export async function createAPIKey(): Promise<{id:string;api_key:string;created_at:string}>{
    const res = await fetch (`${process.env.PESAHOOK_API_URL}/api-keys`,{
        method:"POST",
        headers:{Authorization:`Bearer ${process.env.PESAHOOK_API_KEY}`},
    });

    if (!res.ok){
        throw new Error("failed to create api key");
    }

    return res.json();
}