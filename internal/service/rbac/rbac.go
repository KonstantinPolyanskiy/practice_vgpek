package rbac

import (
	"go.uber.org/zap"
)

type RBACService struct {
	l         *zap.Logger
	actionDAO ActionDAO
	objectDAO ObjectDAO
	roleDAO   RoleDAO
	permDAO   PermissionDAO
}

func New(
	actionDAO ActionDAO,
	objectDAO ObjectDAO,
	roleDAO RoleDAO,
	permDAO PermissionDAO,
	logger *zap.Logger) RBACService {
	return RBACService{
		actionDAO: actionDAO,
		objectDAO: objectDAO,
		roleDAO:   roleDAO,
		permDAO:   permDAO,
		l:         logger,
	}
}
