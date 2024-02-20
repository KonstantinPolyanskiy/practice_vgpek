package account

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"practice_vgpek/internal/model/account"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) SaveAccount(ctx context.Context, savingAcc account.DTO) (account.Entity, error) {
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

	err := r.db.QueryRow(ctx, insertAccQuery, args).Scan(&insertedAccId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, errors.New("сохраненный аккаунт не найден")
		}
		return account.Entity{}, err
	}

	getAccQuery := `
	SELECT * FROM account 
	WHERE account_id = $1
`
	row, err := r.db.Query(ctx, getAccQuery, insertedAccId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account.Entity{}, errors.New("сохраненный аккаунт не найден")
		}
		return account.Entity{}, err
	}

	savedAcc, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account.Entity])
	if err != nil {
		return account.Entity{}, err
	}

	return savedAcc, nil
}
