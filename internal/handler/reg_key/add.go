package reg_key

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/transport/rest"
	"practice_vgpek/pkg/apperr"
	"time"
)

func (h Handler) AddKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	var addingKey dto.NewKeyReq

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.NewKeyOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	err := json.NewDecoder(r.Body).Decode(&addingKey)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.NewKeyOperation,
			Error:  "Преобразование запроса",
		})
		return
	}

	err = validateAddKey(addingKey)
	if err != nil {
		l.Warn(operation.ValidateError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.NewKeyOperation,
			Error:  err.Error(),
		})
		return
	}

	l.Info("попытка создать новый ключ",
		zap.String("группа", addingKey.GroupName),
		zap.Int("роль ключа", addingKey.RoleId),
		zap.Int("макс. кол-во исп-ий", addingKey.MaxCountUsages),
	)

	createdKey, err := h.s.NewKey(ctx, addingKey)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.NewKeyOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			status := http.StatusInternalServerError

			apperr.New(w, r, status, apperr.AppError{
				Action: operation.NewKeyOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("ключ успешно создан",
		zap.Int("id ключа", createdKey.Id),
		zap.String("тело ключа", createdKey.Body),
		zap.Time("время создания", createdKey.CreatedAt),
	)

	render.JSON(w, r, rest.Key{}.DomainToResponse(createdKey))
	return
}

func validateAddKey(addingKey dto.NewKeyReq) error {
	if addingKey.RoleId == 0 {
		return errors.New("roleId не может быть пустым")
	}

	if addingKey.MaxCountUsages <= 0 {
		return errors.New("maxCountUsages не может быть пустым")
	}

	if addingKey.GroupName == "" {
		return errors.New("groupName не может быть пустым")
	}

	return nil
}
