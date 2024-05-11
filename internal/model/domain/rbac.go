package domain

import "time"

const (
	AdminRole   = "ADMIN"
	TeacherRole = "TEACHER"
	StudentRole = "STUDENT"
)

const (
	AddAction    = "ADD"
	GetAction    = "GET"
	EditAction   = "EDIT"
	DeleteAction = "DELETE"
)

const (
	AccountObject = "ACCOUNT"
	KeyObject     = "KEY"
	RBACObject    = "RBAC"
)

type Permissions struct {
	PermissionId int

	Role
	Action
	Object
}

type RolePermission struct {
	Role   Role
	Object ObjectWithActions
}

type ObjectWithActions struct {
	Object  Object
	Actions []Action
}

type RBACPart struct {
	ID int

	Name        string
	Description string

	CreatedAt time.Time

	IsDeleted bool
	DeletedAt *time.Time
}

type Action struct {
	ID int

	Name        string
	Description string

	CreatedAt time.Time

	IsDeleted bool
	DeletedAt *time.Time
}

func (a Action) Deleted() bool {
	return a.IsDeleted
}
func (a Action) Part() RBACPart {
	return RBACPart{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		IsDeleted:   a.IsDeleted,
		DeletedAt:   a.DeletedAt,
	}
}

type Object struct {
	ID int

	Name        string
	Description string

	CreatedAt time.Time

	IsDeleted bool
	DeletedAt *time.Time
}

func (o Object) Deleted() bool {
	return o.IsDeleted
}
func (o Object) Part() RBACPart {
	return RBACPart{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		CreatedAt:   o.CreatedAt,
		IsDeleted:   o.IsDeleted,
		DeletedAt:   o.DeletedAt,
	}
}

type Role struct {
	ID int

	Name        string
	Description string

	CreatedAt time.Time

	IsDeleted bool
	DeletedAt *time.Time
}

func (r Role) Deleted() bool {
	return r.IsDeleted
}
func (r Role) Part() RBACPart {
	return RBACPart{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		IsDeleted:   r.IsDeleted,
		DeletedAt:   r.DeletedAt,
	}
}
