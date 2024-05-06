package rbac

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
	"practice_vgpek/pkg/apperr"
	"time"
)

// @Summary		Назначение доступов
// @Security		ApiKeyAuth
// @Tags			доступы
// @Description	Назначает права действия
// @ID				add-perm
// @Accept			json
// @Produce		json
// @Param			input			body		permissions.AddPermReq	true	"Поля назначении у роли к объекту действий"
// @Success		200				{object}	permissions.AddPermResp	"Возвращает id роли, id объекта и id действий, к ним добавленные"
// @Failure		default			{object}	apperr.AppError
// @Router			/permissions	 [post]
func (h AccessHandler) AddPermission(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingPerm dto.SetPermissionReq

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.AddPermissionOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	err := json.NewDecoder(r.Body).Decode(&addingPerm)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.AddPermissionOperation,
			Error:  "Преобразование запроса на добавление доступа",
		})
		return
	}

	err = h.s.NewPermission(ctx, addingPerm)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.AddPermissionOperation,
				Error:  "таймаут",
			})
			return
		} else {
			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: operation.AddPermissionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("доступы успешно назначены")

	render.JSON(w, r, map[string]string{"result": "ok"})
	return
}
