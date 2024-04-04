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
	"practice_vgpek/pkg/queryutils"
	"strconv"
	"time"
)

func (h AccessHandler) AddObject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingObject permissions.AddObjectReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", rbac.AddObjectOperation),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingObject)
	if err != nil {
		l.Warn("error parse new object request")

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.AddObjectOperation,
			Error:  "Преобразование запроса на добавление объекта",
		})
		return
	}

	added, err := h.s.NewObject(ctx, addingObject)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: rbac.AddObjectOperation,
				Error:  "таймаут",
			})
			return
		} else {
			l.Warn("error add object", zap.String("object name", addingObject.Name))

			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: rbac.AddObjectOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("object successfully added", zap.String("object name", added.Name))

	render.JSON(w, r, &added)
	return
}

func (h AccessHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("operation", rbac.GetObjectOperation),
		zap.String("layer", "handlers"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("error parse get object request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.GetObjectOperation,
			Error:  "Преобразование запроса на получение действия",
		})
		return
	}

	object, err := h.s.ObjectById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: rbac.GetObjectOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, rbac.ErrDontHavePermission) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: rbac.GetObjectOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	render.JSON(w, r, permissions.GetObjectResp{
		Id:   object.Id,
		Name: object.Name,
	})
	return
}

func (h AccessHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.GetObjectsOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	objects, err := h.s.ObjectsByParams(ctx, defaultParams)
	if errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: rbac.GetObjectsOperation,
			Error:  "Таймаут",
		})
		return
	} else if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, rbac.ErrDontHavePermission) {
			code = http.StatusForbidden
		}

		apperr.New(w, r, code, apperr.AppError{
			Action: rbac.GetObjectsOperation,
			Error:  err.Error(),
		})
		return
	}

	var resp permissions.GetObjectsResp
	resp.Objects = make([]permissions.GetObjectResp, 0, len(objects))

	for _, object := range objects {
		resp.Objects = append(resp.Objects, permissions.GetObjectResp{
			Id:   object.Id,
			Name: object.Name,
		})
	}

	render.JSON(w, r, resp)
	return
}
