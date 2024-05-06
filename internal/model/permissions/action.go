package permissions

const (
	GetName    = "GET"
	AddName    = "ADD"
	EditName   = "EDIT"
	DeleteName = "DEL"
)

const (
	AdminRole   = "ADMIN"
	UserRole    = "STUDENT"
	TeacherRole = "TEACHER"
)

const (
	UserObject              = "USER"
	MarkObject              = "MARK"
	IssuedPracticeObject    = "ISSUED_PRACTICE"
	CompletedPracticeObject = "COMPLETED_PRACTICE"
	KeyObject               = "KEY"
	RBACObject              = "RBAC"
)

type GetActionReq struct {
	Id int `json:"action_id"`
}

type GetActionResp struct {
	Id   int    `json:"action_id"`
	Name string `json:"action_name"`
}

type GetActionsResp struct {
	Actions []GetActionResp `json:"actions"`
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
