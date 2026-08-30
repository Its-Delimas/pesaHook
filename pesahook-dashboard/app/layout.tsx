import {Archivo_Narrow,Inter,JetBrains_Mono} from "next/font/google"
import "./globals.css"

const archivoNarrow = Archivo_Narrow({
  subsets:["latin"],
  variable:"--font-display",
  weight:["500","700"]
})

const inter = Inter({
  subsets:["latin"],
  variable:"--font-body",
});

const jetbrainsMono = JetBrains_Mono({
  subsets:["latin"],
  variable:"--font-body"
});

export default function RootLayout({
  children,
}:{
  children:React.ReactNode;
}){
  return (
    <html lang="en">
      <body className={`${archivoNarrow.variable}${inter.variable}${jetbrainsMono.variable}`}>
        {children}
      </body>
    </html>
  );
}