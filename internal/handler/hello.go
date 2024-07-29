package handler

import (
	"log"
	"my-go-template/cmd/web"
	"my-go-template/internal/database"
	"net/http"
)

type HelloWebHandler struct {
	db database.Service
}

func NewHelloWebHandler(db database.Service) http.HandlerFunc {
	h := &HelloWebHandler{db: db}
	return h.ServeHTTP
}

func (h *HelloWebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
