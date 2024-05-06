package account

import (
	"context"
	"go.uber.org/zap"
	"practice_vgpek/internal/model/layer"
	"practice_vgpek/internal/model/operation"
)

func (dao DAO) HardDeleteById(ctx context.Context, id int) error {
	l := dao.logger.With(
		zap.String(operation.Operation, operation.HardDeleteAccountByIdDAO),
		zap.String(layer.Layer, layer.DataLayer),
	)

	deleteQuery := `DELETE FROM account WHERE account_id=$1`

	l.Debug("аргументы запроса", zap.Int("id аккаунта", id))

	_, err := dao.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		l.Error(operation.ExecuteError, zap.Error(err))
		return err
	}

	return nil
}
