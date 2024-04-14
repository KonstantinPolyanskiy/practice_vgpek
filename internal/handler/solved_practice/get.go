package solved_practice

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/solved"
	"practice_vgpek/pkg/apperr"
	"strconv"
	"time"
)

func (h Handler) PracticeById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.GetSolvedPracticeInfoById),
		zap.String("слой", "http обработчики"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetSolvedPracticeInfoById,
			Error:  "Преобразование запроса на получение практической работы",
		})
		return
	}

	practice, err := h.s.ById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetIssuedPracticeInfoById,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetIssuedPracticeInfoById,
				Error:  err.Error(),
			})
			return
		}
	}

	render.JSON(w, r, solved.GetPracticeResp{
		IssuedPracticeId: practice.IssuedPracticeId,
		SolvedPracticeId: practice.SolvedPracticeId,
		SolvedTime:       *practice.SolvedTime,
		Mark:             practice.Mark,
		MarkTime:         *practice.MarkTime,
	})
	return
}
