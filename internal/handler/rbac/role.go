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

func (h AccessHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingRole permissions.AddRoleReq

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", rbac.AddRoleOperation),
		zap.String("layer", "handlers"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingRole)
	if err != nil {
		l.Warn("error parse new role request", zap.String("decoder error", err.Error()))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.AddRoleOperation,
			Error:  "Преобразование запроса на создание роли",
		})
		return
	}

	savedRole, err := h.s.NewRole(ctx, addingRole)
	if err != nil && errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: rbac.AddRoleOperation,
			Error:  "Таймаут",
		})
		return
	} else if err != nil {
		l.Warn("error adding role", zap.Error(err))
		l.Debug("adding data", zap.String("name", addingRole.Name))

		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: rbac.AddRoleOperation,
			Error:  err.Error(),
		})
		return
	}

	l.Info("role successfully added", zap.String("role name", savedRole.Name))

	render.JSON(w, r, &savedRole)
	return
}

func (h AccessHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("endpoint", r.RequestURI),
		zap.String("action", rbac.GetRoleOperation),
		zap.String("layer", "handlers"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("error parse get role request", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.GetRoleOperation,
			Error:  "Преобразование запроса на получение роли",
		})
		return
	}

	role, err := h.s.RoleById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: rbac.GetRoleOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, rbac.ErrDontHavePermission) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: rbac.GetRoleOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	render.JSON(w, r, permissions.GetRoleResp{
		Id:   role.Id,
		Name: role.Name,
	})
	return
}

func (h AccessHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: rbac.GetRolesOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	roles, err := h.s.RolesByParams(ctx, defaultParams)
	if errors.Is(err, context.DeadlineExceeded) {
		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: rbac.GetRolesOperation,
			Error:  "неправильные параметры запроса",
		})
		return
	} else if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, rbac.ErrDontHavePermission) {
			code = http.StatusForbidden
		}

		apperr.New(w, r, code, apperr.AppError{
			Action: rbac.GetRolesOperation,
			Error:  err.Error(),
		})
		return
	}

	var resp permissions.GetRolesResp
	resp.Roles = make([]permissions.GetRoleResp, 0, len(roles))

	for _, role := range roles {
		resp.Roles = append(resp.Roles, permissions.GetRoleResp{
			Id:   role.Id,
			Name: role.Name,
		})
	}

	render.JSON(w, r, resp)
	return
}
