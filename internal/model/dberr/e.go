package dberr

import "errors"

var (
	ErrNotFound          = errors.New("ресурс не найден")
	ErrLoginAlreadyExist = errors.New("логин уже существует")
	ErrManyResult        = errors.New("неоднозначный результат")
)
