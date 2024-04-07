package issued_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/issued"
	"time"
)

type SavePracticeResult struct {
	SavedPractice issued.UploadResp
	Error         error
}

const (
	AddActionName = "ADD"
	ObjectName    = "ISSUED PRACTICE"
)

func (s Service) Save(ctx context.Context, req issued.UploadReq) (issued.UploadResp, error) {
	resCh := make(chan SavePracticeResult)

	l := s.l.With(
		zap.String("операция", operation.UploadIssuedPracticeOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.am.HasAccess(ctx, accountId, ObjectName, AddActionName)
		if err != nil || !hasAccess {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendUploadPracticeResult(resCh, issued.UploadResp{}, permissions.ErrDontHavePerm.Error())
			return
		}

		savedPath, err := s.fs.SaveFile(req.File, "/issued/", ".docx", "test")
		if err != nil || !hasAccess {
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
			DeletedAt:    nil,
		}

		savedPracticeData, err := s.r.Save
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
