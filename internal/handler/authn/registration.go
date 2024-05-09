package authn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/transport/rest"
	"practice_vgpek/pkg/apperr"
	"time"
)

// @Summary		Регистрация пользователя
// @Tags			Authn
// @Description	Создает аккаунт пользователя по ключу регистации. В ключе регистрации зашита группа и роль создаваемого аккаунта.
// @ID				create-account
// @Accept			json
// @Produce		json
// @Param			registrationInput	body		dto.RegistrationReq	true	"Вся поля необходимы, кроме отчества"
// @Success		200					{object}	domain.Person		"Пользователь и его созданный аккаунт"
// @Router			/sign-in [post]
func (h Handler) Registration(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 50000000*time.Second)
	defer cancel()

	l := h.logger.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.RegistrationReq),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	var req dto.RegistrationReq

	l.Info("попытка регистрации пользователя")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.RegistrationReq,
			Error:  "Неизвестная ошибка декодирования запроса",
		})
		return
	}

	err = checkInput(req)
	if err != nil {
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.RegistrationReq,
			Error:  "Ошибка запроса: " + err.Error(),
		})
		return
	}

	user, err := h.personService.NewUser(ctx, req)
	if err != nil {
		apperr.New(w, r, http.StatusUnprocessableEntity, apperr.AppError{
			Action: operation.RegistrationReq,
			Error:  err.Error(),
		})
		return
	}

	l.Info("пользователь успешно зарегистрирован",
		zap.String("UUID пользователя", user.UUID.String()),
		zap.String("логин аккаунта", user.Account.Login),
	)

	render.JSON(w, r, rest.Person{}.DomainToResponse(user))
	return

}

// checkInput проверяет, что все необходимые поля заданны
func checkInput(req dto.RegistrationReq) error {
	var errMsg string

	switch {
	case req.Login == "":
		errMsg = "пустой логин"
	case req.Password == "":
		errMsg = "пустой пароль"
	case req.FirstName == "":
		errMsg = "пустое имя"
	case req.SecondName == "":
		errMsg = "пустая фамилия"
	case req.BodyKey == "":
		errMsg = "пустой ключ регистрации"
	}

	if errMsg == "" {
		return nil
	}

	return fmt.Errorf(errMsg)
}
