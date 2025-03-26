package member

import "context"

type Service interface {
	CreateMember(ctx context.Context, cmd *CreateMemberCommand) (*CreateMemberResult, error)
}
