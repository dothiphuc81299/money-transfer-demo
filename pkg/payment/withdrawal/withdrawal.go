package withdrawal

import "context"

type Service interface {
	CreateWithdrawal(ctx context.Context, cmd *CreateWithdrawalCommand) error
	SearchWithdrawal(ctx context.Context, query *SearchWithdrawalQuery) (*SearchWithdrawalResult, error)
	GetWithdrawalByID(ctx context.Context, id int64) (*WithdrawalDTO, error)
	ApproveWithdrawal(ctx context.Context, cmd *UpdateWithdrawalStatusCommand) error
	RejectWithdrawal(ctx context.Context, cmd *UpdateWithdrawalStatusCommand) error
	TransferWithdrawal(ctx context.Context, cmd *UpdateWithdrawalStatusCommand) error
	ReviewWithdrawal(ctx context.Context, cmd *UpdateWithdrawalStatusCommand) error
	Run(ctx context.Context) error
}
