package memberimpl

import (
	"context"
	"money-transfer-demo/pkg/identity/member"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/util/password"
	"time"
)

func (s *service) LoginMember(ctx context.Context, cmd *member.LoginMemberCommand) (*member.LoginMemberResult, error) {
	membeDto, err := s.store.getMemberByLoginName(ctx, cmd.LoginName)
	if err != nil {
		return nil, err
	}

	if membeDto == nil {
		return nil, member.ErrMemberNotFound
	}

	if membeDto.Status != member.Active {
		return nil, member.ErrMemberInactive
	}

	invalidPassword := password.CheckPassword(cmd.Password, membeDto.Password)
	if !invalidPassword {
		return nil, member.ErrInvalidPassword
	}

	tokenStr, err := token.GenerateJWT(membeDto.ID, membeDto.LoginName, token.Member)
	if err != nil {
		return nil, err
	}

	return &member.LoginMemberResult{
		LoginName:   membeDto.LoginName,
		Currency:    string(membeDto.Currency),
		AccessToken: tokenStr,
		ExpiresIn:   time.Now().Add(2 * time.Hour).Unix(),
	}, nil

}
