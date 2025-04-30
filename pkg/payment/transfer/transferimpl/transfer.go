package transferimpl

import (
	"context"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/transaction"
	"money-transfer-demo/pkg/payment/transfer"
	"money-transfer-demo/pkg/util/generator"
	"time"

	"gorm.io/gorm"
)

type service struct {
	memberAccountSrv memberacc.Service
	store            *store
}

func NewService(memberAccountSrv memberacc.Service, store *store) transfer.Service {
	return &service{memberAccountSrv: memberAccountSrv, store: store}
}

func (s *service) Create(ctx context.Context, cmd *transfer.CreateTransferCommand) (*transfer.CreateTransferResult, error) {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return nil, transfer.ErrUnauthorized
	}

	cmd.FromMemberID = account.ID
	fromAcc, err := s.memberAccountSrv.GetByMemberID(ctx, cmd.FromMemberID)
	if err != nil {
		return nil, err
	}

	if fromAcc == nil {
		return nil, transfer.ErrMemberAccountNotFound
	}

	if fromAcc.Balance < cmd.Amount {
		return nil, transfer.ErrMemberAccountNotEnoughBalance
	}

	toAcc, err := s.memberAccountSrv.GetByMemberID(ctx, cmd.ToMemberID)
	if err != nil {
		return nil, err
	}

	if toAcc == nil {
		return nil, transfer.ErrMemberAccountNotFound
	}

	if fromAcc.Currency != toAcc.Currency {
		return nil, transfer.ErrInvalidCurrency
	}

	cmd.TransactionID = generator.GenerateTransactionID("TRAN")
	entity := &transfer.Transfer{
		FromMemberID:  cmd.FromMemberID,
		FromLoginName: fromAcc.LoginName,
		ToMemberID:    cmd.ToMemberID,
		ToLoginName:   toAcc.LoginName,
		Amount:        cmd.Amount,
		Status:        transfer.Pending,
		TransactionID: cmd.TransactionID,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	var id int64
	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		err = s.memberAccountSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
			MemberID:                  cmd.FromMemberID,
			AdjustedOutstandingAmount: cmd.Amount,
			TransactionType:           transaction.TransferType,
		})

		if err != nil {
			return err
		}

		newID, err := s.store.createTransfer(tx, entity)
		if err != nil {
			return err
		}

		id = newID

		err = s.store.createTransferTimeline(tx, &transfer.TransferTimeline{
			TransferID:    id,
			TransactionID: cmd.TransactionID,
			CreatedBy:     account.LoginName,
			CreatedAt:     time.Now().Format(time.RFC3339),
			Message:       transfer.Message(transfer.Pending, account.LoginName),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &transfer.CreateTransferResult{
		ID: id,
	}, nil
}

func (s *service) UpdateStatus(ctx context.Context, cmd *transfer.UpdateTransferStatusCommand) error {
	return s.store.db.Transaction(func(tx *gorm.DB) error {
		account, ok := ctx.Value("current_account").(*token.AccountData)
		if !ok || account == nil {
			return transfer.ErrUnauthorized
		}
		result, err := s.store.getTransferByID(ctx, cmd.ID)
		if err != nil {
			return err
		}

		if result == nil {
			return transfer.ErrTransferNotFound
		}

		if result.Status != transfer.Pending {
			return transfer.ErrInvalidStatus
		}

		err = s.store.updateTransferStatus(tx, &transfer.Transfer{
			ID:        cmd.ID,
			Status:    cmd.Status,
			UpdatedAt: time.Now().Format(time.RFC3339),
		})
		if err != nil {
			return err
		}

		err = s.store.createTransferTimeline(tx, &transfer.TransferTimeline{
			TransferID:    cmd.ID,
			TransactionID: result.TransactionID,
			CreatedBy:     account.LoginName,
			Note:          cmd.Note,
			CreatedAt:     time.Now().Format(time.RFC3339),
			Message:       transfer.Message(cmd.Status, account.LoginName),
		})
		if err != nil {
			return err
		}

		if cmd.Status == transfer.Successful {
			err = s.memberAccountSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
				MemberID:                  result.FromMemberID,
				AdjustedOutstandingAmount: -result.Amount,
				AdjustedAmount:            -result.Amount,
				TransactionType:           transaction.TransferType,
			})

			if err != nil {
				return err
			}

			err = s.memberAccountSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
				MemberID:        result.ToMemberID,
				AdjustedAmount:  result.Amount,
				TransactionType: transaction.TransferType,
			})

			if err != nil {
				return err
			}
		} else {
			err = s.memberAccountSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
				MemberID:                  result.FromMemberID,
				AdjustedOutstandingAmount: -result.Amount,
				TransactionType:           transaction.TransferType,
			})

			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *service) Search(ctx context.Context, query *transfer.SearchTransferQuery) (*transfer.SearchTransferResult, error) {
	if query.Page == 0 {
		query.Page = 1
	}

	if query.PerPage == 0 {
		query.PerPage = 10
	}

	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return nil, transfer.ErrUnauthorized
	}

	query.FromMemberID = account.ID
	result, err := s.store.search(ctx, query)
	if err != nil {
		return nil, err
	}

	result.Page = query.Page
	result.PerPage = query.PerPage
	return result, nil
}

func (s *service) Get(ctx context.Context, id int64) (*transfer.TransferDTO, error) {
	return s.store.getDetail(ctx, id)
}
