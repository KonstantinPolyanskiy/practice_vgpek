package rbac

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/pkg/apperr"
	"time"
)

// @Summary		Создание ключа регистрации
// @Security ApiKeyAuth
// @Tags			ключ регистрации
// @Description	Создает ключ регистрации
// @ID				create-key
// @Accept			json
// @Produce		json
// @Param			input	body		registration_key.AddReq	true	"Поля необходимые для создания ключа"
// @Success		200		{object}	 registration_key.AddResp				"Возвращает id ключа в системе, его тело, кол-во использований и когда был создан"
// @Failure		default	{object}	apperr.AppError
// @Router			/key	 [post]
func (h AccessHandler) AddPermission(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingPerm permissions.AddPermReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", operation.AddPermissionOperation),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingPerm)
	if err != nil {
		l.Warn("error parse new perm request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.AddPermissionOperation,
			Error:  "Преобразование запроса на добавление доступа",
		})
		return
	}

	_, err = h.s.NewPermission(ctx, addingPerm)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.AddPermissionOperation,
				Error:  "таймаут",
			})
			return
		} else {
			l.Warn("error add permission", zap.Ints("adding id's", addingPerm.ActionsId))

			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: operation.AddPermissionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("доступы успешно назначены")

	render.JSON(w, r, permissions.AddPermResp{AddPermReq: addingPerm})
	return
}
