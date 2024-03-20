package authn

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
	accountRepo "practice_vgpek/internal/repository/account"
	"practice_vgpek/pkg/password"
	"time"
)

const testKey = "uvm23458IKT345fg4TB3"

type PersonRepository interface {
	SavePerson(ctx context.Context, savingPerson person.DTO, accountId int) (person.Entity, error)
}

type KeyRepository interface {
	RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error)
	IncCountUsages(ctx context.Context, keyId int) error
	Invalidate(ctx context.Context, keyId int) error
}

type AccountRepository interface {
	SaveAccount(ctx context.Context, savingAcc account.DTO) (account.Entity, error)
	AccountByLogin(ctx context.Context, login string) (account.Entity, error)
}

type Service struct {
	l  *zap.Logger
	pr PersonRepository
	kr KeyRepository
	ar AccountRepository
}

func NewAuthenticationService(
	repository PersonRepository,
	accountRepository AccountRepository,
	keyRepository KeyRepository,
	logger *zap.Logger) Service {
	return Service{
		l:  logger,
		pr: repository,
		kr: keyRepository,
		ar: accountRepository,
	}
}

type RegistrationResult struct {
	RegisteredPerson person.RegisteredResp
	Error            error
}

type LogInResult struct {
	CreatedToken person.LogInResp
	Error        error
}

func (s Service) NewPerson(ctx context.Context, registering person.RegistrationReq) (person.RegisteredResp, error) {
	resCh := make(chan RegistrationResult)

	l := s.l.With(
		zap.String("action", RegistrationOperation),
		zap.String("layer", "services"),
	)

	go func() {
		// Получаем ключ, по которому зарегистрированн пользователь
		regKey, err := s.kr.RegKeyByBody(ctx, registering.RegistrationKey)
		if err != nil {
			l.Warn("body key error", zap.String("body key", registering.RegistrationKey))
			sendRegistrationResult(resCh, person.RegisteredResp{}, "ошибка с ключем регистрации")
			return
		}

		// Проверка, что ключ валиден, если нет - возвращаем ошибку
		if !regKey.IsValid {
			sendRegistrationResult(resCh, person.RegisteredResp{}, "невалидный ключ")
			return
		}

		// Проверяем, что ключ еще можно использовать, если нет - инвалидируем
		if regKey.CurrentCountUsages >= regKey.MaxCountUsages {
			if err = s.kr.Invalidate(ctx, regKey.RegKeyId); err != nil {
				l.Warn("invalidate key error",
					zap.String("body", regKey.Body),
					zap.Int("key id", regKey.RegKeyId),
					zap.Bool("is valid", regKey.IsValid),
				)

				sendRegistrationResult(resCh, person.RegisteredResp{}, "ошибка деактивирования ключа")
				return
			}
		}

		// Хешируем пароль
		passwordHash, err := password.Hash(registering.Password)
		if err != nil {
			l.Warn("hashing password error",
				zap.String("password", registering.Password),
				zap.Error(err),
			)

			sendRegistrationResult(resCh, person.RegisteredResp{}, "ошибка хеширования пароля")
			return
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

		// Сохраняем регистрируеммый аккаунт пользователя в БД
		savedAcc, err := s.ar.SaveAccount(ctx, dto.Account)
		if err != nil {
			var errMsg string

			l.Warn("error save account in db", zap.String("user login", dto.Account.Login))

			// Проверяем, является ли полученная ошибка - ошибкой сохранения аккаунта
			if errors.Is(err, accountRepo.ErrLoginAlreadyExist) {
				errMsg = "такой логин уже существует"
			} else {
				errMsg = "неизвестная ошибка сохранения аккаунта"
			}

			sendRegistrationResult(resCh, person.RegisteredResp{}, errMsg)
			return
		}

		// Если аккаунт создан, увеличиваем кол-во регистраций по ключу
		err = s.kr.IncCountUsages(ctx, regKey.RegKeyId)
		if err != nil {
			l.Warn("error inc count key",
				zap.Int("key id", regKey.RegKeyId),
				zap.Int("current count", regKey.CurrentCountUsages),
			)

			sendRegistrationResult(resCh, person.RegisteredResp{}, "ошибка обновления ключа")
			return
		}

		// Сохраняем регистируемого пользователя в БД
		savedPerson, err := s.pr.SavePerson(ctx, dto, savedAcc.AccountId)
		if err != nil {
			l.Warn("error save person in db",
				zap.String("full name", dto.FirstName+" "+dto.MiddleName+" "+dto.LastName),
				zap.Int("account id", savedAcc.AccountId),
			)
			sendRegistrationResult(resCh, person.RegisteredResp{}, "ошибка сохранения пользователя")
			return
		}

		// Формируем ответ сервиса
		registeredPerson := person.RegisteredResp{
			Personality: person.Personality{
				FirstName:  savedPerson.FirstName,
				MiddleName: savedPerson.MiddleName,
				LastName:   savedPerson.LastName,
			},
			CreatedAt: savedAcc.CreatedAt,
		}

		// Кладем ответ в канал
		sendRegistrationResult(resCh, registeredPerson, "")
		return
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

func (s Service) NewToken(ctx context.Context, logIn person.LogInReq) (person.LogInResp, error) {
	resCh := make(chan LogInResult)

	l := s.l.With(
		zap.String("action", LoginOperation),
		zap.String("layer", "services"),
	)

	go func() {
		acc, err := s.ar.AccountByLogin(ctx, logIn.Login)
		if err != nil {
			var errMsg string

			l.Warn("error get account from db", zap.String("user login", logIn.Login))

			// Проверяем, является ли полученная ошибка - ошибка получения аккаунта
			if errors.Is(err, accountRepo.ErrAccountNotFound) {
				errMsg = "аккаунт не найден"
			} else {
				errMsg = "неизвестная ошибка получения аккаунта"
			}

			sendCreatedTokenResult(resCh, person.LogInResp{}, errMsg)
			return
		}

		// Если не совпадает - пароль не верен
		if !password.CheckHash(logIn.Password, acc.PasswordHash) {
			l.Debug("user enter incorrect data",
				zap.String("login", logIn.Login),
				zap.String("password", logIn.Password),
			)

			sendCreatedTokenResult(resCh, person.LogInResp{}, "неправильный логин/пароль")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
			"exp":        time.Now().Add(time.Hour).Unix(),
			"account_id": acc.AccountId,
		})

		signedToken, err := token.SignedString([]byte(testKey))
		if err != nil {
			l.Warn("signing token error", zap.Error(err))

			sendCreatedTokenResult(resCh, person.LogInResp{}, "ошибка подписи ключа")
			return
		}

		resp := person.LogInResp{
			Token: signedToken,
		}

		sendCreatedTokenResult(resCh, resp, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return person.LogInResp{}, ctx.Err()
		case result := <-resCh:
			return result.CreatedToken, result.Error
		}
	}
}

func sendRegistrationResult(resCh chan RegistrationResult, resp person.RegisteredResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- RegistrationResult{
		RegisteredPerson: resp,
		Error:            err,
	}
}

func sendCreatedTokenResult(resCh chan LogInResult, resp person.LogInResp, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- LogInResult{
		CreatedToken: resp,
		Error:        err,
	}
}
