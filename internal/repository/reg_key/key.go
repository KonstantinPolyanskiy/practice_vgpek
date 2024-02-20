package reg_key

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"practice_vgpek/internal/model/registration_key"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewKeyRepo(db *pgxpool.Pool) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error) {
	var insertedKeyId int

	insertKeyQuery := `INSERT INTO registration_key (internal_role_id, body_key, max_count_usages, current_count_usages, created_at)  
	VALUES (@RoleId, @BodyKey, @MaxCountUsages, @CurrentCountUsages, @CreatedAt)
	RETURNING reg_key_id
	`
	args := pgx.NamedArgs{
		"RoleId":             key.RoleId,
		"BodyKey":            key.Body,
		"MaxCountUsages":     key.MaxCountUsages,
		"CurrentCountUsages": 0,
		"CreatedAt":          time.Now(),
	}

	// Вставляем полученный ключ в БД и получаем его ID
	err := r.db.QueryRow(ctx, insertKeyQuery, args).Scan(&insertedKeyId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("сохраненный ключ не найден")
		}
		return registration_key.Entity{}, err
	}

	getKeyQuery := `SELECT * FROM registration_key WHERE reg_key_id = $1`

	row, err := r.db.Query(ctx, getKeyQuery, insertedKeyId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("сохраненный ключ не найден")
		}
		return registration_key.Entity{}, err
	}

	savedKey, err := pgx.CollectOneRow(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		return registration_key.Entity{}, err
	}

	return savedKey, nil
}
