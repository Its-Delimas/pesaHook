export async function GET (){
    const res = await fetch (`${process.env.PESAHOOK_API_KEY}/endpoints`,{
        headers:{
            Authorization: `Bearer ${process.env.PESAHOOK_API_KEY}`,
        },
        cache:"no-store"
    });
    if (!res.ok){
        return Response.json({error:"failed to fetch endpoint"},{status:res.status});
    }

    const data = await res.json();
    return Response.json(data)
}
// odo: create the GET /endpoints for the backend