package solved_practice

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"time"
)

type SetMarkResult struct {
	MarkedPractice domain.SolvedPractice
	Error          error
}

func (s Service) SetMark(ctx context.Context, req dto.MarkPracticeReq) (domain.SolvedPractice, error) {
	resCh := make(chan SetMarkResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.SetMarkSolvedPractice),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		accountId := ctx.Value("AccountId").(int)

		markTime := time.Now()

		markedPracticeEntity, err := s.solvedPracticeDAO.Update(ctx, entity.SolvedPracticeUpdate{
			Id:                 req.SolvedPracticeId,
			PerformedAccountId: nil,
			IssuedPracticeId:   nil,
			Mark:               &req.Mark,
			MarkTime:           &markTime,
			SolvedTime:         nil,
			Path:               nil,
			IsDeleted:          nil,
		})
		if err != nil {
			sendSetMarkResult(resCh, domain.SolvedPractice{}, "ошибка при обновлении практической работы")
			return
		}

		practice, err := s.EntityToDomain(ctx, accountId, markedPracticeEntity)
		if err != nil {
			l.Warn("возникла ошибка при переводе сущности БД в сущность логики", zap.Error(err))

			sendSetMarkResult(resCh, domain.SolvedPractice{}, "ошибка формирования практической работы")
			return
		}

		sendSetMarkResult(resCh, practice, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.SolvedPractice{}, ctx.Err()
		case result := <-resCh:
			return result.MarkedPractice, result.Error
		}
	}
}

func sendSetMarkResult(resCh chan SetMarkResult, resp domain.SolvedPractice, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- SetMarkResult{
		MarkedPractice: resp,
		Error:          err,
	}
}
