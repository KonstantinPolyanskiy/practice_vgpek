package reg_key

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
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
		zap.String("операция", operation.NewKeyOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedKeyId int

	insertKeyQuery := `
	INSERT INTO registration_key (internal_role_id, body_key, group_name, max_count_usages, current_count_usages, created_at)  
	VALUES (@RoleId, @BodyKey, @GroupName, @MaxCountUsages, @CurrentCountUsages, @CreatedAt)
	RETURNING reg_key_id
	`

	args := pgx.NamedArgs{
		"RoleId":             key.RoleId,
		"BodyKey":            key.Body,
		"GroupName":          key.GroupName,
		"MaxCountUsages":     key.MaxCountUsages,
		"CurrentCountUsages": 0,
		"CreatedAt":          time.Now(),
	}

	// Вставляем полученный ключ в БД и получаем его ID
	err := r.db.QueryRow(ctx, insertKeyQuery, args).Scan(&insertedKeyId)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("сохраненный ключ не найден")
		}
		return registration_key.Entity{}, err
	}

	getKeyQuery := `SELECT * FROM registration_key WHERE reg_key_id = $1`

	row, err := r.db.Query(ctx, getKeyQuery, insertedKeyId)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("сохраненный ключ не найден")
		}
		return registration_key.Entity{}, err
	}

	savedKey, err := pgx.CollectOneRow(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return registration_key.Entity{}, err
	}

	return savedKey, nil
}

func (r Repository) RegKeyByBody(ctx context.Context, body string) (registration_key.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.GetKeyByBodyOperation),
		zap.String("слой", "репозиторий"),
	)

	findRoleQuery := `
		SELECT * FROM registration_key
		WHERE body_key = @BodyKey 
	`

	args := pgx.NamedArgs{
		"BodyKey": body,
	}

	// Находим ключ по телу (body) - он должен существовать в одном экземпляре
	row, err := r.db.Query(ctx, findRoleQuery, args)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return registration_key.Entity{}, errors.New("ключ регистрации по телу не найден")
		}
		return registration_key.Entity{}, err
	}

	regKey, err := pgx.CollectOneRow(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		l.Warn("ошибка приведения данных к структуре", zap.Error(err))

		return registration_key.Entity{}, err
	}

	return regKey, nil
}

func (r Repository) RegKeyById(ctx context.Context, id int) (registration_key.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.GetKeyByIdOperation),
		zap.String("слой", "репозиторий"),
	)

	findRoleQuery := `SELECT * FROM registration_key WHERE reg_key_id = $1`

	row, err := r.db.Query(ctx, findRoleQuery, id)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return registration_key.Entity{}, err
	}

	key, err := pgx.CollectOneRow(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		l.Warn("ошибка приведения данных к структуре", zap.Error(err))

		return registration_key.Entity{}, err
	}

	return key, nil
}

func (r Repository) IncCountUsages(ctx context.Context, keyId int) error {
	l := r.l.With(
		zap.String("операция", operation.IncCountUsagesOperation),
		zap.String("слой", "репозиторий"),
	)

	incrementCountQuery := `
	UPDATE registration_key
	SET current_count_usages = current_count_usages + 1
	WHERE reg_key_id = $1
`

	_, err := r.db.Exec(ctx, incrementCountQuery, keyId)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return errors.Join(ErrNotUpdate, err)
	}

	return nil
}

func (r Repository) Invalidate(ctx context.Context, keyId int) error {
	l := r.l.With(
		zap.String("операция", operation.InvalidateKeyOperation),
		zap.String("слой", "репозиторий"),
	)

	invalidateKeyQuery := `
	UPDATE registration_key
	SET 
	is_valid = false 
	AND 
	invalidation_time = $2
	WHERE reg_key_id = $1
`

	_, err := r.db.Exec(ctx, invalidateKeyQuery, keyId, time.Now())
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return err
	}

	return nil
}
