package object

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

func (dao DAO) Save(ctx context.Context, object dto.NewRBACPart) (entity.Object, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SaveObjectDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO internal_object 
    				(internal_object_name, description) 
					VALUES 
					(@ObjectName, @Description)
					RETURNING internal_object_id`

	args := pgx.NamedArgs{
		"ObjectName":  object.Name,
		"Description": object.Description,
	}

	l.Debug("аргументы запроса",
		zap.String("название", args["ObjectName"].(string)),
		zap.String("описание", args["Description"].(string)),
	)

	var id int

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&id)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Object{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	getQuery := `SELECT * FROM internal_object WHERE internal_object_id=$1`

	now = time.Now()
	rows, err := dao.db.Query(ctx, getQuery, id)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Object{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	saved, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Object])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Object{}, err
	}

	l.Info(operation.SuccessfullyRecorded, zap.Int("id объекта", id))

	return saved, nil
}
