package reg_key

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/registration_key"
)

type Repository interface {
	SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error)
	RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error)
	IncCountUsages(ctx context.Context, keyId int) error
	Invalidate(ctx context.Context, keyId int) error
}

type Service struct {
	l *zap.Logger
	r Repository
}

func NewKeyService(repository Repository, logger *zap.Logger) Service {
	return Service{
		l: logger,
		r: repository,
	}
}
