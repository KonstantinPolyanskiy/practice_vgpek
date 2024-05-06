package solved_practice

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
)

type SolvedPracticeService interface {
	Save(ctx context.Context, req dto.NewSolvedPracticeReq) (domain.SolvedPractice, error)
	SetMark(ctx context.Context, req dto.MarkPracticeReq) (domain.SolvedPractice, error)

	ById(ctx context.Context, req dto.EntityId) (domain.SolvedPractice, error)
}

type Handler struct {
	l *zap.Logger
	s SolvedPracticeService
}

func NewCompletedPracticeHandler(service SolvedPracticeService, logger *zap.Logger) Handler {
	return Handler{
		l: logger,
		s: service,
	}
}
