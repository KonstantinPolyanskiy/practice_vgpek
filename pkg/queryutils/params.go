package queryutils

import (
	"net/http"
	"practice_vgpek/internal/model/params"
	"strconv"
)

func DefaultParams(r *http.Request, defaultLimit, defaultOffset int) (params.Default, error) {
	var result params.Default

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return result, err
	}

	if limit == 0 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		return result, err
	}

	if offset == 0 {
		offset = defaultOffset
	}

	result.Limit = limit
	result.Offset = offset

	return result, nil
}
