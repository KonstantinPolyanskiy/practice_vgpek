package key

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/params"
)

type DAO interface {
	Update(ctx context.Context, key entity.Key) (entity.Key, error)
	ById(ctx context.Context, id int) (entity.Key, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Key, error)
	Save(ctx context.Context, info dto.NewKeyInfo) (entity.Key, error)
}

type RoleDAO interface {
	ById(ctx context.Context, id int) (entity.Role, error)
}

type AccountMediator interface {
	// RoleByAccountId по Id аккаунта находит его роль
	RoleByAccountId(ctx context.Context, id int) (entity.Role, error)
	// HasAccess проверяет, есть ли у указанной роли доступ к переданному действию к указанному объекту
	HasAccess(ctx context.Context, roleId int, objectName, actionName string) (bool, error)
}

type Service struct {
	l               *zap.Logger
	accountMediator AccountMediator
	keyDAO          DAO
	roleDAO         RoleDAO
}

func New(kd DAO, rd RoleDAO, accountMediator AccountMediator, logger *zap.Logger) Service {
	return Service{
		l:               logger,
		keyDAO:          kd,
		roleDAO:         rd,
		accountMediator: accountMediator,
	}
}
