package reg_key

import "errors"

var (
	NewKeyAction = "создание нового ключа"
)

var (
	ErrDontHavePermission = errors.New("нет доступа")
)
