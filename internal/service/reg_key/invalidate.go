package reg_key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
		zap.String("operation", InvalidateKeyOperation),
		zap.String("layer", "services"),
	)

	go func() {
		if deletingKey.KeyId <= 0 {
			l.Warn("incorrect key id", zap.Int("key id", deletingKey.KeyId))

			sendInvalidateKeyResult(resCh, registration_key.DeleteResp{}, "неправильный id")
			return
		}

		accountId := ctx.Value("AccountId").(int)

		role, err := s.accountMediator.RoleByAccountId(ctx, accountId)
		if err != nil {
			l.Warn("error get role by account id", zap.Error(err))

			sendInvalidateKeyResult(resCh, registration_key.DeleteResp{}, "ошибка проверки доступа")
			return
		}

		hasAccess, err := s.accountMediator.HasAccess(ctx, role.Id, ObjectName, InvalidateActionName)
		if err != nil || !hasAccess {
			l.Warn("ошибка при проверке прав", zap.Error(err))

			sendInvalidateKeyResult(resCh, registration_key.DeleteResp{}, ErrDontHavePermission.Error())
			return
		}

		err = s.r.Invalidate(ctx, deletingKey.KeyId)
		if err != nil {
			l.Warn("ошибка деактивации ключа", zap.Error(err))

			sendInvalidateKeyResult(resCh, registration_key.DeleteResp{}, "ошибка деактивации ключа")
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
