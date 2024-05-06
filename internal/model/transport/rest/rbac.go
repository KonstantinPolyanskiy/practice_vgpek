package rest

import (
	"practice_vgpek/internal/model/domain"
	"time"
)

type RBACAttribute interface {
	Part() domain.RBACPart
}

type RBACPart struct {
	ID int `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`

	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func RBACPartDomainToResponse[T RBACAttribute](part T) RBACPart {
	return RBACPart{
		ID:          part.Part().ID,
		Name:        part.Part().Name,
		Description: part.Part().Description,
		CreatedAt:   part.Part().CreatedAt,
		IsDeleted:   part.Part().IsDeleted,
		DeletedAt:   part.Part().DeletedAt,
	}
}

func RBACPartsDomainToResponse[T RBACAttribute](parts []T) (result []RBACPart) {
	for _, part := range parts {
		result = append(result, RBACPart{
			ID:          part.Part().ID,
			Name:        part.Part().Name,
			Description: part.Part().Description,
			CreatedAt:   part.Part().CreatedAt,
			IsDeleted:   part.Part().IsDeleted,
			DeletedAt:   nil,
		})
	}

	return result
}
