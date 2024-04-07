package issued

import (
	"mime/multipart"
	"time"
)

type UploadReq struct {
	TargetGroups []string        `json:"target_groups"`
	Title        string          `json:"title"`
	Theme        string          `json:"theme"`
	Major        string          `json:"major"`
	File         *multipart.File `json:"file"`
}

type UploadResp struct {
	PracticeId   int       `json:"practice_id"`
	Title        string    `json:"title"`
	TargetGroups []string  `json:"target_groups"`
	UploadAt     time.Time `json:"upload_at"`
}

type DTO struct {
	AccountId           int
	TargetGroups        []string
	Title, Theme, Major string
	Path                string
	UploadAt            time.Time
	DeletedAt           *time.Time
}

type Entity struct {
	PracticeId   int        `db:"issued_practice_id"`
	AccountId    int        `db:"account_id"`
	TargetGroups []string   `db:"target_groups"`
	Title        string     `db:"title"`
	Theme        string     `db:"theme"`
	Major        string     `db:"major"`
	PracticePath string     `db:"practice_path"`
	UploadAt     time.Time  `db:"upload_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
