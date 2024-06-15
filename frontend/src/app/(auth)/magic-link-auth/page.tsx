"use client"

import { useState, useEffect } from "react";
import { siteConfig } from "@/app/siteConfig";

const headerStyle = {
    fontSize: "24px",
    fontWeight: "bold",
    marginBottom: "20px",
  };

export default function MagicLinkAuth() {
    const [token, setToken] = useState("");

    const verifyToken = async (token: string) => {
        try {
            const res = await fetch(`${siteConfig.baseApiUrl}/api/user/getinfo`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
            });

            const data = await res.json();
            console.log("Data:", data);
            if (data.error) {
                alert(data.error);
            } else {
                alert("Token Verified Successfully!");
            }
        } catch (error) {
            console.error("Error:", error);
            alert("Error verifying token. Please try again.");
        }
    }

    // onload, check if "token" param from URL
    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get("token");
        if (token) {
            setToken(token);
        }
    })

    return (
        <div>
            <header style={headerStyle}>
                Verifying Auth Token...
            </header>

            <p style={{ marginBottom: "10px" }}>
                Meanwhile, Feel free to check out how the product works <a href="https://www.youtube.com/watch?v=dQw4w9WgXcQ">here</a>
            </p>
        </div>
    )
}

