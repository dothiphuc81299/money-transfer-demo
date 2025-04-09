package depositimpl

import (
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
	err := tx.Create(deposit).Error
	if err != nil {
		return 0, err
	}

	return deposit.ID, nil
}

func (s *store) getDepositByStatus(tx *gorm.DB, id int64, status deposit.Status) (*deposit.Deposit, error) {
	var deposit deposit.Deposit

	err := tx.Where("id = ? AND status = ?", id, status).First(&deposit).Error
	if err != nil {
		return nil, err
	}

	return &deposit, nil
}

func (s *store) updateDeposit(tx *gorm.DB, deposit *deposit.Deposit) error {
	return tx.Save(deposit).Error
}

func (s *store) createDepositTimeline(tx *gorm.DB, timeline *deposit.DepositTimeline) error {
	return tx.Create(timeline).Error
}
