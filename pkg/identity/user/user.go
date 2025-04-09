package user

import "context"

type Service interface {
	CreateUser(ctx context.Context, cmd *CreateUserCommand) error
	LoginUser(ctx context.Context, cmd *LoginUserCommand) (*LoginUserResult, error)
}
