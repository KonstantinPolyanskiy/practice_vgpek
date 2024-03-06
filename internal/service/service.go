package service

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/internal/repository"
	"practice_vgpek/internal/service/authn"
	"practice_vgpek/internal/service/rbac"
	"practice_vgpek/internal/service/reg_key"
)

type AuthnService interface {
	NewPerson(ctx context.Context, registering person.RegistrationReq) (person.RegisteredResp, error)
}

type RBACService interface {
	NewAction(ctx context.Context, addingAction permissions.AddActionReq) (permissions.AddActionResp, error)
	NewObject(ctx context.Context, addingObject permissions.AddObjectReq) (permissions.AddObjectResp, error)
	NewRole(ctx context.Context, addingRole permissions.AddRoleReq) (permissions.AddRoleResp, error)
}

type KeyService interface {
	NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.AddResp, error)
}

type Service struct {
	AuthnService
	KeyService
	RBACService
}

func New(repository repository.Repository, logger *zap.Logger) Service {
	return Service{
		AuthnService: authn.NewAuthenticationService(repository.PersonRepo, repository.AccountRepo, repository.KeyRepo, logger),
		KeyService:   reg_key.NewKeyService(repository.KeyRepo, logger),
		RBACService:  rbac.NewRBACService(repository.ActionRepo, repository.ObjectRepo, repository.RoleRepo, logger),
	}
}
