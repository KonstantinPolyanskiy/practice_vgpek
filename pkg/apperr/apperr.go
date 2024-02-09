package apperr

type AppError struct {
	Action string `json:"action"`
	Error  string `json:"error"`
}
