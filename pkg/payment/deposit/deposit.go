package deposit

import "context"

type Service interface {
	CreateDepositLBT(ctx context.Context, cmd *CreateDepositCommand) error
	CreateDepositPaypal(ctx context.Context, cmd *CreateDepositPaypalCommand) (*CreateDepositPaypalResult, error)
	VerifyPaypal(ctx context.Context, cmd *VerifyPaypalCommand) error
	ApproveLBT(ctx context.Context, cmd *UpdateDepositStatusCommand) error
	RejectLBT(ctx context.Context, cmd *UpdateDepositStatusCommand) error
}
