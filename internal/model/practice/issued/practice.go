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
	Status string `json:"status"`
}

type DTO struct {
	AccountId           int
	TargetGroups        []string
	Title, Theme, Major string
	Path                string
	UploadAt            time.Time
	DeletedAt           *time.Time
}
