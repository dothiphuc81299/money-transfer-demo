package memberaccimpl

import (
	"context"
	"money-transfer-demo/pkg/payment/memberacc"
	"time"

	"gorm.io/gorm"
)

type service struct {
	store *store
}

func NewService(store *store) memberacc.Service {
	return &service{store: store}
}

func (s *service) Create(ctx context.Context, cmd *memberacc.CreateMemberAccountCommand) error {
	result, err := s.store.getMemberAccount(ctx, cmd.MemberID, cmd.LoginName)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		return memberacc.ErrMemberAccountAlreadyExists
	}

	now := time.Now().UTC().Format(time.RFC3339)
	return s.store.db.Transaction(func(tx *gorm.DB) error {
		err = s.store.createMemberAccount(tx, &memberacc.MemberAccount{
			MemberID:           cmd.MemberID,
			LoginName:          cmd.LoginName,
			Currency:           cmd.Currency,
			Status:             cmd.Status,
			Balance:            0,
			OutstandingBalance: 0,
			CreatedAt:          now,
			UpdatedAt:          now,
		})
		if err != nil {
			return err
		}

		return nil
	})

}

func (s *service) GetByMemberID(ctx context.Context, memberID int64) (*memberacc.MemberAccount, error) {
	result, err := s.store.getMemberAccountByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, memberacc.ErrMemberAccountNotFound
	}

	return result, nil
}

func (s *service) AdjustMemberAccountBalance(ctx context.Context, tx *gorm.DB, cmd *memberacc.AdjustMemberAccountBalanceCommand) error {
	memberAccount, err := s.store.getMemberAccountByMemberID(ctx, cmd.MemberID)
	if err != nil {
		return err
	}

	if memberAccount == nil {
		return memberacc.ErrMemberAccountNotFound
	}

	ok := memberacc.VerifyBalance(memberAccount.Balance, memberAccount.OutstandingBalance, cmd.AdjustedAmount, cmd.AdjustedOutstandingAmount, cmd.TransactionType)
	if !ok {
		return memberacc.ErrMemberAccountNotEnoughBalance
	}

	now := time.Now().UTC().Format(time.RFC3339)
	err = s.store.updateMemberAccountBalance(tx, &memberacc.UpdateMemberAccountBalanceCommand{
		ID:                        memberAccount.ID,
		AdjustedAmount:            cmd.AdjustedAmount,
		AdjustedOutstandingAmount: cmd.AdjustedOutstandingAmount,
		UpdatedAt:                 now,
	})
	if err != nil {
		return err
	}

	return nil

}
