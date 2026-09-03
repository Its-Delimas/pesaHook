import {Geist,Geist_Mono} from "next/font/google";
import './globals.css'
// import { Sidebar } from "@/components/sidebar";
import { Nav } from "@/components/nav";

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
        <Nav />
        {children}
      </body>
    </html>
  )
}