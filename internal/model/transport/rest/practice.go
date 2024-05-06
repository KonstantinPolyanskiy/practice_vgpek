package rest

import (
	"practice_vgpek/internal/model/domain"
	"time"
)

func (p IssuedPractice) DomainToResponse(practice domain.IssuedPractice) IssuedPractice {
	return IssuedPractice{
		Id:           practice.Id,
		AuthorName:   practice.AuthorName,
		AuthorId:     practice.AuthorId,
		TargetGroups: practice.TargetGroups,
		Title:        practice.Title,
		Theme:        practice.Theme,
		Major:        practice.Major,
		UploadAt:     practice.UploadAt,
		IsDeleted:    practice.IsDeleted,
		DeletedAt:    practice.DeletedAt,
	}
}

type IssuedPractice struct {
	Id int `json:"id"`

	AuthorName string `json:"author_name"`
	AuthorId   int    `json:"author_id"`

	TargetGroups []string `json:"target_groups"`

	Title string `json:"title"`
	Theme string `json:"theme"`
	Major string `json:"major"`

	UploadAt time.Time `json:"upload_at"`

	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type SolvedPractice struct {
	Id               int `json:"id"`
	IssuedPracticeId int `json:"issued_practice_id"`

	IssuerName string `json:"issuer_name,omitempty"`

	AuthorName string `json:"author_name"`
	AuthorId   int    `json:"author_id"`

	Mark     int        `json:"mark"`
	MarkTime *time.Time `json:"mark_time,omitempty"`

	SolvedTime time.Time `json:"solved_time,omitempty"`

	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (p SolvedPractice) DomainToResponse(practice domain.SolvedPractice) SolvedPractice {
	return SolvedPractice{}
}

type IssuedPracticeWithLink struct {
	IssuedPractice
	DownloadLink string `json:"download_link"`
}

type SolvedPracticeWithLink struct {
	SolvedPractice
	DownloadLink string `json:"download_link"`
}

func (p IssuedPractice) WithDownloadLink(link string) IssuedPracticeWithLink {
	return IssuedPracticeWithLink{
		IssuedPractice: p,
		DownloadLink:   link,
	}
}

func (p SolvedPractice) WithDownloadLink(link string) SolvedPracticeWithLink {
	return SolvedPracticeWithLink{
		SolvedPractice: p,
		DownloadLink:   link,
	}
}
