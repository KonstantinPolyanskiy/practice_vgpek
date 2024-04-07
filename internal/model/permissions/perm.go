package permissions

type AddPermReq struct {
	RoleId   int `json:"role_id"`
	ObjectId int `json:"object_id"`
	// Слайс действий, которые над объектом может производить роль
	ActionsId []int `json:"actions_id"`
}

type AddPermResp struct {
	AddPermReq `json:"added"`
}

type PermissionEntity struct {
	PermissionId int `db:"role_perm_id"`

	RoleEntity
	ActionEntity
	ObjectEntity
}
