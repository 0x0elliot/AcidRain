"use client"

import React, { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import axios from 'axios';
import { siteConfig } from "@/app/siteConfig";
import { Button } from "@/components/ui/button";
import Link from 'next/link';
import cookies from 'nookies';

import { Copy } from 'lucide-react';
import { useToast } from "@/components/ui/use-toast"

import dynamic from 'next/dynamic';

const ReactJson = dynamic(() => import('react-json-view'), { ssr: false });

export default function CampaignInfo() {
    const params = useParams();
    const [campaignData, setCampaignData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [selectedNotification, setSelectedNotification] = useState(null);

    const { toast } = useToast()

    const handleCopyId = (id) => {
        navigator.clipboard.writeText(id).then(() => {
            toast({
                title: 'ID copied to clipboard',
                variant: 'success',
                description: 'You can now paste the ID where you need it',

            });
        }
        ).catch((err) => {
            console.error('Failed to copy: ', err);
            toast({
                title: 'Failed to copy ID',
                variant: 'error',
                description: 'Please try again',
            });
        });
    };

    useEffect(() => {
        let userinfo = {};
        let userinfoStr = localStorage.getItem('userinfo');
        if (userinfoStr) {
            try {
                userinfo = JSON.parse(userinfoStr);
            } catch (e) {
                window.location.href = "/logout";
            }
        }

        if (userinfo?.current_shop_id === null) {
            setError("You do not have a shop selected");
            setLoading(false);
            return;
        }

        const fetchCampaignDetails = async (shopId) => {
            const accessToken_ = cookies.get(null).access_token;
            try {
                const response = await axios.get(`${siteConfig.baseApiUrl}/api/notification/private/notification-campaigns?shop_id=${shopId}&notification_campaign_id=${params.campaign_id}`, {
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${accessToken_}`
                    }
                });
                setCampaignData(response.data);
                setLoading(false);
            } catch (err) {
                console.error(err);
                setError("Failed to fetch campaign details");
                setLoading(false);
            }
        };

        fetchCampaignDetails(userinfo?.current_shop_id);
    }, [params.campaign_id]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>{error}</div>;
    if (!campaignData) return <div>No campaign found</div>;

    const { campaign, notifications } = campaignData;

    const handleNotificationClick = (notification) => {
        setSelectedNotification(notification);
    };


    return (
        <div className="p-6 max-w-6xl mx-auto">
            <h1 className="text-2xl font-bold mb-4">{campaign.notification_configuration.title || "Untitled Campaign"}</h1>
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
                <h2 className="text-xl font-semibold mb-4">Campaign Details</h2>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <p className="text-sm text-gray-500 dark:text-gray-400">Campaign ID</p>
                        <p className="text-lg font-semibold">{campaign.id}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500 dark:text-gray-400">Campaign Message</p>
                        <p className="text-lg font-semibold">{campaign.notification_configuration.message}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500 dark:text-gray-400">Campaign URL</p>
                        <p className="text-lg font-semibold">{campaign.notification_configuration.url || "No URL"}</p>
                    </div>
                    <div>
                        <p className="text-sm text-gray-500 dark:text-gray-400">Campaign Created At</p>
                        <p className="text-lg font-semibold">{new Date(campaign.created_at).toLocaleString()}</p>
                    </div>

                    <div>
                        <p className="text-sm text-gray-500 dark:text-gray-400">Campaign Icon</p>
                        {campaign.notification_configuration.icon ? (
                            <img src={campaign.notification_configuration.icon} alt="Campaign Icon" className="h-12 w-12" onClick={() => window.open(campaign.notification_configuration.icon)} />
                        ) : <p className="text-lg font-semibold">No Icon</p>}

                    </div>

                    <div>
                        <p className="text-sm text-gray-500 dark:text-gray-400">Campaign Badge</p>
                        {campaign.notification_configuration.badge ? (
                            <img src={campaign.notification_configuration.badge} alt="Campaign Badge" className="h-12 w-12" onClick={() => window.open(campaign.notification_configuration.badge)} />
                        ) : <p className="text-lg font-semibold">No Badge</p>}

                    </div>

                </div>

            </div>

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
                <h2 className="text-xl font-semibold mb-4">Notifications Sent</h2>
                <div className="overflow-x-auto">
                    <table className="min-w-full table-auto">
                        <thead>
                            <tr className="bg-gray-200 dark:bg-gray-700">
                                <th className="px-4 py-2">Status</th>
                                <th className="px-4 py-2">Sent At</th>
                                <th className="px-4 py-2">API Status</th>
                                <th className="px-4 py-2">Actions</th>
                                <th className="px-4 py-2">Copy ID</th>
                            </tr>
                        </thead>
                        <tbody>
                            {notifications.map((notification, index) => (
                                <tr key={notification.id} className={index % 2 === 0 ? 'bg-gray-100 dark:bg-gray-600' : ''}>
                                    <td className="px-4 py-2 text-center">{notification.status}</td>
                                    <td className="px-4 py-2 text-center">{new Date(notification.created_at).toLocaleString()}</td>
                                    <td className="px-4 py-2 text-center">{notification.api_status}</td>
                                    <td className="px-4 py-2 text-center">
                                        <Button onClick={() => handleNotificationClick(notification)}>
                                            View Response
                                        </Button>
                                    </td>
                                    <td className="px-4 py-2 text-center">
                                        <Button
                                            onClick={() => handleCopyId(notification.id)}
                                            variant="outline"
                                            size="icon"
                                            title="Copy Notification ID"
                                        >
                                            <Copy className="h-4 w-4" />
                                        </Button>
                                    </td>

                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {selectedNotification && (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
                    <h2 className="text-xl font-semibold mb-4">API Response</h2>
                    <ReactJson
                        src={JSON.parse(selectedNotification.api_response)}
                        theme="monokai"
                        displayDataTypes={false}
                        enableClipboard={false}
                    />
                </div>
            )}

            <Link href="/campaigns/history" passHref>
                <Button className="mt-6">
                    Back to Campaign History
                </Button>
            </Link>
        </div>
    );
}