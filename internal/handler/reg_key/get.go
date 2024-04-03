package reg_key

import (
	"context"
	"github.com/go-chi/render"
	"net/http"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/service/reg_key"
	"practice_vgpek/pkg/apperr"
	"practice_vgpek/pkg/queryutils"
	"time"
)

func (h Handler) GetKeys(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	defaultParams, err := queryutils.DefaultParams(r, 10, 0)
	if err != nil {
		apperr.New(w, r, http.StatusBadRequest, apperr.AppError{
			Action: reg_key.GetKeysOperation,
			Error:  "Неправильные параметры запроса",
		})
		return
	}

	keyParams := getKeyParams(r, defaultParams)

	keys, err := h.s.Keys(ctx, keyParams)
	if err != nil {
		apperr.New(w, r, http.StatusInternalServerError, apperr.AppError{
			Action: reg_key.GetKeysOperation,
			Error:  err.Error(),
		})
	}

	render.JSON(w, r, keys)
	return
}

func getKeyParams(r *http.Request, defaultParams params.Default) params.Key {
	isValid := true

	v := r.URL.Query().Get("valid")
	if v == "false" {
		isValid = false
	}

	return params.Key{
		IsValid: isValid,
		Default: defaultParams,
	}

}
