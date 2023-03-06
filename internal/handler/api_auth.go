package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itksb/go-mart/api"
	"github.com/itksb/go-mart/internal/service/auth"
	"io"
	"net/http"
)

// APIAuthRegister - регистрация пользователя.
// Возможные коды ответа:
// 200 — пользователь успешно зарегистрирован и аутентифицирован;
// 400 — неверный формат запроса;
// 409 — логин уже занят;
// 500 — внутренняя ошибка сервера.
func (h *Handler) APIAuthRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.apiSendJSONError(
			w,
			api.ErrorBadRequest{Error: "error while reading the request body"},
			http.StatusInternalServerError,
		)
		h.logger.Errorf("APIAuthRegister. Error while reading the request body %s", err.Error())
		return
	}

	request := api.SignUPRequest{
		Login:    "",
		Password: "",
	}

	if err := json.Unmarshal(reqBytes, &request); err != nil {
		h.apiSendJSONError(
			w,
			api.ErrorBadRequest{Error: "bad request type: json error."},
			http.StatusBadRequest,
		)
		h.logger.Errorf("APIAuthRegister. Wrong input. Json decoding error: %s", err.Error())
		return
	}

	if request.Login == "" || request.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tok, err := h.auth.SignUp(r.Context(), auth.ClientCredential{
		Login:    request.Login,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, auth.ErrDuplicateKeyValue) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		h.logger.Errorf("auth.SignUp error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tok))
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) APIAuthLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.apiSendJSONError(
			w,
			api.ErrorBadRequest{Error: "error while reading the request body"},
			http.StatusInternalServerError,
		)
		h.logger.Errorf("APIAuthRegister. Error while reading the request body %s", err.Error())
		return
	}

	request := api.AuthRequest{
		Login:    "",
		Password: "",
	}

	if err := json.Unmarshal(reqBytes, &request); err != nil {
		h.apiSendJSONError(
			w,
			api.ErrorBadRequest{Error: "bad request type: json error."},
			http.StatusBadRequest,
		)
		h.logger.Errorf("APIAuthRegister. Wrong input. Json decoding error: %s", err.Error())
		return
	}

	if request.Login == "" || request.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tok, err := h.auth.SignIn(r.Context(), auth.ClientCredential{
		Login:    request.Login,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.logger.Errorf("auth.SignIn error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tok))
	w.WriteHeader(http.StatusOK)

}
