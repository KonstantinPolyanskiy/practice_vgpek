package entity

type Permissions struct {
	PermissionId int `db:"role_perm_id"`

	Role
	Action
	Object
}
