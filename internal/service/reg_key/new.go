package reg_key

import (
	"context"
	"errors"
	"log"
	"practice_vgpek/internal/model/registration_key"
	"practice_vgpek/pkg/rndutils"
)

type CreatingKeyResult struct {
	CreatedKey registration_key.AddResp
	Error      error
}

func (s Service) NewKey(ctx context.Context, req registration_key.AddReq) (registration_key.AddResp, error) {
	resCh := make(chan CreatingKeyResult)

	go func() {
		if req.MaxCountUsages < 0 {
			resCh <- CreatingKeyResult{
				CreatedKey: registration_key.AddResp{},
				Error:      errors.New("неправильное кол-во использований ключа"),
			}
		}

		dto := registration_key.DTO{
			RoleId:         req.RoleId,
			Body:           rndutils.RandNumberString(5) + rndutils.RandString(5),
			MaxCountUsages: req.MaxCountUsages,
		}
		savedKey, err := s.r.SaveKey(ctx, dto)
		log.Println(err)
		if err != nil {
			resCh <- CreatingKeyResult{
				CreatedKey: registration_key.AddResp{},
				Error:      errors.New("ошибка сохранения ключа"),
			}
		}

		resp := registration_key.AddResp{
			RegKeyId:           savedKey.RegKeyId,
			MaxCountUsages:     savedKey.MaxCountUsages,
			CurrentCountUsages: savedKey.CurrentCountUsages,
			Body:               savedKey.Body,
			CreatedAt:          savedKey.CreatedAt,
		}

		resCh <- CreatingKeyResult{
			CreatedKey: resp,
			Error:      nil,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return registration_key.AddResp{}, ctx.Err()
		case result := <-resCh:
			return result.CreatedKey, result.Error
		}
	}
}
