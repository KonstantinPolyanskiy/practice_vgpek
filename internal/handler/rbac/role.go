package rbac

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/pkg/apperr"
	"time"
)

func (h AccessHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingRole permissions.AddRoleReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", "добавление роли"),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingRole)
	if err != nil {
		l.Warn("error parse new role request", zap.String("decoder error", err.Error()))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: "добавление роли",
			Error:  "Преобразование запроса на создание роли",
		})
		return
	}

	savedRole, err := h.s.NewRole(ctx, addingRole)
	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: "добавление роли",
			Error:  "Таймаут",
		})
		return
	} else if err != nil {
		l.Warn("error adding role", zap.Error(err))
		l.Debug("adding data", zap.String("name", addingRole.Name))

		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: "добавление роли",
			Error:  err.Error(),
		})
		return
	}

	l.Info("role successfully added", zap.String("role name", savedRole.Name))

	render.JSON(w, r, &savedRole)
	return
}
