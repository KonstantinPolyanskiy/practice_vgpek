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

func (h Handler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var deletingKey registration_key.DeleteReq

	err := json.NewDecoder(r.Body).Decode(&deletingKey)
	if err != nil {
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: "Удаление ключа",
			Error:  "Преобразование запроса",
		})

		return
	}

	deletedKey, err := h.s.InvalidateKey(ctx, deletingKey)
	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: "Удаление ключа",
			Error:  "Таймаут",
		})

		return
	} else if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, reg_key.ErrDontHavePermission) {
			status = http.StatusForbidden
		}

		apperr.New(w, r, status, apperr.AppError{
			Action: "Удаление ключа",
			Error:  err.Error(),
		})

		return
	}

	render.JSON(w, r, deletedKey)
	return
}
