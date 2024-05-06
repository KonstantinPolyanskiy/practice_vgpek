package dao

import (
	"context"
	"github.com/google/uuid"
	"practice_vgpek/internal/model/dto"
	"practice_vgpek/internal/model/entity"
	"practice_vgpek/internal/model/params"
)

type ActionDAO interface {
	ById(ctx context.Context, id int) (entity.Action, error)
	SoftDeleteById(ctx context.Context, id int, info dto.DeleteInfo) error
	Save(ctx context.Context, action dto.NewRBACPart) (entity.Action, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Action, error)
}

type RoleDAO interface {
	ById(ctx context.Context, id int) (entity.Role, error)
	SoftDeleteById(ctx context.Context, id int, info dto.DeleteInfo) error
	Save(ctx context.Context, role dto.NewRBACPart) (entity.Role, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Role, error)
}

type ObjectDAO interface {
	ById(ctx context.Context, id int) (entity.Object, error)
	SoftDeleteById(ctx context.Context, id int, info dto.DeleteInfo) error
	Save(ctx context.Context, role dto.NewRBACPart) (entity.Object, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Object, error)
}

type PersonDAO interface {
	Save(ctx context.Context, data dto.PersonRegistrationData) (entity.Person, error)

	ByUUID(ctx context.Context, uid uuid.UUID) (entity.Person, error)
	ByAccountId(ctx context.Context, accountId int) (entity.Person, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Person, error)
}

type AccountDAO interface {
	Save(ctx context.Context, data dto.AccountRegistrationData) (entity.Account, error)

	ById(ctx context.Context, id int) (entity.Account, error)
	ByLogin(ctx context.Context, login string) (entity.Account, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Account, error)

	HardDeleteById(ctx context.Context, id int) error
}

type KeyDAO interface {
	Save(ctx context.Context, data dto.NewKeyInfo) (entity.Key, error)

	ById(ctx context.Context, id int) (entity.Key, error)
	ByBody(ctx context.Context, body string) (entity.Key, error)
	ByParams(ctx context.Context, p params.Default) ([]entity.Key, error)

	Update(ctx context.Context, old entity.Key) (entity.Key, error)
}

type PermissionDAO interface {
	ByRoleId(ctx context.Context, roleId int) ([]entity.Permissions, error)
	Save(ctx context.Context, roleId, objectId int, actionsId []int) error
}

type IssuedPracticeDAO interface {
	Save(ctx context.Context, data dto.NewIssuedPractice) (entity.IssuedPractice, error)
	ById(ctx context.Context, id int) (entity.IssuedPractice, error)
}

type SolvedPracticeDAO interface {
	Save(ctx context.Context, data dto.NewSolvedPractice) (entity.SolvedPractice, error)
	ById(ctx context.Context, id int) (entity.SolvedPractice, error)
	Update(ctx context.Context, old entity.SolvedPracticeUpdate) (entity.SolvedPractice, error)
}
