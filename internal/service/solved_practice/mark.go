package solved_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
	"practice_vgpek/internal/model/practice/solved"
	"time"
)

type SetMarkResult struct {
	MarkedPractice solved.SetMarkResp
	Error          error
}

func (s Service) SetMark(ctx context.Context, req solved.SetMarkReq) (solved.SetMarkResp, error) {
	resCh := make(chan SetMarkResult)

	l := s.l.With(
		zap.String("операция", operation.SetMarkSolvedPractice),
		zap.String("слой", "сервисы"),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		hasAccess, err := s.am.HasAccess(ctx, accountId, MarkObjectName, AddActionName)
		if err != nil {
			l.Warn("возникла ошибка при проверке прав", zap.Error(err))

			sendSetMarkResult(resCh, solved.SetMarkResp{}, permissions.ErrCheckAccess.Error())
			return
		}

		if !hasAccess {
			l.Warn("попытка поставить оценку без прав", zap.Int("id аккаунта", accountId))

			sendSetMarkResult(resCh, solved.SetMarkResp{}, permissions.ErrDontHavePerm.Error())
			return
		}

		currentPractice, err := s.spr.ById(ctx, req.SolvedPracticeId)
		if err != nil {
			sendSetMarkResult(resCh, solved.SetMarkResp{}, "ошибка при получении практической работы")
			return
		}

		markTime := time.Now()

		updatedPractice := solved.Entity{
			SolvedPracticeId:   currentPractice.SolvedPracticeId,
			PerformedAccountId: currentPractice.PerformedAccountId,
			IssuedPracticeId:   currentPractice.IssuedPracticeId,
			Mark:               req.Mark,
			MarkTime:           &markTime,
			SolvedTime:         currentPractice.SolvedTime,
			Path:               currentPractice.Path,
			IsDeleted:          currentPractice.IsDeleted,
		}

		markedPractice, err := s.spr.Update(ctx, updatedPractice)
		if err != nil {
			sendSetMarkResult(resCh, solved.SetMarkResp{}, "ошибка при обновлении практической работы")
			return
		}

		resp := solved.SetMarkResp{
			SolvedPracticeId: markedPractice.SolvedPracticeId,
			Mark:             markedPractice.Mark,
			MarkTime:         *markedPractice.MarkTime,
		}

		sendSetMarkResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return solved.SetMarkResp{}, ctx.Err()
		case result := <-resCh:
			return result.MarkedPractice, result.Error
		}
	}
}

func sendSetMarkResult(resCh chan SetMarkResult, resp solved.SetMarkResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- SetMarkResult{
		MarkedPractice: resp,
		Error:          err,
	}
}
