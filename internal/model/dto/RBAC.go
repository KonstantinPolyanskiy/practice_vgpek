package dto

import "time"

type NewRBACPart struct {
	Name        string
	Description string
	CreatedAt   time.Time
}

type SetPermissionReq struct {
	RoleId    int   `json:"roleId"`
	ObjectId  int   `json:"objectId"`
	ActionsId []int `json:"actionsId"`
}

type NewRBACReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeleteInfo struct {
	DeleteTime time.Time
}
