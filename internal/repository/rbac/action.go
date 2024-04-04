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
	ManyActionErr     = errors.New("неоднозначный результат")
	ActionNotFoundErr = errors.New("действие не найдено")
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
		zap.String("executing query name", "save action"),
		zap.String("layer", "repo"),
	)

	var insertedActionId int

	insertActionQuery := `
	INSERT INTO internal_action (internal_action_name)
	VALUES (@ActionName)
	RETURNING internal_action_id
`

	l.Debug("insert action", zap.String("query", insertActionQuery))

	args := pgx.NamedArgs{
		"ActionName": savingAction.Name,
	}

	l.Debug("args in insert action query", zap.Any("name action", args["ActionName"]))

	// Вставляем действие в БД
	err := ar.db.QueryRow(ctx, insertActionQuery, args).Scan(&insertedActionId)
	if err != nil {
		l.Warn("error insert action", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ActionEntity{}, errors.New("сохраненное действие не найдено")
		}
		return permissions.ActionEntity{}, err
	}

	getActionQuery := `
	SELECT * FROM internal_action WHERE internal_action_id=$1
`
	l.Debug("get inserted action", zap.String("query", getActionQuery))

	row, err := ar.db.Query(ctx, getActionQuery, insertedActionId)
	defer row.Close()
	if err != nil {
		l.Warn("error get inserted action", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ActionEntity{}, errors.New("сохраненное действие не найдено")
		}
		return permissions.ActionEntity{}, err
	}

	savedAction, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ActionEntity])
	if err != nil {
		l.Warn("error collect action in struct", zap.Error(err))

		return permissions.ActionEntity{}, err
	}

	return savedAction, nil
}

func (ar ActionRepository) ActionByName(ctx context.Context, name string) (permissions.ActionEntity, error) {
	l := ar.l.With(
		zap.String("executing query name", "get action by name"),
		zap.String("layer", "repo"),
	)

	var action permissions.ActionEntity

	getActionQuery := `SELECT * FROM internal_action WHERE internal_action_name=$1`

	err := ar.db.QueryRow(ctx, getActionQuery, name).Scan(&action)
	if err != nil {
		l.Warn("error get action by name",
			zap.String("action name", name),
			zap.Error(err),
		)

		if errors.Is(err, pgx.ErrTooManyRows) {
			return permissions.ActionEntity{}, ManyActionErr
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ActionEntity{}, ActionNotFoundErr
		}

		return permissions.ActionEntity{}, errors.New("unknown error")
	}

	return action, nil
}

func (ar ActionRepository) ActionById(ctx context.Context, id int) (permissions.ActionEntity, error) {
	l := ar.l.With(
		zap.String("executing query name", "get action by id"),
		zap.String("layer", "repo"),
	)

	getActionQuery := `SELECT * FROM internal_action WHERE internal_action_id=$1`

	row, err := ar.db.Query(ctx, getActionQuery, id)
	defer row.Close()
	if err != nil {
		l.Warn("error get action", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return permissions.ActionEntity{}, ActionNotFoundErr
		}

		return permissions.ActionEntity{}, errors.New("unknown error")
	}

	action, err := pgx.CollectOneRow(row, pgx.RowToStructByName[permissions.ActionEntity])
	if err != nil {
		l.Warn("error collect action to struct", zap.Error(err))

		return permissions.ActionEntity{}, errors.New("unknown error")
	}

	return action, nil
}

func (ar ActionRepository) ActionsByParams(ctx context.Context, params params.Default) ([]permissions.ActionEntity, error) {
	l := ar.l.With(
		zap.String("operation", "get actions by params"),
		zap.String("layer", "repo"),
	)

	getActionsQuery := squirrel.Select("*").From("internal_action").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	q, args, err := getActionsQuery.ToSql()
	if err != nil {
		l.Warn("error build sql", zap.Error(err))

		return nil, err
	}

	row, err := ar.db.Query(ctx, q, args...)
	defer row.Close()
	if err != nil {
		l.Warn("error get action by params", zap.Error(err))

		return nil, err
	}

	actions, err := pgx.CollectRows(row, pgx.RowToStructByName[permissions.ActionEntity])
	if err != nil {
		l.Warn("error collect action to struct", zap.Error(err))

		return nil, err
	}

	return actions, nil
}
