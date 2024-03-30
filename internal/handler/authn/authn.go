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
	NewToken(ctx context.Context, logIn person.LogInReq) (person.LogInResp, error)
	ParseToken(token string) (int, error)
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
	ctx, cancel := context.WithTimeout(r.Context(), 10000*time.Second)
	defer cancel()

	var registering person.RegistrationReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", authn.RegistrationOperation),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&registering)
	if err != nil {
		l.Warn("error parse new person request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: authn.RegistrationOperation,
			Error:  "Преобразование запроса на регистрацию",
		})
		return
	}

	registered, err := h.s.NewPerson(ctx, registering)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: authn.RegistrationOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			l.Warn("error registering user", zap.String("Ключ регистрации", registering.RegistrationKey))

			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: authn.RegistrationOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("user successfully registered",
		zap.String("first name", registered.FirstName),
		zap.Time("registration time", registered.CreatedAt),
	)

	render.JSON(w, r, &registered)
	return
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var logIn person.LogInReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", authn.LoginOperation),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&logIn)
	if err != nil {
		l.Warn("error parse login request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: authn.LoginOperation,
			Error:  "Преобразование запроса на вход",
		})
		return
	}

	token, err := h.s.NewToken(ctx, logIn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: authn.LoginOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			l.Warn("error login user", zap.String("login", logIn.Login))

			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: authn.LoginOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	render.JSON(w, r, token)
}
