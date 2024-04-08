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

// @Summary		Создание роли
// @Security		ApiKeyAuth
// @Tags			Роль
// @Description	Создает роль в системе
// @ID				create-role
// @Accept			json
// @Produce		json
// @Param			input	body		permissions.AddRoleReq	true	"Поля необходимые для создания роли"
// @Success		200		{object}	permissions.AddRoleResp	"Возвращает название созданной роли"
// @Failure		default	{object}	apperr.AppError
// @Router			/role	 [post]
func (h AccessHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var addingRole permissions.AddRoleReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.AddRoleOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingRole)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.String("decoder error", err.Error()))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.AddRoleOperation,
			Error:  "Преобразование запроса на создание роли",
		})
		return
	}

	savedRole, err := h.s.NewRole(ctx, addingRole)
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

	l.Info("роль успешно добавлена", zap.String("название роли", savedRole.Name))

	render.JSON(w, r, &savedRole)
	return
}

// @Summary		Получение роли
// @Security		ApiKeyAuth
// @Tags			Роль
// @Description	Получение роли по id
// @ID				get-role
// @Accept			json
// @Produce		json
// @Param			id		query		int						true	"ID ключа"
// @Success		200		{object}	permissions.GetRoleResp	"Возвращает id и название роли"
// @Failure		default	{object}	apperr.AppError
// @Router			/key	 [get]
func (h AccessHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.GetRoleOperation),
		zap.String("слой", "http обработчики"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetRoleOperation,
			Error:  "Преобразование запроса на получение роли",
		})
		return
	}

	role, err := h.s.RoleById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("роль успешно отдана", zap.Int("id роли", role.Id))

	render.JSON(w, r, permissions.GetRoleResp{
		Id:   role.Id,
		Name: role.Name,
	})
	return
}

// @Summary		Получение ролей по параметрам
// @Security		ApiKeyAuth
// @Tags			Роль
// @Description	Получение ролей
// @ID				get-roles
// @Accept			json
// @Produce		json
// @Param			limit			query		int							false	"Сколько выдать ролей"
// @Param			offset			query		int							false	"С какой позиции выдать роли"
// @Param			id				query		int							true	"ID ключа"
// @Success		200				{object}	permissions.GetRolesResp	"Возвращает id и названия ролей"
// @Failure		default			{object}	apperr.AppError
// @Router			/role/params	 [get]
func (h AccessHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.GetRolesOperation),
		zap.String("слой", "http обработчики"),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetRolesOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	roles, err := h.s.RolesByParams(ctx, defaultParams)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetRolesOperation,
				Error:  "неправильные параметры запроса",
			})
			return
		} else {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetRolesOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	var resp permissions.GetRolesResp
	resp.Roles = make([]permissions.GetRoleResp, 0, len(roles))

	for _, role := range roles {
		resp.Roles = append(resp.Roles, permissions.GetRoleResp{
			Id:   role.Id,
			Name: role.Name,
		})
	}

	l.Info("роли успешно отданы")

	render.JSON(w, r, resp)
	return
}
