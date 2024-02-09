package authn

import (
	"context"
	"fmt"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/pkg/password"
)

type Repository interface {
	SavePerson(ctx context.Context, person person.DTO) (person.Entity, error)
}

type Service struct {
	r Repository
}

func NewAuthenticationService(repository Repository) Service {
	return Service{
		r: repository,
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

		// Формируем DTO
		dto := person.DTO{
			Personality: person.Personality{
				FirstName:  registering.FirstName,
				MiddleName: registering.MiddleName,
				LastName:   registering.LastName,
			},
			Login:        registering.Login,
			PasswordHash: passwordHash,
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
