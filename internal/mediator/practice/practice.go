package practice

import (
	"context"
	"practice_vgpek/internal/model/entity"
)

type AccountDAO interface {
	ById(ctx context.Context, id int) (entity.Account, error)
}

type IssuedPracticeDAO interface {
	ById(ctx context.Context, id int) (entity.IssuedPractice, error)
}

type KeyDAO interface {
	ById(ctx context.Context, id int) (entity.Key, error)
}

type Mediator struct {
	accountDAO        AccountDAO
	issuedPracticeDAO IssuedPracticeDAO
	keyDAO            KeyDAO
}

func NewIssuedPracticeMediator(accountDAO AccountDAO, issuedPracticeDAO IssuedPracticeDAO, keyDAO KeyDAO) Mediator {
	return Mediator{
		accountDAO:        accountDAO,
		issuedPracticeDAO: issuedPracticeDAO,
		keyDAO:            keyDAO,
	}
}

// IssuedGroupMatch проверяет, что группа ключа совпадает с группой практики
func (m Mediator) IssuedGroupMatch(ctx context.Context, accountId, practiceId int) (bool, error) {
	var match bool

	acc, err := m.accountDAO.ById(ctx, accountId)
	if err != nil {
		return false, err
	}

	practice, err := m.issuedPracticeDAO.ById(ctx, practiceId)
	if err != nil {
		return false, err
	}

	key, err := m.keyDAO.ById(ctx, acc.KeyId)
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
