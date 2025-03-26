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
