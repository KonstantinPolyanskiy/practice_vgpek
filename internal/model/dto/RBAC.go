package dto

import "time"

type NewRBACPart struct {
	Name        string
	Description string
	CreatedAt   time.Time
}

type DeleteInfo struct {
	DeleteTime time.Time
}
