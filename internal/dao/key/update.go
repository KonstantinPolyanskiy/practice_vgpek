package key

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"log"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Update(ctx context.Context, new entity.KeyUpdate) (entity.Key, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.UpdateKeyDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	update, count := updateQ("registration_key", new)

	// Предикат WHERE reg_key_id должен идти последним, по этому нам нужно вычислить кол-во измененных полей и добавить к ним единицу
	update = update.Where(fmt.Sprintf("reg_key_id = $%d", count+1), new.Id)

	updateQuery, args, err := update.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		l.Error("ошибка сборки запроса", zap.Error(err))
		return entity.Key{}, err
	}

	log.Println(updateQuery)

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

// updateQ возвращает builder, и кол-во полей, которые запрос изменит
func updateQ(table string, key entity.KeyUpdate) (squirrel.UpdateBuilder, int) {
	updateBuilder := squirrel.Update(table)
	var count int

	if key.RoleId != nil {
		updateBuilder = updateBuilder.Set("internal_role_id", key.RoleId)
		count++
	}
	if key.Body != nil {
		updateBuilder = updateBuilder.Set("body_key", key.Body)
		count++
	}
	if key.MaxCountUsages != nil {
		updateBuilder = updateBuilder.Set("max_count_usages", key.MaxCountUsages)
		count++
	}
	if key.CurrentCountUsages != nil {
		updateBuilder = updateBuilder.Set("current_count_usages", key.CurrentCountUsages)
		count++
	}
	if key.CreatedAt != nil {
		updateBuilder = updateBuilder.Set("created_at", key.CreatedAt)
		count++
	}
	if key.IsValid != nil {
		updateBuilder = updateBuilder.Set("is_valid", &key.IsValid)
		count++
	}
	if key.InvalidationTime != nil {
		updateBuilder = updateBuilder.Set("invalidation_time", &key.InvalidationTime)
		count++
	}
	if key.GroupName != nil {
		updateBuilder = updateBuilder.Set("group_name", key.GroupName)
		count++
	}

	return updateBuilder, count
}
