package token

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
)

type AccountDAO interface {
	ByLogin(ctx context.Context, login string) (entity.Account, error)
}

type Service struct {
	accountDAO AccountDAO

	logger *zap.Logger

	signingKey string
}

func New(accountDAO AccountDAO, key string, logger *zap.Logger) Service {
	return Service{
		accountDAO: accountDAO,
		logger:     logger,
		signingKey: key,
	}
}
