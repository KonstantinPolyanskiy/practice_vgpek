package reg_key

import (
	"context"
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

func (h Handler) GetKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetKeyByIdOperation),
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

	l.Info("попытка получить ключ регистрации", zap.Int("id", id))

	hasAccess, err := h.accountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.AccountObject, domain.GetAction)
	if err != nil {
		l.Warn("ошибка проверки доступа", zap.Error(err))

		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetKeyByIdOperation,
			Error:  "Ошибка проверки доступа",
		})
		return
	}

	if !hasAccess {
		apperr.New(w, r, http.StatusForbidden, apperr.AppError{
			Action: operation.GetKeyByIdOperation,
			Error:  "Недостаточно прав",
		})
		return
	}

	key, err := h.s.KeyById(ctx, dto.EntityId{Id: id})
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

	l.Info("ключ успешно получен", zap.String("тело", key.Body))

	render.JSON(w, r, rest.Key{}.DomainToResponse(key))
	return
}

func (h Handler) GetKeys(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetKeysOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Warn("ошибка получения параметров запроса", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetKeysOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	stateParams := queryutils.StateParams(r, defaultParams)

	l.Info("попытка получить ключи регистрации",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("лимит", stateParams.Limit),
		zap.Int("оффсет", stateParams.Offset),
		zap.String("состояние", stateParams.State),
	)

	keys, err := h.s.KeysByParams(ctx, stateParams)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetKeysOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetKeysOperation,
				Error:  err.Error(),
			})
			return
		}

	}

	l.Info("ключи регистрации успешно получены")

	render.JSON(w, r, rest.Keys{}.DomainToResponse(keys))
	return
}
