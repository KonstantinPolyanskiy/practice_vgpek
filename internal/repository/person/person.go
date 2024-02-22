package person

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
		zap.String("action query", "save person"),
		zap.String("layer", "repo"),
	)

	var insertedPersonId int

	insertPersonQuery := `
	INSERT INTO person (account_id, first_name, middle_name, last_name)
	VALUES (@AccountId, @FirstName, @MiddleName, @LastName)
	RETURNING person_id
`
	l.Debug("insert person", zap.String("query", insertPersonQuery))

	args := pgx.NamedArgs{
		"AccountId":  accountId,
		"FirstName":  savingPerson.FirstName,
		"MiddleName": savingPerson.MiddleName,
		"LastName":   savingPerson.LastName,
	}

	l.Debug("args in query",
		zap.Int("account id", accountId),
		zap.String("first name", savingPerson.FirstName),
		zap.String("middle name", savingPerson.MiddleName),
		zap.String("last name", savingPerson.LastName),
	)

	// Если запрос не возвращает Id, то пользователь не создан
	err := r.db.QueryRow(ctx, insertPersonQuery, args).Scan(&insertedPersonId)
	if err != nil {
		l.Warn("error insert person", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return person.Entity{}, errors.New("сохраненный пользователь не найден")
		}
		return person.Entity{}, err
	}

	getPersonQuery := `
	SELECT * FROM person
	WHERE account_id=$1
`

	l.Debug("get person", zap.String("query", getPersonQuery))

	row, err := r.db.Query(ctx, getPersonQuery, insertedPersonId)
	if err != nil {
		l.Warn("error get inserted person", zap.Error(err))

		if errors.Is(err, pgx.ErrNoRows) {
			return person.Entity{}, errors.New("сохраненный пользователь не найдет")
		}
		return person.Entity{}, err
	}

	savedPerson, err := pgx.CollectOneRow(row, pgx.RowToStructByName[person.Entity])
	if err != nil {
		l.Warn("error collect person in struct", zap.Error(err))

		return person.Entity{}, err
	}

	return savedPerson, nil
}
