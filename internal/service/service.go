package service

import (
	"context"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/repository"
	"practice_vgpek/internal/service/authn"
)

type AuthnService interface {
	NewPerson(ctx context.Context, registering person.RegistrationReq) (person.RegisteredResp, error)
}

type Service struct {
	AuthnService
}

func New(repository repository.Repository) Service {
	return Service{
		AuthnService: authn.NewAuthenticationService(repository.PersonRepo),
	}
}
