package account

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/account"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
)

var duplicateKeyCodeError = "23505"

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
		zap.String("операция", operation.NewAccountOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedAccId int

	insertAccQuery := `
	INSERT INTO account (login, password_hash, internal_role_id, reg_key_id) 
	VALUES (@Login, @PasswordHash, @RoleId, @RegKeyId)
	RETURNING account_id
`

	args := pgx.NamedArgs{
		"Login":        savingAcc.Login,
		"PasswordHash": savingAcc.PasswordHash,
		"RoleId":       savingAcc.RoleId,
		"RegKeyId":     savingAcc.RegKeyId,
	}

	// Если запрос не возвращает Id, то аккаунт не создан
	err := r.db.QueryRow(ctx, insertAccQuery, args).Scan(&insertedAccId)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		var pgErr *pgconn.PgError

		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, dberr.ErrNotFound
		} else if errors.As(err, &pgErr) {
			if pgErr.Code == duplicateKeyCodeError {
				return account.Entity{}, dberr.ErrLoginAlreadyExist
			}
		}
		return account.Entity{}, err
	}

	getAccQuery := `
	SELECT * FROM account 
	WHERE account_id = $1
`

	row, err := r.db.Query(ctx, getAccQuery, insertedAccId)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, dberr.ErrNotFound
		}
		return account.Entity{}, err
	}

	savedAcc, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return account.Entity{}, err
	}

	return savedAcc, nil
}

func (r Repository) AccountByLogin(ctx context.Context, login string) (account.Entity, error) {
	l := r.l.With(
		zap.String("запрос к базе данных", "получение аккаунта по логину&паролю"),
		zap.String("слой", "репозиторий"),
	)

	getAccountQuery := `SELECT * FROM account WHERE login=$1`

	row, err := r.db.Query(ctx, getAccountQuery, login)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, dberr.ErrNotFound
		}

		return account.Entity{}, err
	}

	acc, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return account.Entity{}, err
	}

	return acc, nil
}

func (r Repository) AccountById(ctx context.Context, id int) (account.Entity, error) {
	l := r.l.With(
		zap.String("запрос к базе данных", operation.GetAccountOperation),
		zap.String("слой", "репозиторий"),
	)

	getAccountQuery := `SELECT * FROM account WHERE account.account_id=$1`

	row, err := r.db.Query(ctx, getAccountQuery, id)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, dberr.ErrNotFound
		}

		return account.Entity{}, err
	}

	acc, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return account.Entity{}, err
	}

	return acc, nil
}
