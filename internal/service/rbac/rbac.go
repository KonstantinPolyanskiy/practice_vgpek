package rbac

import "go.uber.org/zap"

type RBACService struct {
	l  *zap.Logger
	ar ActionRepository
	or ObjectRepository
	rr RoleRepository
	pr PermissionRepo
}

func NewRBACService(
	actionRepo ActionRepository,
	objectRepo ObjectRepository,
	roleRepo RoleRepository,
	permRepo PermissionRepo,
	logger *zap.Logger) RBACService {
	return RBACService{
		ar: actionRepo,
		or: objectRepo,
		rr: roleRepo,
		pr: permRepo,
		l:  logger,
	}
}
