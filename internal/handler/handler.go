package handler

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	_ "practice_vgpek/docs" // docs are generated by Swag CLI, you have to import it.
	"practice_vgpek/internal/handler/authn"
	"practice_vgpek/internal/handler/issued_practice"
	"practice_vgpek/internal/handler/rbac"
	"practice_vgpek/internal/handler/reg_key"
	"practice_vgpek/internal/handler/solved_practice"
	"practice_vgpek/internal/handler/user"
	"practice_vgpek/internal/mediator/account"
	"practice_vgpek/internal/service"
)

type AuthnHandler interface {
	Registration(w http.ResponseWriter, r *http.Request)

	Login(w http.ResponseWriter, r *http.Request)
	Identity(next http.Handler) http.Handler
}

type UserHandler interface {
	GetAccount(w http.ResponseWriter, r *http.Request)
	GetAccountsByParam(w http.ResponseWriter, r *http.Request)
	GetPersonsByParam(w http.ResponseWriter, r *http.Request)
}

type KeyHandler interface {
	AddKey(w http.ResponseWriter, r *http.Request)
	DeleteKey(w http.ResponseWriter, r *http.Request)

	GetKey(w http.ResponseWriter, r *http.Request)
	GetKeys(w http.ResponseWriter, r *http.Request)
}

type RBACHandler interface {
	AddAction(w http.ResponseWriter, r *http.Request)
	DeleteAction(w http.ResponseWriter, r *http.Request)

	AddObject(w http.ResponseWriter, r *http.Request)
	DeleteObject(w http.ResponseWriter, r *http.Request)

	AddRole(w http.ResponseWriter, r *http.Request)
	DeleteRole(w http.ResponseWriter, r *http.Request)

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

	PracticeById(w http.ResponseWriter, r *http.Request)
	PracticeByParams(w http.ResponseWriter, r *http.Request)

	Download(w http.ResponseWriter, r *http.Request)
}

type SolvedPracticeHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)

	PracticeById(w http.ResponseWriter, r *http.Request)

	SetMark(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	l *zap.Logger

	AuthnHandler

	UserHandler

	KeyHandler

	RBACHandler

	IssuedPracticeHandler
	SolvedPracticeHandler
}

func New(service service.Service, logger *zap.Logger) Handler {
	accountMediator := account.NewAccountMediator(service.PersonService, service.KeyService, service.RBACService, service.RBACService)
	return Handler{
		l:                     logger,
		AuthnHandler:          authn.NewAuthenticationHandler(service.PersonService, service.TokenService, service.RBACService, logger),
		KeyHandler:            reg_key.NewKeyHandler(service.KeyService, accountMediator, logger),
		RBACHandler:           rbac.NewAccessHandler(service.RBACService, accountMediator, logger),
		IssuedPracticeHandler: issued_practice.NewIssuedPracticeHandler(service.IssuedPracticeService, logger),
		SolvedPracticeHandler: solved_practice.NewCompletedPracticeHandler(service.SolvedPracticeService, logger),
		UserHandler:           user.New(service.PersonService, service.PersonService, accountMediator, logger),
	}
}

func (h Handler) Init() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/person", func(r chi.Router) {
		r.Post("/", h.AuthnHandler.Registration)

		r.Route("/", func(r chi.Router) {
			r.Use(h.AuthnHandler.Identity)

			r.Get("/", h.GetPersonsByParam)
		})
		r.Route("/account", func(r chi.Router) {
			r.Use(h.AuthnHandler.Identity)

			r.Get("/", h.GetAccount)
			r.Get("/params", h.GetAccountsByParam)
		})
	})

	r.Route("/login", func(r chi.Router) {
		r.Post("/", h.AuthnHandler.Login)
	})

	r.Route("/key", func(r chi.Router) {
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.KeyHandler.AddKey)
		r.Delete("/", h.KeyHandler.DeleteKey)

		r.Get("/", h.KeyHandler.GetKey)
		r.Get("/params", h.KeyHandler.GetKeys)
	})

	r.Route("/action", func(r chi.Router) {
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.RBACHandler.AddAction)

		r.Get("/", h.RBACHandler.GetAction)
		r.Get("/params", h.RBACHandler.GetActions)

		r.Delete("/", h.RBACHandler.DeleteAction)
	})

	r.Route("/object", func(r chi.Router) {
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.RBACHandler.AddObject)

		r.Get("/", h.RBACHandler.GetObject)
		r.Get("/params", h.RBACHandler.GetObjects)

		r.Delete("/", h.RBACHandler.DeleteObject)
	})

	r.Route("/role", func(r chi.Router) {
		r.Use(h.AuthnHandler.Identity)

		r.Post("/", h.RBACHandler.AddRole)

		r.Get("/", h.RBACHandler.GetRole)
		r.Get("/params", h.RBACHandler.GetRoles)

		r.Delete("/", h.RBACHandler.DeleteRole)
	})

	r.Route("/permissions", func(r chi.Router) {
		r.Post("/", h.RBACHandler.AddPermission)
	})

	r.Route("/practice", func(r chi.Router) {
		r.Route("/issued", func(r chi.Router) {
			r.Use(h.AuthnHandler.Identity)

			r.Post("/", h.IssuedPracticeHandler.Upload)

			r.Get("/", h.IssuedPracticeHandler.PracticeById)
			r.Get("/download", h.IssuedPracticeHandler.Download)
			r.Get("/params", h.IssuedPracticeHandler.PracticeByParams)

		})
		r.Route("/solved", func(r chi.Router) {
			r.Use(h.AuthnHandler.Identity)

			r.Post("/", h.SolvedPracticeHandler.Upload)

			r.Post("/mark", h.SolvedPracticeHandler.SetMark)

			r.Get("/", h.SolvedPracticeHandler.PracticeById)
		})
	})

	return r
}
