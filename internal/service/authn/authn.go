package authn

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/person"
	"practice_vgpek/internal/model/registration_key"
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

type authClaims struct {
	jwt.RegisteredClaims
	AccountId int `json:"acc_id,omitempty"`
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
		zap.String("операция", operation.RegistrationOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		// Получаем ключ, по которому зарегистрированн пользователь
		regKey, err := s.kr.RegKeyByBody(ctx, registering.RegistrationKey)
		if err != nil {
			l.Warn("ошибка получения ключа регистрации", zap.String("тело ключа при регистрации", registering.RegistrationKey))

			// TODO: написать проверку, найден ли ключ или другая ошибка
			sendRegistrationResult(resCh, person.RegisteredResp{}, "Ошибка с ключем регистрации / ключ регистрации не найден")
			return
		}

		// Проверка, что ключ валиден, если нет - возвращаем ошибку
		if !regKey.IsValid {
			l.Warn("попытка регистрации по неактивному ключу", zap.String("тело ключа регистрации", regKey.Body))
			sendRegistrationResult(resCh, person.RegisteredResp{}, "Ключ регистрации неактивен")
			return
		}

		// Проверяем, что ключ еще можно использовать, если нет - инвалидируем
		if regKey.CurrentCountUsages >= regKey.MaxCountUsages {
			if err = s.kr.Invalidate(ctx, regKey.RegKeyId); err != nil {
				l.Warn("ошибка деактивации ключа",
					zap.String("тело ключа", regKey.Body),
					zap.Int("id ключа", regKey.RegKeyId),
					zap.Bool("валиден", regKey.IsValid),
				)

				sendRegistrationResult(resCh, person.RegisteredResp{}, "Ошибка с ключем регистрации")
				return
			}
		}

		// Хешируем пароль
		passwordHash, err := password.Hash(registering.Password)
		if err != nil {
			l.Warn("ошибка хеширования пароля",
				zap.String("пароль", registering.Password),
				zap.Error(err),
			)

			sendRegistrationResult(resCh, person.RegisteredResp{}, "Ошибка с введенным паролем")
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
			errMsg := "Ошибка создания аккаунта"

			l.Warn("ошибка сохранения аккаунта",
				zap.String("логин", dto.Account.Login),
				zap.Int("id ключа регистрации", dto.Account.RegKeyId),
			)

			if errors.Is(err, dberr.ErrLoginAlreadyExist) {
				errMsg = "Введенный логин уже занят"
			}

			sendRegistrationResult(resCh, person.RegisteredResp{}, errMsg)
			return
		}

		// Если аккаунт создан, увеличиваем кол-во регистраций по ключу
		err = s.kr.IncCountUsages(ctx, regKey.RegKeyId)
		if err != nil {
			l.Warn("ошибка увеличения текущего кол-ва регистраций по ключу",
				zap.Int("id ключа", regKey.RegKeyId),
				zap.Bool("валиден", regKey.IsValid),
			)

			sendRegistrationResult(resCh, person.RegisteredResp{}, "Ошибка при работе с ключем регистрации")
			return
		}

		// Сохраняем регистируемого пользователя в БД
		savedPerson, err := s.pr.SavePerson(ctx, dto, savedAcc.AccountId)
		if err != nil {
			l.Warn("ошибка сохранения пользователя",
				zap.String("имя", dto.FirstName),
				zap.String("фамилия", dto.LastName),
				zap.String("отчество", dto.MiddleName),
				zap.Int("id аккаунта", savedAcc.AccountId),
			)

			sendRegistrationResult(resCh, person.RegisteredResp{}, "Ошибка создания пользователя")
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
		zap.String("операция", operation.LoginOperation),
		zap.String("слой", "сервисы"),
	)

	go func() {
		acc, err := s.ar.AccountByLogin(ctx, logIn.Login)
		if err != nil {
			errMsg := "Ошибка получения аккаунта"

			l.Warn("ошибка получения аккаунта", zap.String("логин аккаунта", logIn.Login))

			if errors.Is(err, dberr.ErrNotFound) {
				errMsg = "Аккаунт не найден"
			}

			sendCreatedTokenResult(resCh, person.LogInResp{}, errMsg)
			return
		}

		// Если не совпадает - пароль не верен
		if !password.CheckHash(logIn.Password, acc.PasswordHash) {
			l.Warn("вход по некорректным данным",
				zap.String("логин", logIn.Login),
				zap.String("пароль", logIn.Password),
			)

			sendCreatedTokenResult(resCh, person.LogInResp{}, "Неправильный логин или пароль")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &authClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			AccountId: acc.AccountId,
		})

		signedToken, err := token.SignedString([]byte(testKey))
		if err != nil {
			l.Warn("ошибка подписи токена", zap.Error(err))

			sendCreatedTokenResult(resCh, person.LogInResp{}, "Ошибка создания токена авторизации")
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

func (s Service) ParseToken(token string) (int, error) {
	l := s.l.With(
		zap.String("операция", "расшифровка токена"),
		zap.String("слой", "сервисы"),
	)

	t, err := jwt.ParseWithClaims(token, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			l.Warn("неправильная подпись токена",
				zap.String("токен", token.Raw),
				zap.String("ожидаем", jwt.SigningMethodHS256.Name),
				zap.String("текущий", token.Method.Alg()),
			)
			return nil, errors.New("Неправильный метод подписи")
		}

		return []byte(testKey), nil
	})
	if err != nil {
		l.Warn("ошибка расшифровки токена", zap.Error(err))

		return 0, errors.New("Ошибка расшифровки токена")
	}

	c, ok := t.Claims.(*authClaims)
	if !ok {
		return 0, errors.New("Ошибка получения полей токена")
	}

	return c.AccountId, err
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
