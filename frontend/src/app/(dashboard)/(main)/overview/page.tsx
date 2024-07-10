"use client"

import { siteConfig } from "@/app/siteConfig";
import React from "react"
import { useEffect } from "react"
import cookies from 'nookies';

export default function Overview() {
  const [shopExists, setShopExists] = React.useState(false);

  const [showPopup, setShowPopup] = React.useState(false);
  const [shopName, setShopName] = React.useState("");

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
          WhatsApp Automations
        </h1>
      </section>

      {/* Have a button that says "Add a shop" that convinces the user to add a shop if they haven't already */}
      {/* The button should be convincing, and should be the first thing the user sees */}
      {/* Clicking on it should open a popup text box that asks for the shop name */}
      {(!shopExists && !showPopup) && (
      <div className="flex items-center justify-center h-screen">
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        {/* <div className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center"> */}
          <div className="relative p-4 bg-white border border-gray-200 rounded-lg shadow-sm dark:bg-gray-800">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">
                Add a shop
              </h3>
            </div>
            <div className="mt-2">
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Add a shop to get started with WhatsApp automations.
              </p>
            </div>
            <div className="mt-4">
              <button
                type="button"
                className="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                onClick={() => setShowPopup(true)}
              >
                Add a shop
              </button>
            </div>
            </div>
        </div>
      </div>
      )}
      
      {showPopup && (
        <div className="flex items-center justify-center h-screen">
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div className="absolute top-0 left-0 right-0 bottom-0 flex items-center justify-center">
            <div className="p-6 rounded-lg shadow-md">
              <h3 className="mb-4 text-lg font-semibold text-gray-900 dark:text-gray-50">
                Enter Shop Name
              </h3>
              <input
                type="text"
                className="w-full px-3 py-2 mb-4 border rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                value={shopName}
                onChange={(e) => setShopName(e.target.value)}
                placeholder="Shop Name"
              />
              <button
                type="button"
                className="inline-flex items-center px-4 py-2 mr-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
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
      )}
{/*   
      {shopExists ? null : (
        <button
          type="button"
          className="inline-flex items-center px-4 py-2 mt-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          onClick={() => {
            window.location.href = siteConfig.baseApiUrl + "/api/user/shopify-oauth";
          }}
        >
          Add a shop
        </button>
      )} */}

      

    </>
  )
}
