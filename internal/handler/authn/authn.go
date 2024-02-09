package authn

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/pkg/apperr"
	"time"
)

type Service interface {
	NewPerson(ctx context.Context, registering person.RegistrationReq) (person.RegisteredResp, error)
}

type Handler struct {
	s Service
}

func NewAuthenticationHandler(service Service) Handler {
	return Handler{
		s: service,
	}
}

func (h Handler) Registration(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var registering person.RegistrationReq

	err := json.NewDecoder(r.Body).Decode(&registering)
	if err != nil {
		log.Printf("Ошибка в unmarshall - %s\n", err)
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: "Регистрация пользователя",
			Error:  "Преобразование запроса",
		})
		return
	}

	registered, err := h.s.NewPerson(ctx, registering)
	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		log.Printf("Ошибка в регистрации - %s\n", err)
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: "Регистрация пользователя",
			Error:  "Таймаут",
		})
		return
	} else if err != nil {
		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: "Регистрация пользователя",
			Error:  err.Error(),
		})
	}

	render.JSON(w, r, &registered)
}
