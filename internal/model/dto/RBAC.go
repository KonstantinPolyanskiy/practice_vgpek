package dto

import "time"

type NewRBACPart struct {
	Name        string
	Description string
	CreatedAt   time.Time
}

type SetPermissionReq struct {
	RoleId    int   `json:"role_id"`
	ObjectId  int   `json:"object_id"`
	ActionsId []int `json:"actions_id"`
}

type NewRBACReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeleteInfo struct {
	DeleteTime time.Time
}
