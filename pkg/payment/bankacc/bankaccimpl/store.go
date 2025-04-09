package bankaccimpl

import (
	"context"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/bankacc"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database postgres.DBConnector) *store {
	return &store{db: database.GetDB()}
}

func (s *store) isBankAccountTaken(ctx context.Context, bankCode, accountNo string) ([]*bankacc.BankAccount, error) {
	var bankAccounts []*bankacc.BankAccount

	err := s.db.WithContext(ctx).Where("bank_code = ? AND account_no = ?", bankCode, accountNo).Find(&bankAccounts).Error
	return bankAccounts, err
}

func (s *store) createBankAccount(tx *gorm.DB, bankAccount *bankacc.BankAccount) error {
	return tx.Create(bankAccount).Error
}

func (s *store) getBankAccount(ctx context.Context, id int64) (*bankacc.BankAccount, error) {
	var bankAccount bankacc.BankAccount

	err := s.db.WithContext(ctx).Where("id = ?", id).First(&bankAccount).Error
	return &bankAccount, err
}

func (s *store) adjustBankAccountBalance(tx *gorm.DB, cmd *bankacc.AdjustBankAccountBalanceCommand) error {
	result := tx.Model(&bankacc.BankAccount{}).
		Where("id = ?", cmd.BankAccountID).
		Updates(map[string]interface{}{
			"balance":    gorm.Expr("balance + ?", cmd.ChangedAmount),
			"updated_at": cmd.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
