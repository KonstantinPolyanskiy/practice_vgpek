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

func (h AccessHandler) AddAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingAction permissions.AddActionReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", "добавление действия"),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingAction)
	if err != nil {
		l.Warn("error parse new action request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			//TODO: заменить action на константу
			Action: "добавление действия",
			Error:  "Преобразование запроса на добавление действия",
		})
		return
	}

	added, err := h.s.NewAction(ctx, addingAction)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: "добавление действия",
				Error:  "Таймаут",
			})
			return
		} else {
			l.Warn("error add action", zap.String("action name", addingAction.Name))

			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: "добавление действия",
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("action successfully added", zap.String("action name", added.Name))

	render.JSON(w, r, &added)
	return
}
