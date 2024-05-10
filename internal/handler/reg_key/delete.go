package reg_key

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/transport/rest"
	"practice_vgpek/pkg/apperr"
	"strconv"
	"time"
)

func (h Handler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.NewKeyOperation),
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

	l.Info("попытка инвалидировать ключ регистрации",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("id ключа", id),
	)

	deletedKey, err := h.s.InvalidateKey(ctx, dto.EntityId{Id: id})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.InvalidateKeyOperation,
				Error:  "Таймаут",
			})
			return
		} else {
			status := http.StatusInternalServerError

			apperr.New(w, r, status, apperr.AppError{
				Action: operation.InvalidateKeyOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("ключ успешно инвалидирован",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("id ключа", deletedKey.Id),
	)

	render.JSON(w, r, rest.InvalidatedKey{}.DomainToResponse(deletedKey))
	return
}
