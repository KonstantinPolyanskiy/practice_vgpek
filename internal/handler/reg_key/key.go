package reg_key

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/registration_key"
)

type Service interface {
	NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.AddResp, error)
	InvalidateKey(ctx context.Context, req registration_key.DeleteReq) (registration_key.DeleteResp, error)

	KeysByParams(ctx context.Context, keyParams params.Key) ([]registration_key.Entity, error)
}

type Handler struct {
	l *zap.Logger
	s Service
}

func NewRegKeyHandler(service Service, logger *zap.Logger) Handler {
	return Handler{
		l: logger,
		s: service,
	}
}
