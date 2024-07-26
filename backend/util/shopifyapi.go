package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
    "log"
)

type Customer struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}

type CustomerResponse struct {
    Customer Customer `json:"customer"`
}

func GetCustomer(customerID string, accessToken string, shopName string) (Customer, error) {
    // curl -X GET "https://your-development-store.myshopify.com/admin/api/2024-01/customers/207119551.json" \
    // -H "X-Shopify-Access-Token: {access_token}"

    url := fmt.Sprintf("https://%s/admin/api/2024-01/customers/%s.json", shopName, customerID)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return Customer{}, err
    }

    log.Printf("[DEBUG] Fetching customer data from: %s", url)

    log.Printf("[DEBUG] Access token: %s", accessToken)

    req.Header.Set("X-Shopify-Access-Token", accessToken)
    req.Header.Set("Host", fmt.Sprintf("%s", shopName))
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return Customer{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return Customer{}, fmt.Errorf("failed to fetch customer data: %s", resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return Customer{}, err
    }

    var customerResponse CustomerResponse
    if err := json.Unmarshal(body, &customerResponse); err != nil {
        return Customer{}, err
    }

    return customerResponse.Customer, nil
}