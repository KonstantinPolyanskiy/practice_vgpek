package permissions

type AddPermReq struct {
	RoleId   int `json:"role_id"`
	ObjectId int `json:"object_id"`
	// Слайс действий, которые над объектом может производить роль
	ActionsId []int `json:"actions_id"`
}

type AddPermResp struct {
	Success string `json:"success"`
}
