package solved_practice

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
	"strings"
	"time"
)

func (h Handler) PracticeById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetSolvedPracticeInfoById),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("ошибка получения параметров запроса", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetSolvedPracticeInfoById,
			Error:  "Преобразование запроса на получение практической работы",
		})
		return
	}

	practice, err := h.s.ById(ctx, dto.EntityId{Id: id})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetIssuedPracticeInfoById,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetIssuedPracticeInfoById,
				Error:  err.Error(),
			})
			return
		}
	}

	link := r.Host + r.URL.String()
	link = strings.Replace(link, "?", "/download?", 1)

	render.JSON(w, r, rest.SolvedPractice{}.DomainToResponse(practice).WithDownloadLink(link))
	return
}
