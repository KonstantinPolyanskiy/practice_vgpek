package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"practice_vgpek/internal/handler/authn"
	"practice_vgpek/internal/service"
)

type AuthnHandler interface {
	Registration(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	AuthnHandler
}

func New(service service.Service) Handler {
	return Handler{
		//TODO: подставить service
		AuthnHandler: authn.NewAuthenticationHandler(service.AuthnService),
	}
}

func (h Handler) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/registration", func(r chi.Router) {
		r.Post("/", h.AuthnHandler.Registration)
	})

	return r
}
