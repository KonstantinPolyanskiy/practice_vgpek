package entity

import "time"

type IssuedPractice struct {
	Id int `db:"issued_practice_id"`

	AccountId int `db:"account_id"`

	TargetGroups []string `db:"target_groups"`

	Title string `db:"title"`
	Theme string `db:"theme"`
	Major string `db:"major"`

	Path string `db:"practice_path"`

	UploadAt  time.Time  `db:"upload_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type SolvedPractice struct {
	Id int `db:"solved_practice_id"`

	PerformedAccountId int `db:"performed_account_id"`
	IssuedPracticeId   int `db:"issued_practice_id"`

	Mark     int `db:"mark"`
	MarkTime *time.Time

	SolvedTime *time.Time

	Path      string `db:"path"`
	IsDeleted *time.Time
}

// SolvedPracticeUpdate структура для обновления записи. Если поле nil - поле в запрос не попадает
type SolvedPracticeUpdate struct {
	Id int

	PerformedAccountId *int
	IssuedPracticeId   *int

	Mark     *int
	MarkTime *time.Time

	SolvedTime *time.Time

	Path      *string
	IsDeleted *time.Time
}
