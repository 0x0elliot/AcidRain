"use client"

import React from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export default function Campaigns() {

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
                <TabsTrigger value="sms">SMS</TabsTrigger>
                <TabsTrigger value="email">Email</TabsTrigger>
                <TabsTrigger value="whatsapp">WhatsApp</TabsTrigger>
                </TabsList>


                <TabsContent value="sms">
                <div className="p-4">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                    SMS Campaigns
                    </h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        Try out SMS campaigns for your shop.
                    </p>
                </div>
                </TabsContent>

                <TabsContent value="email">
                <div className="p-4">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                    Email Campaigns
                    </h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        Try out Email campaigns for your shop.
                    </p>
                </div>
                </TabsContent>

                <TabsContent value="whatsapp">
                <div className="p-4">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                    WhatsApp Campaigns
                    </h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        Try out WhatsApp campaigns for your shop.
                    </p>
                </div>
                </TabsContent>
            </Tabs>
        </div>


        </>
    )
}

