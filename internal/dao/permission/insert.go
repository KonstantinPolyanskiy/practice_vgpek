package permission

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Save(ctx context.Context, roleId, objectId int, actionsId []int) error {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SavePermissionsDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO role_permission 
    				(internal_role_id, internal_action_id, internal_object_id) 
					VALUES 
					($1, $2, $3)`

	now := time.Now()
	for _, actionId := range actionsId {
		_, err := dao.db.Exec(ctx, insertQuery, roleId, actionId, objectId)
		if err != nil {
			l.Warn(operation.ExecuteError, zap.Error(err))

			return err
		}
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))
	l.Info(operation.SuccessfullyRecorded)

	return nil
}
