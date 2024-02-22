package account

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/account"
)

type Repository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool, logger *zap.Logger) Repository {
	return Repository{
		l:  logger,
		db: db,
	}
}

func (r Repository) SaveAccount(ctx context.Context, savingAcc account.DTO) (account.Entity, error) {
	l := r.l.With(
		zap.String("action query", "save account"),
		zap.String("layer", "repo"),
	)

	var insertedAccId int

	insertAccQuery := `
	INSERT INTO account (login, password_hash, internal_role_id, reg_key_id) 
	VALUES (@Login, @PasswordHash, @RoleId, @RegKeyId)
	RETURNING account_id
`
	l.Debug("insert account", zap.String("query", insertAccQuery))

	args := pgx.NamedArgs{
		"Login":        savingAcc.Login,
		"PasswordHash": savingAcc.PasswordHash,
		"RoleId":       savingAcc.RoleId,
		"RegKeyId":     savingAcc.RegKeyId,
	}

	l.Debug("args in query",
		zap.String("login", savingAcc.Login),
		zap.Int("role id", savingAcc.RoleId),
		zap.Int("key id", savingAcc.RegKeyId),
	)

	// Если запрос не возвращает Id, то аккаунт не создан
	err := r.db.QueryRow(ctx, insertAccQuery, args).Scan(&insertedAccId)
	if err != nil {
		l.Warn("error insert account", zap.Error(err))

		var pgErr *pgconn.PgError

		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, errors.New("сохраненный аккаунт не найден")
		} else if errors.As(err, &pgErr) {
			if pgErr.Code == duplicateKeyCodeError {
				return account.Entity{}, ErrLoginAlreadyExist
			}
		}
		return account.Entity{}, err
	}

	getAccQuery := `
	SELECT * FROM account 
	WHERE account_id = $1
`
	l.Debug("get account", zap.String("query", getAccQuery))

	row, err := r.db.Query(ctx, getAccQuery, insertedAccId)
	if err != nil {
		l.Warn("error get inserted account", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, errors.New("сохраненный аккаунт не найден")
		}
		return account.Entity{}, err
	}

	savedAcc, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account.Entity])
	if err != nil {
		l.Warn("error collect account in struct", zap.Error(err))

		return account.Entity{}, err
	}

	return savedAcc, nil
}
