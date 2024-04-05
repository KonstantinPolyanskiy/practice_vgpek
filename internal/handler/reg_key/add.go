package reg_key

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/internal/service/reg_key"
	"practice_vgpek/pkg/apperr"
	"time"
)

// AddKey REST хэндлер для создания нового ключа регистрации
func (h Handler) AddKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	var addingKey registration_key.AddReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", reg_key.NewKeyOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingKey)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: reg_key.NewKeyOperation,
			Error:  "Преобразование запроса",
		})
		return
	}

	l.Info("попытка создать новый ключ",
		zap.Int("роль ключа", addingKey.RoleId),
		zap.Int("макс. кол-во исп-ий", addingKey.MaxCountUsages),
	)

	createdKey, err := h.s.NewKey(ctx, addingKey)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: reg_key.NewKeyOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			status := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				status = http.StatusForbidden
			}

			apperr.New(w, r, status, apperr.AppError{
				Action: reg_key.NewKeyOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("ключ успешно создан",
		zap.Int("id ключа", createdKey.RegKeyId),
		zap.String("тело ключа", createdKey.Body),
		zap.Time("время создания", createdKey.CreatedAt),
	)

	render.JSON(w, r, createdKey)
	return
}
