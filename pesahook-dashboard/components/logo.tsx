export function Logo ({size = 24}:{size?:number}){
    return (
        <div className="rounded-lg bg-accent flex items-center justify-center shrink-0" style={{width:size,height:size}}>
            <svg width={size*0.6} height={size*0.6} viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M18 16.98h-5.99c-1.1 0-1.95.94-2.48 1.9A4 4 0 0 1 2 17c.01-.7.2-1.4.57-2M6 6.3a4 4 0 0 1 6.6 4.34" />
                <path d="M9 12.8v2.7m0-13.5v3M6.8 15H2m14.3-9.3-1.9 1.9M18 6.3l-2.5-2.4M7.5 6.8 6.3 8" />
                <circle cx="12" cy="7" r="2.5" />
                <circle cx="4" cy="17" r="2.5" />
                <circle cx="19" cy="17" r="2.5" />
            </svg>
        </div>
    )
}