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

func (h AccessHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingRole dto.NewRBACReq

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.AddRoleOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	err := json.NewDecoder(r.Body).Decode(&addingRole)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.AddRoleOperation,
			Error:  "Преобразование запроса на создание роли",
		})
		return
	}

	role, err := h.s.NewRole(ctx, addingRole)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.AddRoleOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
				Action: operation.AddRoleOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("роль успешно добавлена", zap.String("название роли", role.Name))

	render.JSON(w, r, rest.RBACPartDomainToResponse(role))
	return
}

func (h AccessHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.SoftDeleteRoleById),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.SoftDeleteRoleById,
			Error:  "Преобразование запроса на получение действия",
		})
		return
	}

	hasAccess, err := h.accountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.RBACObject, domain.DeleteAction)
	if err != nil {
		l.Warn("ошибка проверки доступа", zap.Error(err))

		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.SoftDeleteRoleById,
			Error:  "Ошибка проверки доступа",
		})
		return
	}

	if !hasAccess {
		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.SoftDeleteRoleById,
			Error:  "Недостаточно прав",
		})
		return
	}

	deletedRole, err := h.s.DeleteRoleById(ctx, dto.EntityId{Id: id})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.SoftDeleteRoleById,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.SoftDeleteRoleById,
				Error:  err.Error(),
			})
			return
		}
	}

	render.JSON(w, r, rest.RBACPartDomainToResponse(deletedRole))
	return
}

func (h AccessHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetRoleOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetRoleOperation,
			Error:  "Преобразование запроса на получение роли",
		})
		return
	}

	role, err := h.s.RoleById(ctx, dto.EntityId{Id: id})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("роль успешно отдана", zap.Int("id роли", role.ID))

	render.JSON(w, r, rest.RBACPartDomainToResponse(role))
	return
}

func (h AccessHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetRolesOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Warn("ошибка получение параметров запроса", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetRolesOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	stateParams := queryutils.StateParams(r, defaultParams)

	roles, err := h.s.RolesByParams(ctx, stateParams)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetRolesOperation,
				Error:  "неправильные параметры запроса",
			})
			return
		} else {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetRolesOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("роли успешно отданы")

	render.JSON(w, r, rest.RBACPartsDomainToResponse(roles))
	return
}
