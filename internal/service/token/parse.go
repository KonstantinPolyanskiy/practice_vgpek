package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

type ParseTokenResult struct {
	AccountId int
	Error     error
}

func (s Service) ParseToken(ctx context.Context, token string) (int, error) {
	resCh := make(chan ParseTokenResult)

	l := s.logger.With(
		zap.String(operation.Operation, operation.ParseToken),
		zap.String(layer.Layer, layer.ServiceLayer),
	)

	go func() {
		t, err := jwt.ParseWithClaims(token, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				l.Warn("неправильная подпись токена",
					zap.String("токен", token.Raw),
					zap.String("ожидаем", jwt.SigningMethodHS256.Name),
					zap.String("текущий", token.Method.Alg()),
				)
				return nil, errors.New("неправильный метод подписи")
			}

			return []byte(s.signingKey), nil
		})
		if err != nil {
			l.Warn("ошибка расшифровки токена", zap.Error(err))

			sendParseTokenResult(resCh, 0, "ошибка расшифровки токена")
			return
		}

		c, ok := t.Claims.(*authClaims)
		if !ok {
			l.Warn("ошибка получения полей токена")

			sendParseTokenResult(resCh, 0, "ошибка расшифровки токена")
			return
		}

		sendParseTokenResult(resCh, c.AccountId, "")
		return
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case result := <-resCh:
			return result.AccountId, result.Error
		}
	}
}

func sendParseTokenResult(resCh chan ParseTokenResult, resp int, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- ParseTokenResult{
		AccountId: resp,
		Error:     err,
	}
}
