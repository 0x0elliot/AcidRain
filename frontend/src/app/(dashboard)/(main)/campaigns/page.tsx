"use client"

import React from "react"
import { siteConfig } from "@/app/siteConfig";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import WebPushNotificationsCampaignSettings from "@/components/WebPushNotificationsCampaignSettings"
import WebPushNotificationThemeSettings from "@/components/WebPushNotificationThemeSettings"

import cookies from 'nookies';

export default function Campaigns() {
    const [shops, setShops] = React.useState([]);

    React.useEffect(() => {
        document.title = "Quick Campaigns | Dashboard"

        let accessToken_ = cookies.get(null).access_token;

        // get /api/shop/private/all
        fetch(`${siteConfig.baseApiUrl}/api/shop/private/all`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${accessToken_}`,
            },
        }).then((res) => {
            if (res.status === 200) {
                res.json().then((data) => {
                    setShops(data.shops);
                });
            }
        });
    }, [])

    const [notificationTypes, setNotificationTypes] = React.useState(["WhatsApp", "SMS", "Email", "Push"]);

    return (   
        <>
            <section aria-labelledby="flows-title">
                <h1
                    id="overall-title"
                    className="scroll-mt-10 text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50"
                >
                        Quick Campaigns
                </h1>
            </section>

            <section aria-labelledby="flows-description">
                <p
                    id="overall-description"
                    className="text-sm text-gray-500 dark:text-gray-400"
                >
                    Select a campaign type and a segment to send out your first quick message!
                </p>
            </section>

        <div className="my-tabs mt-4">
            <Tabs defaultValue="push" className="w-[400px]">
                
                <TabsList>
                    {notificationTypes.map((type) => (
                        <TabsTrigger key={type} value={type.toLowerCase()}>
                            {type}
                        </TabsTrigger>
                    ))}
                </TabsList>

                {notificationTypes.map((type) => (
                    <TabsContent key={type} value={type.toLowerCase()}>
                        <div className="p-4" style={{width: "max-content" }} >
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                                {type} Campaigns
                            </h2>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                Try out {type} campaigns for your shop.
                            </p>

                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
                                {type === "Push" && <WebPushNotificationsCampaignSettings />}
                                {type === "Push" && <WebPushNotificationThemeSettings storeName={shops[0]?.shop_identifier} />}
                            </div>
                        </div>
                    </TabsContent>
                ))}

            </Tabs>
        </div>


        </>
    )
}

