package deposit

import "context"

type Service interface {
	CreateDeposit(ctx context.Context, cmd *CreateDepositCommand) error
	ApproveLBT(ctx context.Context, cmd *UpdateDepositStatusCommand) error
	RejectLBT(ctx context.Context, cmd *UpdateDepositStatusCommand) error
}
