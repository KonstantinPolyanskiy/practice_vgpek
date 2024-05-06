package person

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
)

type KeyDAO interface {
	ByBody(ctx context.Context, body string) (entity.Key, error)
	Update(ctx context.Context, key entity.Key) (entity.Key, error)
}

type AccountMediator interface {
	PermByAccountId(ctx context.Context, id int) (domain.RolePermission, error)
}

type KeyService interface {
	InvalidateKey(ctx context.Context, id int) (domain.InvalidatedKey, error)
	Increment(ctx context.Context, key entity.Key) (entity.Key, error)
}

type RoleDAO interface {
	ById(ctx context.Context, id int) (entity.Role, error)
}

type PersonDAO interface {
	Save(ctx context.Context, data dto.PersonRegistrationData) (entity.Person, error)
}

type AccountDAO interface {
	Save(ctx context.Context, data dto.AccountRegistrationData) (entity.Account, error)
	HardDeleteById(ctx context.Context, id int) error
}

type Service struct {
	logger *zap.Logger

	keyDAO     KeyDAO
	personDAO  PersonDAO
	accountDAO AccountDAO
	roleDAO    RoleDAO

	keyService KeyService

	accountMediator AccountMediator
}

func New(kd KeyDAO, pd PersonDAO, ad AccountDAO, rd RoleDAO, keyService KeyService, am AccountMediator, logger *zap.Logger) Service {
	return Service{
		logger: logger,

		keyDAO:     kd,
		personDAO:  pd,
		accountDAO: ad,
		roleDAO:    rd,

		keyService: keyService,

		accountMediator: am,
	}
}
