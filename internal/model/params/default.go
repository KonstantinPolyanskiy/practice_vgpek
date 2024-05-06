package params

type Default struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

const (
	All        = "all"
	Deleted    = "deleted"
	NotDeleted = "not_deleted"
)

type State struct {
	// Состояние объекта: удален или нет, или все
	State string `json:"state"`
	Default
}

type IssuedPractice struct {
	IsSolved string `json:"is_solved"`
	Default
}
