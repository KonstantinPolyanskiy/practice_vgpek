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
	"practice_vgpek/pkg/queryutils"
	"strconv"
	"time"
)

// @Summary		Создание объекта действия
// @Security		ApiKeyAuth
// @Tags			Объект действия
// @Description	Создает объект действия в системе
// @ID				create-object
// @Accept			json
// @Produce		json
// @Param			input	body		permissions.AddObjectReq	true	"Поля необходимые для создания объекта"
// @Success		200		{object}	permissions.AddObjectResp	"Возвращает название созданной роли"
// @Failure		default	{object}	apperr.AppError
// @Router			/object	 [post]
func (h AccessHandler) AddObject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingObject permissions.AddObjectReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.AddObjectOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingObject)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.AddObjectOperation,
			Error:  "Преобразование запроса на добавление объекта",
		})
		return
	}

	added, err := h.s.NewObject(ctx, addingObject)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.AddObjectOperation,
				Error:  "таймаут",
			})
			return
		} else {
			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: operation.AddObjectOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("объект действия успешно добавлен", zap.String("название объекта", added.Name))

	render.JSON(w, r, &added)
	return
}

// @Summary		Получение объекта действия
// @Security		ApiKeyAuth
// @Tags			Объект действия
// @Description	Получение объекта действия по id
// @ID				get-object
// @Accept			json
// @Produce		json
// @Param			id		query		int							true	"ID объекта"
// @Success		200		{object}	permissions.GetObjectResp	"Возвращает id и название объекта"
// @Failure		default	{object}	apperr.AppError
// @Router			/object	 [get]
func (h AccessHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("operation", operation.GetObjectOperation),
		zap.String("layer", "handlers"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("error parse get object request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetObjectOperation,
			Error:  "Преобразование запроса на получение действия",
		})
		return
	}

	object, err := h.s.ObjectById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetObjectOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetObjectOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("объект успешно отдан", zap.Int("id объекта", object.Id))

	render.JSON(w, r, permissions.GetObjectResp{
		Id:   object.Id,
		Name: object.Name,
	})
	return
}

// @Summary		Получение объектов по параметрам
// @Security		ApiKeyAuth
// @Tags			Объект действия
// @Description	Получение объектов
// @ID				get-objects
// @Accept			json
// @Produce		json
// @Param			limit			query		int							false	"Сколько выдать ролей"
// @Param			offset			query		int							false	"С какой позиции выдать роли"
// @Success		200				{object}	permissions.GetObjectsResp	"Возвращает id и названия объектов"
// @Failure		default			{object}	apperr.AppError
// @Router			/object/params	 [get]
func (h AccessHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("операция", r.RequestURI),
		zap.String("операция", operation.GetObjectsOperation),
		zap.String("слой", "http обработчики"),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Info("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetObjectsOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	objects, err := h.s.ObjectsByParams(ctx, defaultParams)
	if errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: operation.GetObjectsOperation,
			Error:  "Таймаут",
		})
		return
	} else if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, permissions.ErrDontHavePerm) {
			code = http.StatusForbidden
		}

		apperr.New(w, r, code, apperr.AppError{
			Action: operation.GetObjectsOperation,
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

	l.Info("объекты успешно отданы")

	render.JSON(w, r, resp)
	return
}
