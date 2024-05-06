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
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/transport/rest"
	"practice_vgpek/pkg/apperr"
	"time"
)

// @Summary		удаление ключа регистрации
// @Security		ApiKeyAuth
// @Tags			ключ регистрации
// @Description	Удаляет ключ регистрации
// @ID				delete-key
// @Accept			json
// @Produce		json
// @Param			input	body		registration_key.DeleteReq	true	"Поля необходимые для создания ключа"
// @Success		200		{object}	registration_key.DeleteResp	"Возвращает id удаленного ключа"
// @Failure		default	{object}	apperr.AppError
// @Router			/key	 [delete]
func (h Handler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var deletingKey dto.EntityId

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.NewKeyOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	err := json.NewDecoder(r.Body).Decode(&deletingKey)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.InvalidateKeyOperation,
			Error:  "Преобразование запроса",
		})
		return
	}

	l.Info("попытка инвалидировать ключ регистрации",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("id ключа", deletingKey.Id),
	)

	deletedKey, err := h.s.InvalidateKey(ctx, deletingKey.Id)
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
		zap.Int("id ключа", deletedKey.Id),
	)

	render.JSON(w, r, rest.InvalidatedKey{}.DomainToResponse(deletedKey))
	return
}
