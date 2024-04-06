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

	added, err := h.s.NewPermission(ctx, addingPerm)
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

	l.Info("permission successfully added")

	render.JSON(w, r, &added)
	return
}
