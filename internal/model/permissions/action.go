package permissions

type GetActionReq struct {
	Id int `json:"action_id"`
}

type GetActionResp struct {
	Id   int    `json:"action_id"`
	Name string `json:"action_name"`
}

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
