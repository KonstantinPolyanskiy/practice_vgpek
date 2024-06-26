package permission

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) ByRoleId(ctx context.Context, roleId int) ([]entity.Permissions, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectPermByRoleIdDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := `
	SELECT role_perm_id,
       ir.internal_role_id, ir.role_name as role_name, ir.description as description, ir.created_at as created_at, ir.is_deleted as is_deleted,
       ia.internal_action_id, ia.internal_action_name as internal_action_name, ia.description as description, ia.created_at as created_at, ia.is_deleted as is_deleted,
       io.internal_object_id, io.internal_object_name as internal_object_name, io.description as description, io.created_at as created_at, io.is_deleted as is_deleted
FROM role_permission rp
         JOIN internal_role ir ON rp.internal_role_id = ir.internal_role_id
         JOIN internal_action ia ON rp.internal_action_id = ia.internal_action_id
         JOIN internal_object io ON rp.internal_object_id = io.internal_object_id
WHERE ir.internal_role_id = $1;`

	now := time.Now()
	rows, err := dao.db.Query(ctx, selectQuery, roleId)
	defer rows.Close()
	if err != nil {
		l.Warn(operation.ExecuteError, zap.Error(err))

		return nil, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	perm, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Permissions])
	if err != nil {
		l.Warn(operation.CollectError, zap.Error(err))

		return nil, err
	}

	l.Info(operation.SuccessfullyReceived)

	return perm, nil
}
