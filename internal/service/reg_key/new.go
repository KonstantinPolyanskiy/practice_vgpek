package reg_key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/pkg/rndutils"
)

type CreatingKeyResult struct {
	CreatedKey registration_key.AddResp
	Error      error
}

func (s Service) NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.AddResp, error) {
	resCh := make(chan CreatingKeyResult)

	l := s.l.With(
		zap.String("action", NewKeyAction),
		zap.String("layer", "services"),
	)

	go func() {
		// Проверяем что указанное количество использований ключа больше 0
		if req.MaxCountUsages < 0 {
			l.Warn("incorrect max count usages", zap.Int("body key", req.MaxCountUsages))

			sendNewKeyResult(resCh, registration_key.AddResp{}, "неправильное кол-во использований ключа")
			return
		}

		// Формируем DTO
		dto := registration_key.DTO{
			RoleId:         req.RoleId,
			Body:           rndutils.RandNumberString(5) + rndutils.RandString(5),
			MaxCountUsages: req.MaxCountUsages,
		}

		// Сохраняем ключ
		savedKey, err := s.r.SaveKey(ctx, dto)
		if err != nil {
			l.Warn("error save key in db",
				zap.String("body", dto.Body),
				zap.Int("max count usages", dto.MaxCountUsages),
				zap.Int("role id", dto.RoleId),
			)

			sendNewKeyResult(resCh, registration_key.AddResp{}, "ошибка сохранения ключа")
			return
		}

		// Формируем и отправляем ответ
		resp := registration_key.AddResp{
			RegKeyId:           savedKey.RegKeyId,
			MaxCountUsages:     savedKey.MaxCountUsages,
			CurrentCountUsages: savedKey.CurrentCountUsages,
			Body:               savedKey.Body,
			CreatedAt:          savedKey.CreatedAt,
		}

		sendNewKeyResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return registration_key.AddResp{}, ctx.Err()
		case result := <-resCh:
			return result.CreatedKey, result.Error
		}
	}
}

func sendNewKeyResult(resCh chan CreatingKeyResult, resp registration_key.AddResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- CreatingKeyResult{
		CreatedKey: resp,
		Error:      err,
	}
}
