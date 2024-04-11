package issued_practice

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/practice/issued"
)

type IssuedPracticeService interface {
	Save(ctx context.Context, req issued.UploadReq) (issued.UploadResp, error)
	ById(ctx context.Context, id int) (issued.Entity, error)
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
