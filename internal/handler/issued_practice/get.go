package issued_practice

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/transport/rest"
	"practice_vgpek/pkg/apiutils"
	"practice_vgpek/pkg/apperr"
	"practice_vgpek/pkg/queryutils"
	"strconv"
	"strings"
	"time"
)

func (h Handler) PracticeById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetIssuedPracticeInfoById),
		zap.String(layer.Layer, layer.HTTPLayer),
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

	render.JSON(w, r, rest.IssuedPractice{}.DomainToResponse(practice).WithDownloadLink(link))
	return
}

func (h Handler) Download(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.DownloadIssuedPractice),
		zap.String(layer.Layer, layer.HTTPLayer),
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

	practice, err := h.s.ById(ctx, dto.EntityId{Id: id})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
				Action: operation.DownloadIssuedPractice,
				Error:  "Таймаут",
			})
			return
		} else {
			code := http.StatusInternalServerError

			apperr.New(w, r, code, apperr.AppError{
				Action: operation.DownloadIssuedPractice,
				Error:  err.Error(),
			})
			return
		}
	}

	path := practice.Path

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

func (h Handler) PracticeByParams(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	l := h.l.With(
		zap.String(layer.Endpoint, r.RequestURI),
		zap.String(operation.Operation, operation.GetIssuedPracticeInfoByParams),
		zap.String(layer.Layer, layer.HTTPLayer),
	)

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		l.Warn("ошибка получени параметров запроса", zap.Error(err))

		apperr.New(w, r, http.StatusRequestTimeout, apperr.AppError{
			Action: operation.GetIssuedPracticeInfoByParams,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	practiceParams := getPracticeParams(r, defaultParams)

	l.Info("попытка получить практические задания",
		zap.Int("id аккаунта", r.Context().Value("AccountId").(int)),
		zap.Int("лимит", practiceParams.Limit),
		zap.Int("оффсет", practiceParams.Offset),
		zap.String("статус решения", practiceParams.IsSolved),
	)

}

func getPracticeParams(r *http.Request, defaultParams params.Default) params.IssuedPractice {
	var isSolved string

	v := r.URL.Query().Get("solved")

	switch v {
	case "all":
		isSolved = "all"
	case "yes":
		isSolved = "yes"
	case "no":
		isSolved = "no"
	// по умолчанию получают только не решенные практические
	default:
		isSolved = "no"
	}

	return params.IssuedPractice{
		IsSolved: isSolved,
		Default:  defaultParams,
	}
}
