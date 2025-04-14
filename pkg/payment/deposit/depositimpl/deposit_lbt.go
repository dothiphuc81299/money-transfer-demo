package depositimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/bankacc"
	"money-transfer-demo/pkg/payment/deposit"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/transaction"
	"time"

	"gorm.io/gorm"
)

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
			return fmt.Errorf("user not authenticated")
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
		cmd.Amount = dp.Amount
		cmd.MemberID = dp.MemberID

		err = s.updateDetail(tx, cmd)
		if err != nil {
			return err
		}

		err = s.bankAccSrv.AdjustBalance(ctx, tx, &bankacc.AdjustBankAccountBalanceCommand{
			BankAccountID: dp.BankAccountID,
			ChangedAmount: cmd.Amount,
			UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
			MemberID:        cmd.MemberID,
			UpdatedBy:       cmd.UpdatedBy,
			AdjustedAmount:  cmd.Amount,
			TransactionID:   cmd.TransactionID,
			TransactionType: transaction.DepositType,
			Note:            cmd.Note,
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
			return fmt.Errorf("user not authenticated")
		}

		cmd.UpdatedBy = account.LoginName

		dp, err := s.store.getDepositByStatus(tx, cmd.ID, deposit.Processing)
		if err != nil {
			return err
		}

		if dp == nil {
			return deposit.ErrDepositNotFound
		}

		err = s.updateDetail(tx, cmd)
		if err != nil {
			return err
		}

		return nil
	})
}
