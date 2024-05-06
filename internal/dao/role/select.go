package role

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/params"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) ById(ctx context.Context, id int) (entity.Role, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectRoleById),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM internal_role WHERE internal_role_id=@RoleId`

	args := pgx.NamedArgs{
		"RoleId": id,
	}

	l.Debug("аргументы запроса", zap.Int("id роли", args["RoleId"].(int)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Role{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Role])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Role{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("id роли", role.Id))

	return role, nil
}

func (dao DAO) ByParams(ctx context.Context, p params.Default) ([]entity.Role, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectRoleByParams),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := squirrel.Select("*").From("internal_role").
		Limit(uint64(p.Limit)).
		Offset(uint64(p.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	q, args, err := selectQuery.ToSql()
	if err != nil {
		l.Warn("ошибка подготовки запроса", zap.Error(err))

		return nil, err
	}

	l.Debug("аргументы запроса",
		zap.Int("лимит", p.Limit),
		zap.Int("смещение", p.Offset),
	)

	now := time.Now()
	rows, err := dao.db.Query(ctx, q, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return nil, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Role])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return nil, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("количество ролей", len(roles)))

	return roles, nil
}
