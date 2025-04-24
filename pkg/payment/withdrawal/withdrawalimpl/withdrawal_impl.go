package withdrawalimpl

import (
	"context"
	"encoding/json"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/bankacc"
	"money-transfer-demo/pkg/payment/config"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/memberpayacc"
	"money-transfer-demo/pkg/payment/transaction"
	"money-transfer-demo/pkg/payment/withdrawal"
	"money-transfer-demo/pkg/util/generator"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type service struct {
	store                 *store
	memberAccSrv          memberacc.Service
	bankAccSrv            bankacc.Service
	memberPaymentAccStore memberpayacc.Store
	cfg                   *config.Config
}

func NewService(store *store, memberAccSrv memberacc.Service, bankAccSrv bankacc.Service, memberPaymentAccStore memberpayacc.Store, cfg *config.Config) withdrawal.Service {
	return &service{store: store, memberAccSrv: memberAccSrv, bankAccSrv: bankAccSrv, memberPaymentAccStore: memberPaymentAccStore, cfg: cfg}
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

	cmd.MemberCurrency = string(member.Currency)
	err = s.getWithdrawalDetails(ctx, cmd)
	if err != nil {
		return err
	}

	if cmd.PaymentMethodCode == string(withdrawal.LBT) {
		cmd.Currency = string(member.Currency)
	}

	return s.store.db.Transaction(func(tx *gorm.DB) error {
		var adjustAmount float64
		if cmd.PaymentMethodCode == string(withdrawal.PAYPAL) && cmd.Currency != string(member.Currency) {
			adjustAmount = cmd.Amount * s.cfg.ExchangeVNDRate
		} else {
			adjustAmount = cmd.Amount
		}

		err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
			MemberID:                  cmd.MemberID,
			AdjustedOutstandingAmount: adjustAmount,
			TransactionID:             cmd.TransactionID,
			TransactionType:           transaction.WithdrawalType,
		})
		if err != nil {
			return err
		}

		entity := &withdrawal.Withdrawal{
			MemberID:               cmd.MemberID,
			LoginName:              cmd.LoginName,
			TransactionID:          cmd.TransactionID,
			PaymentMethodCode:      cmd.PaymentMethodCode,
			MemberPaymentAccountID: cmd.MemberPaymentAccountID,
			GrossAmount:            cmd.Amount,
			NetAmount:              cmd.Amount,
			Detail:                 cmd.DetailStr,
			Status:                 withdrawal.Pending,
			CreatedAt:              time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:              time.Now().UTC().Format(time.RFC3339),
			Currency:               cmd.Currency,
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
			PaypalEmail:      cmd.PaypalEmail,
		})
		if err != nil {
			return err
		}
		err = s.store.createWithdrawalTimeline(tx, &withdrawal.WithdrawalTimeline{
			WithdrawalID:      id,
			Message:           withdrawal.Message(withdrawal.Pending, cmd.LoginName),
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

func (s *service) getWithdrawalDetails(ctx context.Context, cmd *withdrawal.CreateWithdrawalCommand) error {
	var detail withdrawal.WithdrawalDetail

	mpa, err := s.memberPaymentAccStore.Get(ctx, cmd.MemberPaymentAccountID)
	if err != nil {
		return err
	}

	detail.MemberBankCode = mpa.MemberBankCode
	detail.MemberAccountNo = mpa.MemberAccountNo
	detail.MemberAccountName = mpa.MemberAccountName
	detail.MemberFullName = mpa.MemberFullName
	detail.PaypalEmail = mpa.PaypalEmail
	detail.MemberCurrency = cmd.MemberCurrency

	cmd.MemberBankCode = detail.MemberBankCode
	cmd.PaypalEmail = detail.PaypalEmail

	dt, err := json.Marshal(detail)
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

func (s *service) ApproveWithdrawal(ctx context.Context, cmd *withdrawal.UpdateWithdrawalStatusCommand) error {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return withdrawal.ErrUnauthorized
	}

	cmd.UpdatedBy = account.LoginName

	result, err := s.store.getWithdrawal(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if result == nil {
		return withdrawal.ErrWithdrawalNotFound
	}

	previousStatuses := map[withdrawal.Status]struct{}{
		withdrawal.Transferring: {},
	}

	err = s.validateStatus(result, previousStatuses)
	if err != nil {
		return err
	}

	cmd.TransactionID = result.TransactionID
	cmd.BankAccountID = result.BankAccountID
	cmd.NetAmount = result.NetAmount
	cmd.ChargeAmount = result.ChargeAmount

	return s.store.db.Transaction(func(tx *gorm.DB) error {
		err = s.updateWithdrawal(tx, result, cmd)
		if err != nil {
			return err
		}

		if cmd.BankAccountID != 0 {
			err := s.bankAccSrv.AdjustBalance(ctx, tx, &bankacc.AdjustBankAccountBalanceCommand{
				BankAccountID: cmd.BankAccountID,
				Currency:      result.Currency,
				ChangedAmount: -result.NetAmount,
				UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
			})
			if err != nil {
				return err
			}
		}

		err = s.memberPaymentAccStore.UpdateVerifyStatus(tx, &memberpayacc.MemberPayAccount{
			ID:           result.MemberPaymentAccountID,
			VerifyStatus: memberpayacc.Verified,
			UpdatedBy:    cmd.UpdatedBy,
			UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
			MemberID:                  result.MemberID,
			UpdatedBy:                 cmd.UpdatedBy,
			AdjustedAmount:            -result.GrossAmount,
			AdjustedOutstandingAmount: -result.GrossAmount,
			TransactionID:             cmd.TransactionID,
			TransactionType:           transaction.WithdrawalType,
		})
		if err != nil {
			return err
		}

		return nil
	})

}

func (s *service) RejectWithdrawal(ctx context.Context, cmd *withdrawal.UpdateWithdrawalStatusCommand) error {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return withdrawal.ErrUnauthorized
	}

	cmd.UpdatedBy = account.LoginName

	result, err := s.store.getWithdrawal(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if result == nil {
		return withdrawal.ErrWithdrawalNotFound
	}

	previousStatuses := map[withdrawal.Status]struct{}{
		withdrawal.Transferring: {},
		withdrawal.Pending:      {},
		withdrawal.Reviewing:    {},
	}

	err = s.validateStatus(result, previousStatuses)
	if err != nil {
		return err
	}

	cmd.TransactionID = result.TransactionID
	cmd.BankAccountID = result.BankAccountID
	cmd.NetAmount = result.NetAmount
	cmd.ChargeAmount = result.ChargeAmount

	return s.store.db.Transaction(func(tx *gorm.DB) error {
		err = s.updateWithdrawal(tx, result, cmd)
		if err != nil {
			return err
		}

		err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
			MemberID:                  result.MemberID,
			AdjustedOutstandingAmount: -result.GrossAmount,
			UpdatedBy:                 cmd.UpdatedBy,
			TransactionID:             cmd.TransactionID,
			TransactionType:           transaction.WithdrawalType,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) TransferWithdrawal(ctx context.Context, cmd *withdrawal.UpdateWithdrawalStatusCommand) error {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return withdrawal.ErrUnauthorized
	}

	cmd.UpdatedBy = account.LoginName

	result, err := s.store.getWithdrawal(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if result == nil {
		return withdrawal.ErrWithdrawalNotFound
	}

	previousStatuses := map[withdrawal.Status]struct{}{
		withdrawal.Reviewing: {},
	}

	err = s.validateStatus(result, previousStatuses)
	if err != nil {
		return err
	}

	cmd.TransactionID = result.TransactionID
	if cmd.ChargeAmount > 0 {
		cmd.NetAmount = result.GrossAmount - cmd.ChargeAmount
	}

	if cmd.BankAccountID > 0 {
		_, err := s.bankAccSrv.GetBankAccountByID(ctx, cmd.BankAccountID)
		if err != nil {
			return err
		}
	}

	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		if result.PaymentMethodCode == string(withdrawal.PAYPAL) {
			err = s.createSinglePaypal(ctx, tx, result, cmd)
			if err != nil {
				return err
			}
		}

		err = s.updateWithdrawal(tx, result, cmd)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) ReviewWithdrawal(ctx context.Context, cmd *withdrawal.UpdateWithdrawalStatusCommand) error {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return withdrawal.ErrUnauthorized
	}

	cmd.UpdatedBy = account.LoginName
	result, err := s.store.getWithdrawal(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if result == nil {
		return withdrawal.ErrWithdrawalNotFound
	}

	previousStatuses := map[withdrawal.Status]struct{}{
		withdrawal.Pending: {},
	}

	err = s.validateStatus(result, previousStatuses)
	if err != nil {
		return err
	}

	cmd.TransactionID = result.TransactionID
	cmd.NetAmount = result.NetAmount

	detail := &withdrawal.WithdrawalDetail{
		MemberBankCode:    result.MemberBankCode,
		MemberAccountNo:   result.MemberAccountNo,
		MemberAccountName: result.MemberAccountName,
		MemberFullName:    result.MemberFullName,
		MemberCurrency:    result.Currency,
	}

	dt, err := json.Marshal(detail)
	if err != nil {
		return err
	}

	cmd.Detail = string(dt)

	return s.store.db.Transaction(func(tx *gorm.DB) error {
		err = s.updateWithdrawal(tx, result, cmd)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *service) validateStatus(w *withdrawal.WithdrawalDTO, previousStatuses map[withdrawal.Status]struct{}) error {
	_, exists := previousStatuses[w.Status]
	if exists {
		return nil
	}

	return withdrawal.ErrWithdrawalNotFound
}

func (s *service) updateWithdrawal(tx *gorm.DB, entity *withdrawal.WithdrawalDTO, cmd *withdrawal.UpdateWithdrawalStatusCommand) error {
	now := time.Now().UTC().Format(time.RFC3339)

	wDetail := &withdrawal.WithdrawalDetail{
		MemberBankCode:    entity.MemberBankCode,
		MemberAccountNo:   entity.MemberAccountNo,
		MemberAccountName: entity.MemberAccountName,
		MemberFullName:    entity.MemberFullName,
		MemberCurrency:    entity.Currency,
		PaypalEmail:       entity.PaypalEmail,
		PayoutBatchID:     entity.PayoutBatchID,
	}

	dt, err := json.Marshal(wDetail)
	if err != nil {
		return err
	}

	cmd.Detail = string(dt)
	err = s.store.updateWithdrawal(tx, &withdrawal.Withdrawal{
		ID:            cmd.ID,
		BankAccountID: cmd.BankAccountID,
		Detail:        cmd.Detail,
		ChargeAmount:  cmd.ChargeAmount,
		NetAmount:     cmd.NetAmount,
		Status:        cmd.Status,
		UpdatedAt:     now,
	})
	if err != nil {
		return err
	}

	message := withdrawal.Message(cmd.Status, cmd.UpdatedBy)
	detail, err := json.Marshal(&withdrawal.TimelineDetail{
		WithdrawalStatus: cmd.Status,
		Note:             cmd.Note,
		TransactionID:    cmd.TransactionID,
	})
	if err != nil {
		return err
	}

	err = s.store.createWithdrawalTimeline(tx, &withdrawal.WithdrawalTimeline{
		WithdrawalID:      cmd.ID,
		Message:           message,
		AdditionalContent: datatypes.JSON(detail),
		CreatedAt:         time.Now().UTC().Format(time.RFC3339Nano),
		CreatedBy:         cmd.UpdatedBy,
	})
	if err != nil {
		return err
	}

	return nil
}
