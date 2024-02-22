package account

import "errors"

var (
	duplicateKeyCodeError = "23505"
)

var (
	ErrLoginAlreadyExist = errors.New("login already exist")
)
