package bankaccimpl

import (
	"context"
	"money-transfer-demo/pkg/payment/bankacc"
	"time"

	"gorm.io/gorm"
)

type service struct {
	store *store
}

func NewService(store *store) bankacc.Service {
	return &service{store: store}
}

func (s *service) CreateBankAccount(ctx context.Context, cmd *bankacc.CreateBankAccountCommand) error {
	exist, err := s.store.isBankAccountTaken(ctx, cmd.BankCode, cmd.AccountNo)
	if err != nil {
		return err
	}

	if len(exist) > 0 {
		return bankacc.ErrBankAccountAlreadyExists
	}

	return s.store.db.Transaction(func(tx *gorm.DB) error {
		entity := &bankacc.BankAccount{
			BankCode:           cmd.BankCode,
			AccountNo:          cmd.AccountNo,
			Balance:            cmd.Balance,
			Status:             bankacc.Active,
			Currency:           cmd.Currency,
			OutstandingBalance: 0,
			CreatedAt:          time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:          time.Now().UTC().Format(time.RFC3339),
		}

		err := s.store.createBankAccount(tx, entity)
		if err != nil {

			return err
		}
		return nil
	})

}

func (s *service) AdjustBalance(ctx context.Context, tx *gorm.DB, cmd *bankacc.AdjustBankAccountBalanceCommand) error {
	result, err := s.store.getBankAccount(ctx, cmd.BankAccountID)
	if err != nil {
		return err
	}

	if result == nil {
		return bankacc.ErrBankAccountNotFound
	}

	if result.Status != bankacc.Active {
		return bankacc.ErrBankAccountNotActive
	}

	if result.Currency != cmd.Currency {
		return bankacc.ErrInvalidCurrency
	}

	if result.Balance < -cmd.ChangedAmount {
		return bankacc.ErrInsufficientBalance
	}
	err = s.store.adjustBankAccountBalance(tx, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetBankAccountByID(ctx context.Context, id int64) (*bankacc.BankAccount, error) {
	result, err := s.store.getBankAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, bankacc.ErrBankAccountNotFound
	}

	if result.Status != bankacc.Active {
		return nil, bankacc.ErrBankAccountNotActive
	}

	return result, nil
}
