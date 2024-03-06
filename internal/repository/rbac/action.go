package rbac

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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