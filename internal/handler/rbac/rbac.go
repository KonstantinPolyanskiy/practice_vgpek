package rbac

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/domain"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/params"
)

type RBACService interface {
	ActionById(ctx context.Context, req dto.EntityId) (domain.Action, error)
	DeleteActionById(ctx context.Context, req dto.EntityId) (domain.Action, error)
	ActionsByParams(ctx context.Context, params params.State) ([]domain.Action, error)

	ObjectById(ctx context.Context, req dto.EntityId) (domain.Object, error)
	ObjectsByParams(ctx context.Context, params params.State) ([]domain.Object, error)

	RoleById(ctx context.Context, req dto.EntityId) (domain.Role, error)
	RolesByParams(ctx context.Context, params params.State) ([]domain.Role, error)

	NewAction(ctx context.Context, addingAction dto.NewRBACReq) (domain.Action, error)
	NewObject(ctx context.Context, addingObject dto.NewRBACReq) (domain.Object, error)
	NewRole(ctx context.Context, addingRole dto.NewRBACReq) (domain.Role, error)
	NewPermission(ctx context.Context, req dto.SetPermissionReq) error
}

type AccountMediator interface {
	HasAccess(ctx context.Context, accountId int, objectName, actionName string) (bool, error)
}

type AccessHandler struct {
	l               *zap.Logger
	s               RBACService
	accountMediator AccountMediator
}

func NewAccessHandler(service RBACService, accountMediator AccountMediator, logger *zap.Logger) AccessHandler {
	return AccessHandler{
		l:               logger,
		s:               service,
		accountMediator: accountMediator,
	}
}
