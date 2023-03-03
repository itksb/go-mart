package handler

import (
	"encoding/json"
	"github.com/itksb/go-mart/api"
	"io"
	"net/http"
)

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

	/*encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	w.Header().Add("Content-Type", "application/json")*/
	/*if errors.Is(err, shortener.ErrDuplicate) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}*/

	//response := api.ShortenResponse{Result: createShortenURL(sURLId, h.cfg.ShortBaseURL)}

}
