package reg_key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
)

const (
	InvalidateActionName = "DEL"
)

type InvalidateKeyResult struct {
	DeletedKey registration_key.DeleteResp
	Error      error
}

func (s Service) InvalidateKey(ctx context.Context, deletingKey registration_key.DeleteReq) (registration_key.DeleteResp, error) {
	resCh := make(chan InvalidateKeyResult)

	l := s.l.With(
		zap.String("операция", InvalidateKeyOperation),
		zap.String("слой", "services"),
	)

	go func() {
		if deletingKey.KeyId <= 0 {
			l.Warn("неправильный id ключа", zap.Int("id ключа", deletingKey.KeyId))

			sendInvalidateKeyResult(resCh, registration_key.DeleteResp{}, "Неправильный id")
			return
		}

		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.accountMediator.HasAccess(ctx, accountId, ObjectName, InvalidateActionName)
		if err != nil || !hasAccess {
			l.Warn("ошибка при проверке прав", zap.Error(err))

			sendInvalidateKeyResult(resCh, registration_key.DeleteResp{}, permissions.ErrDontHavePerm.Error())
			return
		}

		err = s.r.Invalidate(ctx, deletingKey.KeyId)
		if err != nil {
			l.Warn("ошибка деактивации ключа", zap.Error(err))

			sendInvalidateKeyResult(resCh, registration_key.DeleteResp{}, "Ошибка деактивации ключа")
			return
		}

		resp := registration_key.DeleteResp{KeyId: deletingKey.KeyId}

		sendInvalidateKeyResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return registration_key.DeleteResp{}, ctx.Err()
		case result := <-resCh:
			return result.DeletedKey, result.Error
		}
	}
}

func sendInvalidateKeyResult(resCh chan InvalidateKeyResult, resp registration_key.DeleteResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- InvalidateKeyResult{
		DeletedKey: resp,
		Error:      err,
	}
}
