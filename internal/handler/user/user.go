package user

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/params"
)

type AccountService interface {
	EntityAccountById(ctx context.Context, id dto.EntityId) (entity.Account, error)
	EntityAccountByParam(ctx context.Context, p params.State) ([]entity.Account, error)
}

type AccountMediator interface {
	HasAccess(ctx context.Context, accountId int, objectName, actionName string) (bool, error)
}

type Handler struct {
	logger *zap.Logger
	AccountService
	AccountMediator
}

func New(accountService AccountService, accountMediator AccountMediator, logger *zap.Logger) Handler {
	return Handler{
		logger:          logger,
		AccountService:  accountService,
		AccountMediator: accountMediator,
	}
}
