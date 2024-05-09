package person

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
)

type GetAccountResult struct {
	Account entity.Account
	Error   error
}

type GetAccountsResult struct {
	Accounts []entity.Account
	Error    error
}

func (s Service) EntityAccountById(ctx context.Context, req dto.EntityId) (entity.Account, error) {
	resCh := make(chan GetAccountResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.GetAccountOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		account, err := s.accountDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetAccountResult(resCh, entity.Account{}, "ошибка получения аккаунта")
			return
		}

		l.Info("получен аккаунт", zap.Int("id", account.Id))

		sendGetAccountResult(resCh, account, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return entity.Account{}, ctx.Err()
		case result := <-resCh:
			return result.Account, result.Error
		}
	}
}

func (s Service) EntityAccountByParam(ctx context.Context, p params.State) ([]entity.Account, error) {
	resCh := make(chan GetAccountsResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.GetAccountsByParamsOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		rawAccounts, err := s.accountDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendGetAccountsResult(resCh, nil, "ошибка получения аккаунтов")
			return
		}

		accounts := make([]entity.Account, 0)

		l.Info("запрос на получение аккаунтов", zap.String("состояние", p.State))
		switch p.State {
		case params.All:
			accounts = append(accounts, rawAccounts...)
		case params.Deleted:
			for _, account := range rawAccounts {
				if account.DeactivateTime != nil {
					accounts = append(accounts, account)
				}
			}
		case params.NotDeleted:
			for _, account := range rawAccounts {
				if account.DeactivateTime == nil {
					accounts = append(accounts, account)
				}
			}
		}

		sendGetAccountsResult(resCh, accounts, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Accounts, result.Error
		}
	}
}

func sendGetAccountsResult(resCh chan GetAccountsResult, resp []entity.Account, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetAccountsResult{
		Accounts: resp,
		Error:    err,
	}
}

func sendGetAccountResult(resCh chan GetAccountResult, resp entity.Account, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetAccountResult{
		Account: resp,
		Error:   err,
	}
}
