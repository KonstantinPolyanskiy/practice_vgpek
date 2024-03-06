package permissions

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
