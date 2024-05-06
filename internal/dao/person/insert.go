package person

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/pkg/timeutils"
	"time"
)

func (dao DAO) Save(ctx context.Context, data dto.PersonRegistrationData) (entity.Person, error) {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.SavePersonDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	insertQuery := `INSERT INTO 
    				person (person_uuid, account_id, first_name, middle_name, last_name) 
					VALUES (@PersonUUID, @AccountId, @FirstName, @MiddleName, @LastName)
					RETURNING person_uuid`

	args := pgx.NamedArgs{
		"PersonUUID": data.UUID,
		"AccountId":  data.AccountId,
		"FirstName":  data.FirstName,
		"MiddleName": data.SecondName,
		"LastName":   data.LastName,
	}

	l.Debug("аргументы запроса",
		zap.String("uuid пользователя", args["PersonUUID"].(uuid.UUID).String()),
		zap.Int("id аккаунта", args["AccountId"].(int)),
		zap.String("имя", args["FirstName"].(string)),
		zap.String("фамилия", args["MiddleName"].(string)),
		zap.String("отчество", args["LastName"].(string)),
	)

	var personUUID uuid.UUID

	now := time.Now()
	err := dao.db.QueryRow(ctx, insertQuery, args).Scan(&personUUID)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Person{}, err
	}

	l.Debug(operation.Insert, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	selectQuery := `SELECT * FROM person WHERE person_uuid=@PersonUIID`

	args = pgx.NamedArgs{
		"PersonUIID": personUUID,
	}

	now = time.Now()
	rows, err := dao.db.Query(ctx, selectQuery, args)
	defer rows.Close()
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return entity.Person{}, err
	}

	l.Debug(operation.Select, zap.Duration("время выполнения", timeutils.TrackTime(now)))

	person, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Person])
	if err != nil {
		l.Error(operation.CollectError, zap.Error(err))
		return entity.Person{}, err
	}

	l.Info(operation.SuccessfullyReceived, zap.String("uuid пользователя", person.UUID.String()))

	return person, nil
}
