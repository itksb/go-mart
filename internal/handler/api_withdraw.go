package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itksb/go-mart/api"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/itksb/go-mart/internal/middleware"
	"github.com/itksb/go-mart/internal/service/order"
	"net/http"
)

// APIWithdraw POST
func (h *Handler) APIWithdraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawRequest := api.BalanceWithdrawRequest{}

	err := json.NewDecoder(r.Body).Decode(&withdrawRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("withdraw json decoding error %s", err.Error())
		return
	}

	_, err = h.withdrawService.Create(
		r.Context(),
		withdrawRequest.Order,
		withdrawRequest.Sum,
		userID,
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWithdrawNotEnoughBalance):
			w.WriteHeader(http.StatusPaymentRequired)
		case errors.Is(err, order.ErrOrderIncorrectOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		h.logger.Errorf("Failed to create withdraw", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) APIWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.withdrawService.FindAllByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("no user withdrawals %s", err.Error())
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	bWithdrawals, err := json.Marshal(withdrawals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("marshaller user withdrawals error, %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bWithdrawals))
}
