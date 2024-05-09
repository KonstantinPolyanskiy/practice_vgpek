package user

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

func (h Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 300000*time.Second)
	defer cancel()

	var req dto.EntityId

	l := h.logger.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetAccountOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("ошибка получения параметров запроса", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Преобразование запроса на получение аккаунта",
		})
		return
	}

	req.Id = id

	hasAccess, err := h.AccountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.AccountObject, domain.GetAction)
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

	account, err := h.AccountService.EntityAccountById(ctx, req)
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

	l.Info("аккаунт успешно отдан", zap.Int("id аккаунта", account.Id))

	render.JSON(w, r, rest.AccountEntity{}.EntityToResponse(account))
	return
}

func (h Handler) GetAccountsByParam(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.logger.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetAccountsByParamsOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Warn("ошибка получения параметров запроса", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetAccountOperation,
			Error:  "Преобразование запроса на получение аккаунта",
		})
		return
	}

	stateParams := queryutils.StateParams(r, defaultParams)

	hasAccess, err := h.AccountMediator.HasAccess(ctx, ctx.Value("AccountId").(int), domain.AccountObject, domain.GetAction)
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

	accounts, err := h.AccountService.EntityAccountByParam(ctx, stateParams)
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

	l.Info("аккаунт успешно отдан", zap.Int("кол-во аккаунтов", len(accounts)))

	render.JSON(w, r, rest.AccountsEntity{}.EntityToResponse(accounts))
	return
}
