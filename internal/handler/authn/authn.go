package authn

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/service/authn"
	"practice_vgpek/pkg/apperr"
	"time"
)

type Service interface {
	NewPerson(ctx context.Context, registering person.RegistrationReq) (person.RegisteredResp, error)
}

type Handler struct {
	l *zap.Logger
	s Service
}

func NewAuthenticationHandler(service Service, logger *zap.Logger) Handler {
	return Handler{
		l: logger,
		s: service,
	}
}

func (h Handler) Registration(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var registering person.RegistrationReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", authn.RegistrationAction),
	)

	err := json.NewDecoder(r.Body).Decode(&registering)
	if err != nil {
		l.Warn("error parse new person request",
			zap.String("decoder error", err.Error()),
		)

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: authn.RegistrationAction,
			Error:  "Преобразование запроса на регистрацию",
		})
		return
	}

	registered, err := h.s.NewPerson(ctx, registering)
	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: authn.RegistrationAction,
			Error:  "Таймаут",
		})
		return
	} else if err != nil {
		l.Warn("error registering user",
			zap.String("Ключ регистрации", registering.RegistrationKey),
		)
		l.Debug("registering data",
			zap.String("full name", registering.FirstName+" "+registering.MiddleName+" "+registering.MiddleName),
			zap.String("login", registering.Login),
		)

		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: authn.RegistrationAction,
			Error:  err.Error(),
		})
		return
	}

	l.Info("user successfully registered",
		zap.String("first name", registered.FirstName),
		zap.Time("registration time", registered.CreatedAt),
	)

	render.JSON(w, r, &registered)
	return
}
