package rbac

import (
	"fmt"
	"practice_vgpek/internal/model/domain"
)

type partResult struct {
	part  Part
	error error
}

type partsResult struct {
	parts []Part
	error error
}

func sendPartResult[T Part](resCh chan partResult, resp T, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	resCh <- partResult{
		part:  resp,
		error: err,
	}
}

func sendPartsResult[T Part](resCh chan partsResult, resp []T, errMsg string) {
	var err error

	if errMsg != "" {
		err = fmt.Errorf(errMsg)
	}

	parts := make([]Part, 0)

	for _, a := range resp {
		parts = append(parts, a)
	}

	resCh <- partsResult{
		parts: parts,
		error: err,
	}
}

type Deletable interface {
	Deleted() bool
}

type Part interface {
	Part() domain.RBACPart
}

// filterDeleted возвращает только удаленные элементы
func filterDeleted[T Deletable](items []T) (result []T) {
	for _, item := range items {
		if item.Deleted() {
			result = append(result, item)
		}
	}
	return result
}

// filterNotDeleted возвращает только не удаленные элементы
func filterNotDeleted[T Deletable](items []T) (result []T) {
	for _, item := range items {
		if !item.Deleted() {
			result = append(result, item)
		}
	}
	return result
}
