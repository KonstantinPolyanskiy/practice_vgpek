package reg_key

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/params"
)

type Service interface {
	NewKey(ctx context.Context, req dto.NewKeyReq) (domain.Key, error)
	InvalidateKey(ctx context.Context, id int) (domain.InvalidatedKey, error)

	KeysByParams(ctx context.Context, keyParams params.State) ([]domain.Key, error)
}

type Handler struct {
	l *zap.Logger
	s Service
}

func NewKeyHandler(service Service, logger *zap.Logger) Handler {
	return Handler{
		l: logger,
		s: service,
	}
}
