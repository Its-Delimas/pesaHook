"use server";

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