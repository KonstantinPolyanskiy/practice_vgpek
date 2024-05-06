package apperr

import (
	"github.com/go-chi/render"
	"net/http"
)

type AppError struct {
	Action string `json:"action"`
	Error  string `json:"error"`
}

func New(w http.ResponseWriter, r *http.Request, code int, ae AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	render.JSON(w, r, ae)
}
