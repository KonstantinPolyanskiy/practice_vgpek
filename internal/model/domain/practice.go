package domain

import "time"

type IssuedPractice struct {
	Id int

	AuthorName string
	AuthorId   int

	TargetGroups []string

	Title string
	Theme string
	Major string

	Path string

	UploadAt time.Time

	IsDeleted bool
	DeletedAt *time.Time
}

type SolvedPractice struct {
	Id               int
	IssuedPracticeId int

	// IssuerName - ФИО того, что выслал практическое задание
	IssuerName string

	// AuthorName - ФИО того, что сделал практическую работу
	AuthorName string
	AuthorId   int

	Mark     int
	MarkTime *time.Time

	SolvedTime time.Time

	Path string

	IsDeleted bool
	DeletedAt *time.Time
}
