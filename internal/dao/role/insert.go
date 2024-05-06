package role

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

func (dao DAO) Save(ctx context.Context, role dto.NewRBACPart) (entity.Role, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SaveRoleDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO internal_role 
    				(role_name, description, created_at) 
					VALUES 
					(@RoleName, @Description, created_at)
					RETURNING internal_role_id`

	args := pgx.NamedArgs{
		"RoleName":    role.Name,
		"Description": role.Description,
		"CreatedAt":   role.CreatedAt,
	}

	l.Debug("аргументы запроса",
		zap.String("название", args["RoleName"].(string)),
		zap.String("описание", args["Description"].(string)),
		zap.Time("время создания", args["CreatedAt"].(time.Time)),
	)

	var id int

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&id)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Role{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	getQuery := `SELECT * FROM internal_role WHERE internal_role_id=$1`

	now = time.Now()
	rows, err := dao.db.Query(ctx, getQuery, id)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Role{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	saved, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Role])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Role{}, err
	}

	l.Info(operation.SuccessfullyRecorded, zap.Int("id роли", id))

	return saved, nil
}
