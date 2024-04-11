package practice

import (
	"context"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/practice/issued"
	"practice_vgpek/internal/model/registration_key"
)

type AccountRepo interface {
	AccountById(ctx context.Context, id int) (account.Entity, error)
}

type IssuedPracticeRepo interface {
	ById(ctx context.Context, id int) (issued.Entity, error)
}

type KeyRepo interface {
	RegKeyById(ctx context.Context, id int) (registration_key.Entity, error)
}

type Mediator struct {
	accountRepo AccountRepo
	issuedRepo  IssuedPracticeRepo
	keyRepo     KeyRepo
}

func NewPracticeMediator(accountRepo AccountRepo, issuedRepo IssuedPracticeRepo, keyRepo KeyRepo) Mediator {
	return Mediator{
		accountRepo: accountRepo,
		issuedRepo:  issuedRepo,
		keyRepo:     keyRepo,
	}
}

func (m Mediator) IssuedGroupMatch(ctx context.Context, accountId, practiceId int) (bool, error) {
	var match bool

	acc, err := m.accountRepo.AccountById(ctx, accountId)
	if err != nil {
		return false, err
	}

	practice, err := m.issuedRepo.ById(ctx, practiceId)
	if err != nil {
		return false, err
	}

	key, err := m.keyRepo.RegKeyById(ctx, acc.RegKeyId)
	if err != nil {
		return false, err
	}

	for _, group := range practice.TargetGroups {
		if group == key.GroupName {
			match = true
			break
		}
	}

	return match, nil
}
