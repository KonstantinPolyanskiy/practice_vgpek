package permissions

type GetRoleResp struct {
	Id   int    `json:"role_id"`
	Name string `json:"role_name"`
}

type GetRolesResp struct {
	Roles []GetRoleResp `json:"roles"`
}

type AddRoleReq struct {
	Name string `json:"name"`
}

type AddRoleResp struct {
	Name string `json:"name"`
}

type RoleDTO struct {
	Name string
}

type RoleEntity struct {
	Id   int    `db:"internal_role_id"`
	Name string `db:"role_name"`
}
