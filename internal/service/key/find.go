package key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
)

type GetKeysResult struct {
	Keys  []domain.Key
	Error error
}

type GetKeyResult struct {
	Key   domain.Key
	Error error
}

func (s Service) KeyById(ctx context.Context, req dto.EntityId) (domain.Key, error) {
	resCh := make(chan GetKeyResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetKeyByIdOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		keyEntity, err := s.keyDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetKeyResult(resCh, domain.Key{}, "Ошибка получения ключа")
			return
		}

		role, err := s.roleDAO.ById(ctx, keyEntity.RoleId)
		if err != nil {
			sendGetKeyResult(resCh, domain.Key{}, "Ошибка получения роли")
			return
		}

		key := domain.Key{
			Id:             keyEntity.Id,
			RoleId:         role.Id,
			RoleName:       role.Name,
			Body:           keyEntity.Body,
			MaxCountUsages: keyEntity.MaxCountUsages,
			CountUsages:    keyEntity.CurrentCountUsages,
			CreatedAt:      keyEntity.CreatedAt,
			Group:          keyEntity.GroupName,
			IsValid:        keyEntity.IsValid,
		}

		l.Info("ключ найден", zap.Int("id", key.Id))

		sendGetKeyResult(resCh, key, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Key{}, ctx.Err()
		case result := <-resCh:
			return result.Key, result.Error

		}
	}
}

func (s Service) KeysByParams(ctx context.Context, keyParams params.State) ([]domain.Key, error) {
	resCh := make(chan GetKeysResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.GetKeysOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		keys, err := s.keyDAO.ByParams(ctx, keyParams.Default)
		if err != nil {
			sendGetKeysResult(resCh, nil, "Ошибка получения ключей")
			return
		}

		domainKeys := make([]domain.Key, 0, 10)

		for _, key := range keys {
			role, err := s.roleDAO.ById(ctx, key.RoleId)
			if err != nil {
				l.Warn("ошибка получения роли",
					zap.Int("id роли", key.RoleId),
					zap.Error(err),
				)
				continue
			}

			domainKeys = append(domainKeys, domain.Key{
				Id:             key.Id,
				RoleId:         role.Id,
				RoleName:       role.Name,
				Body:           key.Body,
				MaxCountUsages: key.MaxCountUsages,
				CountUsages:    key.CurrentCountUsages,
				CreatedAt:      key.CreatedAt,
				Group:          key.GroupName,
				IsValid:        key.IsValid,
			})
		}

		keysResp := make([]domain.Key, 0)

		switch keyParams.State {
		case params.All:
			keysResp = append(keysResp, domainKeys...)
		case params.Deleted:
			keysResp = filterDeleted(domainKeys)
		case params.NotDeleted:
			keysResp = filterNotDeleted(domainKeys)
		default:
			l.Warn("ошибка при фильтрации",
				zap.String("состояние", keyParams.State),
				zap.Int("кол-во ключей из бд", len(keys)),
			)

			sendGetKeysResult(resCh, nil, "ошибка фильтрации")
			return
		}

		sendGetKeysResult(resCh, keysResp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Keys, result.Error

		}
	}
}

// filterDeleted возвращает ключи, которые помечены как удаленные
func filterDeleted(keys []domain.Key) []domain.Key {
	result := make([]domain.Key, 0, 10)

	for _, key := range keys {
		if !key.IsValid {
			result = append(result, key)
		}
	}

	return result
}

// filterNotDeleted возвращает ключи, которые не помечены как удаленные
func filterNotDeleted(keys []domain.Key) []domain.Key {
	result := make([]domain.Key, 0, 10)

	for _, key := range keys {
		if key.IsValid {
			result = append(result, key)
		}
	}

	return result
}

func sendGetKeysResult(resCh chan GetKeysResult, resp []domain.Key, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetKeysResult{
		Keys:  resp,
		Error: err,
	}
}

func sendGetKeyResult(resCh chan GetKeyResult, resp domain.Key, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetKeyResult{
		Key:   resp,
		Error: err,
	}
}
