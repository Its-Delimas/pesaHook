import {Geist,Geist_Mono} from "next/font/google";
import './globals.css'
import { Sidebar } from "@/components/sidebar";

const geist = Geist({
  subsets:["latin"],
  variable:"--font-sans",
});

const geistMono = Geist_Mono ({
  subsets:["latin"],
  variable:"--font-mono",
})

export default function RootLayout ({children}:{children:React.ReactNode;}){
  return (
    <html lang="en">
      <body className={`${geist.variable} ${geistMono.variable} font-sans antialiased`}>
        <Sidebar />
        <div className="flex-1 min-w-0">{children}</div>
      </body>
    </html>
  )
}