package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) apiSendJSONError(w http.ResponseWriter, errorMsg interface{}, status int) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := encoder.Encode(errorMsg); err != nil {
		h.logger.Errorf("error while encoding errorMsg: %s", err.Error())
	}
}
