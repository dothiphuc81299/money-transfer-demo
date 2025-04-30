package depositimpl

import (
	"context"
	"encoding/json"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/bankacc"
	"money-transfer-demo/pkg/payment/deposit"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/transaction"
	"money-transfer-demo/pkg/util/generator"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (s *service) CreateDepositLBT(ctx context.Context, cmd *deposit.CreateDepositCommand) error {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return deposit.ErrUnauthorized
	}

	if account.AccountType != token.Member {
		return deposit.ErrUnauthorizedUserIsNotMember
	}

	cmd.MemberID = account.ID

	memberAcc, err := s.memberAccSrv.GetByMemberID(ctx, cmd.MemberID)
	if err != nil {
		return err
	}

	cmd.LoginName = memberAcc.LoginName
	cmd.TransactionID = generator.GenerateTransactionID(string(cmd.PaymentMethodCode))

	ba, err := s.bankAccSrv.GetBankAccountByID(ctx, cmd.BankAccountID)
	if err != nil {
		return err
	}

	if ba == nil {
		return deposit.ErrBankAccountNotFound
	}

	err = s.getLBTDetails(cmd)
	if err != nil {
		return err
	}

	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		entity := deposit.Deposit{
			RefCode:           cmd.RefCode,
			PaymentMethodCode: string(cmd.PaymentMethodCode),
			TransactionID:     cmd.TransactionID,
			Status:            deposit.Processing,
			MemberID:          cmd.MemberID,
			LoginName:         cmd.LoginName,
			Currency:          string(memberAcc.Currency),
			GrossAmount:       cmd.Amount,
			NetAmount:         cmd.Amount,
			BankAccountID:     cmd.BankAccountID,
			CreatedAt:         time.Now().UTC().Format(time.RFC3339),
			CreatedBy:         cmd.LoginName,
			UpdatedAt:         time.Now().UTC().Format(time.RFC3339),
		}

		entity.Detail = datatypes.JSON(cmd.DetailStr)
		id, err := s.store.createDeposit(tx, &entity)
		if err != nil {
			return err
		}

		timeline, err := json.Marshal(&deposit.TimelineDetail{
			TransactionID: cmd.TransactionID,
			Amount:        cmd.Amount,
			PaymentMethod: string(cmd.PaymentMethodCode),
			BankAccount:   ba.BankCode + "-" + ba.AccountNo,
			RefCode:       cmd.RefCode,
		})
		if err != nil {
			return err
		}

		err = s.store.createDepositTimeline(tx, &deposit.DepositTimeline{
			DepositID:         id,
			Message:           deposit.Message(false, deposit.Processing, cmd.LoginName),
			AdditionalContent: timeline,
			CreatedAt:         time.Now().UTC().Format(time.RFC3339),
			CreatedBy:         cmd.LoginName,
		})
		if err != nil {
			return err
		}

		return nil

	})

	return nil
}

func (s *service) getLBTDetails(cmd *deposit.CreateDepositCommand) error {
	details := &deposit.LBTDetail{}

	err := details.ValidateLBT(cmd.Detail)
	if err != nil {
		return err
	}

	dt, err := json.Marshal(details)
	if err != nil {
		return err
	}

	cmd.DetailStr = string(dt)
	return nil
}

func (s *service) ApproveLBT(ctx context.Context, cmd *deposit.UpdateDepositStatusCommand) error {
	return s.store.db.Transaction(func(tx *gorm.DB) error {
		account, ok := ctx.Value("current_account").(*token.AccountData)
		if !ok || account == nil {
			return deposit.ErrUnauthorized
		}

		cmd.UpdatedBy = account.LoginName

		dp, err := s.store.getDepositByStatus(tx, cmd.ID, deposit.Processing)
		if err != nil {
			return err
		}

		if dp == nil {
			return deposit.ErrDepositNotFound
		}

		cmd.TransactionID = dp.TransactionID
		cmd.GrossAmount = dp.GrossAmount
		cmd.NetAmount = dp.NetAmount
		cmd.MemberID = dp.MemberID

		err = s.updateDetail(tx, cmd)
		if err != nil {
			return err
		}

		err = s.bankAccSrv.AdjustBalance(ctx, tx, &bankacc.AdjustBankAccountBalanceCommand{
			BankAccountID: dp.BankAccountID,
			ChangedAmount: cmd.GrossAmount,
			Currency:      dp.Currency,
			UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
			MemberID:        cmd.MemberID,
			UpdatedBy:       cmd.UpdatedBy,
			AdjustedAmount:  cmd.GrossAmount,
			TransactionType: transaction.DepositType,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) RejectLBT(ctx context.Context, cmd *deposit.UpdateDepositStatusCommand) error {
	return s.store.db.Transaction(func(tx *gorm.DB) error {
		account, ok := ctx.Value("current_account").(*token.AccountData)
		if !ok || account == nil {
			return deposit.ErrUnauthorized
		}

		cmd.UpdatedBy = account.LoginName

		dp, err := s.store.getDepositByStatus(tx, cmd.ID, deposit.Processing)
		if err != nil {
			return err
		}

		if dp == nil {
			return deposit.ErrDepositNotFound
		}

		cmd.GrossAmount = dp.GrossAmount
		cmd.NetAmount = dp.NetAmount
		cmd.MemberID = dp.MemberID
		err = s.updateDetail(tx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}
