package reg_key

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
)

const ObjectName = "KEY"

type Repository interface {
	SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error)
	KeysByParams(ctx context.Context, params params.Key) ([]registration_key.Entity, error)
	RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error)
	IncCountUsages(ctx context.Context, keyId int) error
	Invalidate(ctx context.Context, keyId int) error
}

type AccountMediator interface {
	// RoleByAccountId по Id аккаунта находит его роль
	RoleByAccountId(ctx context.Context, id int) (permissions.RoleEntity, error)
	// HasAccess проверяет, есть ли у указанной роли доступ к переданному действию к указанному объекту
	HasAccess(ctx context.Context, roleId int, objectName, actionName string) (bool, error)
}

type Service struct {
	l               *zap.Logger
	r               Repository
	accountMediator AccountMediator
}

func NewKeyService(repository Repository, logger *zap.Logger, accountMediator AccountMediator) Service {
	return Service{
		l:               logger,
		r:               repository,
		accountMediator: accountMediator,
	}
}
