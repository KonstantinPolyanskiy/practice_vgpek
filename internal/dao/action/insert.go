package action

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Save(ctx context.Context, action dto.NewRBACPart) (entity.Action, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SaveActionDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO internal_action 
    				(internal_action_name, description) 
					VALUES 
					(@ActionName, @Description)
					RETURNING internal_action_id`

	args := pgx.NamedArgs{
		"ActionName":  action.Name,
		"Description": action.Description,
	}

	l.Debug("аргументы запроса",
		zap.String("название", args["ActionName"].(string)),
		zap.String("описание", args["Description"].(string)),
	)

	var id int

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&id)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Action{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	getQuery := `SELECT * FROM internal_action WHERE internal_action_id=$1`

	now = time.Now()
	rows, err := dao.db.Query(ctx, getQuery, id)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Action{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	saved, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Action])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Action{}, err
	}

	l.Info(operation.SuccessfullyRecorded, zap.Int("id действия", id))

	return saved, nil
}
