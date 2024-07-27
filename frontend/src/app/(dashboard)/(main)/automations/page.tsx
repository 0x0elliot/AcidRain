"use client"

import { siteConfig } from "@/app/siteConfig";
import React from "react"
import { useEffect } from "react"
import cookies from 'nookies';

import axios from 'axios';

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export default function Automations() {
  const [shopExists, setShopExists] = React.useState(true);

  const [showPopup, setShowPopup] = React.useState(false);
  const [shopName, setShopName] = React.useState("");

  const [notificationTypes, setNotificationTypes] = React.useState(["WhatsApp", "SMS", "Email", "Push"]);

  const GetRedirectUrl = async (shopName: string) => {
    const response = await fetch(`${siteConfig.baseApiUrl}/api/user/shopify-oauth?shop=${shopName}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    const data = await response.json();
    return data;
  }

  const handleAddShop = (shopName: string) => {
    // window.location.href = `${siteConfig.baseApiUrl}/api/user/shopify-oauth?shop=${shopName}`;
    GetRedirectUrl(shopName).then((data) => {
      window.location.href = data.auth_url;
    });
  };

  // make a request to API + /api/shop/priv/all
  // get the response, check if data.shops is empty.
  // if it is, then show the "Add a shop" button
  useEffect(() => {
    let accessToken = cookies.get(null).access_token;

    // use axios to make a request to the API
    axios.get(`${siteConfig.baseApiUrl}/api/user/private/getinfo`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
      },
    }).then((response) => {
      localStorage.setItem('userinfo', JSON.stringify(response.data));
    });

    fetch(`${siteConfig.baseApiUrl}/api/shop/private/all`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
      },
    })
      .then(response => response.json())
      .then(data => {
        if (data.shops.length > 0) {
          setShopExists(true);
        } else {
          setShopExists(false);
        }
      });
  }
  , []);

  return (
    <>
      <section aria-labelledby="flows-title">
        <h1
          id="overall-title"
          className="scroll-mt-10 text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50"
        >
          Configure Automations
        </h1>
      </section>

      <section aria-labelledby="flows-description">
        <p
          id="overall-description"
          className="text-sm text-gray-500 dark:text-gray-400"
        >
          Configure automations for your shop.
        </p>
      </section>

      {/* Have a button that says "Add a shop" that convinces the user to add a shop if they haven't already */}
      {/* The button should be convincing, and should be the first thing the user sees */}
      {/* Clicking on it should open a popup text box that asks for the shop name */}
      {(!shopExists && !showPopup) && (
        <div className="relative flex items-center justify-center h-screen">
        <div className="relative p-4 bg-white border border-gray-200 rounded-lg shadow-sm dark:bg-gray-800 z-50">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
              Add a shop
            </h3>
          </div>
          <div className="mt-2">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Add a shop to get started with your automations.
            </p>
          </div>
          <div className="mt-4">
            <button
              type="button"
              className="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-orange-600 border border-transparent rounded-md shadow-sm hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500"
              onClick={() => setShowPopup(true)}
            >
              Add a shop
            </button>
          </div>
        </div>
      </div>
      

    )} 
      
      {showPopup && (
        <div className="flex items-center justify-center h-screen">
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center">
            <div className="p-6 rounded-lg shadow-md w-96"> {/* Adjust width here */}
              <h3 className="mb-4 text-lg font-semibold text-gray-900 dark:text-gray-50">
                Enter Shop Name
              </h3>

              <div className="relative flex w-full border rounded-md">
                <input
                  type="text"
                  className="w-full px-3 py-2 border-none rounded-l-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                  value={shopName}
                  onChange={(e) => setShopName(e.target.value)}
                  placeholder="Shop Name"
                />
                <div className="absolute inset-y-0 right-0 flex items-center pr-3 bg-gray-200 border-l border-gray-300 rounded-r-md">
                  <span className="text-gray-500">.myshopify.com</span>
                </div>
              </div>

            <div className="mt-4">

              <button
                type="button"
                className="inline-flex items-center px-4 py-2 mr-2 text-sm font-medium text-white bg-orange-600 border border-transparent rounded-md shadow-sm hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500"
                onClick={ () => handleAddShop(shopName) }
              >
                Submit
              </button>
              <button
                type="button"
                className="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-gray-300 border border-transparent rounded-md shadow-sm hover:bg-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
                onClick={() => setShowPopup(false)}
              >
                Cancel
              </button>
              </div>
            </div>
          </div>
        </div>
        </div>
      )} 

      {/* Add two clickable "tabs" from shadcn */}
      {/* The first tab should be "Order Automations" */}
      {/* The second tab should be "Promotion Automations" */}
      <div className="my-tabs mt-4">


      <Tabs defaultValue="push" className="w-[400px]">
        <TabsList>
          {notificationTypes.map((type) => (
            <TabsTrigger key={type} value={type.toLowerCase()}>{type}</TabsTrigger>
          ))}
        </TabsList>

        {notificationTypes.map((type) => (
          <TabsContent key={type} value={type.toLowerCase()}>
            <div className="p-4">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                {type} Automations
              </h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Configure {type} automations for your shop.
              </p>
            </div>
          </TabsContent>
        ))}


      </Tabs>
      </div>

    </>
  )
}
