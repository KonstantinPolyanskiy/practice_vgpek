package reg_key

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/internal/model/registration_key"
)

func (r Repository) KeysByParams(ctx context.Context, params params.Key) ([]registration_key.Entity, error) {
	_ = r.l.With(
		zap.String("action", "get keys by params"),
		zap.String("layer", "repo"),
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
		return nil, err
	}

	row, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	keys, err := pgx.CollectRows(row, pgx.RowToStructByName[registration_key.Entity])
	if err != nil {
		return nil, err
	}

	return keys, nil
}
