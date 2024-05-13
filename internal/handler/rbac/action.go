package rbac

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/transport/rest"
	"practice_vgpek/pkg/apperr"
	"practice_vgpek/pkg/queryutils"
	"strconv"
	"time"
)

func (h AccessHandler) AddAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingAction dto.NewRBACReq

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.AddActionOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	err := json.NewDecoder(r.Body).Decode(&addingAction)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.AddActionOperation,

			Error: "Преобразование запроса на добавление действия",
		})
		return
	}

	hasAccess, err := h.accountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.RBACObject, domain.AddAction)
	if err != nil {
		l.Warn("ошибка проверки доступа", zap.Error(err))

		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Ошибка проверки доступа",
		})
		return
	}

	if !hasAccess {
		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Недостаточно прав",
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

	render.JSON(w, r, rest.RBACPartDomainToResponse(added))
	return
}

func (h AccessHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.SoftDeleteActionById),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetActionOperation,
			Error:  "Преобразование запроса на получение действия",
		})
		return
	}

	hasAccess, err := h.accountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.RBACObject, domain.DeleteAction)
	if err != nil {
		l.Warn("ошибка проверки доступа", zap.Error(err))

		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Ошибка проверки доступа",
		})
		return
	}

	if !hasAccess {
		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Недостаточно прав",
		})
		return
	}

	deletedAction, err := h.s.DeleteActionById(ctx, dto.EntityId{Id: id})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetActionOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetActionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	render.JSON(w, r, rest.RBACPartDomainToResponse(deletedAction))
	return
}

func (h AccessHandler) GetAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	var req dto.EntityId

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetActionOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetActionOperation,
			Error:  "Преобразование запроса на получение действия",
		})
		return
	}

	hasAccess, err := h.accountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.RBACObject, domain.GetAction)
	if err != nil {
		l.Warn("ошибка проверки доступа", zap.Error(err))

		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Ошибка проверки доступа",
		})
		return
	}

	if !hasAccess {
		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Недостаточно прав",
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

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetActionOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("действие успешно отдано", zap.String("название действия", action.Name))

	render.JSON(w, r, rest.RBACPartDomainToResponse(action))
	return
}

func (h AccessHandler) GetActions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetActionsOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Warn("возникла ошибка при получении параметров", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetActionsOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	stateParams := queryutils.StateParams(r, defaultParams)

	hasAccess, err := h.accountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.RBACObject, domain.GetAction)
	if err != nil {
		l.Warn("ошибка проверки доступа", zap.Error(err))

		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Ошибка проверки доступа",
		})
		return
	}

	if !hasAccess {
		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Недостаточно прав",
		})
		return
	}

	actions, err := h.s.ActionsByParams(ctx, stateParams)
	if errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: operation.GetActionsOperation,
			Error:  "неправильные параметры запроса",
		})
		return
	} else if err != nil {
		code := http.StatusInternalServerError

		apperr.New(w, r, code, apperr.AppError{
			Action: operation.GetActionsOperation,
			Error:  err.Error(),
		})
		return
	}

	render.JSON(w, r, rest.RBACPartsDomainToResponse(actions))
	return
}
