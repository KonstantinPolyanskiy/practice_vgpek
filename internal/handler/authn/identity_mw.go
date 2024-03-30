package authn

import (
	"context"
	"net/http"
	"practice_vgpek/pkg/apperr"
	"strings"
)

func (h Handler) Identity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			apperr.New(w, r, http.StatusUnauthorized, apperr.AppError{
				Action: "Авторизация",
				Error:  "пустой хедер Authorization",
			})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			apperr.New(w, r, http.StatusUnauthorized, apperr.AppError{
				Action: "Авторизация",
				Error:  "невалидный хедер Authorization",
			})
			return
		}

		id, err := h.s.ParseToken(headerParts[1])
		if err != nil {
			apperr.New(w, r, http.StatusUnauthorized, apperr.AppError{
				Action: "Авторизация",
				Error:  err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), "AccountId", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
