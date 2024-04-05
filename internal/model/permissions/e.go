package permissions

import "errors"

var (
	ErrDontHavePerm = errors.New("недостаточно прав")
)
