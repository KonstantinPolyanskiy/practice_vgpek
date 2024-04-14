package issued_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/issued"
	"practice_vgpek/pkg/rndutils"
	"strings"
	"time"
)

type SavePracticeResult struct {
	SavedPractice issued.UploadResp
	Error         error
}

func (s Service) Save(ctx context.Context, req issued.UploadReq) (issued.UploadResp, error) {
	resCh := make(chan SavePracticeResult)

	l := s.l.With(
		zap.String("операция", operation.UploadIssuedPracticeOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.am.HasAccess(ctx, accountId, IssuedObjectName, AddActionName)
		if err != nil || !hasAccess {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendUploadPracticeResult(resCh, issued.UploadResp{}, permissions.ErrDontHavePerm.Error())
			return
		}

		// Формируем название, добавляем в конце набор случайных символов для уникальности
		name := fmt.Sprintf("%s_%s", req.Title, rndutils.RandString(5))
		name = strings.Replace(name, " ", "_", -1)

		savedPath, err := s.fs.SaveFile(ctx, req.File, "issued", ".docx", name)
		if err != nil {
			l.Warn("возникла ошибка при сохрании файла", zap.Error(err))

			sendUploadPracticeResult(resCh, issued.UploadResp{}, "Не удалось сохранить файл")
			return
		}

		dto := issued.DTO{
			AccountId:    accountId,
			TargetGroups: req.TargetGroups,
			Title:        req.Title,
			Theme:        req.Theme,
			Major:        req.Major,
			Path:         savedPath,
			UploadAt:     time.Now(),
			DeletedAt:    &time.Time{},
		}

		savedPracticeData, err := s.r.Save(ctx, dto)
		if err != nil {
			sendUploadPracticeResult(resCh, issued.UploadResp{}, "Не удалось сохранить практическое задание")
			return
		}

		resp := issued.UploadResp{
			PracticeId:   savedPracticeData.PracticeId,
			Title:        savedPracticeData.Title,
			TargetGroups: savedPracticeData.TargetGroups,
			UploadAt:     savedPracticeData.UploadAt,
		}

		sendUploadPracticeResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return issued.UploadResp{}, ctx.Err()
		case result := <-resCh:
			return result.SavedPractice, result.Error
		}
	}
}

func sendUploadPracticeResult(resCh chan SavePracticeResult, resp issued.UploadResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- SavePracticeResult{
		SavedPractice: resp,
		Error:         err,
	}
}
