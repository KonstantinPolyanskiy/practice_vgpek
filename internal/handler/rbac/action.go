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

// @Summary		Создание действия
// @Security		ApiKeyAuth
// @Tags			Действие
// @Description	Создает действие в системе
// @ID				create-action
// @Accept			json
// @Produce		json
// @Param			input	body		permissions.AddActionReq	true	"Поля необходимые для создания действия"
// @Success		200		{object}	permissions.AddObjectResp	"Возвращает название созданного действия"
// @Failure		default	{object}	apperr.AppError
// @Router			/action	 [post]
func (h AccessHandler) AddAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingAction permissions.AddActionReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.AddActionOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingAction)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.AddActionOperation,
			Error:  "Преобразование запроса на добавление действия",
		})
		return
	}

	added, err := h.s.NewAction(ctx, addingAction)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.AddActionOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: operation.AddActionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("действие успешно добавлено", zap.String("название действия", added.Name))

	render.JSON(w, r, &added)
	return
}

// @Summary		Получение действия
// @Security		ApiKeyAuth
// @Tags			Действие
// @Description	Получение объекта действия по id
// @ID				get-action
// @Accept			json
// @Produce		json
// @Param			id		query		int							true	"ID действия"
// @Success		200		{object}	permissions.GetActionResp	"Возвращает id и название действия"
// @Failure		default	{object}	apperr.AppError
// @Router			/object	 [get]
func (h AccessHandler) GetAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	var req permissions.GetActionReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.GetActionOperation),
		zap.String("слой", "http обработчики"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetActionOperation,
			Error:  "Преобразование запроса на получение действия",
		})
		return
	}

	req.Id = id

	action, err := h.s.ActionById(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetActionOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetActionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("действие успешно отдано", zap.String("название действия", action.Name))

	render.JSON(w, r, permissions.GetActionResp{
		Id:   action.Id,
		Name: action.Name,
	})
	return
}

// @Summary		Получение действий по параметрам
// @Security		ApiKeyAuth
// @Tags			Действие
// @Description	Получение действий
// @ID				get-actions
// @Accept			json
// @Produce		json
// @Param			limit			query		int							false	"Сколько выдать действия"
// @Param			offset			query		int							false	"С какой позиции выдать действия"
// @Success		200				{object}	permissions.GetActionsResp	"Возвращает id и названия действия"
// @Failure		default			{object}	apperr.AppError
// @Router			/action/params	 [get]
func (h AccessHandler) GetActions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetActionsOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	actions, err := h.s.ActionsByParams(ctx, defaultParams)
	if errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: operation.GetActionsOperation,
			Error:  "неправильные параметры запроса",
		})
		return
	} else if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, permissions.ErrDontHavePerm) {
			code = http.StatusForbidden
		}

		apperr.New(w, r, code, apperr.AppError{
			Action: operation.GetActionsOperation,
			Error:  err.Error(),
		})
		return
	}

	var resp permissions.GetActionsResp
	resp.Actions = make([]permissions.GetActionResp, 0, len(actions))

	for _, action := range actions {
		resp.Actions = append(resp.Actions, permissions.GetActionResp{
			Id:   action.Id,
			Name: action.Name,
		})
	}

	render.JSON(w, r, resp)
	return
}
