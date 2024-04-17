package action

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) ById(ctx context.Context, id int) (entity.Action, error) {
	l := dao.logger.With(
		zap.String("операция", operation.SelectActionById),
		zap.String("слой", operation.DataLayer),
	)

	getQuery := `SELECT * FROM internal_action WHERE internal_action_id=@ActionId`

	args := pgx.NamedArgs{
		"ActionId": id,
	}

	l.Debug("аргументы запроса", zap.Int("id действия", args["ActionId"].(int)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Action{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	action, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Action])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Action{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("id действия", action.Id))

	return action, nil
}
