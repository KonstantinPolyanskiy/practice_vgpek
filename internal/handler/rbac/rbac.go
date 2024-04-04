package rbac

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
)

type RBACService interface {
	NewAction(ctx context.Context, addingAction permissions.AddActionReq) (permissions.AddActionResp, error)

	ActionById(ctx context.Context, req permissions.GetActionReq) (permissions.ActionEntity, error)
	ActionsByParams(ctx context.Context, params params.Default) ([]permissions.ActionEntity, error)

	ObjectById(ctx context.Context, id int) (permissions.ObjectEntity, error)
	ObjectsByParams(ctx context.Context, params params.Default) ([]permissions.ObjectEntity, error)

	RoleById(ctx context.Context, id int) (permissions.RoleEntity, error)
	RolesByParams(ctx context.Context, params params.Default) ([]permissions.RoleEntity, error)

	NewObject(ctx context.Context, addingObject permissions.AddObjectReq) (permissions.AddObjectResp, error)
	NewRole(ctx context.Context, addingRole permissions.AddRoleReq) (permissions.AddRoleResp, error)
	NewPermission(ctx context.Context, addingPermission permissions.AddPermReq) (permissions.AddPermResp, error)
}

type AccessHandler struct {
	l *zap.Logger
	s RBACService
}

func NewAccessHandler(service RBACService, logger *zap.Logger) AccessHandler {
	return AccessHandler{
		l: logger,
		s: service,
	}
}
