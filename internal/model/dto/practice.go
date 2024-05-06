package dto

import (
	"mime/multipart"
	"time"
)

type NewIssuedPracticeReq struct {
	TargetGroups []string `json:"target_groups"`

	Title string `json:"title"`
	Theme string `json:"theme"`
	Major string `json:"major"`

	File *multipart.File `json:"file"`
}

type NewIssuedPractice struct {
	AccountId int

	TargetGroups []string

	Title string
	Theme string
	Major string

	Path     string
	UploadAt time.Time
}

type NewSolvedPracticeReq struct {
	PerformedAccountId int `json:"performed_account_id"`
	IssuedPracticeId   int `json:"issued_practice_id"`

	File *multipart.File `json:"file"`
}

type NewSolvedPractice struct {
	PerformedAccountId int
	IssuedPracticeId   int

	Mark     int
	MarkTime *time.Time

	SolvedTime *time.Time

	Path string

	IsDeleted *time.Time
}

type MarkPracticeReq struct {
	SolvedPracticeId int `json:"solved_practice_id"`
	Mark             int `json:"mark"`
}

type MarkPractice struct {
	// SolvedPracticeId id практической работы, которую необходимо оценить
	SolvedPracticeId int

	Mark     int
	MarkTime time.Time
}
