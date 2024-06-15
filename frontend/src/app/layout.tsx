"use client"

import { ThemeProvider } from "next-themes"
import { Inter } from "next/font/google"
import "./globals.css"

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter",
})

import { Sidebar } from "@/components/ui/navigation/sidebar"
import { siteConfig } from "./siteConfig"
import { isLastDayOfMonth } from "date-fns";

// export const metadata: Metadata = {
//   metadataBase: new URL("https://yoururl.com"),
//   title: siteConfig.name,
//   description: siteConfig.description,
//   keywords: [],
//   authors: [
//     {
//       name: "yourname",
//       url: "",
//     },
//   ],
//   creator: "yourname",
//   openGraph: {
//     type: "website",
//     locale: "en_US",
//     url: siteConfig.url,
//     title: siteConfig.name,
//     description: siteConfig.description,
//     siteName: siteConfig.name,
//   },
//   icons: {
//     icon: "/favicon.ico",
//   },
// }

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body
        className={`${inter.className} overflow-y-scroll scroll-auto antialiased selection:bg-indigo-100 selection:text-indigo-700 dark:bg-gray-950`}
        suppressHydrationWarning
      >
        <ThemeProvider defaultTheme="system" attribute="class">
            <main>{children}</main>
          </ThemeProvider>
      </body>
    </html>
  )
}
