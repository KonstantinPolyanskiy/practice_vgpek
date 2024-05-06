package account

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

func (dao DAO) Save(ctx context.Context, data dto.AccountRegistrationData) (entity.Account, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SaveAccountDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO
					account (login, password_hash, internal_role_id, reg_key_id) 
					VALUES  (@Login, @PasswordHash, @RoleId, @KeyId)
					RETURNING account_id`

	args := pgx.NamedArgs{
		"Login":        data.Login,
		"PasswordHash": data.PasswordHash,
		"RoleId":       data.RoleId,
		"KeyId":        data.KeyId,
	}

	l.Debug("аргументы запроса",
		zap.String("login", args["Login"].(string)),
		zap.Int("id роли", args["RoleId"].(int)),
		zap.Int("id ключа", args["KeyId"].(int)),
	)

	var id int

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&id)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Account{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	selectQuery := `SELECT * FROM account WHERE account_id=@AccountId`

	args = pgx.NamedArgs{
		"AccountId": id,
	}

	now = time.Now()
	rows, err := dao.db.Query(ctx, selectQuery, args)
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

	return account, nil
}
