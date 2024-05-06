package issued_practice

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
	Practice domain.IssuedPractice
	Error    error
}

func (s Service) ById(ctx context.Context, req dto.EntityId) (domain.IssuedPractice, error) {
	// необходимо проверить id, кто запрашивает
	// если это студент и его целевая группа совпадает и id верен - отдаем ее,
	// в ином случае, если доступ есть - отдаем по id
	resCh := make(chan GetPracticeResult)

	_ = s.logger.With(
		zap.String(operation.Operation, operation.GetIssuedPracticeInfoById),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		practiceEntity, err := s.issuedPracticeDAO.ById(ctx, req.Id)
		if err != nil {
			sendGetPracticeResult(resCh, domain.IssuedPractice{}, "нет практического задания с таким id")
			return
		}

		person, err := s.personDAO.ByAccountId(ctx, practiceEntity.AccountId)
		if err != nil {
			sendGetPracticeResult(resCh, domain.IssuedPractice{}, "ошибка получения автора задания")
			return
		}

		var isDeleted bool

		if practiceEntity.DeletedAt != nil {
			isDeleted = true
		}

		practice := domain.IssuedPractice{
			Id:           practiceEntity.Id,
			AuthorName:   fmt.Sprintf("%s %s %s", person.LastName, person.FirstName, person.MiddleName),
			AuthorId:     practiceEntity.AccountId,
			TargetGroups: practiceEntity.TargetGroups,
			Title:        practiceEntity.Title,
			Theme:        practiceEntity.Theme,
			Major:        practiceEntity.Major,
			Path:         practiceEntity.Path,
			UploadAt:     practiceEntity.UploadAt,
			IsDeleted:    isDeleted,
			DeletedAt:    practiceEntity.DeletedAt,
		}

		sendGetPracticeResult(resCh, practice, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.IssuedPractice{}, ctx.Err()
		case result := <-resCh:
			return result.Practice, result.Error
		}
	}
}

func sendGetPracticeResult(resCh chan GetPracticeResult, practice domain.IssuedPractice, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- GetPracticeResult{
		Practice: practice,
		Error:    err,
	}
}
