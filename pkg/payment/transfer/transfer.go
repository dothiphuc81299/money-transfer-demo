package transfer

import "context"

type Service interface {
	Create(ctx context.Context, cmd *CreateTransferCommand) (*CreateTransferResult, error)
	UpdateStatus(ctx context.Context, cmd *UpdateTransferStatusCommand) error
	Search(ctx context.Context, query *SearchTransferQuery) (*SearchTransferResult, error)
	Get(ctx context.Context, id int64) (*TransferDTO, error)
}
