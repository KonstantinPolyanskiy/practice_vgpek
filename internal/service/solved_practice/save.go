package solved_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/solved"
	"practice_vgpek/pkg/rndutils"
)

type SavePracticeResult struct {
	SavedPractice solved.UploadResp
	Error         error
}

func (s Service) Save(ctx context.Context, req solved.UploadReq) (solved.UploadResp, error) {
	resCh := make(chan SavePracticeResult)

	l := s.l.With(
		zap.String("операция", operation.UploadSolvedPracticeOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		// Проверка: есть ли доступ к загрузке практических работ
		hasAccess, err := s.am.HasAccess(ctx, accountId, SolvedObjectName, AddActionName)
		if err != nil {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendSavePracticeResult(resCh, solved.UploadResp{}, permissions.ErrCheckAccess.Error())
			return
		}

		// Проверка самого доступа
		if !hasAccess {
			l.Warn("попытка загрузить работу без прав", zap.Int("id аккаунта", accountId))

			sendSavePracticeResult(resCh, solved.UploadResp{}, permissions.ErrDontHavePerm.Error())
			return
		}

		// Проверка: совпадает ли целевая группа выполненного задания с группой студента
		groupMatch, err := s.ipm.IssuedGroupMatch(ctx, accountId, req.IssuedPracticeId)
		if err != nil {
			l.Warn("ошибка при проверке совпадания группы", zap.Error(err))

			sendSavePracticeResult(resCh, solved.UploadResp{}, "Ошибка при проверке целевой группы задания")
			return
		}

		if !groupMatch {
			l.Warn("попытка загрузить работу не с целевой группой", zap.Int("id аккаунта", accountId))

			sendSavePracticeResult(resCh, solved.UploadResp{}, "Некорректная целевая группа у практического задания")
			return
		}

		// Формируем случайное название практического задания
		// TODO: возможно стоит сделать его более осмысленным
		name := rndutils.RandString(10)

		// Сохраняем файл выполненной практической работы
		savedPath, err := s.fs.SaveFile(ctx, req.File, "solved", ".docx", name)
		if err != nil {
			l.Warn("возникла ошибка при сохранении файла", zap.Error(err))

			sendSavePracticeResult(resCh, solved.UploadResp{}, "Не удалось сохранить файл практического задания")
			return
		}

		dto := solved.DTO{
			IssuedPracticeId:   req.IssuedPracticeId,
			PerformedAccountId: accountId,
			Path:               savedPath,
		}

		savedPracticeInfo, err := s.spr.Save(ctx, dto)
		if err != nil {
			sendSavePracticeResult(resCh, solved.UploadResp{}, "Не удалось сохранить информацию о практическом задании")
			return
		}

		issuedPractice, err := s.ipr.ById(ctx, savedPracticeInfo.IssuedPracticeId)
		if err != nil {
			sendSavePracticeResult(resCh, solved.UploadResp{}, "Не удалось получить информаци по практическому заданию")
			return
		}

		resp := solved.UploadResp{
			SolvedPracticeId:   savedPracticeInfo.SolvedPracticeId,
			SolvedTime:         *savedPracticeInfo.SolvedTime,
			IssuedPracticeName: issuedPractice.Theme,
		}

		sendSavePracticeResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return solved.UploadResp{}, ctx.Err()
		case result := <-resCh:
			return result.SavedPractice, result.Error
		}
	}
}

func sendSavePracticeResult(resCh chan SavePracticeResult, resp solved.UploadResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- SavePracticeResult{
		SavedPractice: resp,
		Error:         err,
	}
}
