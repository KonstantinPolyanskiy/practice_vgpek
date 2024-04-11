package authn

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"log"
	"net/http"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/person"
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

// @Summary		Регистрация
// @Tags			авторизация
// @Description	Cоздает аккаунт по ключу регистрации
// @ID				create-person
// @Accept			json
// @Produce		json
// @Param			input	body		person.RegistrationReq	true	"Поля необходимые для регистрации"
// @Success		200		{object}	person.RegisteredResp	"На кого зарегестрирован аккаунт и кем"
// @Failure		default	{object}	apperr.AppError
// @Router			/registration [post]
func (h Handler) Registration(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10000*time.Second)
	defer cancel()

	var registering person.RegistrationReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.RegistrationOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&registering)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.RegistrationOperation,
			Error:  "Преобразование запроса на регистрацию",
		})
		return
	}

	l.Info("попытка регистрации пользователя",
		zap.String("имя", registering.Personality.FirstName),
		zap.String("фамилия", registering.Personality.LastName),
		zap.String("отчество", registering.Personality.MiddleName),
		zap.String("ключ регистрации", registering.RegistrationKey),
		zap.String("логин", registering.Credentials.Login),
		zap.String("пароль", registering.Credentials.Password),
	)

	registered, err := h.s.NewPerson(ctx, registering)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.RegistrationOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: operation.RegistrationOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("пользователь успешно зарегистрирован",
		zap.String("имя", registered.FirstName),
		zap.String("фамилия", registered.LastName),
		zap.String("отчество", registered.MiddleName),
		zap.Time("дата регистрации", registered.CreatedAt),
	)

	render.JSON(w, r, &registered)
	return
}

// @Summary		Авторизация
// @Tags			авторизация
// @Description	Вход в систему (возвращает jwt bearer token)
// @ID				login
// @Accept			json
// @Produce		json
// @Param			input	body		person.LogInReq		true	"Поля необходимые для авторизации"
// @Success		200		{object}	person.LogInResp	"Token для авторизации"
// @Failure		default	{object}	apperr.AppError
// @Router			/login [post]
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30000*time.Second)
	defer cancel()

	var logIn person.LogInReq
	log.Println(r.RequestURI)
	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.LoginOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&logIn)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.LoginOperation,
			Error:  "Преобразование запроса на вход",
		})
		return
	}

	l.Info("попытка входа",
		zap.String("логин", logIn.Login),
		zap.String("пароль", logIn.Password),
	)

	token, err := h.s.NewToken(ctx, logIn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.LoginOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: operation.LoginOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("пользователь успешно вошел", zap.String("логин", logIn.Login))

	render.JSON(w, r, token)
	return
}
