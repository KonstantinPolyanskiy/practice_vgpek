package key

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/params"
)

type DAO interface {
	Update(ctx context.Context, key entity.KeyUpdate) (entity.Key, error)
	ById(ctx context.Context, id int) (entity.Key, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Key, error)
	Save(ctx context.Context, info dto.NewKeyInfo) (entity.Key, error)
}

type RoleDAO interface {
	ById(ctx context.Context, id int) (entity.Role, error)
}

type Service struct {
	l       *zap.Logger
	keyDAO  DAO
	roleDAO RoleDAO
}

func New(kd DAO, rd RoleDAO, logger *zap.Logger) Service {
	return Service{
		l:       logger,
		keyDAO:  kd,
		roleDAO: rd,
	}
}
