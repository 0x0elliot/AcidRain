"use client"

import React from "react"
import { siteConfig } from "@/app/siteConfig";
import { Button } from "@/components/ui/button"
import cookies from 'nookies';

import axios from 'axios';
import Link from "next/link";

export default function CampaignHistory() {
    const [campaigns, setCampaigns] = React.useState([]);
    const [loading, setLoading] = React.useState(true);
    const [error, setError] = React.useState("");
    const [userinfo, setUserinfo] = React.useState({});

    const [shopOwned, setShopOwned] = React.useState(true);


    React.useEffect(() => {
        document.title = "Campaign History | Dashboard"

        let accessToken_ = cookies.get(null).access_token;
        let userinfoLocal;
        let userInfo = localStorage.getItem('userinfo');

        if (userInfo) {
            try {
                let userinfoJSON = JSON.parse(userInfo);
                setUserinfo(userinfoJSON);
                userinfoLocal = userinfoJSON;
            } catch (e) {
                window.location.href = "/logout";
            }
        }

        if (userinfoLocal?.current_shop_id === null || userinfoLocal?.current_shop_id === undefined) {
            setShopOwned(false);
            return;
        }

        axios.get(`${siteConfig.baseApiUrl}/api/notification/private/notification-campaigns`, {
            params: {
                shop_id: userinfoLocal?.current_shop_id
            },
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${accessToken_}`
            }
        })
            .then(response => {
                setCampaigns(response.data.campaigns);
                setLoading(false);
            })
            .catch(error => {
                if (error.response) {
                    // The request was made and the server responded with a status code
                    // that falls out of the range of 2xx
                    setError("Failed to fetch campaigns");
                } else if (error.request) {
                    // The request was made but no response was received
                    setError("No response received from server");
                } else {
                    // Something happened in setting up the request that triggered an Error
                    setError("An error occurred while fetching campaigns");
                }
                setLoading(false);
            });

    }, [])

    return (
        <>
            <section aria-labelledby="history-title">
                <h1
                    id="history-title"
                    className="scroll-mt-10 text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50"
                >
                    Campaign History
                </h1>
            </section>

            <section aria-labelledby="history-description">
                <p
                    id="history-description"
                    className="text-sm text-gray-500 dark:text-gray-400"
                >
                    View all your past campaigns and their performance.
                </p>
            </section>

            {shopOwned ? (
                <div className="mt-6">
                    {loading ? (
                        <p className="text-gray-600 dark:text-gray-400">Loading campaigns...</p>
                    ) : error ? (
                        <p className="text-red-500">{error}</p>
                    ) : (
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
                            {campaigns.map((campaign) => (
                                <div key={campaign.id} className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
                                    <div className="p-6">
                                        <div className="flex justify-between items-start mb-4">
                                            <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
                                                {campaign.notification_configuration.title || "Untitled Campaign"}
                                            </h3>
                                            <span className="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">
                                                {new Date(campaign.created_at).toLocaleDateString()}
                                            </span>
                                        </div>
                                        <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                                            {campaign.notification_configuration.message || "No message provided"}
                                        </p>
                                        <div className="flex flex-wrap gap-2 mb-4">
                                            <span className="px-2 py-1 text-xs font-semibold text-blue-800 bg-blue-100 rounded-full">
                                                Shop: {campaign.shop.name}
                                            </span>
                                            <span className="px-2 py-1 text-xs font-semibold text-purple-800 bg-purple-100 rounded-full">
                                                Platform: {campaign.shop.platform}
                                            </span>
                                        </div>
                                        <div className="text-sm text-gray-500 dark:text-gray-400 mb-4">
                                            <p>store URL: <a href={campaign.notification_configuration.url} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">{campaign.notification_configuration.url}</a></p>
                                        </div>
                                        <div className="flex justify-end">
                                            <Link
                                                href={`/campaigns/history/${campaign.id}/info`}
                                                passHref
                                            >
                                                <Button
                                                    className="text-white bg-black hover:bg-gray-800 focus:ring-2 focus:ring-gray-500 font-medium text-sm px-5 py-2.5 text-center dark:bg-white dark:text-black dark:hover:bg-gray-200 dark:focus:ring-gray-400"
                                                    onClick={() => {
                                                        // Add logic to view campaign details
                                                        console.log("View details for campaign", campaign.id);
                                                    }}
                                                >
                                                    View details
                                                </Button>
                                            </Link>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            ) : (
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