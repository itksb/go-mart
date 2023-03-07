package handler

import (
	"encoding/json"
	"fmt"
	"github.com/itksb/go-mart/internal/middleware"
	"net/http"
)

func (h *Handler) APIGetUserSum(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	summary, err := h.balanceService.GetSummaryForUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("error getting user balance summary, %s", err)
		return
	}

	bSummary, err := json.Marshal(summary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("error marshaller user balance summary", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bSummary))

}
