package rbac

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/permissions"
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
		zap.String("операция", operation.AddObjectOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedObjectId int

	insertedObjectQuery := `
	INSERT INTO internal_object (internal_object_name)
	VALUES (@ObjectName)
	RETURNING internal_object_id
`

	args := pgx.NamedArgs{
		"ObjectName": savingObject.Name,
	}

	// Вставляем объект в БД
	err := or.db.QueryRow(ctx, insertedObjectQuery, args).Scan(&insertedObjectId)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ObjectEntity{}, dberr.ErrNotFound
		}
		return permissions.ObjectEntity{}, err
	}

	getObjectQuery := `
	SELECT * FROM internal_object WHERE internal_object_id=$1
`

	row, err := or.db.Query(ctx, getObjectQuery, insertedObjectId)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполениня запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ObjectEntity{}, dberr.ErrNotFound
		}
		return permissions.ObjectEntity{}, err
	}

	savedObject, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ObjectEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return permissions.ObjectEntity{}, err
	}

	return savedObject, nil
}

func (or ObjectRepository) ObjectById(ctx context.Context, id int) (permissions.ObjectEntity, error) {
	l := or.l.With(
		zap.String("операция", operation.GetObjectOperation),
		zap.String("слой", "репозиторий"),
	)

	getObjectQuery := `SELECT * FROM internal_object WHERE internal_object.internal_object_id=$1`

	row, err := or.db.Query(ctx, getObjectQuery, id)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ObjectEntity{}, dberr.ErrNotFound
		}

		return permissions.ObjectEntity{}, err
	}

	object, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ObjectEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return permissions.ObjectEntity{}, err
	}

	return object, nil
}

func (or ObjectRepository) ObjectsByParams(ctx context.Context, params params.Default) ([]permissions.ObjectEntity, error) {
	l := or.l.With(
		zap.String("операция", operation.GetObjectsOperation),
		zap.String("слой", "репозиторий"),
	)

	getObjectsQuery := squirrel.Select("*").From("internal_object").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	q, args, err := getObjectsQuery.ToSql()
	if err != nil {
		l.Warn("ошибка подготовки запроса", zap.Error(err))

		return nil, err
	}

	row, err := or.db.Query(ctx, q, args...)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return nil, err
	}

	objects, err := pgx.CollectRows(row, pgx.RowToStructByName[permissions.ObjectEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return nil, err
	}

	return objects, nil
}
