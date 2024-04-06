package reg_key

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/registration_key"
)

func (r Repository) KeysByParams(ctx context.Context, params params.Key) ([]registration_key.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.GetKeysOperation),
		zap.String("слой", "репозиторий"),
	)

	findKeysQuery := squirrel.Select("*").From("registration_key").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	if params.IsValid {
		findKeysQuery = findKeysQuery.Where("is_valid = true")
	} else {
		findKeysQuery = findKeysQuery.Where("is_valid = false")
	}

	q, args, err := findKeysQuery.ToSql()
	if err != nil {
		l.Warn("ошибка подготовки запроса", zap.Error(err))
		return nil, err
	}

	row, err := r.db.Query(ctx, q, args...)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))
		return nil, err
	}

	keys, err := pgx.CollectRows(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		l.Warn("ошибка приведения данных к структуре", zap.Error(err))
		return nil, err
	}

	return keys, nil
}
