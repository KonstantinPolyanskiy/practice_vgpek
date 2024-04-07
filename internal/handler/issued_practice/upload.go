package issued_practice

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/issued"
	"practice_vgpek/pkg/apperr"
	"time"
)

type IssuedPracticeService interface {
	Save(ctx context.Context, req issued.UploadReq) (issued.UploadResp, error)
}

type Handler struct {
	l *zap.Logger
	s IssuedPracticeService
}

func NewIssuedPracticeHandler(service IssuedPracticeService, logger *zap.Logger) Handler {
	return Handler{
		s: service,
		l: logger,
	}
}

func (h Handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.UploadIssuedPracticeOperation),
		zap.String("слой", "http обработчики"),
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

	req := issued.UploadReq{
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
	resp, err := h.s.Save(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  "Таймаут",
			})
			return
		} else if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.GetRoleOperation,
				Error:  err.Error(),
			})
			return
		}
	}
	l.Info("практическое задание успешно загружено")

	render.JSON(w, r, resp)
	return
}
