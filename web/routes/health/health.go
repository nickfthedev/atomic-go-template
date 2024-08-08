package health

import (
	"atomic-go-template/internal/config"
	"atomic-go-template/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
)

type Health struct {
	db     database.Service
	config *config.Config
}

func New(db database.Service, config *config.Config) *Health {
	return &Health{db: db, config: config}
}

func (h *Health) GET(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Config: %+v", h.config)
	jsonResp, _ := json.Marshal(h.db.Health())

	_, _ = w.Write(jsonResp)
}
