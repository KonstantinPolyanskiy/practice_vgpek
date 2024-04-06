package reg_key

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/pkg/apperr"
	"time"
)

func (h Handler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var deletingKey registration_key.DeleteReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.NewKeyOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&deletingKey)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.InvalidateKeyOperation,
			Error:  "Преобразование запроса",
		})
		return
	}

	l.Info("попытка инвалидировать ключ регистрации",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("id ключа", deletingKey.KeyId),
	)

	deletedKey, err := h.s.InvalidateKey(ctx, deletingKey)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.InvalidateKeyOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			status := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				status = http.StatusForbidden
			}

			apperr.New(w, r, status, apperr.AppError{
				Action: operation.InvalidateKeyOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("ключ успешно инвалидирован",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("id ключа", deletedKey.KeyId),
	)

	render.JSON(w, r, deletedKey)
	return
}
