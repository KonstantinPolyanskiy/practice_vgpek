package reg_key

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/pkg/rndutils"
)

const (
	AddActionName = "ADD"
)

type CreatingKeyResult struct {
	CreatedKey registration_key.AddResp
	Error      error
}

func (s Service) NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.AddResp, error) {
	resCh := make(chan CreatingKeyResult)

	l := s.l.With(
		zap.String("операция", operation.NewKeyOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		// Проверяем что указанное количество использований ключа больше 0
		if req.MaxCountUsages <= 0 {
			l.Warn("неправильное макс. кол-во использований ключа", zap.Int("макс кол-во использований", req.MaxCountUsages))

			sendNewKeyResult(resCh, registration_key.AddResp{}, "Неправильное кол-во использований ключа")
			return
		}

		// Получаем ID аккаунта
		accountId := ctx.Value("AccountId").(int)

		// Получаем роль по id аккаунта
		role, err := s.accountMediator.RoleByAccountId(ctx, accountId)
		if err != nil {
			sendNewKeyResult(resCh, registration_key.AddResp{}, err.Error())
			return
		}

		// Проверяем, есть ли доступ
		hasAccess, err := s.accountMediator.HasAccess(ctx, role.Id, ObjectName, AddActionName)
		if err != nil || !hasAccess {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendNewKeyResult(resCh, registration_key.AddResp{}, permissions.ErrDontHavePerm.Error())
			return
		}

		// Формируем DTO
		dto := registration_key.DTO{
			RoleId:         req.RoleId,
			Body:           rndutils.RandNumberString(5) + rndutils.RandString(5),
			GroupName:      req.GroupName,
			MaxCountUsages: req.MaxCountUsages,
		}

		// Сохраняем ключ
		savedKey, err := s.r.SaveKey(ctx, dto)
		if err != nil {
			sendNewKeyResult(resCh, registration_key.AddResp{}, "Ошибка сохранения ключа")
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
