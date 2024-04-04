package rbac

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

const (
	ObjectName    = "RBAC"
	GetActionName = "GET"
)

var (
	ErrDontHavePermission = errors.New("нет доступа")
	ErrCheckAccess        = errors.New("ошибка проверки доступа")
)

type AccountMediator interface {
	// HasAccess проверяет, есть ли у указанной роли доступ к переданному действию к указанному объекту
	HasAccess(ctx context.Context, roleId int, objectName, actionName string) (bool, error)
}

type RBACService struct {
	l               *zap.Logger
	ar              ActionRepository
	or              ObjectRepository
	rr              RoleRepository
	pr              PermissionRepo
	accountMediator AccountMediator
}

func NewRBACService(
	actionRepo ActionRepository,
	objectRepo ObjectRepository,
	roleRepo RoleRepository,
	permRepo PermissionRepo,
	accountMediator AccountMediator,
	logger *zap.Logger) RBACService {
	return RBACService{
		ar:              actionRepo,
		or:              objectRepo,
		rr:              roleRepo,
		pr:              permRepo,
		accountMediator: accountMediator,
		l:               logger,
	}
}
