package withdrawal

import "context"

type Service interface {
	CreateWithdrawal(ctx context.Context, cmd *CreateWithdrawalCommand) error
	SearchWithdrawal(ctx context.Context, query *SearchWithdrawalQuery) (*SearchWithdrawalResult, error)
	GetWithdrawalByID(ctx context.Context, id int64) (*WithdrawalDTO, error)
	ApproveWithdrawal(ctx context.Context, id int64) error
	RejectWithdrawal(ctx context.Context, id int64) error
	TransferWithdrawal(ctx context.Context, id int64) error
}
