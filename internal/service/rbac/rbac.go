package rbac

import (
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
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

type Deletable interface {
	Deleted() bool
}

type Part interface {
	Part() domain.RBACPart
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
