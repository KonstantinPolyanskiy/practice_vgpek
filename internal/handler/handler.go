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
	Identity(next http.Handler) http.Handler
}

type KeyHandler interface {
	AddKey(w http.ResponseWriter, r *http.Request)
	DeleteKey(w http.ResponseWriter, r *http.Request)
	GetKeys(w http.ResponseWriter, r *http.Request)
}

type RBACHandler interface {
	AddAction(w http.ResponseWriter, r *http.Request)

	AddObject(w http.ResponseWriter, r *http.Request)

	AddRole(w http.ResponseWriter, r *http.Request)

	AddPermission(w http.ResponseWriter, r *http.Request)

	GetActions(w http.ResponseWriter, r *http.Request)
	GetAction(w http.ResponseWriter, r *http.Request)

	GetObject(w http.ResponseWriter, r *http.Request)
	GetObjects(w http.ResponseWriter, r *http.Request)

	GetRole(w http.ResponseWriter, r *http.Request)
	GetRoles(w http.ResponseWriter, r *http.Request)
}

type IssuedPracticeHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	l *zap.Logger
	AuthnHandler
	KeyHandler
	RBACHandler
	IssuedPracticeHandler
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
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.KeyHandler.AddKey)
		r.Delete("/", h.KeyHandler.DeleteKey)

		r.Get("/params", h.KeyHandler.GetKeys)
	})

	r.Route("/action", func(r chi.Router) {
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.RBACHandler.AddAction)

		r.Get("/", h.RBACHandler.GetAction)
		r.Get("/params", h.RBACHandler.GetActions)
	})

	r.Route("/object", func(r chi.Router) {
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.RBACHandler.AddObject)

		r.Get("/", h.RBACHandler.GetObject)
		r.Get("/params", h.RBACHandler.GetObjects)
	})

	r.Route("/role", func(r chi.Router) {
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.RBACHandler.AddRole)

		r.Get("/", h.RBACHandler.GetRole)
		r.Get("/params", h.RBACHandler.GetRoles)
	})

	r.Route("/permissions", func(r chi.Router) {
		r.Post("/", h.RBACHandler.AddPermission)
	})

	r.Route("/practice", func(r chi.Router) {
		r.Route("/issued", func(r chi.Router) {
			r.Post("/", h.IssuedPracticeHandler.Upload)
		})
	})

	return r
}
