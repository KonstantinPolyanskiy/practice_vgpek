package permissions

import "errors"

var (
	ErrDontHavePerm = errors.New("недостаточно прав")
)

var (
	ErrCheckAccess = errors.New("ошибка проверки доступа")
)
