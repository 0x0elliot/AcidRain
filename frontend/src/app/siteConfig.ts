export const siteConfig = {
  name: "zappush",
  url: "https://dashboard.tremor.so",
  description: "The only dashboard you will ever need.",
  baseLinks: {
    home: "/",
    dashboard: "/dashboard",
    campaigns: "/campaigns",
    campaign_history: "/campaigns/history",
  },
  externalLink: {
    blocks: "https://blocks.tremor.so/templates#dashboard",
  },
  baseApiUrl: "http://localhost:5002",
}

export type siteConfig = typeof siteConfig
