package handler

import (
	"context"
	"fmt"
	"my-go-template/internal/model"
	"net/http"
	"os"
	"strings"

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

	// TODO: Store the token and shop name in the database, take a look at GoShopifyTest Repo for an example.
	db := h.db.GetDB()

	// Update or create the user's Shopify information
	var user model.User
	result := db.Where("shopify_shop = ?", shopName).FirstOrCreate(&user, model.User{
		Email:              shop.Email,
		Username:           shop.Name,
		ShopifyAccessToken: &token,
		ShopifyShop:        &shopName,
	})

	if result.Error != nil {
		http.Error(w, "Failed to update or create user with Shopify information", http.StatusInternalServerError)
		return
	}

	// Update the user if it already existed
	if result.RowsAffected == 0 {
		result = db.Model(&user).Updates(model.User{
			Email:              shop.Email,
			Username:           shop.Name,
			ShopifyAccessToken: &token,
			ShopifyShop:        &shop.MyshopifyDomain,
		})

		if result.Error != nil {
			http.Error(w, "Failed to update existing user with Shopify information", http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("OAuth completeted. Redirecting...")
	// https://admin.shopify.com/store/freshstoretest12122/apps/my-go-template
	redirectURL := fmt.Sprintf("https://admin.shopify.com/store/%s/apps/%s", strings.Split(shop.MyshopifyDomain, ".")[0], os.Getenv("APP_HANDLE"))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
