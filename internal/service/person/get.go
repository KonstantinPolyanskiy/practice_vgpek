package person

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
)

type GetAccountEntityResult struct {
	Account entity.Account
	Error   error
}

type GetAccountsEntityResult struct {
	Accounts []entity.Account
	Error    error
}

type GetPersonsEntityResult struct {
	Persons []entity.Person
	Error   error
}

type GetAccountResult struct {
	Account domain.Account
	Error   error
}

func (s Service) AccountById(ctx context.Context, req dto.EntityId) (domain.Account, error) {
	resCh := make(chan GetAccountResult)

	_ = s.logger.With(
		zap.String(operation.Operation, operation.GetAccountOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		account, err := s.accountDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetAccountResult(resCh, domain.Account{}, "ошибка получения аккаунта")
			return
		}

		role, err := s.roleDAO.ById(ctx, account.RoleId)
		if err != nil {
			sendGetAccountResult(resCh, domain.Account{}, "ошибка получения роли")
			return
		}

		var isDeleted bool

		if account.DeactivateTime != nil {
			isDeleted = true
		}

		acc := domain.Account{
			Login:          account.Login,
			IsActive:       isDeleted,
			DeactivateTime: account.DeactivateTime,
			RoleName:       role.Name,
			RoleId:         role.Id,
			KeyId:          account.KeyId,
			CreatedAt:      account.CreatedAt,
		}

		sendGetAccountResult(resCh, acc, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Account{}, ctx.Err()
		case result := <-resCh:
			return result.Account, result.Error
		}
	}
}

func (s Service) EntityAccountById(ctx context.Context, req dto.EntityId) (entity.Account, error) {
	resCh := make(chan GetAccountEntityResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.GetAccountOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		account, err := s.accountDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetAccountEntityResult(resCh, entity.Account{}, "ошибка получения аккаунта")
			return
		}

		l.Info("получен аккаунт", zap.Int("id", account.Id))

		sendGetAccountEntityResult(resCh, account, "")
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
	resCh := make(chan GetAccountsEntityResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.GetAccountsByParamsOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		rawAccounts, err := s.accountDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendGetAccountsEntityResult(resCh, nil, "ошибка получения аккаунтов")
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

		sendGetAccountsEntityResult(resCh, accounts, "")
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

func (s Service) EntityPersonByParam(ctx context.Context, p params.State) ([]entity.Person, error) {
	resCh := make(chan GetPersonsEntityResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.GetPersonsByParams),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		l.Info("параметры запроса",
			zap.String("статус", p.State),
			zap.Int("лимит", p.Default.Limit),
			zap.Int("оффсет", p.Default.Offset),
		)

		rawPersons, err := s.personDAO.ByParams(ctx, p.Default)
		if err != nil {
			sendGetPersonsEntityResult(resCh, nil, "ошибка получения пользователей")
			return
		}

		sendGetPersonsEntityResult(resCh, rawPersons, "")
		return

	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resCh:
			return result.Persons, result.Error
		}
	}
}

func sendGetPersonsEntityResult(resCh chan GetPersonsEntityResult, resp []entity.Person, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetPersonsEntityResult{
		Persons: resp,
		Error:   err,
	}
}

func sendGetAccountsEntityResult(resCh chan GetAccountsEntityResult, resp []entity.Account, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetAccountsEntityResult{
		Accounts: resp,
		Error:    err,
	}
}

func sendGetAccountEntityResult(resCh chan GetAccountEntityResult, resp entity.Account, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetAccountEntityResult{
		Account: resp,
		Error:   err,
	}
}

func sendGetAccountResult(resCh chan GetAccountResult, resp domain.Account, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetAccountResult{
		Account: resp,
		Error:   err,
	}
}
