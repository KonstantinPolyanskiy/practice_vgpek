package reg_key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/registration_key"
)

const (
	GetActionName = "GET"
)

type GetKeysResult struct {
	Keys  registration_key.GetKeysResp
	Error error
}

func (s Service) Keys(ctx context.Context, keyParams params.Key) (registration_key.GetKeysResp, error) {
	resCh := make(chan GetKeysResult)

	l := s.l.With(
		zap.String("operation", GetKeysOperation),
		zap.String("layer", "services"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		role, err := s.accountMediator.RoleByAccountId(ctx, accountId)
		if err != nil {
			sendGetKeysResult(resCh, registration_key.GetKeysResp{}, err.Error())
			return
		}

		hasAccess, err := s.accountMediator.HasAccess(ctx, role.Id, ObjectName, AddActionName)
		if err != nil || !hasAccess {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendGetKeysResult(resCh, registration_key.GetKeysResp{}, ErrDontHavePermission.Error())
			return
		}

		keys, err := s.r.KeysByParams(ctx, keyParams)
		if err != nil {
			l.Warn("error get keys by params",
				zap.Int("limit", keyParams.Limit),
				zap.Int("offset", keyParams.Offset),
				zap.Error(err),
			)

			sendGetKeysResult(resCh, registration_key.GetKeysResp{}, "ошибка получения ключей")
			return
		}

		resp := registration_key.GetKeysResp{Keys: keys}

		sendGetKeysResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return registration_key.GetKeysResp{}, ctx.Err()
		case result := <-resCh:
			return result.Keys, result.Error

		}
	}
}

func sendGetKeysResult(resCh chan GetKeysResult, resp registration_key.GetKeysResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetKeysResult{
		Keys:  resp,
		Error: err,
	}
}
