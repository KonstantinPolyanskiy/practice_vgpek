package reg_key

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/internal/service/reg_key"
	"practice_vgpek/pkg/apperr"
	"practice_vgpek/pkg/queryutils"
	"time"
)

func (h Handler) GetKeys(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", reg_key.GetKeysOperation),
		zap.String("слой", "http обработчики"),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Warn("ошибка получения параметров запроса", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: reg_key.GetKeysOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	keyParams := getKeyParams(r, defaultParams)

	l.Info("попытка получить ключи регистрации",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("лимит", keyParams.Limit),
		zap.Int("оффсет", keyParams.Offset),
		zap.Bool("только валидные", keyParams.IsValid),
	)

	keys, err := h.s.KeysByParams(ctx, keyParams)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: reg_key.GetKeysOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: reg_key.GetKeysOperation,
				Error:  err.Error(),
			})
			return
		}

	}

	var resp registration_key.GetKeysResp
	resp.Keys = make([]registration_key.GetKeyResp, 0, len(keys))

	for _, key := range keys {
		resp.Keys = append(resp.Keys, registration_key.GetKeyResp{
			RegKeyId:           key.RegKeyId,
			RoleId:             key.RoleId,
			Body:               key.Body,
			MaxCountUsages:     key.MaxCountUsages,
			CurrentCountUsages: key.CurrentCountUsages,
			CreatedAt:          key.CreatedAt,
			IsValid:            key.IsValid,
			InvalidationTime:   key.InvalidationTime,
		})
	}

	l.Info("ключи регистарции успешно получены")

	render.JSON(w, r, resp)
	return
}

func getKeyParams(r *http.Request, defaultParams params.Default) params.Key {
	isValid := true

	v := r.URL.Query().Get("valid")
	if v == "false" {
		isValid = false
	}

	return params.Key{
		IsValid: isValid,
		Default: defaultParams,
	}

}
