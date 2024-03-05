package rbac

import "go.uber.org/zap"

type RBACService struct {
	l  *zap.Logger
	ar ActionRepository
	or ObjectRepository
}

func NewRBACService(
	actionRepo ActionRepository,
	objectRepo ObjectRepository,
	logger *zap.Logger) RBACService {
	return RBACService{
		ar: actionRepo,
		or: objectRepo,
		l:  logger,
	}
}
