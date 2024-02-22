package handler

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/handler/authn"
	"practice_vgpek/internal/handler/reg_key"
	"practice_vgpek/internal/service"
)

type AuthnHandler interface {
	Registration(w http.ResponseWriter, r *http.Request)
}

type KeyHandler interface {
	AddKey(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	l *zap.Logger
	AuthnHandler
	KeyHandler
}

func New(service service.Service, logger *zap.Logger) Handler {
	return Handler{
		l:            logger,
		AuthnHandler: authn.NewAuthenticationHandler(service.AuthnService, logger),
		KeyHandler:   reg_key.NewRegKeyHandler(service.KeyService, logger),
	}
}

func (h Handler) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/registration", func(r chi.Router) {
		r.Post("/", h.AuthnHandler.Registration)
	})

	r.Route("/key", func(r chi.Router) {
		r.Post("/", h.KeyHandler.AddKey)
	})

	return r
}
