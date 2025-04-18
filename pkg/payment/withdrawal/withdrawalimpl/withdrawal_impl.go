package withdrawalimpl

import (
	"context"
	"encoding/json"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/bankacc"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/transaction"
	"money-transfer-demo/pkg/payment/withdrawal"
	"money-transfer-demo/pkg/util/generator"
	"time"

	"gorm.io/gorm"
)

type service struct {
	store        *store
	memberAccSrv memberacc.Service
	bankAccSrv   bankacc.Service
}

func NewService(store *store, memberAccSrv memberacc.Service, bankAccSrv bankacc.Service) withdrawal.Service {
	return &service{store: store, memberAccSrv: memberAccSrv, bankAccSrv: bankAccSrv}
}

func (s *service) CreateWithdrawal(ctx context.Context, cmd *withdrawal.CreateWithdrawalCommand) error {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return withdrawal.ErrUnauthorized
	}

	if account.AccountType != token.Member {
		return withdrawal.ErrUnauthorizedUserIsNotMember
	}

	cmd.CreatedBy = account.LoginName
	cmd.MemberID = account.ID

	member, err := s.memberAccSrv.GetByMemberID(ctx, cmd.MemberID)
	if err != nil {
		return err
	}

	cmd.LoginName = member.LoginName
	cmd.TransactionID = generator.GenerateTransactionID(string(cmd.PaymentMethodCode))

	err = s.getWithdrawalDetails(cmd)
	if err != nil {
		return err
	}

	return s.store.db.Transaction(func(tx *gorm.DB) error {
		err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
			MemberID:                  cmd.MemberID,
			AdjustedOutstandingAmount: -cmd.Amount,
			AdjustedAmount:            cmd.Amount,
			TransactionID:             cmd.TransactionID,
			TransactionType:           transaction.WithdrawalType,
		})
		if err != nil {
			return err
		}

		entity := &withdrawal.Withdrawal{
			MemberID:          cmd.MemberID,
			LoginName:         cmd.LoginName,
			TransactionID:     cmd.TransactionID,
			PaymentMethodCode: cmd.PaymentMethodCode,
			Amount:            cmd.Amount,
			Detail:            cmd.DetailStr,
			Status:            withdrawal.Pending,
			CreatedAt:         time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:         time.Now().UTC().Format(time.RFC3339),
			Currency:          string(member.Currency),
		}
		id, err := s.store.createWithdrawal(tx, entity)
		if err != nil {
			return err
		}

		timeline, err := json.Marshal(&withdrawal.TimelineDetail{
			WithdrawalStatus: withdrawal.Pending,
			TransactionID:    cmd.TransactionID,
			PaymentMethod:    cmd.PaymentMethodCode,
			Amount:           cmd.Amount,
			Bank:             cmd.MemberBankCode,
		})
		if err != nil {
			return err
		}
		err = s.store.createWithdrawalTimeline(tx, &withdrawal.WithdrawalTimeline{
			WithdrawalID:      id,
			Message:           withdrawal.Message(false, withdrawal.Pending, cmd.LoginName),
			AdditionalContent: timeline,
			CreatedAt:         time.Now().UTC().Format(time.RFC3339),
			CreatedBy:         cmd.LoginName,
		})

		if err != nil {
			return err
		}
		return nil
	})

}

func (s *service) getWithdrawalDetails(cmd *withdrawal.CreateWithdrawalCommand) error {
	var details withdrawal.WithdrawalDetail

	err := details.ValidateWithdrawal(cmd.Detail)
	if err != nil {
		return err
	}

	cmd.MemberBankCode = details.MemberBankCode
	dt, err := json.Marshal(details)
	if err != nil {
		return err
	}

	cmd.DetailStr = string(dt)
	return nil
}

func (s *service) SearchWithdrawal(ctx context.Context, query *withdrawal.SearchWithdrawalQuery) (*withdrawal.SearchWithdrawalResult, error) {
	if query.Page == 0 {
		query.Page = 1
	}

	if query.PerPage == 0 {
		query.PerPage = 10
	}

	result, err := s.store.searchWithdrawal(ctx, query)
	if err != nil {
		return nil, err
	}

	result.Page = query.Page
	result.PerPage = query.PerPage
	return result, nil
}

func (s *service) GetWithdrawalByID(ctx context.Context, id int64) (*withdrawal.WithdrawalDTO, error) {
	result, err := s.store.getWithdrawal(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, withdrawal.ErrWithdrawalNotFound
	}

	return result, nil
}

func (s *service) ApproveWithdrawal(ctx context.Context, id int64) error {
	panic("TOD")
}

func (s *service) RejectWithdrawal(ctx context.Context, id int64) error {
	panic("TOD")
}

func (s *service) TransferWithdrawal(ctx context.Context, id int64) error {
	panic("TOD")
}
