package solved_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

type GetPracticeResult struct {
	Practice domain.SolvedPractice
	Error    error
}

func (s Service) ById(ctx context.Context, req dto.EntityId) (domain.SolvedPractice, error) {
	// необходимо проверить id, кто запрашивает
	// если это студент и его целевая группа совпадает и id верен - отдаем ее,
	// в ином случае, если доступ есть - отдаем по id
	resCh := make(chan GetPracticeResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.GetSolvedPracticeInfoById),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		groupMatch, err := s.issuedPracticeMediator.IssuedGroupMatch(ctx, accountId, req.Id)
		if err != nil {
			l.Warn("возникла ошибка при проверке совпадений")
		}

		if !groupMatch {
			l.Warn("попытка получить практическую студентом с неправильной группой")

			sendGetPracticeResult(resCh, domain.SolvedPractice{}, "нет доступа к практическим группы")
			return
		}

		solvedPracticeEntity, err := s.solvedPracticeDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetPracticeResult(resCh, domain.SolvedPractice{}, "нет практической работы с таким id")
			return
		}

		practice, err := s.EntityToDomain(ctx, accountId, solvedPracticeEntity)
		if err != nil {
			l.Warn("возникла ошибка при переводе сущности БД в сущность логики", zap.Error(err))

			sendGetPracticeResult(resCh, domain.SolvedPractice{}, "ошибка формирования практической работы")
			return
		}

		sendGetPracticeResult(resCh, practice, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.SolvedPractice{}, ctx.Err()
		case result := <-resCh:
			return result.Practice, result.Error
		}
	}
}

func sendGetPracticeResult(resCh chan GetPracticeResult, practice domain.SolvedPractice, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetPracticeResult{
		Practice: practice,
		Error:    err,
	}
}
