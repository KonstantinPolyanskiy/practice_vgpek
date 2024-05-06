package issued_practice

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
	"time"
)

func (h Handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10000*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.UploadIssuedPracticeOperation),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	// Максимальный размер файла - 10 мб
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		l.Warn("попытка загрузить большой файл", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.UploadIssuedPracticeOperation,
			Error:  "Слишком большой файл",
		})
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		l.Warn("ошибка чтения файла из формы", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.UploadIssuedPracticeOperation,
			Error:  "Ошибка чтения файла",
		})
		return
	}
	defer file.Close()

	req := dto.NewIssuedPracticeReq{
		TargetGroups: r.MultipartForm.Value["target_groups"],
		Title:        r.FormValue("title"),
		Theme:        r.FormValue("theme"),
		Major:        r.FormValue("major"),
		File:         &file,
	}

	l.Info("попытка загрузить практическое задание",
		zap.String("тема", req.Theme),
		zap.Strings("целевые группы", req.TargetGroups),
	)

	practice, err := h.s.Save(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  err.Error(),
			})
			return
		}
	}
	l.Info("практическое задание успешно загружено")

	render.JSON(w, r, rest.IssuedPractice{}.DomainToResponse(practice))
	return
}
