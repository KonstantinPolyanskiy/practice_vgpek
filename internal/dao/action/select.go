package action

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) ById(ctx context.Context, id int) (entity.Action, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectActionById),
		zap.String(layer.Layer, layer.DataLayer),
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

func (dao DAO) ByParams(ctx context.Context, p params.Default) ([]entity.Action, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectActionByParams),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := squirrel.Select("*").From("internal_action").
		Limit(uint64(p.Limit)).
		Offset(uint64(p.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	q, args, err := selectQuery.ToSql()
	if err != nil {
		l.Warn("ошибка подготовки запроса", zap.Error(err))

		return nil, err
	}

	l.Debug("аргументы запроса",
		zap.Int("лимит", p.Limit),
		zap.Int("смещение", p.Offset),
	)

	now := time.Now()
	rows, err := dao.db.Query(ctx, q, args...)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return nil, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	actions, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Action])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return nil, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("количество действий", len(actions)))

	return actions, nil
}
