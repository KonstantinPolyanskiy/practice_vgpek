package key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

type IncrementResult struct {
	Key   entity.Key
	Error error
}

func (s Service) Increment(ctx context.Context, key entity.Key) (entity.Key, error) {
	resCh := make(chan IncrementResult)

	_ = s.l.With(
		zap.String(operation.Operation, operation.IncrementKey),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	current := key.CurrentCountUsages + 1

	go func() {
		incremented, err := s.keyDAO.Update(ctx, entity.KeyUpdate{
			Id:                 key.Id,
			RoleId:             nil,
			Body:               nil,
			MaxCountUsages:     nil,
			CurrentCountUsages: &current,
			CreatedAt:          nil,
			IsValid:            nil,
			InvalidationTime:   nil,
			GroupName:          nil,
		})
		if err != nil {
			sendIncrementResult(resCh, entity.Key{}, "Ошибка инкрементирования ключа")
			return
		}

		sendIncrementResult(resCh, incremented, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return entity.Key{}, ctx.Err()
		case result := <-resCh:
			return result.Key, result.Error
		}
	}
}

func sendIncrementResult(resCh chan IncrementResult, resp entity.Key, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- IncrementResult{
		Key:   resp,
		Error: err,
	}
}
