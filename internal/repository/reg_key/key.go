package reg_key

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/registration_key"
	"time"
)

type Repository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewKeyRepo(db *pgxpool.Pool, logger *zap.Logger) Repository {
	return Repository{
		l:  logger,
		db: db,
	}
}

func (r Repository) SaveKey(ctx context.Context, key registration_key.DTO) (registration_key.Entity, error) {
	l := r.l.With(
		zap.String("action query", "save key"),
		zap.String("layer", "repo"),
	)

	var insertedKeyId int

	insertKeyQuery := `
	INSERT INTO registration_key (internal_role_id, body_key, max_count_usages, current_count_usages, created_at)  
	VALUES (@RoleId, @BodyKey, @MaxCountUsages, @CurrentCountUsages, @CreatedAt)
	RETURNING reg_key_id
	`

	l.Debug("insert key", zap.String("query", insertKeyQuery))

	args := pgx.NamedArgs{
		"RoleId":             key.RoleId,
		"BodyKey":            key.Body,
		"MaxCountUsages":     key.MaxCountUsages,
		"CurrentCountUsages": 0,
		"CreatedAt":          time.Now(),
	}

	l.Debug("args in query",
		zap.Int("rbac id", key.RoleId),
		zap.String("body key", key.Body),
		zap.Int("max count", key.MaxCountUsages),
		zap.Any("current count", args["CurrentCountUsages"]),
		zap.Any("created at", args["CreatedAt"]),
	)

	// Вставляем полученный ключ в БД и получаем его ID
	err := r.db.QueryRow(ctx, insertKeyQuery, args).Scan(&insertedKeyId)
	if err != nil {
		l.Warn("error insert key", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("сохраненный ключ не найден")
		}
		return registration_key.Entity{}, err
	}

	getKeyQuery := `SELECT * FROM registration_key WHERE reg_key_id = $1`

	l.Debug("get inserted key", zap.String("query", getKeyQuery))

	row, err := r.db.Query(ctx, getKeyQuery, insertedKeyId)
	if err != nil {
		l.Warn("error get inserted key", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("сохраненный ключ не найден")
		}
		return registration_key.Entity{}, err
	}

	savedKey, err := pgx.CollectOneRow(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		l.Warn("error collect key in struct", zap.Error(err))

		return registration_key.Entity{}, err
	}

	return savedKey, nil
}

func (r Repository) RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error) {
	l := r.l.With(
		zap.String("action", "get key by body"),
		zap.String("layer", "repo"),
	)

	findRoleQuery := `
		SELECT * FROM registration_key
		WHERE body_key = @BodyKey 
	`

	l.Debug("get key", zap.String("query", findRoleQuery))

	args := pgx.NamedArgs{
		"BodyKey": body,
	}

	l.Debug("args in query", zap.String("body", body))

	// Находим ключ по телу (body) - он должен существовать в одном экземпляре
	row, err := r.db.Query(ctx, findRoleQuery, args)
	if err != nil {
		l.Warn("error get key by body", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("ключ регистрации по телу не найден")
		}
		return registration_key.Entity{}, err
	}

	regKey, err := pgx.CollectOneRow(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		l.Warn("error collect key in struct")

		return registration_key.Entity{}, err
	}

	return regKey, nil
}

func (r Repository) IncCountUsages(ctx context.Context, keyId int) error {
	l := r.l.With(
		zap.String("action", "increment key count"),
		zap.String("layer", "repo"),
	)

	incrementCountQuery := `
	UPDATE registration_key
	SET current_count_usages = current_count_usages + 1
	WHERE reg_key_id = $1
`

	l.Debug("increment key", zap.String("query", incrementCountQuery))

	_, err := r.db.Exec(ctx, incrementCountQuery, keyId)
	if err != nil {
		l.Warn("error increment key", zap.Error(err))

		return errors.Join(ErrNotUpdate, err)
	}

	return nil
}

func (r Repository) Invalidate(ctx context.Context, keyId int) error {
	l := r.l.With(
		zap.String("action", "invalidate key"),
		zap.String("layer", "repo"),
	)

	invalidateKeyQuery := `
	UPDATE registration_key
	SET 
	is_valid = false 
	AND 
	invalidation_time = $2
	WHERE reg_key_id = $1
`

	l.Debug("invalidate key", zap.String("query", invalidateKeyQuery))

	_, err := r.db.Exec(ctx, invalidateKeyQuery, keyId, time.Now())
	if err != nil {
		l.Warn("error invalidate key", zap.Error(err))

		return err
	}

	return nil
}
