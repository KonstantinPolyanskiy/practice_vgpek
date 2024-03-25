package rbac

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/permissions"
)

type RBACService interface {
	NewAction(ctx context.Context, addingAction permissions.AddActionReq) (permissions.AddActionResp, error)
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
