package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(h.db.Health())
	_, _ = w.Write(jsonResp)
}
