package key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/rndutils"
	"time"
)

type CreatingKeyResult struct {
	CreatedKey domain.Key
	Error      error
}

func (s Service) NewKey(ctx context.Context, req dto.NewKeyReq) (domain.Key, error) {
	resCh := make(chan CreatingKeyResult)

	l := s.l.With(
		zap.String(operation.Operation, operation.NewKeyOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		// Проверяем что указанное количество использований ключа больше 0
		if req.MaxCountUsages <= 0 {
			l.Warn("неправильное макс. кол-во использований ключа", zap.Int("макс кол-во использований", req.MaxCountUsages))

			sendNewKeyResult(resCh, domain.Key{}, "Неправильное кол-во использований ключа")
			return
		}

		if req.GroupName == "" {
			req.GroupName = "unknown"
		}

		// Формируем DTO
		info := dto.NewKeyInfo{
			RoleId:         req.RoleId,
			Body:           rndutils.RandString(7),
			MaxCountUsages: req.MaxCountUsages,
			CreatedAt:      time.Now(),
			Group:          req.GroupName,
		}

		// Сохраняем ключ
		savedKey, err := s.keyDAO.Save(ctx, info)
		if err != nil {
			sendNewKeyResult(resCh, domain.Key{}, "Ошибка сохранения ключа")
			return
		}

		role, err := s.roleDAO.ById(ctx, savedKey.RoleId)
		if err != nil {
			sendNewKeyResult(resCh, domain.Key{}, "Ошибка получения роли")
			return
		}

		// Формируем и отправляем ответ
		key := domain.Key{
			Id:             savedKey.Id,
			RoleId:         role.Id,
			RoleName:       role.Name,
			Body:           savedKey.Body,
			MaxCountUsages: savedKey.MaxCountUsages,
			CountUsages:    savedKey.CurrentCountUsages,
			CreatedAt:      savedKey.CreatedAt,
			Group:          savedKey.GroupName,
			IsValid:        savedKey.IsValid,
		}

		sendNewKeyResult(resCh, key, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Key{}, ctx.Err()
		case result := <-resCh:
			return result.CreatedKey, result.Error
		}
	}
}

func sendNewKeyResult(resCh chan CreatingKeyResult, resp domain.Key, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- CreatingKeyResult{
		CreatedKey: resp,
		Error:      err,
	}
}
