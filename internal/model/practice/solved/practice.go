package solved

import (
	"mime/multipart"
	"time"
)

type UploadReq struct {
	IssuedPracticeId int `json:"issued_practice_id"`
	File             *multipart.File
}

type UploadResp struct {
	SolvedPracticeId   int       `json:"solved_practice_id"`
	SolvedTime         time.Time `json:"solved_time"`
	IssuedPracticeName string    `json:"issued_practice_name"`
}

type DTO struct {
	IssuedPracticeId   int
	PerformedAccountId int

	Path string
}

type Entity struct {
	SolvedPracticeId int `db:"solved_practice_id"`

	PerformedAccountId int `db:"performed_account_id"`
	IssuedPracticeId   int `db:"issued_practice_id"`

	Mark     int        `db:"mark"`
	MarkTime *time.Time `db:"mark_time"`

	SolvedTime *time.Time `db:"solved_time"`

	Path string `db:"path"`

	IsDeleted *time.Time `db:"is_deleted"`
}
