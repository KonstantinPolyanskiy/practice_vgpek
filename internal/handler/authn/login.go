package authn

import (
	"context"
	"encoding/json"
	"errors"
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

	var cred dto.Credentials

	l := h.logger.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.LoginOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.LoginOperation,
			Error:  "Преобразование запроса на вход",
		})
		return
	}

	l.Info("попытка входа",
		zap.String("логин", cred.Login),
		zap.String("пароль", cred.Password),
	)

	token, err := h.tokenService.CreateToken(ctx, cred)
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

	l.Info("пользователь успешно вошел", zap.String("логин", cred.Login))

	render.JSON(w, r, rest.Token{}.TokenToResponse(token))
	return
}
