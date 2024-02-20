package reg_key

import (
	"context"
	"practice_vgpek/internal/model/registration_key"
)

type Repository interface {
	SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error)
}

type Service struct {
	r Repository
}

func NewKeyService(repository Repository) Service {
	return Service{
		r: repository,
	}
}
