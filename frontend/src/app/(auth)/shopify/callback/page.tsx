"use client"

// generate a page that takes the params and sends it to localhost:3000/api/shopify/callback as a GET request

import { useEffect } from "react";
import { siteConfig } from "@/app/siteConfig";
import { setCookie } from "nookies";

export default function ShopifyCallback() {
    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);

        // take all the params and send them to the backend
        fetch(`${siteConfig.baseApiUrl}/api/user/shopify/callback?${urlParams.toString()}`)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.error(data.error);
                    return;
                } else {
                    // Set the cookie with the token
                    window.location.href = "/overview"
                }
            });
    }, []);

    return null;
}
