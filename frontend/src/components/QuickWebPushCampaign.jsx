import React from "react";
import { useState } from "react";
import { Button } from "@/components/ui/button"

import cookies from 'nookies';

export default function QuickWebPushCampaign() {
    const [segmentFinalized, setSegmentFinalized] = useState(false);

    return (
        <div className="flex bg-white dark:bg-gray-900">
            <div className="flex-1 p-4">
                <h1 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Launch a quick campaign</h1>
                <p className="text-sm text-gray-600 dark:text-gray-400">Select from your segment</p>
                <div className="flex justify-center p-4 mt-10 space-x-4" style={{ marginTop: "2rem" }}>
                    <Button
                        className="text-white bg-black hover:bg-gray-800 focus:ring-2 focus:ring-gray-500 font-medium text-sm px-5 py-2.5 text-center dark:bg-white dark:text-black dark:hover:bg-gray-200 dark:focus:ring-gray-400 disabled:opacity-50"
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
        </div>
    )
}
