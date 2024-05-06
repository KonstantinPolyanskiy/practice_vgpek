package person

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

type PermResult struct {
	Perm  domain.RolePermission
	Error error
}

func (s Service) PermByAccountId(ctx context.Context, id dto.EntityId) (domain.RolePermission, error) {
	resCh := make(chan PermResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.GetPermByAccountIdOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		accId := ctx.Value("accountId").(int)

		perm, err := s.accountMediator.PermByAccountId(ctx, accId)
		if err != nil {
			l.Warn("ошибка получения доступов по аккаунту", zap.Error(err), zap.Int("accountId", accId))
			sendGetPermResult(resCh, domain.RolePermission{}, "ошибка получения доступов по аккаунту")
			return
		}

		sendGetPermResult(resCh, perm, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.RolePermission{}, ctx.Err()
		case result := <-resCh:
			return result.Perm, result.Error
		}
	}
}

func sendGetPermResult(ch chan PermResult, perm domain.RolePermission, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	ch <- PermResult{
		Perm:  perm,
		Error: err,
	}
}
