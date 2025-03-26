package memberacc

import "context"

type Service interface {
	Create(ctx context.Context, cmd *CreateMemberAccountCommand) error
}
