package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/password"
	"time"
)

type LogInResult struct {
	CreatedToken string
	Error        error
}

type authClaims struct {
	jwt.RegisteredClaims
	AccountId int `json:"acc_id,omitempty"`
}

func (s Service) CreateToken(ctx context.Context, cred dto.Credentials) (string, error) {
	resCh := make(chan LogInResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.LoginOperation),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		// Находим аккаунт по введенному логину
		acc, err := s.accountDAO.ByLogin(ctx, cred.Login)
		if err != nil {
			errMsg := "Ошибка получения аккаунта"

			l.Warn("ошибка получения аккаунта", zap.String("логин аккаунта", cred.Login))

			if errors.Is(err, pgx.ErrNoRows) {
				errMsg = "Аккаунт не найден"
			}

			sendCreatedTokenResult(resCh, "", errMsg)
			return
		}

		// Если не совпадает - пароль не верен
		if !password.CheckHash(cred.Password, acc.PasswordHash) {
			l.Warn("вход по некорректным данным",
				zap.String("логин", cred.Login),
				zap.String("пароль", cred.Password),
			)

			sendCreatedTokenResult(resCh, "", "Неправильный логин или пароль")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &authClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			AccountId: acc.Id,
		})

		signedToken, err := token.SignedString([]byte(s.signingKey))
		if err != nil {
			l.Warn("ошибка подписи токена", zap.Error(err))

			sendCreatedTokenResult(resCh, "", "Ошибка создания токена авторизации")
			return
		}

		sendCreatedTokenResult(resCh, signedToken, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case result := <-resCh:
			return result.CreatedToken, result.Error
		}
	}
}

func sendCreatedTokenResult(resCh chan LogInResult, token string, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- LogInResult{
		CreatedToken: token,
		Error:        err,
	}
}
