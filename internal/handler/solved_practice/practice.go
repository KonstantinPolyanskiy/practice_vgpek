package solved_practice

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/practice/solved"
)

type SolvedPracticeService interface {
	Save(ctx context.Context, req solved.UploadReq) (solved.UploadResp, error)
	SetMark(ctx context.Context, req solved.SetMarkReq) (solved.SetMarkResp, error)

	ById(ctx context.Context, id int) (solved.Entity, error)
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
