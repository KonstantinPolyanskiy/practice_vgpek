package person

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/password"
	"time"
)

type NewUserResult struct {
	User  domain.Person
	Error error
}

func (s Service) NewUser(ctx context.Context, registration dto.RegistrationReq) (domain.Person, error) {
	resCh := make(chan NewUserResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.RegistrationOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		// Получаем ключ по указанному телу
		key, err := s.keyDAO.ByBody(ctx, registration.BodyKey)
		if err != nil {
			sendRegistrationResult(resCh, domain.Person{}, "Ошибка получения ключа регистрации")
			return
		}

		// Проверяем, валиден ли он
		if !key.IsValid {
			l.Warn("попытка зарегистрироваться по невалидному ключу", zap.Int("id ключа", key.Id))
			sendRegistrationResult(resCh, domain.Person{}, "Невалидный ключ регистрации")
			return
		}

		// Если текущее кол-во регистраций больше или равно допустимому - инвалидируем
		if key.CurrentCountUsages >= key.MaxCountUsages {
			l.Warn("превышено кол-во попыток регистрации по ключу", zap.Int("id ключа", key.Id))
			_, err = s.keyService.InvalidateKey(ctx, dto.EntityId{Id: key.Id})
			if err != nil {
				sendRegistrationResult(resCh, domain.Person{}, "Ошибка инвалидирования ключа регистрации")
				return
			}
			sendRegistrationResult(resCh, domain.Person{}, "Превышено кол-во регистраций по ключу")
			return
		}

		pHash, err := password.Hash(registration.Password)
		if err != nil {
			l.Warn("ошибка хеширования пароля", zap.Error(err))
			sendRegistrationResult(resCh, domain.Person{}, "Ошибка хеширования пароля")
			return
		}

		// Сохраняем сущность Аккаунт
		accountEntity, err := s.accountDAO.Save(ctx, dto.AccountRegistrationData{
			Login:        registration.Login,
			PasswordHash: pHash,
			CreatedAt:    time.Now(),
			RoleId:       key.RoleId,
			KeyId:        key.Id,
		})
		if err != nil {
			sendRegistrationResult(resCh, domain.Person{}, "Ошибка создания пользователя")
			return
		}

		// Сохраняем сущность Пользователь
		personEntity, err := s.personDAO.Save(ctx, dto.PersonRegistrationData{
			UUID:       uuid.New(),
			FirstName:  registration.FirstName,
			SecondName: registration.SecondName,
			LastName:   registration.LastName,
			AccountId:  accountEntity.Id,
		})
		if err != nil {
			accErr := s.accountDAO.HardDeleteById(ctx, accountEntity.Id)
			if accErr != nil {
				l.Warn("ошибка удаления аккаунта", zap.Error(accErr))
			}
			sendRegistrationResult(resCh, domain.Person{}, "Ошибка создания пользователя")
			return
		}

		// Если ключ хороший и регистрация успешная, увеличиваем кол-во регистраций по нему
		_, err = s.keyService.Increment(ctx, key)
		if err != nil {
			sendRegistrationResult(resCh, domain.Person{}, err.Error())
			return
		}

		roleEntity, err := s.roleDAO.ById(ctx, accountEntity.RoleId)
		if err != nil {
			sendRegistrationResult(resCh, domain.Person{}, "Ошибка получения роли пользователя")
			return
		}

		person := domain.Person{
			UUID:       personEntity.UUID,
			FirstName:  personEntity.FirstName,
			MiddleName: personEntity.MiddleName,
			LastName:   personEntity.LastName,
			Account: domain.Account{
				Login:          accountEntity.Login,
				IsActive:       accountEntity.IsActive,
				DeactivateTime: accountEntity.DeactivateTime,
				RoleName:       roleEntity.Name,
				RoleId:         roleEntity.Id,
				KeyId:          key.Id,
				CreatedAt:      accountEntity.CreatedAt,
			},
		}

		sendRegistrationResult(resCh, person, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return domain.Person{}, ctx.Err()
		case result := <-resCh:
			return result.User, result.Error
		}
	}

}

func sendRegistrationResult(resCh chan NewUserResult, resp domain.Person, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- NewUserResult{
		User:  resp,
		Error: err,
	}
}
