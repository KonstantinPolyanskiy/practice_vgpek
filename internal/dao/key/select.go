package key

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

func (dao DAO) ById(ctx context.Context, id int) (entity.Key, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectKeyById),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM registration_key WHERE reg_key_id=@KeyId`

	args := pgx.NamedArgs{
		"KeyId": id,
	}

	l.Debug("аргументы запроса", zap.Int("id ключа", args["KeyId"].(int)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	key, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Key])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("id ключа", key.Id))

	return key, nil
}

func (dao DAO) ByBody(ctx context.Context, body string) (entity.Key, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectKeyByBody),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM registration_key WHERE body_key=@Body`

	args := pgx.NamedArgs{
		"Body": body,
	}

	l.Debug("аргументы запроса", zap.String("тело ключа", args["Body"].(string)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	key, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Key])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Key{}, err
	}

	return key, nil
}

func (dao DAO) ByParams(ctx context.Context, p params.Default) ([]entity.Key, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectKeysByParams),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := squirrel.Select("*").From("registration_key").
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

	keys, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Key])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return nil, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("количество ключей", len(keys)))

	return keys, nil
}
