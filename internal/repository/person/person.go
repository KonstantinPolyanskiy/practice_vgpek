package person

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"practice_vgpek/internal/model/person"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewPersonRepo(db *pgxpool.Pool) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) SavePerson(ctx context.Context, savingPerson person.DTO, accountId int) (person.Entity, error) {
	var insertedPersonId int

	insertPersonQuery := `
	INSERT INTO person (account_id, first_name, middle_name, last_name)
	VALUES (@AccountId, @FirstName, @MiddleName, @LastName)
	RETURNING person_id
`
	args := pgx.NamedArgs{
		"AccountId":  accountId,
		"FirstName":  savingPerson.FirstName,
		"MiddleName": savingPerson.MiddleName,
		"LastName":   savingPerson.LastName,
	}

	err := r.db.QueryRow(ctx, insertPersonQuery, args).Scan(&insertedPersonId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return person.Entity{}, errors.New("сохраненный пользователь не найден")
		}
		return person.Entity{}, err
	}

	getPersonQuery := `
	SELECT * FROM person
	WHERE account_id=$1
`
	row, err := r.db.Query(ctx, getPersonQuery, insertedPersonId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return person.Entity{}, errors.New("сохраненный пользователь не найдет")
		}
		return person.Entity{}, err
	}

	savedPerson, err := pgx.CollectOneRow(row, pgx.RowToStructByName[person.Entity])
	if err != nil {
		return person.Entity{}, err
	}

	return savedPerson, nil
}
