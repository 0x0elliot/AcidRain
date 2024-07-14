"use client"

import React from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export default function Campaigns() {

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
            <Tabs defaultValue="sms" className="w-[400px]">
                
                <TabsList>
                    {notificationTypes.map((type) => (
                        <TabsTrigger key={type} value={type.toLowerCase()}>
                            {type}
                        </TabsTrigger>
                    ))}
                </TabsList>

                {notificationTypes.map((type) => (
                    <TabsContent key={type} value={type.toLowerCase()}>
                        <div className="p-4">
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                                {type} Campaigns
                            </h2>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                Try out {type} campaigns for your shop.
                            </p>
                        </div>
                    </TabsContent>
                ))}

            </Tabs>
        </div>


        </>
    )
}

