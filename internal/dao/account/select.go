package account

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

func (dao DAO) ById(ctx context.Context, id int) (entity.Account, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectAccountByIdDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM account WHERE account_id=@AccountId`

	args := pgx.NamedArgs{
		"AccountId": id,
	}

	l.Debug("аргументы запроса", zap.Int("id аккаунта", args["AccountId"].(int)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Account{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Account])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Account{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("id аккаунта", account.Id))

	return account, nil
}

func (dao DAO) ByLogin(ctx context.Context, login string) (entity.Account, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectAccountByLoginDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	getQuery := `SELECT * FROM account WHERE login=@Login`

	args := pgx.NamedArgs{
		"Login": login,
	}

	l.Debug("аргументы запроса", zap.String("логин", args["Login"].(string)))

	now := time.Now()
	rows, err := dao.db.Query(ctx, getQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Account{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Account])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Account{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("id аккаунта", account.Id))

	return account, nil
}

func (dao DAO) ByParams(ctx context.Context, p params.Default) ([]entity.Account, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SelectAccountsByParamsDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	selectQuery := squirrel.Select("*").From("account").
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
	rows, err := dao.db.Query(ctx, q, args...)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return nil, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	persons, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Account])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return nil, err
	}

	l.Info(operation.SuccessfullyReceived, zap.Int("количество аккаунтов", len(persons)))

	return persons, nil
}
