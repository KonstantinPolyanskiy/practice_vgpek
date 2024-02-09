package reg_key

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/pkg/apperr"
	"time"
)

type Service interface {
	NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.Entity, error)
	GetKeyById(ctx context.Context)
	GetKeyByRoleId(ctx context.Context)
	InvalidateKey(ctx context.Context)
}

type Handler struct {
	s Service
}

func NewRegKeyHandler(service Service) Handler {
	return Handler{
		s: service,
	}
}

func (h Handler) AddKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
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
		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: "Добавление ключа",
			Error:  err.Error(),
		})
		return
	}

	render.JSON(w, r, createdKey)
}
