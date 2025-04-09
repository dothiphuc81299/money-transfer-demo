package userimpl

import (
	"context"
	"money-transfer-demo/pkg/identity/member"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/identity/user"
	"money-transfer-demo/pkg/util/password"
	"time"
)

func (s *service) LoginUser(ctx context.Context, cmd *user.LoginUserCommand) (*user.LoginUserResult, error) {
	userDTO, err := s.store.getUserByLoginName(ctx, cmd.LoginName)
	if err != nil {
		return nil, err
	}

	if userDTO == nil {
		return nil, user.ErrUserNotFound
	}

	if userDTO.Status != user.Active {
		return nil, user.ErrUserInactive
	}

	invalidPassword := password.CheckPassword(cmd.Password, userDTO.Password)
	if !invalidPassword {
		return nil, member.ErrInvalidPassword
	}

	tokenStr, err := token.GenerateJWT(userDTO.ID, userDTO.LoginName, token.User)
	if err != nil {
		return nil, err
	}

	return &user.LoginUserResult{
		LoginName:   userDTO.LoginName,
		AccessToken: tokenStr,
		ExpiresIn:   time.Now().Add(2 * time.Hour).Unix(),
	}, nil

}
