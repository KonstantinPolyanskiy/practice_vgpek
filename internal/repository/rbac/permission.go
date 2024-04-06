package rbac

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/permissions"
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
		zap.String("операция", operation.AddPermissionOperation),
		zap.String("layer", "repo"),
	)

	insertPermQuery := `
	INSERT INTO role_permission (internal_role_id, internal_action_id, internal_object_id) 
	VALUES ($1, $2, $3) 
`
	for _, actionId := range actionsId {
		_, err := r.db.Exec(ctx, insertPermQuery, roleId, actionId, objectId)
		if err != nil {
			l.Warn("ошибка выполнения запроса", zap.Error(err))

			return err
		}
	}

	return nil
}

func (r PermissionRepository) PermissionsByRoleId(ctx context.Context, roleId int) ([]permissions.PermissionEntity, error) {
	l := r.l.With(
		zap.String("операция", operation.GetPermissionOperation),
		zap.String("слой", "репозиторий"),
	)
	getPermissionsQuery := `
	SELECT role_perm_id, ir.internal_role_id, ir.role_name, ia.internal_action_id, ia.internal_action_name, io.internal_object_id, io.internal_object_name
	FROM role_permission rp
	JOIN internal_role ir ON rp.internal_role_id = ir.internal_role_id
	JOIN internal_action ia ON rp.internal_action_id = ia.internal_action_id
	JOIN internal_object io ON rp.internal_object_id = io.internal_object_id
	WHERE ir.internal_role_id = $1;
`
	rows, err := r.db.Query(ctx, getPermissionsQuery, roleId)
	defer rows.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberr.ErrNotFound
		}
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[permissions.PermissionEntity])
	if err != nil {
		r.l.Info("ошибка записи данных в структуру", zap.Error(err))
		return nil, err
	}

	return result, err
}
