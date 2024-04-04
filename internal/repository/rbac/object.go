package rbac

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
)

var (
	ManyObjectErr     = errors.New("неоднозначный результат")
	ObjectNotFoundErr = errors.New("объект действия не найден")
)

type ObjectRepository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewObjectRepo(db *pgxpool.Pool, logger *zap.Logger) ObjectRepository {
	return ObjectRepository{
		l:  logger,
		db: db,
	}
}

func (or ObjectRepository) SaveObject(ctx context.Context, savingObject permissions.ObjectDTO) (permissions.ObjectEntity, error) {
	l := or.l.With(
		zap.String("executing query name", "save object"),
		zap.String("layer", "repo"),
	)

	var insertedObjectId int

	insertedObjectQuery := `
	INSERT INTO internal_object (internal_object_name)
	VALUES (@ObjectName)
	RETURNING internal_object_id
`

	l.Debug("insert object", zap.String("query", insertedObjectQuery))

	args := pgx.NamedArgs{
		"ObjectName": savingObject.Name,
	}

	l.Debug("args in insert object query", zap.Any("name object", args["ObjectName"]))

	// Вставляем объект в БД
	err := or.db.QueryRow(ctx, insertedObjectQuery, args).Scan(&insertedObjectId)
	if err != nil {
		l.Warn("error insert action", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ObjectEntity{}, errors.New("сохраненный объект не найден")
		}
		return permissions.ObjectEntity{}, err
	}

	getObjectQuery := `
	SELECT * FROM internal_object WHERE internal_object_id=$1
`

	l.Debug("get inserted object", zap.String("query", getObjectQuery))

	row, err := or.db.Query(ctx, getObjectQuery, insertedObjectId)
	if err != nil {
		l.Warn("error get inserted object", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ObjectEntity{}, errors.New("сохраненный объект не найден")
		}
		return permissions.ObjectEntity{}, err
	}

	savedObject, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ObjectEntity])
	if err != nil {
		l.Warn("error collect object in struct", zap.Error(err))

		return permissions.ObjectEntity{}, err
	}

	return savedObject, nil
}

func (or ObjectRepository) ObjectById(ctx context.Context, id int) (permissions.ObjectEntity, error) {
	l := or.l.With(
		zap.String("executing query name", "get object by id"),
		zap.String("layer", "repo"),
	)

	getObjectQuery := `SELECT * FROM internal_object WHERE internal_object.internal_object_id=$1`

	row, err := or.db.Query(ctx, getObjectQuery, id)
	if err != nil {
		l.Warn("error get object", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ObjectEntity{}, ActionNotFoundErr
		}

		return permissions.ObjectEntity{}, errors.New("unknown error")
	}

	object, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ObjectEntity])
	if err != nil {
		l.Warn("error collect object to struct", zap.Error(err))

		return permissions.ObjectEntity{}, errors.New("unknown error")
	}

	return object, nil
}

func (or ObjectRepository) ObjectsByParams(ctx context.Context, params params.Default) ([]permissions.ObjectEntity, error) {
	l := or.l.With(
		zap.String("operation", "get objects by params"),
		zap.String("layer", "repo"),
	)

	getObjectsQuery := squirrel.Select("*").From("internal_object").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	q, args, err := getObjectsQuery.ToSql()
	if err != nil {
		l.Warn("error build sql", zap.Error(err))

		return nil, err
	}

	row, err := or.db.Query(ctx, q, args...)
	if err != nil {
		l.Warn("error get objects by params", zap.Error(err))

		return nil, err
	}

	objects, err := pgx.CollectRows(row, pgx.RowToStructByName[permissions.ObjectEntity])
	if err != nil {
		l.Warn("error collect objects to struct", zap.Error(err))

		return nil, err
	}

	return objects, nil
}
