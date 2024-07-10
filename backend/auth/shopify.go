package auth 

import (
	"fmt"
	"net/url"
	"os"

	"log"
)

var ShopifyAPIKey string = os.Getenv("ACIDRAIN_SHOPIFY_CLIENT_ID")
var ShopifyAPISecret string = os.Getenv("ACIDRAIN_SHOPIFY_CLIENT_SECRET")
var ShopifyScope string = "read_products,write_products"
var ShopifyRedirectURI string = "http://localhost:3000/api/user/shopify/callback"

func GenerateAuthURL(shopName string) string {
	log.Printf("[DEBUG] Crdentials: %s, %s, %s, %s", ShopifyAPIKey, ShopifyAPISecret, ShopifyScope, ShopifyRedirectURI)

	params := url.Values{}
	params.Add("client_id", ShopifyAPIKey)
	params.Add("scope", ShopifyScope)
	params.Add("redirect_uri", ShopifyRedirectURI)
	params.Add("state", "nonce") // Add state to prevent CSRF attacks
	params.Add("grant_options[]", "per-user") // Optional: for offline access

	return fmt.Sprintf("https://%s.myshopify.com/admin/oauth/authorize?%s", shopName, params.Encode())
}