package rbac

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PermissionRepository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewPermissionRepository(db *pgxpool.Pool, logger *zap.Logger) PermissionRepository {
	return PermissionRepository{
		l:  logger,
		db: db,
	}
}

func (r PermissionRepository) SavePermission(ctx context.Context, roleId, objectId int, actionsId []int) error {
	l := r.l.With(
		zap.String("executing query name", "save permission"),
		zap.String("layer", "repo"),
	)

	insertPermQuery := `
	INSERT INTO role_permission (internal_role_id, internal_action_id, internal_object_id) 
	VALUES ($1, $2, $3) 
`
	for _, actionId := range actionsId {
		_, err := r.db.Exec(ctx, insertPermQuery, roleId, actionId, objectId)
		if err != nil {
			l.Warn("ошибка сохранения доступа",
				zap.Int("Id действия", actionId),
				zap.Int("Id объекта", objectId),
				zap.Int("Id роли", roleId),
				zap.Error(err),
			)

			return errors.New("ошибка сохранения доступа")
		}
	}

	return nil
}
