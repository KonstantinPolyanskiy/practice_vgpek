package params

type Default struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Key struct {
	IsValid bool `json:"is_valid"`
	Default
}
