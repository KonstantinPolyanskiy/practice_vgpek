package entity

import "github.com/google/uuid"

type Person struct {
	UUID      uuid.UUID `db:"person_uuid"`
	AccountId int       `db:"account_id"`

	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
	LastName   string `db:"last_name"`
}
