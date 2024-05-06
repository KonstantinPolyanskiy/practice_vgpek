package key

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Save(ctx context.Context, info dto.NewKeyInfo) (entity.Key, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SaveKeyDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO registration_key
					(internal_role_id, body_key, max_count_usages, current_count_usages, created_at, group_name) 
					VALUES 
					(@RoleId, @Body, @MaxUsages, @CurrentUsages, @CreatedAt, @GroupName)
					RETURNING reg_key_id`

	args := pgx.NamedArgs{
		"RoleId":        info.RoleId,
		"Body":          info.Body,
		"MaxUsages":     info.MaxCountUsages,
		"CurrentUsages": 0,
		"CreatedAt":     info.CreatedAt,
		"GroupName":     info.Group,
	}

	l.Debug("аргументы запроса",
		zap.Int("id роли", args["RoleId"].(int)),
		zap.String("тело ключа", args["Body"].(string)),
		zap.Int("макс. кол-во исп-ий", args["MaxUsages"].(int)),
		zap.Int("тек. кол-во исп-ий", args["CurrentUsages"].(int)),
		zap.Time("время создания", args["CreatedAt"].(time.Time)),
		zap.String("группа", args["GroupName"].(string)),
	)

	var id int

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&id)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	getQuery := `SELECT * FROM registration_key WHERE reg_key_id=$1`

	now = time.Now()
	rows, err := dao.db.Query(ctx, getQuery, id)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	saved, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Key])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Key{}, err
	}

	l.Info(operation.SuccessfullyRecorded, zap.Int("id ключа", id))

	return saved, nil
}
