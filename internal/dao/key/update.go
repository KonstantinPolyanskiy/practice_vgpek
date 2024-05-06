package key

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Update(ctx context.Context, new entity.Key) (entity.Key, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.UpdateKeyDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	update := updateQ("registration_key", new)

	update = update.Where("reg_key_id = $2", new.Id)

	updateQuery, args, err := update.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		l.Error("ошибка сборки запроса", zap.Error(err))
		return entity.Key{}, err
	}

	now := time.Now()
	_, err = dao.db.Exec(ctx, updateQuery, args...)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Debug(operation.Update, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	getQuery := `SELECT * FROM registration_key WHERE reg_key_id=$1`

	now = time.Now()
	rows, err := dao.db.Query(ctx, getQuery, new.Id)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Key])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Key{}, err
	}

	return updated, nil
}

func updateQ(table string, key entity.Key) squirrel.UpdateBuilder {
	updateBuilder := squirrel.Update(table)

	if key.RoleId != 0 {
		updateBuilder = updateBuilder.Set("internal_role_id", key.RoleId)
	}
	if key.Body != "" {
		updateBuilder = updateBuilder.Set("body_key", key.Body)
	}
	if key.MaxCountUsages != 0 {
		updateBuilder = updateBuilder.Set("max_count_usages", key.MaxCountUsages)
	}
	if key.CurrentCountUsages != 0 {
		updateBuilder = updateBuilder.Set("current_count_usages", key.CurrentCountUsages)
	}
	if !key.CreatedAt.IsZero() {
		updateBuilder = updateBuilder.Set("created_at", key.CreatedAt)
	}
	if key.IsValid {
		updateBuilder = updateBuilder.Set("is_valid", key.IsValid)
	}
	if key.InvalidationTime != nil {
		updateBuilder = updateBuilder.Set("invalidation_time", *key.InvalidationTime)
	}
	if key.GroupName != "" {
		updateBuilder = updateBuilder.Set("group_name", key.GroupName)
	}

	return updateBuilder
}
