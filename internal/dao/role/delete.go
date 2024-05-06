package role

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) SoftDeleteById(ctx context.Context, id int, info dto.DeleteInfo) error {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SoftDeleteRoleById),
		zap.String(layer.Layer, layer.DataLayer),
	)

	deleteQuery := `UPDATE public.internal_role SET is_deleted = @DeleteTime WHERE internal_role_id = @RoleId`

	args := pgx.NamedArgs{
		"RoleId":     id,
		"DeleteTime": info.DeleteTime,
	}

	l.Debug("аргументы запроса", zap.Time("время удаления", args["DeleteTime"].(time.Time)))

	now := time.Now()
	_, err := dao.db.Exec(ctx, deleteQuery, args)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return err
	}

	l.Debug(operation.Update, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	l.Info(operation.SuccessfullyUpdated)

	return nil
}
