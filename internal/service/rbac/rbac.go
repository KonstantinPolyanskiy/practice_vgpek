package rbac

import "go.uber.org/zap"

type RBACService struct {
	l  *zap.Logger
	ar ActionRepository
	or ObjectRepository
	rr RoleRepository
}

func NewRBACService(
	actionRepo ActionRepository,
	objectRepo ObjectRepository,
	roleRepo RoleRepository,
	logger *zap.Logger) RBACService {
	return RBACService{
		ar: actionRepo,
		or: objectRepo,
		rr: roleRepo,
		l:  logger,
	}
}
