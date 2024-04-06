package person

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/dberr"
	"practice_vgpek/internal/model/operation"
	"practice_vgpek/internal/model/person"
)

type Repository struct {
	l  *zap.Logger
	db *pgxpool.Pool
}

func NewPersonRepo(db *pgxpool.Pool, logger *zap.Logger) Repository {
	return Repository{
		l:  logger,
		db: db,
	}
}

func (r Repository) SavePerson(ctx context.Context, savingPerson person.DTO, accountId int) (person.Entity, error) {
	l := r.l.With(
		zap.String("операция", operation.NewPersonOperation),
		zap.String("слой", "репозиторий"),
	)

	var insertedPersonUUID uuid.UUID

	insertPersonQuery := `
	INSERT INTO person (person_uuid, account_id, first_name, middle_name, last_name)
	VALUES (@PersonUUID, @AccountId, @FirstName, @MiddleName, @LastName)
	RETURNING person_uuid
`

	args := pgx.NamedArgs{
		"PersonUUID": uuid.New(),
		"AccountId":  accountId,
		"FirstName":  savingPerson.FirstName,
		"MiddleName": savingPerson.MiddleName,
		"LastName":   savingPerson.LastName,
	}

	// Если запрос не возвращает Id, то пользователь не создан
	err := r.db.QueryRow(ctx, insertPersonQuery, args).Scan(&insertedPersonUUID)
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		return person.Entity{}, err
	}

	getPersonQuery := `
	SELECT * FROM person
	WHERE person_uuid=$1
`

	row, err := r.db.Query(ctx, getPersonQuery, insertedPersonUUID)
	defer row.Close()
	if err != nil {
		l.Warn("ошибка выполнения запроса", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return person.Entity{}, dberr.ErrNotFound
		}
		return person.Entity{}, err
	}

	savedPerson, err := pgx.CollectOneRow(row, pgx.RowToStructByName[person.Entity])
	if err != nil {
		l.Warn("ошибка записи данных в структуру", zap.Error(err))

		return person.Entity{}, err
	}

	return savedPerson, nil
}
