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

// @Summary		Создание ключа регистрации
// @Security		ApiKeyAuth
// @Tags			ключ регистрации
// @Description	Создает ключ регистрации
// @ID				create-key
// @Accept			json
// @Produce		json
// @Param			input	body		registration_key.AddReq		true	"Поля необходимые для создания ключа"
// @Success		200		{object}	registration_key.AddResp	"Возвращает id ключа в системе, его тело, кол-во использований и когда был создан"
// @Failure		default	{object}	apperr.AppError
// @Router			/key	 [post]
func (h Handler) AddKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Second)
	defer cancel()

	var addingKey registration_key.AddReq

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.NewKeyOperation),
		zap.String("слой", "http обработчики"),
	)

	err := json.NewDecoder(r.Body).Decode(&addingKey)
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.NewKeyOperation,
			Error:  "Преобразование запроса",
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

			if errors.Is(err, permissions.ErrDontHavePerm) {
				status = http.StatusForbidden
			}

			apperr.New(w, r, status, apperr.AppError{
				Action: operation.NewKeyOperation,
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
