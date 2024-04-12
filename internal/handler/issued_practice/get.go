package issued_practice

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/issued"
	"practice_vgpek/pkg/apiutils"
	"practice_vgpek/pkg/apperr"
	"strconv"
	"strings"
	"time"
)

func (h Handler) PracticeById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.GetIssuedPracticeInfoById),
		zap.String("слой", "http обработчики"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.GetIssuedPracticeInfoById,
			Error:  "Преобразование запроса на получение практического задания",
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

	link := r.Host + r.URL.String()
	link = strings.Replace(link, "?", "/download?", 1)

	render.JSON(w, r, issued.GetPracticeResp{
		PracticeId:   practice.PracticeId,
		AuthorId:     practice.AccountId,
		Title:        practice.Title,
		Theme:        practice.Theme,
		Major:        practice.Major,
		DownloadLink: link,
	})
	return
}

func (h Handler) Download(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String("адрес", r.RequestURI),
		zap.String("операция", operation.DownloadIssuedPractice),
		zap.String("слой", "http обработчики"),
	)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		l.Warn("ошибка декодирования данных", zap.Error(err))

		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: operation.DownloadIssuedPractice,
			Error:  "Преобразование запроса на получение практического задания",
		})
		return
	}

	practice, err := h.s.ById(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.DownloadIssuedPractice,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			if errors.Is(err, permissions.ErrDontHavePerm) {
				code = http.StatusForbidden
			}

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.DownloadIssuedPractice,
				Error:  err.Error(),
			})
			return
		}
	}

	path := practice.PracticePath

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		l.Warn("ошибка открытия файла", zap.Error(err))
		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: operation.DownloadIssuedPractice,
			Error:  "не удалось найти файл",
		})
		return
	}

	i, err := f.Stat()
	if err != nil {
		l.Warn("ошибка чтения метаинформации файла", zap.Error(err))
		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: operation.DownloadIssuedPractice,
			Error:  "не удалось найти файл",
		})
		return
	}

	apiutils.SetDownloadHeaders(w, path, strconv.Itoa(int(i.Size())))
	w.WriteHeader(http.StatusOK)

	log.Println(r.Host + r.URL.String())

	_, err = io.Copy(w, f)
	if err != nil {
		l.Warn("ошибка выдачи файла", zap.Error(err))
		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: operation.GetIssuedPracticeInfoById,
			Error:  "не удалось выгрузить файл",
		})
		return
	}
}