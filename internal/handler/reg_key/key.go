package reg_key

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/params"
)

type Service interface {
	NewKey(ctx context.Context, req dto.NewKeyReq) (domain.Key, error)
	InvalidateKey(ctx context.Context, req dto.EntityId) (domain.InvalidatedKey, error)

	KeysByParams(ctx context.Context, keyParams params.State) ([]domain.Key, error)
	KeyById(ctx context.Context, req dto.EntityId) (domain.Key, error)
}

type AccountMediator interface {
	HasAccess(ctx context.Context, accountId int, objectName, actionName string) (bool, error)
}

type Handler struct {
	l *zap.Logger
	s Service

	accountMediator AccountMediator
}

func NewKeyHandler(service Service, accountMediator AccountMediator, logger *zap.Logger) Handler {
	return Handler{
		l:               logger,
		s:               service,
		accountMediator: accountMediator,
	}
}
