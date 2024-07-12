import Automations from "./(dashboard)/(main)/automations/page"

export const siteConfig = {
  name: "AcidRain",
  url: "https://dashboard.tremor.so",
  description: "The only dashboard you will ever need.",
  baseLinks: {
    home: "/",
    automations: "/automations",
    campaigns: "/campaigns",
  },
  externalLink: {
    blocks: "https://blocks.tremor.so/templates#dashboard",
  },
  baseApiUrl: "http://localhost:5002",
}

export type siteConfig = typeof siteConfig
