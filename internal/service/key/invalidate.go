package key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"time"
)

const (
	InvalidateActionName = "DEL"
)

type InvalidateKeyResult struct {
	DeletedKey domain.InvalidatedKey
	Error      error
}

func (s Service) InvalidateKey(ctx context.Context, req dto.EntityId) (domain.InvalidatedKey, error) {
	resCh := make(chan InvalidateKeyResult)

	_ = s.l.With(
		zap.String(operation.Operation, operation.InvalidateKeyOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	isValid := false

	go func() {
		deleteTime := time.Now()
		invalidated, err := s.keyDAO.Update(ctx, entity.KeyUpdate{
			Id:                 req.Id,
			RoleId:             nil,
			Body:               nil,
			MaxCountUsages:     nil,
			CurrentCountUsages: nil,
			CreatedAt:          nil,
			IsValid:            &isValid,
			InvalidationTime:   &deleteTime,
			GroupName:          nil,
		})
		if err != nil {
			sendInvalidateKeyResult(resCh, domain.InvalidatedKey{}, "ошибка инвалидирования ключа")
			return
		}

		resp := domain.InvalidatedKey{
			Id:               invalidated.Id,
			RoleId:           invalidated.RoleId,
			CreatedAt:        invalidated.CreatedAt,
			IsValid:          invalidated.IsValid,
			InvalidationTime: *invalidated.InvalidationTime,
		}

		sendInvalidateKeyResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.InvalidatedKey{}, ctx.Err()
		case result := <-resCh:
			return result.DeletedKey, result.Error
		}
	}
}

func sendInvalidateKeyResult(resCh chan InvalidateKeyResult, resp domain.InvalidatedKey, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- InvalidateKeyResult{
		DeletedKey: resp,
		Error:      err,
	}
}
