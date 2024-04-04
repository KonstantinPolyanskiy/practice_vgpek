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

func (h AccessHandler) AddAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingAction permissions.AddActionReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", rbac.AddActionOperation),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingAction)
	if err != nil {
		l.Warn("error parse new action request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			//TODO: заменить action на константу
			Action: rbac.AddActionOperation,
			Error:  "Преобразование запроса на добавление действия",
		})
		return
	}

	added, err := h.s.NewAction(ctx, addingAction)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: rbac.AddActionOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			l.Warn("error add action", zap.String("action name", addingAction.Name))

			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: rbac.AddActionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("action successfully added", zap.String("action name", added.Name))

	render.JSON(w, r, &added)
	return
}

func (h AccessHandler) GetAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	var req permissions.GetActionReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", rbac.GetActionOperation),
		zap.String("layer", "handlers"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("error parse get action request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.GetActionOperation,
			Error:  "Преобразование запроса на получение действия",
		})
		return
	}

	req.Id = id

	action, err := h.s.ActionById(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: rbac.GetActionOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, rbac.ErrDontHavePermission) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: rbac.GetActionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	render.JSON(w, r, permissions.GetActionResp{
		Id:   action.Id,
		Name: action.Name,
	})
	return
}

func (h AccessHandler) GetActions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.GetActionsOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	actions, err := h.s.ActionsByParams(ctx, defaultParams)
	if errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: rbac.GetActionsOperation,
			Error:  "неправильные параметры запроса",
		})
		return
	} else if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, rbac.ErrDontHavePermission) {
			code = http.StatusForbidden
		}

		apperr.New(w, r, code, apperr.AppError{
			Action: rbac.GetActionsOperation,
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
