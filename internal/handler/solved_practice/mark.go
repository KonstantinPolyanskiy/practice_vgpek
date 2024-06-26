package solved_practice

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

func (h Handler) SetMark(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.SetMarkSolvedPractice),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	var markBody dto.MarkPracticeReq
	err := json.NewDecoder(r.Body).Decode(&markBody)
	if err != nil {
		l.Warn(operation.DecodeError, zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.SetMarkSolvedPractice,
			Error:  "Преобразование запроса",
		})
		return
	}

	l.Info("попытка выставить оценку за практическую работу",
		zap.Int("id работы", markBody.SolvedPracticeId),
		zap.Int("оценка", markBody.Mark),
	)

	practice, err := h.s.SetMark(ctx, markBody)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.SetMarkSolvedPractice,
				Error:  "Таймаут",
			})
			return
		} else {
			status := http.StatusInternalServerError

			apperr.New(w, r, status, apperr.AppError{
				Action: operation.SetMarkSolvedPractice,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("оценка успешно выставлена", zap.Int("id работы", practice.Mark))

	render.JSON(w, r, rest.SolvedPractice{}.DomainToResponse(practice))
	return
}
