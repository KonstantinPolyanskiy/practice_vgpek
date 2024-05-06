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

// StateParams возвращает параметры состояния, если указаны неверное - возвращает params.All
func StateParams(r *http.Request, p params.Default) params.State {
	var result params.State

	result.Default = p

	state := r.URL.Query().Get("state")

	if state == "" {
		result.State = params.All
		return result
	}

	result.State = state

	return result
}
