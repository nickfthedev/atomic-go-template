package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"my-go-template/cmd/web"
	mw "my-go-template/internal/middleware"
	"my-go-template/internal/model"
	"net/http"
)

func (h *Handler) HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	component := web.HelloPost(name)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}

func (h *Handler) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	// Sample for getting userID from context when user is loggedin
	userID, ok := r.Context().Value(mw.UserIDKey).(string)
	if !ok {
		log.Println("UserID not found in context")
		// Handle the error appropriately, e.g., return an error response
		return
	}
	fmt.Println("UserID", userID)

	// Sample for getting user from context when user is loggedin
	user, ok := r.Context().Value(mw.UserKey).(model.User)
	if !ok {
		log.Println("User not found in context")
		// Handle the error appropriately, e.g., return an error response
		return
	}

	fmt.Printf("%+v\n", user)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
