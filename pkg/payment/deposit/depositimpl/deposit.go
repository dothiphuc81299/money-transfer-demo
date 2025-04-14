package depositimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/bankacc"
	"money-transfer-demo/pkg/payment/deposit"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/util/generator"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type service struct {
	store        *store
	memberAccSrv memberacc.Service
	bankAccSrv   bankacc.Service
}

func NewService(store *store, memberAccSrv memberacc.Service, bankAccSrv bankacc.Service) deposit.Service {
	return &service{store: store, memberAccSrv: memberAccSrv, bankAccSrv: bankAccSrv}
}

func (s *service) CreateDeposit(ctx context.Context, cmd *deposit.CreateDepositCommand) error {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return fmt.Errorf("user not authenticated")
	}

	if account.AccountType != token.Member {
		return fmt.Errorf("unauthorized: user is not a member")
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
			Amount:            cmd.Amount,
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

func (s *service) updateDetail(tx *gorm.DB, cmd *deposit.UpdateDepositStatusCommand) error {
	now := time.Now().UTC().Format(time.RFC3339)
	err := s.store.updateDeposit(tx, &deposit.Deposit{
		ID:        cmd.ID,
		Status:    cmd.Status,
		UpdatedBy: cmd.UpdatedBy,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	detailsJSON, err := json.Marshal(deposit.TimelineDetail{
		Remark: cmd.Note,
	})
	if err != nil {
		return err
	}

	timeline := &deposit.DepositTimeline{
		DepositID:         cmd.ID,
		Message:           deposit.Message(false, cmd.Status, cmd.UpdatedBy),
		AdditionalContent: detailsJSON,
		CreatedAt:         now,
		CreatedBy:         cmd.UpdatedBy,
	}

	if len(cmd.Note) != 0 {
		timeline.AdditionalContent = detailsJSON
	}

	err = s.store.createDepositTimeline(tx, timeline)
	if err != nil {
		return err
	}

	return nil
}
