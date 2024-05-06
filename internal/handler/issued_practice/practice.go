package issued_practice

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
)

type IssuedPracticeService interface {
	Save(ctx context.Context, req dto.NewIssuedPracticeReq) (domain.IssuedPractice, error)
	ById(ctx context.Context, req dto.EntityId) (domain.IssuedPractice, error)
}

type Handler struct {
	l *zap.Logger
	s IssuedPracticeService
}

func NewIssuedPracticeHandler(service IssuedPracticeService, logger *zap.Logger) Handler {
	return Handler{
		s: service,
		l: logger,
	}
}
