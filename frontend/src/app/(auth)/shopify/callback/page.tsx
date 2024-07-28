"use client"

// generate a page that takes the params and sends it to localhost:3000/api/shopify/callback as a GET request

import { useEffect } from "react";
import { siteConfig } from "@/app/siteConfig";
import cookies from 'nookies';

export default function ShopifyCallback() {
    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);

        let accessToken = cookies.get(null).access_token;

        // take all the params and send them to the backend
        fetch(`${siteConfig.baseApiUrl}/api/user/private/shopify/callback?${urlParams.toString()}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${accessToken}`,
            },
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.error(data.error);
                    return;
                } else {
                    // Set the cookie with the token
                    window.location.href = "/dashboard";
                }
            });
    }, []);

    return null;
}
