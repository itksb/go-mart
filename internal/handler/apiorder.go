package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itksb/go-mart/api"
	"github.com/itksb/go-mart/internal/middleware"
	"github.com/itksb/go-mart/internal/service/order"
	"io"
	"net/http"
)

func (h *Handler) APIOrderLoad(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("APICreateOrder. Error while reading the request body %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := api.OrderLoadRequest{Number: string(reqBytes)}

	if len(request.Number) == 0 {
		h.logger.Infof("APICreateOrder. Bad request %v", request)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currnetUserID, ok := middleware.GetUserID(r.Context())

	if !ok {
		h.logger.Errorf("APICreateOrder. Unauthorized access!")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = h.orderService.Create(r.Context(), request.Number, currnetUserID)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrOrderIncorrectOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, order.ErrOrderAlreadyCreatedByCurUser):
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, order.ErrOrderAlreadyCreatedByAnotherUser):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		h.logger.Infof("error while creating new order: %s", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

func (h *Handler) APIGetOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.FindAllByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("Failed to find user orders, %s", err.Error())
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	bOrders, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("Json marshaller for user orders error %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bOrders))
}
