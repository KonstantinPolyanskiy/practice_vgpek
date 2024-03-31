package reg_key

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/internal/service/reg_key"
	"practice_vgpek/pkg/apperr"
	"time"
)

// AddKey REST хэндлер для создания ключа регистрации
func (h Handler) AddKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	var addingKey registration_key.AddReq

	err := json.NewDecoder(r.Body).Decode(&addingKey)
	if err != nil {
		//TODO: логгирование
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: "Добавление ключа",
			Error:  "Преобразование запроса",
		})
		return
	}

	createdKey, err := h.s.NewKey(ctx, addingKey)
	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: "Добавление ключа",
			Error:  "Таймаут",
		})
		return
	} else if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, reg_key.ErrDontHavePermission) {
			status = http.StatusForbidden
		}

		apperr.New(w, r, status, apperr.AppError{
			Action: "Добавление ключа",
			Error:  err.Error(),
		})
		return
	}

	render.JSON(w, r, createdKey)
	return
}
