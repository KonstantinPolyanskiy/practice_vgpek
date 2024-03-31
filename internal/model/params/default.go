package params

type Default struct {
	NotDeleted bool `json:"not_deleted"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
}
