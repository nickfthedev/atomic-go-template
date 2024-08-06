package health

import (
	"encoding/json"
	"fmt"
	"my-go-template/internal/config"
	"my-go-template/internal/database"
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
