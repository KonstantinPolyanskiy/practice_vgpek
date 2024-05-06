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
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/transport/rest"
	"practice_vgpek/pkg/apperr"
	"strconv"
	"time"
)

func (h Handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.UploadSolvedPracticeOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	// Максимальный размер файла - 10 мб
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		l.Warn("попытка загрузить большой файл", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.UploadSolvedPracticeOperation,
			Error:  "Слишком большой файл",
		})
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		l.Warn("ошибка чтения файла из формы", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.UploadSolvedPracticeOperation,
			Error:  "Ошибка чтения файла",
		})
		return
	}
	defer file.Close()

	issuedId, err := strconv.Atoi(r.FormValue("issued_practice_id"))
	if err != nil {
		l.Warn("некоректный id решаемой работы", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.UploadSolvedPracticeOperation,
			Error:  "Некорректный id практического задания",
		})
		return
	}

	req := dto.NewSolvedPracticeReq{
		PerformedAccountId: ctx.Value("AccountId").(int),
		IssuedPracticeId:   issuedId,
		File:               &file,
	}

	l.Info("попытка загрузить практическое задание",
		zap.Int("практическая работа", req.IssuedPracticeId),
	)

	practice, err := h.s.Save(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.UploadSolvedPracticeOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.UploadSolvedPracticeOperation,
				Error:  err.Error(),
			})
			return
		}
	}

	l.Info("практическая работа успешно загружена")

	render.JSON(w, r, rest.SolvedPractice{}.DomainToResponse(practice))
	return
}
