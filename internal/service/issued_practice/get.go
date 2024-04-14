package issued_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/issued"
)

type GetPracticeResult struct {
	Practice issued.Entity
	Error    error
}

func (s Service) ById(ctx context.Context, id int) (issued.Entity, error) {
	// необходимо проверить id, кто запрашивает
	// если это студент и его целевая группа совпадает и id верен - отдаем ее,
	// в ином случае, если доступ есть - отдаем по id
	resCh := make(chan GetPracticeResult)

	l := s.l.With(
		zap.String("операция", operation.GetIssuedPracticeInfoById),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.am.HasAccess(ctx, accountId, IssuedObjectName, GetActionName)
		if err != nil {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendGetPracticeResult(resCh, issued.Entity{}, permissions.ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendGetPracticeResult(resCh, issued.Entity{}, permissions.ErrDontHavePerm.Error())
			return
		}

		groupMatch, err := s.pm.IssuedGroupMatch(ctx, accountId, id)
		if err != nil {
			l.Warn("возникла ошибка при проверке совпадений")
		}

		if !groupMatch {
			l.Warn("попытка получить практическую студентом с неправильной группой")

			sendGetPracticeResult(resCh, issued.Entity{}, "нет доступа к практическим группы")
			return
		}

		practice, err := s.r.ById(ctx, id)
		if err != nil {
			sendGetPracticeResult(resCh, issued.Entity{}, "нет практического задания с таким id")
			return
		}

		sendGetPracticeResult(resCh, practice, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return issued.Entity{}, ctx.Err()
		case result := <-resCh:
			return result.Practice, result.Error
		}
	}
}

func sendGetPracticeResult(resCh chan GetPracticeResult, practice issued.Entity, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetPracticeResult{
		Practice: practice,
		Error:    err,
	}
}
