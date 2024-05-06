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
}

type TokenService interface {
	// CreateToken создает JWT токен с вшитым id аккаунта
	CreateToken(ctx context.Context, cred dto.Credentials) (string, error)

	// ParseToken возвращает ID аккаунта пользователя
	ParseToken(ctx context.Context, token string) (int, error)
}

type Handler struct {
	logger        *zap.Logger
	personService PersonService
	tokenService  TokenService
}

func NewAuthenticationHandler(personService PersonService, tokenService TokenService, logger *zap.Logger) Handler {
	return Handler{
		logger:        logger,
		personService: personService,
		tokenService:  tokenService,
	}
}
