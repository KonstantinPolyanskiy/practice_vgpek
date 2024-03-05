package rbac

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/service/rbac"
	"practice_vgpek/pkg/apperr"
	"time"
)

func (h AccessHandler) AddObject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingObject permissions.AddObjectReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", rbac.AddObjectAction),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingObject)
	if err != nil {
		l.Warn("error parse new object request")

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.AddObjectAction,
			Error:  "Преобразование запроса на добавление объекта",
		})
		return
	}

	added, err := h.s.NewObject(ctx, addingObject)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: rbac.AddObjectAction,
				Error:  "таймаут",
			})
			return
		} else {
			l.Warn("error add object", zap.String("object name", addingObject.Name))

			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: rbac.AddObjectAction,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("object successfully added", zap.String("object name", added.Name))

	render.JSON(w, r, &added)
	return
}
