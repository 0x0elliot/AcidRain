import React, { useEffect } from "react";
import { useState } from "react";
import { Button } from "@/components/ui/button"

import cookies from 'nookies';
import { siteConfig } from "@/app/siteConfig";

export default function QuickWebPushCampaign() {
    const [segmentFinalized, setSegmentFinalized] = useState(false);
    const [segmentSelectionModalOpen, setSegmentSelectionModalOpen] = useState(false);
    const [webPushSegments, setWebPushSegments] = useState([]);

    const [userinfo, setUserinfo] = useState({});

    const [accessToken, setAccessToken] = useState(null);

    useEffect(() => {
        const token  = cookies.get(null).access_token;
        setAccessToken(token);

        const userInfoString = localStorage.getItem("userinfo") || "{}";
        try {
            // setUserinfo(JSON.parse(userInfoString));
            json = JSON.parse(userInfoString);
            if (json.id) {
                setUserinfo(json);
            }
        } catch (error) {
            setUserinfo({});
        }
    }, []);

    useEffect(() => {
        if (segmentSelectionModalOpen) {
            getWebPushSegments();
        }
    }, [segmentSelectionModalOpen]);

    const getWebPushSegments = async () => {
        // check local storage, userinfo. if it has a "CurrentShop" key, then use that shop identifier
        // if not, then ask API again.
        let localUserinfo = userinfo

        if (!localUserinfo.current_shop_id) {
            const userInfoResponse = await fetch(`${siteConfig.baseApiUrl}/api/user/private/getinfo`, {
                method: 'GET',
                headers: {
                    "Authorization": `Bearer ${accessToken}`,
                },
            });

            const userInfoData = await userInfoResponse.json();
            if (userInfoResponse.status !== 200) {
                console.error(userInfoData.message);
                return;
            }

            localStorage.setItem('userinfo', JSON.stringify(userInfoData));
            setUserinfo(userInfoData);
            localUserinfo = userInfoData;
        }

        if (!localUserinfo.current_shop_id) {
            console.error("No shop identifier found. Please select a shop first.");
            return;
        }

        const response = await fetch(`${siteConfig.baseApiUrl}/api/notification/private/push-subscribers?shop_identifier=${localUserinfo.current_shop.shop_identifier}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${accessToken}`
            },
        });

        const data = await response.json();
        if (response.status !== 200) {
            console.error(data.message);
            return;
        }

        setWebPushSegments(data.subscriptions);
    }

    return (
        <div className="flex bg-white dark:bg-gray-900">
            <div className="flex-1 p-4">
                <h1 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Launch a quick campaign</h1>
                <p className="text-sm text-gray-600 dark:text-gray-400">Select from your segment</p>
                <div className="flex justify-center p-4 mt-10 space-x-4" style={{ marginTop: "2rem" }}>
                    <Button
                        className="text-white bg-black hover:bg-gray-800 focus:ring-2 focus:ring-gray-500 font-medium text-sm px-5 py-2.5 text-center dark:bg-white dark:text-black dark:hover:bg-gray-200 dark:focus:ring-gray-400 disabled:opacity-50"
                        onClick={() => segmentSelectionModalOpen ? setSegmentSelectionModalOpen(false) : setSegmentSelectionModalOpen(true)}
                    >
                        Select a segment
                    </Button>

                    <Button
                        className="text-white bg-black hover:bg-gray-800 focus:ring-2 focus:ring-gray-500 font-medium text-sm px-5 py-2.5 text-center dark:bg-white dark:text-black dark:hover:bg-gray-200 dark:focus:ring-gray-400 disabled:opacity-50"
                        disabled={!segmentFinalized}
                        >
                        Launch
                    </Button>
                </div>
            </div>

            {(segmentSelectionModalOpen) && (
                <div className="fixed inset-0 bg-opacity-50 z-50 flex items-center justify-center">
                    <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-lg">
                        <h1 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Select a segment</h1>
                        <p className="text-sm text-gray-600 dark:text-gray-400">Select a segment to send the campaign to</p>
                        <div className="flex flex-col items-center overflow-y-auto p-4 rounded-lg">
                        {/* Iterate over all the subscribers and show them here as a checklist */}
                            {webPushSegments.map((segment) => (
                                // this list should be vertically scrollable and centered and pretty. 
                                // after each element, new line should start
                                <div key={segment.id} className="flex items-center space-x-4 mb-2">
                                    <input type="checkbox" id={segment.id} name={segment.id} value={segment.id} />
                                    <label htmlFor={segment.id}>{segment.owner_id ? "Your test subscription" : segment.customer_ids }</label>
                                </div>
                            ))}
                        </div>

                        <div className="flex justify-center">
                            <Button
                                className="text-white bg-black hover:bg-gray-800 focus:ring-2 focus:ring-gray-500 font-medium text-sm px-5 py-2.5 text-center dark:bg-white dark:text-black dark:hover:bg-gray-200 dark:focus:ring-gray-400 disabled:opacity-50"
                                onClick={() => setSegmentSelectionModalOpen(false)}
                            >
                                Select
                            </Button>
                        </div>
                    </div>
                </div>
            )}


        </div>



    )
}
