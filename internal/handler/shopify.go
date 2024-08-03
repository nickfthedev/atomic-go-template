package handler

import (
	"context"
	"fmt"
	"net/http"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

// Create an oauth-authorize url for the app and redirect to it.
// In some request handler, you probably want something like this:
func (h *Handler) MyHandler(w http.ResponseWriter, r *http.Request) {
	shopName := r.URL.Query().Get("shop")
	if shopName == "" {
		http.Error(w, "Shop name is required. Please add ?shop=yourshop (without myshopify.com) to the URL", http.StatusBadRequest)
		return
	}
	state := "nonce"
	authUrl, err := h.shopifyApp.AuthorizeUrl(shopName, state)
	if err != nil {
		http.Error(w, "Failed to generate authorization URL", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, authUrl, http.StatusFound)
}

// Fetch a permanent access token in the callback
func (h *Handler) MyCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check that the callback signature is valid
	if ok, _ := h.shopifyApp.VerifyAuthorizationURL(r.URL); !ok {
		http.Error(w, "Invalid Signature", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	shopName := query.Get("shop")
	code := query.Get("code")
	token, err := h.shopifyApp.GetAccessToken(context.Background(), shopName, code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	fmt.Println(token)

	// Create a new API client
	client, err := goshopify.NewClient(*h.shopifyApp, shopName, token)
	if err != nil {
		http.Error(w, "Failed to create API client", http.StatusInternalServerError)
		return
	}
	shop, _ := client.Shop.Get(context.Background(), nil)
	fmt.Printf("%+v \n \n MyShopifyDomain: %s \n", shop, shop.MyshopifyDomain)

	// Do something with the token, like store it in a DB.
}
