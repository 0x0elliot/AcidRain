"use client"

import React from "react"
import { siteConfig } from "@/app/siteConfig";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import WebPushNotificationsCampaignSettings from "@/components/WebPushNotificationsCampaignSettings"
import WebPushNotificationThemeSettings from "@/components/WebPushNotificationThemeSettings"
import QuickWebPushCampaign from "@/components/QuickWebPushCampaign"

import CampaignHistory from "@/components/ui/popups/CampaignHistory"

import { Button } from "@/components/ui/button"

import cookies from 'nookies';

export default function Campaigns() {
    const [shops, setShops] = React.useState([]);
    const [userinfo, setUserinfo] = React.useState({});

    const [historyModal, setHistoryModal] = React.useState(false);

    React.useEffect(() => {
        document.title = "Quick Campaigns | Dashboard"

        let accessToken_ = cookies.get(null).access_token;

        let userinfo = {};
        let userinfoStr = localStorage.getItem('userinfo');
        if (userinfoStr) {
            try {
                userinfo = JSON.parse(userinfoStr);
            } catch (e) {
                window.location.href = "/logout";
            }
        }

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

    // const [notificationTypes, setNotificationTypes] = React.useState(["WhatsApp", "SMS", "Email", "Push"]);
    const [notificationTypes, setNotificationTypes] = React.useState(["Push"]);

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
                    {/* Select a campaign type and a segment to send out your first quick message! */}
                    Set up and send out your first quick campaign!
                </p>
            </section>
            {/* Open CampaignHistory modal */}
            { historyModal && <CampaignHistory onClose={() => setHistoryModal(false)} isOpen={historyModal} /> }
            
            { shops.length !== 0 ? (
            <div className="my-tabs mt-4">
                <div className="w-[400px]">
                    {notificationTypes.map((type) => (
                        <div key={type} className="mb-6">
                            <div className="p-4" style={{ width: "max-content" }}>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                                    {type} Campaigns
                                </h2>
                                <p className="text-sm text-gray-500 dark:text-gray-400">
                                    Try out {type} campaigns for your shop.
                                </p>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
                                    {type === "Push" && <WebPushNotificationsCampaignSettings />}
                                    {type === "Push" && <WebPushNotificationThemeSettings storeName={shops[0]?.shop_identifier} />}
                                    {type === "Push" && <QuickWebPushCampaign />}
                                    {type === "Push" && (
                                        <>
                                            <div className="flex bg-white dark:bg-gray-900">
                                                <div className="flex-1 p-4">
                                                    <h1 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Campaign History</h1>
                                                    <p style={{ marginBottom: 10 }} className="text-sm text-gray-600 dark:text-gray-400">All the campaigns launched in the past.</p>
                                                    {/* <div className="flex justify-center p-4 mt-10 space-x-4" style={{ marginTop: "2rem" }}> */}
                                                    <div className="flex flex-col space-y-4">

                                                        <Button
                                                            className="text-white bg-black hover:bg-gray-800 focus:ring-2 focus:ring-gray-500 font-medium text-sm px-5 py-2.5 text-center dark:bg-white dark:text-black dark:hover:bg-gray-200 dark:focus:ring-gray-400 disabled:opacity-50"
                                                            // on click,  open CampaignHistory
                                                            onClick={() => {
                                                                window.location.href = "/campaigns/history";
                                                            }}
                                                        >
                                                            View Campaigns
                                                        </Button>
                                                    </div>
                                                </div>
                                            </div>
                                            </>
                                    )}
                                        </div>
                                </div>
                            </div>
                    ))}
                        </div>
            </div>) : (
                <div className="flex flex-col items-center justify-center h-[calc(100vh-200px)]">
                    <div className="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md text-center">
                        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
                            Create a shop first!
                        </h2>
                        <p className="text-gray-600 dark:text-gray-300 mb-6">
                            You need to create a shop from the dashboard before you can view campaign history.
                        </p>
                        <Button
                            className="text-white bg-black hover:bg-gray-800 focus:ring-2 focus:ring-gray-500 font-medium text-sm px-5 py-2.5 text-center dark:bg-white dark:text-black dark:hover:bg-gray-200 dark:focus:ring-gray-400"
                            onClick={() => {
                                window.location.href = "/dashboard";
                            }}
                        >
                            Go to Dashboard
                        </Button>
                    </div>
                </div>
            )}


            </>
            )
}

