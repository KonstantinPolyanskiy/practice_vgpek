package permissions

type ActionEntity struct {
	Id   int    `db:"internal_action_id"`
	Name string `db:"internal_action_name"`
}

type ActionDTO struct {
	Name string
}

type AddActionReq struct {
	Name string `json:"name"`
}

type AddActionResp struct {
	Name string `json:"name"`
}
