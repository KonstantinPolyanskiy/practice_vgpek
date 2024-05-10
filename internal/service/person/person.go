package person

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/params"
)

type KeyDAO interface {
	ByBody(ctx context.Context, body string) (entity.Key, error)
	Update(ctx context.Context, key entity.KeyUpdate) (entity.Key, error)
}

type KeyService interface {
	InvalidateKey(ctx context.Context, id dto.EntityId) (domain.InvalidatedKey, error)
	Increment(ctx context.Context, key entity.Key) (entity.Key, error)
}

type RoleDAO interface {
	ById(ctx context.Context, id int) (entity.Role, error)
}

type RoleService interface {
	RoleById(ctx context.Context, req dto.EntityId) (domain.Role, error)
}

type PermDAO interface {
	ByRoleId(ctx context.Context, roleId int) ([]entity.Permissions, error)
}

type PersonDAO interface {
	Save(ctx context.Context, data dto.PersonRegistrationData) (entity.Person, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Person, error)
}

type AccountDAO interface {
	Save(ctx context.Context, data dto.AccountRegistrationData) (entity.Account, error)
	ById(ctx context.Context, id int) (entity.Account, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Account, error)
	HardDeleteById(ctx context.Context, id int) error
}

type Service struct {
	logger *zap.Logger

	keyDAO      KeyDAO
	personDAO   PersonDAO
	accountDAO  AccountDAO
	roleDAO     RoleDAO
	roleService RoleService

	permDAO PermDAO

	keyService KeyService
}

func New(roleService RoleService, permDAO PermDAO, kd KeyDAO, pd PersonDAO, ad AccountDAO, rd RoleDAO, keyService KeyService, logger *zap.Logger) Service {
	return Service{
		logger: logger,

		keyDAO:     kd,
		personDAO:  pd,
		accountDAO: ad,
		roleDAO:    rd,
		permDAO:    permDAO,

		roleService: roleService,
		keyService:  keyService,
	}
}
