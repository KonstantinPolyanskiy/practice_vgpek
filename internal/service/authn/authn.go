package authn

import (
	"context"
	"fmt"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/pkg/password"
)

type Repository interface {
	SavePerson(ctx context.Context, person person.DTO) (person.Entity, error)
}

type KeyRepository interface {
	RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error)
}

type Service struct {
	r  Repository
	kr KeyRepository
}

func NewAuthenticationService(repository Repository, keyRepository KeyRepository) Service {
	return Service{
		r:  repository,
		kr: keyRepository,
	}
}

type RegistrationResult struct {
	RegisteredPerson person.RegisteredResp
	Error            error
}

func (s Service) NewPerson(ctx context.Context, registering person.RegistrationReq) (person.RegisteredResp, error) {
	resCh := make(chan RegistrationResult)

	go func() {
		// Хешируем пароль
		passwordHash, err := password.Hash(registering.Password)
		if err != nil {
			resCh <- RegistrationResult{
				RegisteredPerson: person.RegisteredResp{},
				Error:            fmt.Errorf("ошибка в хешировании пароля - %s\n", err.Error()),
			}
		}

		regKey, err := s.kr.RegKeyByBody(ctx, registering.RegistrationKey)
		if err != nil {
			resCh <- RegistrationResult{
				RegisteredPerson: person.RegisteredResp{},
				Error:            fmt.Errorf("ошибка с ключем регистрации"),
			}
		}

		// Формируем DTO
		dto := person.DTO{
			Personality: person.Personality{
				FirstName:  registering.FirstName,
				MiddleName: registering.MiddleName,
				LastName:   registering.LastName,
			},
			Account: account.DTO{
				Login:        registering.Login,
				PasswordHash: passwordHash,
				RoleId:       regKey.RoleId,
				RegKeyId:     regKey.RegKeyId,
			},
		}

		// Сохраняем регистируемого пользователя в БД
		savedPerson, err := s.r.SavePerson(ctx, dto)
		if err != nil {
			resCh <- RegistrationResult{
				RegisteredPerson: person.RegisteredResp{},
				Error:            err,
			}
		}

		// Формируем ответ сервиса
		registeredPerson := person.RegisteredResp{
			Personality: person.Personality{
				FirstName:  savedPerson.FirstName,
				MiddleName: savedPerson.MiddleName,
				LastName:   savedPerson.LastName,
			},
			CreatedAt: savedPerson.CreatedAt,
		}

		// Кладем ответ в канал
		resCh <- RegistrationResult{
			RegisteredPerson: registeredPerson,
			Error:            nil,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return person.RegisteredResp{}, ctx.Err()
		case result := <-resCh:
			return result.RegisteredPerson, result.Error
		}
	}

}
