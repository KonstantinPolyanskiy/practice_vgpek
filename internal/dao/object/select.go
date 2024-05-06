package object

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

func (dao DAO) ById(ctx context.Context, id int) (entity.Object, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectObjectById),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM internal_object WHERE internal_object_id=@ObjectId`

	args := pgx.NamedArgs{
		"ObjectId": id,
	}

	l.Debug("аргументы запроса", zap.Int("id объекта", args["ObjectId"].(int)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Object{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	object, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Object])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Object{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("id действия", object.Id))

	return object, nil
}

func (dao DAO) ByParams(ctx context.Context, p params.Default) ([]entity.Object, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectObjectByParams),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := squirrel.Select("*").From("internal_object").
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
	rows, err := dao.db.Query(ctx, q, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return nil, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	objects, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Object])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return nil, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("количество объектов", len(objects)))

	return objects, nil
}
