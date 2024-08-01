package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Config: %+v", h.config)
	jsonResp, _ := json.Marshal(h.db.Health())
	_, _ = w.Write(jsonResp)
}
