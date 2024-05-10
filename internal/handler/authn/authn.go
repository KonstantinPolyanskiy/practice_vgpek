package authn

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
)

type PersonService interface {
	// NewUser создает нового пользователя (Person и Account)
	NewUser(ctx context.Context, registration dto.RegistrationReq) (domain.Person, error)
	AccountById(ctx context.Context, req dto.EntityId) (domain.Account, error)
}

type TokenService interface {
	// CreateToken создает JWT токен с вшитым id аккаунта
	CreateToken(ctx context.Context, cred dto.Credentials) (string, error)

	// ParseToken возвращает ID аккаунта пользователя
	ParseToken(ctx context.Context, token string) (int, error)
}

type RBACService interface {
	RoleById(ctx context.Context, req dto.EntityId) (domain.Role, error)
}

type Handler struct {
	logger *zap.Logger

	personService PersonService
	tokenService  TokenService
	RBACService   RBACService
}

func NewAuthenticationHandler(personService PersonService, tokenService TokenService, rbacService RBACService, logger *zap.Logger) Handler {
	return Handler{
		logger:        logger,
		personService: personService,
		tokenService:  tokenService,
		RBACService:   rbacService,
	}
}
