package member

import "context"

type Service interface {
	CreateMember(ctx context.Context, cmd *CreateMemberCommand) (*CreateMemberResult, error)
	LoginMember(ctx context.Context, cmd *LoginMemberCommand) (*LoginMemberResult, error)
	GetMemberByID(ctx context.Context, id string) (*Member, error)
}
