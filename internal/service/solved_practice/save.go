package solved_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/rndutils"
	"time"
)

type SavePracticeResult struct {
	SavedPractice domain.SolvedPractice
	Error         error
}

func (s Service) Save(ctx context.Context, req dto.NewSolvedPracticeReq) (domain.SolvedPractice, error) {
	resCh := make(chan SavePracticeResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.UploadSolvedPracticeOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		// Проверка: совпадает ли целевая группа выполненного задания с группой студента
		groupMatch, err := s.issuedPracticeMediator.IssuedGroupMatch(ctx, accountId, req.IssuedPracticeId)
		if err != nil {
			l.Warn("ошибка при проверке совпадения группы", zap.Error(err))

			sendSavePracticeResult(resCh, domain.SolvedPractice{}, "Ошибка при проверке целевой группы задания")
			return
		}

		if !groupMatch {
			l.Warn("попытка загрузить работу не с целевой группой", zap.Int("id аккаунта", accountId))

			sendSavePracticeResult(resCh, domain.SolvedPractice{}, "Некорректная целевая группа у практического задания")
			return
		}

		// Формируем случайное название практического задания
		// TODO: возможно стоит сделать его более осмысленным
		name := rndutils.RandString(10)

		// Сохраняем файл выполненной практической работы
		savedPath, err := s.fileStorage.SaveFile(ctx, req.File, "solved", ".docx", name)
		if err != nil {
			l.Warn("возникла ошибка при сохранении файла", zap.Error(err))

			sendSavePracticeResult(resCh, domain.SolvedPractice{}, "Не удалось сохранить файл практического задания")
			return
		}

		solvedTime := time.Now()

		data := dto.NewSolvedPractice{
			PerformedAccountId: accountId,
			IssuedPracticeId:   req.IssuedPracticeId,
			Mark:               0,
			MarkTime:           nil,
			SolvedTime:         &solvedTime,
			Path:               savedPath,
			IsDeleted:          nil,
		}

		savedPracticeEntity, err := s.solvedPracticeDAO.Save(ctx, data)
		if err != nil {
			sendSavePracticeResult(resCh, domain.SolvedPractice{}, "Не удалось сохранить информацию о практическом задании")
			return
		}

		practice, err := s.EntityToDomain(ctx, accountId, savedPracticeEntity)
		if err != nil {
			l.Warn("возникла ошибка при переводе сущности БД в сущность логики", zap.Error(err))

			sendSavePracticeResult(resCh, domain.SolvedPractice{}, "ошибка формирования практической работы")
			return
		}

		sendSavePracticeResult(resCh, practice, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.SolvedPractice{}, ctx.Err()
		case result := <-resCh:
			return result.SavedPractice, result.Error
		}
	}
}

func sendSavePracticeResult(resCh chan SavePracticeResult, resp domain.SolvedPractice, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- SavePracticeResult{
		SavedPractice: resp,
		Error:         err,
	}
}
