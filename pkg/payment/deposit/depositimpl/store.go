package depositimpl

import (
	"context"
	"fmt"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/deposit"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database postgres.DBConnector) *store {
	return &store{db: database.GetDB()}
}

func (s *store) createDeposit(tx *gorm.DB, deposit *deposit.Deposit) (int64, error) {
	err := tx.Table("deposit").Create(deposit).Error

	if err != nil {
		return 0, err
	}

	return deposit.ID, nil
}

func (s *store) getDepositByStatus(tx *gorm.DB, id int64, status deposit.Status) (*deposit.Deposit, error) {
	var deposit deposit.Deposit

	err := tx.Where("id = ? AND status = ?", id, status).First(&deposit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &deposit, nil
}


func (s *store) updateDeposit(tx *gorm.DB, entity *deposit.Deposit) error {
	if err := tx.Model(&deposit.Deposit{}).Where("id = ?", entity.ID).Updates(map[string]interface{}{
		"status":        entity.Status,
		"gross_amount":  entity.GrossAmount,
		"net_amount":    entity.NetAmount,
		"charge_amount": entity.ChargeAmount,
		"detail":        entity.Detail,
		"updated_at":    entity.UpdatedAt,
		"updated_by":    entity.UpdatedBy,
	}).Error; err != nil {
		return fmt.Errorf("failed to update deposit: %w", err)
	}

	return nil
}

func (s *store) createDepositTimeline(tx *gorm.DB, timeline *deposit.DepositTimeline) error {
	return tx.Table("deposit_timeline").Create(timeline).Error
}

func (s *store) getDepositByTransactionID(ctx context.Context, transactionID string) (*deposit.Deposit, error) {
	var deposit deposit.Deposit

	err := s.db.WithContext(ctx).Table("deposit").Where("transaction_id = ?", transactionID).First(&deposit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &deposit, nil
}
