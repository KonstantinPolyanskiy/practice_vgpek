package handler

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"practice_vgpek/internal/handler/authn"
	"practice_vgpek/internal/handler/rbac"
	"practice_vgpek/internal/handler/reg_key"
	"practice_vgpek/internal/service"
)

type AuthnHandler interface {
	Registration(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type KeyHandler interface {
	AddKey(w http.ResponseWriter, r *http.Request)
}

type RBACHandler interface {
	AddAction(w http.ResponseWriter, r *http.Request)
	AddObject(w http.ResponseWriter, r *http.Request)
	AddRole(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	l *zap.Logger
	AuthnHandler
	KeyHandler
	RBACHandler
}

func New(service service.Service, logger *zap.Logger) Handler {
	return Handler{
		l:            logger,
		AuthnHandler: authn.NewAuthenticationHandler(service.AuthnService, logger),
		KeyHandler:   reg_key.NewRegKeyHandler(service.KeyService, logger),
		RBACHandler:  rbac.NewAccessHandler(service.RBACService, logger),
	}
}

func (h Handler) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/registration", func(r chi.Router) {
		r.Post("/", h.AuthnHandler.Registration)
	})

	r.Route("/login", func(r chi.Router) {
		r.Post("/", h.AuthnHandler.Login)
	})

	r.Route("/key", func(r chi.Router) {
		r.Post("/", h.KeyHandler.AddKey)
	})

	r.Route("/action", func(r chi.Router) {
		r.Post("/", h.RBACHandler.AddAction)
	})

	r.Route("/object", func(r chi.Router) {
		r.Post("/", h.RBACHandler.AddObject)
	})

	r.Route("/role", func(r chi.Router) {
		r.Post("/", h.RBACHandler.AddRole)
	})

	return r
}
