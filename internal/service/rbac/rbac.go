package rbac

import (
	"context"
	"go.uber.org/zap"
)

const (
	ObjectName    = "RBAC"
	GetActionName = "GET"
)

type AccountMediator interface {
	// HasAccess проверяет, есть ли у указанной роли доступ к переданному действию к указанному объекту
	HasAccess(ctx context.Context, roleId int, objectName, actionName string) (bool, error)
}

type RBACService struct {
	l               *zap.Logger
	actionDAO       ActionDAO
	objectDAO       ObjectDAO
	roleDAO         RoleDAO
	permDAO         PermissionDAO
	accountMediator AccountMediator
}

func New(
	actionDAO ActionDAO,
	objectDAO ObjectDAO,
	roleDAO RoleDAO,
	permDAO PermissionDAO,
	accountMediator AccountMediator,
	logger *zap.Logger) RBACService {
	return RBACService{
		actionDAO:       actionDAO,
		objectDAO:       objectDAO,
		roleDAO:         roleDAO,
		permDAO:         permDAO,
		accountMediator: accountMediator,
		l:               logger,
	}
}

type Deletable interface {
	Deleted() bool
}

// filterDeleted возвращает только удаленные элементы
func filterDeleted[T Deletable](items []T) (result []T) {
	for _, item := range items {
		if item.Deleted() {
			result = append(result, item)
		}
	}
	return result
}

// filterNotDeleted возвращает только не удаленные элементы
func filterNotDeleted[T Deletable](items []T) (result []T) {
	for _, item := range items {
		if !item.Deleted() {
			result = append(result, item)
		}
	}
	return result
}
