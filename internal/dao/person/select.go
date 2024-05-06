package person

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) ByUUID(ctx context.Context, uid uuid.UUID) (entity.Person, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectPersonByUIIDDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM person WHERE person_uuid=@PersonUUID`

	args := pgx.NamedArgs{
		"PersonUUID": uid,
	}

	l.Debug("аргументы запроса", zap.String("uuid пользователя", args["PersonUUID"].(uuid.UUID).String()))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Person{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	person, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Person])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Person{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.String("uuid пользователя", person.UUID.String()))

	return person, nil
}

func (dao DAO) ByParams(ctx context.Context, p params.Default) ([]entity.Person, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectPersonByUIIDDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := squirrel.Select("*").From("person").
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

	persons, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Person])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return nil, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("количество пользователей", len(persons)))

	return persons, nil
}

func (dao DAO) ByAccountId(ctx context.Context, id int) (entity.Person, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectPersonByAccIdDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM person WHERE account_id=@AccountId`

	args := pgx.NamedArgs{
		"AccountId": id,
	}

	l.Debug("аргументы запроса", zap.Int("id пользователя", args["AccountId"].(int)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Person{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	person, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Person])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Person{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.String("uuid пользователя", person.UUID.String()))

	return person, nil
}
