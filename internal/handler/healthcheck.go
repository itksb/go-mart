package handler

import (
	"encoding/json"
	"net/http"
)

// HealthCheck - for monitoring stuff
func (h *Handler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
