"use client"

import * as React from "react"
import { TrendingUp } from "lucide-react"
import { Label, Pie, PieChart, Cell } from "recharts"

import Link from 'next/link';

import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    ChartConfig,
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart"
const chartData = [
    {
        os: "linux",
        visitors: 1000,
        fill: "#E2366F"
    },
    {
        os: "macOS",
        visitors: 2000,
        fill: "#2761D8",
    },
    {
        os: "windows",
        visitors: 3000,
        fill: "#AF56DB",
    },
    {
        os: "other",
        visitors: 4000,
        fill: "#2DB78A"
    },
]

const chartConfig = {
    visitors: {
        label: "Visitors",
    },
    linux: {
        label: "Linux",
    },
    macOS: {
        label: "Mac",
    },
    windows: {
        label: "Windows",
    },
    other: {
        label: "Other",
    },
} satisfies ChartConfig

const notificationCampaignData = [
    {
        url: "https://example.com",
        visitors: 100,
        notification_campaign_id: "1",
        fill: "#E2366F"
    },
    {
        url: "https://example.com",
        visitors: 200,
        notification_campaign_id: "2",
        fill: "#2761D8",
    },
    {
        url: "https://example.com",
        visitors: 300,
        notification_campaign_id: "3",
        fill: "#AF56DB",
    },
    {
        url: "https://example.com",
        visitors: 400,
        notification_campaign_id: "4",
        fill: "#2DB78A"
    },
    {
        url: "https://example.com",
        visitors: 500,
        notification_campaign_id: "5",
        fill: "#FF8C00"
    }
]


export function DashboardClickerOS() {
    const totalVisitors = React.useMemo(() => {
        return chartData.reduce((acc, curr) => acc + curr.visitors, 0)
    }, [])

    return (
        <Card className="flex flex-col">
            <CardHeader className="items-center pb-0">
                <CardTitle>About your clickers..</CardTitle>
                <CardDescription>All time data</CardDescription>
            </CardHeader>
            <div className="flex flex-row gap-4">
                <CardContent className="flex-1 pb-0">
                    <ChartContainer
                        config={chartConfig}
                        className="mx-auto aspect-square max-h-[250px]"
                    >
                        <PieChart>
                            <ChartTooltip
                                cursor={false}
                                content={<ChartTooltipContent hideLabel />}
                            />
                            <Pie
                                data={chartData}
                                dataKey="visitors"
                                nameKey="os"
                                innerRadius={60}
                                strokeWidth={5}
                            >
                                <Label
                                    content={({ viewBox }) => {
                                        if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                                            return (
                                                <text
                                                    x={viewBox.cx}
                                                    y={viewBox.cy}
                                                    textAnchor="middle"
                                                    dominantBaseline="middle"
                                                >
                                                    <tspan
                                                        x={viewBox.cx}
                                                        y={viewBox.cy}
                                                        className="fill-foreground text-3xl font-bold"
                                                    >
                                                        {totalVisitors.toLocaleString()}
                                                    </tspan>
                                                    <tspan
                                                        x={viewBox.cx}
                                                        y={(viewBox.cy || 0) + 24}
                                                        className="fill-muted-foreground"
                                                    >
                                                        Visitors
                                                    </tspan>
                                                </text>
                                            )
                                        }
                                    }}
                                />
                            </Pie>
                        </PieChart>
                    </ChartContainer>
                </CardContent>

                <CardContent className="flex-1 pb-0">
                    <ChartContainer
                        config={chartConfig}
                        className="mx-auto aspect-square max-h-[250px]"
                    >
                        <PieChart>
                            <ChartTooltip
                                cursor={false}
                                content={<ChartTooltipContent hideLabel />}
                            />
                            <Pie
                                data={notificationCampaignData}
                                dataKey="visitors"
                                nameKey="notification_campaign_id"
                                innerRadius={60}
                                strokeWidth={5}
                                onClick={(data) => {
                                    if (data && data.notification_campaign_id) {
                                        window.open(`/campaigns/history/${data.notification_campaign_id}/info`)
                                    }
                                }}
                            >
                                {notificationCampaignData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} style={{ cursor: 'pointer' }} />
                                ))}
                                <Label
                                    content={({ viewBox }) => {
                                        if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                                            const totalCampaignVisitors = notificationCampaignData.reduce((sum, item) => sum + item.visitors, 0);
                                            return (
                                                <text
                                                    x={viewBox.cx}
                                                    y={viewBox.cy}
                                                    textAnchor="middle"
                                                    dominantBaseline="middle"
                                                >
                                                    <tspan
                                                        x={viewBox.cx}
                                                        y={viewBox.cy}
                                                        className="fill-foreground text-3xl font-bold"
                                                    >
                                                        {totalCampaignVisitors.toLocaleString()}
                                                    </tspan>
                                                    <tspan
                                                        x={viewBox.cx}
                                                        y={(viewBox.cy || 0) + 24}
                                                        className="fill-muted-foreground"
                                                    >
                                                        Campaign Visitors
                                                    </tspan>
                                                </text>
                                            )
                                        }
                                    }}
                                />
                            </Pie>
                        </PieChart>
                    </ChartContainer>
                </CardContent>
            </div>
            <CardFooter className="flex-col gap-2 text-sm">
                <Link
                    href="https://calendly.com/aditya-zappush"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 font-medium leading-none text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-200 transition-colors"
                >
                    <div className="flex items-center gap-2">
                        Trend up more by talking to us <TrendingUp className="h-4 w-4" /> We probably can help
                    </div>
                </Link>
                <div className="leading-none text-muted-foreground">
                    Showing total visitors for the last 6 months
                </div>
            </CardFooter>
        </Card>
    )
}
