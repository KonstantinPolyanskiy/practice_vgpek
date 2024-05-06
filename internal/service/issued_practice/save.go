package issued_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/rndutils"
	"strings"
	"time"
)

type SavePracticeResult struct {
	SavedPractice domain.IssuedPractice
	Error         error
}

func (s Service) Save(ctx context.Context, req dto.NewIssuedPracticeReq) (domain.IssuedPractice, error) {
	resCh := make(chan SavePracticeResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.UploadIssuedPracticeOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		// Формируем название, добавляем в конце набор случайных символов для уникальности
		name := fmt.Sprintf("%s_%s", req.Title, rndutils.RandString(5))
		name = strings.Replace(name, " ", "_", -1)

		// Сохраняем файл практического задания
		savedPath, err := s.fileStorage.SaveFile(ctx, req.File, "issued", ".docx", name)
		if err != nil {
			l.Warn("возникла ошибка при сохранении файла", zap.Error(err))

			sendUploadPracticeResult(resCh, domain.IssuedPractice{}, "Не удалось сохранить файл")
			return
		}

		data := dto.NewIssuedPractice{
			AccountId:    accountId,
			TargetGroups: req.TargetGroups,
			Title:        req.Title,
			Theme:        req.Theme,
			Major:        req.Major,
			Path:         savedPath,
			UploadAt:     time.Now(),
		}

		savedPracticeData, err := s.issuedPracticeDAO.Save(ctx, data)
		if err != nil {
			sendUploadPracticeResult(resCh, domain.IssuedPractice{}, "Не удалось сохранить практическое задание")
			return
		}

		person, err := s.personDAO.ByAccountId(ctx, savedPracticeData.AccountId)
		if err != nil {
			sendUploadPracticeResult(resCh, domain.IssuedPractice{}, "Не удалось получить данные пользователя")
			return
		}

		var isDeleted bool

		if savedPracticeData.DeletedAt != nil {
			isDeleted = true
		}

		practice := domain.IssuedPractice{
			Id:           savedPracticeData.Id,
			AuthorName:   fmt.Sprintf("%s %s %s", person.LastName, person.FirstName, person.MiddleName),
			AuthorId:     savedPracticeData.AccountId,
			TargetGroups: savedPracticeData.TargetGroups,
			Title:        savedPracticeData.Title,
			Theme:        savedPracticeData.Theme,
			Major:        savedPracticeData.Major,
			Path:         savedPracticeData.Path,
			UploadAt:     savedPracticeData.UploadAt,
			IsDeleted:    isDeleted,
			DeletedAt:    savedPracticeData.DeletedAt,
		}

		sendUploadPracticeResult(resCh, practice, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.IssuedPractice{}, ctx.Err()
		case result := <-resCh:
			return result.SavedPractice, result.Error
		}
	}
}

func sendUploadPracticeResult(resCh chan SavePracticeResult, resp domain.IssuedPractice, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- SavePracticeResult{
		SavedPractice: resp,
		Error:         err,
	}
}
