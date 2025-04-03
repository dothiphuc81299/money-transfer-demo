package bankacc

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	CreateBankAccount(ctx context.Context, cmd *CreateBankAccountCommand) error
	GetBankAccountByID(ctx context.Context, id int64) (*BankAccount, error)
	AdjustBalance(ctx context.Context, tx *gorm.DB, cmd *AdjustBankAccountBalanceCommand) error
}
