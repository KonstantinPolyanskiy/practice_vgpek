package service

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/mediator/account"
	"practice_vgpek/internal/model/params"
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
	NewToken(ctx context.Context, logIn person.LogInReq) (person.LogInResp, error)
	// ParseToken в случае, если токен распарсен - возвращает id аккаунта
	ParseToken(token string) (int, error)
}

type RBACService interface {
	NewAction(ctx context.Context, addingAction permissions.AddActionReq) (permissions.AddActionResp, error)
	ActionById(ctx context.Context, req permissions.GetActionReq) (permissions.ActionEntity, error)
	ActionsByParams(ctx context.Context, params params.Default) ([]permissions.ActionEntity, error)

	NewObject(ctx context.Context, addingObject permissions.AddObjectReq) (permissions.AddObjectResp, error)
	ObjectById(ctx context.Context, id int) (permissions.ObjectEntity, error)
	ObjectsByParams(ctx context.Context, params params.Default) ([]permissions.ObjectEntity, error)

	NewRole(ctx context.Context, addingRole permissions.AddRoleReq) (permissions.AddRoleResp, error)
	RoleById(ctx context.Context, id int) (permissions.RoleEntity, error)
	RolesByParams(ctx context.Context, params params.Default) ([]permissions.RoleEntity, error)

	NewPermission(ctx context.Context, addingPerm permissions.AddPermReq) (permissions.AddPermResp, error)
}

type KeyService interface {
	NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.AddResp, error)
	InvalidateKey(ctx context.Context, deletingKey registration_key.DeleteReq) (registration_key.DeleteResp, error)
	Keys(ctx context.Context, keyParams params.Key) (registration_key.GetKeysResp, error)
}

type Service struct {
	AuthnService
	KeyService
	RBACService
}

func New(repository repository.Repository, logger *zap.Logger) Service {
	am := account.NewAccountMediator(repository.AccountRepo, repository.KeyRepo, repository.RoleRepo, repository.PermissionRepo)

	return Service{
		AuthnService: authn.NewAuthenticationService(repository.PersonRepo, repository.AccountRepo, repository.KeyRepo, logger),
		KeyService:   reg_key.NewKeyService(repository.KeyRepo, logger, am),
		RBACService:  rbac.NewRBACService(repository.ActionRepo, repository.ObjectRepo, repository.RoleRepo, repository.PermissionRepo, am, logger),
	}
}
