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

type ActionRepository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewActionRepo(db *pgxpool.Pool, logger *zap.Logger) ActionRepository {
	return ActionRepository{
		l:  logger,
		db: db,
	}
}

func (ar ActionRepository) SaveAction(ctx context.Context, savingAction permissions.ActionDTO) (permissions.ActionEntity, error) {
	l := ar.l.With(
		zap.String("операция", operation.AddActionOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedActionId int

	insertActionQuery := `
	INSERT INTO internal_action (internal_action_name)
	VALUES (@ActionName)
	RETURNING internal_action_id
`

	args := pgx.NamedArgs{
		"ActionName": savingAction.Name,
	}

	// Вставляем действие в БД
	err := ar.db.QueryRow(ctx, insertActionQuery, args).Scan(&insertedActionId)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ActionEntity{}, dberr.ErrNotFound
		}
		return permissions.ActionEntity{}, err
	}

	getActionQuery := `
	SELECT * FROM internal_action WHERE internal_action_id=$1
`

	row, err := ar.db.Query(ctx, getActionQuery, insertedActionId)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ActionEntity{}, dberr.ErrNotFound
		}
		return permissions.ActionEntity{}, err
	}

	savedAction, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ActionEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return permissions.ActionEntity{}, err
	}

	return savedAction, nil
}

func (ar ActionRepository) ActionById(ctx context.Context, id int) (permissions.ActionEntity, error) {
	l := ar.l.With(
		zap.String("операция", operation.GetActionOperation),
		zap.String("layer", "repo"),
	)

	getActionQuery := `SELECT * FROM internal_action WHERE internal_action_id=$1`

	row, err := ar.db.Query(ctx, getActionQuery, id)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ActionEntity{}, dberr.ErrNotFound
		}

		return permissions.ActionEntity{}, errors.New("unknown error")
	}

	action, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ActionEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return permissions.ActionEntity{}, err
	}

	return action, nil
}

func (ar ActionRepository) ActionsByParams(ctx context.Context, params params.Default) ([]permissions.ActionEntity, error) {
	l := ar.l.With(
		zap.String("операция", operation.GetActionsOperation),
		zap.String("слой", "репозиторий"),
	)

	getActionsQuery := squirrel.Select("*").From("internal_action").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	q, args, err := getActionsQuery.ToSql()
	if err != nil {
		l.Warn("ошибка подготовки запроса", zap.Error(err))

		return nil, err
	}

	row, err := ar.db.Query(ctx, q, args...)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return nil, err
	}

	actions, err := pgx.CollectRows(row, pgx.RowToStructByName[permissions.ActionEntity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return nil, err
	}

	return actions, nil
}
