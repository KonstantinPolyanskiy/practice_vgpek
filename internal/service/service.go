package service

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/internal/repository"
	"practice_vgpek/internal/service/authn"
	"practice_vgpek/internal/service/reg_key"
)

type AuthnService interface {
	NewPerson(ctx context.Context, registering person.RegistrationReq) (person.RegisteredResp, error)
}

type KeyService interface {
	NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.AddResp, error)
}

type Service struct {
	AuthnService
	KeyService
}

func New(repository repository.Repository, logger *zap.Logger) Service {
	return Service{
		AuthnService: authn.NewAuthenticationService(repository.PersonRepo, repository.AccountRepo, repository.KeyRepo, logger),
		KeyService:   reg_key.NewKeyService(repository.KeyRepo, logger),
	}
}
