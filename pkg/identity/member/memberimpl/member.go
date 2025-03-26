package memberimpl

import (
	"context"
	"money-transfer-demo/pkg/identity/member"
	"money-transfer-demo/pkg/util/password"
	"strings"
	"time"
)

const (
	TimestampFormat = time.RFC3339
)

type service struct {
	store *store
}

func NewService(store *store) member.Service {
	return &service{store: store}
}

func (s *service) CreateMember(ctx context.Context, cmd *member.CreateMemberCommand) (*member.CreateMemberResult, error) {
	result, err := s.store.isMemberTaken(ctx, 0, cmd.LoginName, cmd.Email, cmd.Phone)
	if err != nil {
		return nil, err
	}

	if len(result) != 0 {
		if result[0].Email == cmd.Email {
			return nil, member.ErrEmailExists
		}

		if result[0].Phone == cmd.Phone {
			return nil, member.ErrPhoneNumberExists
		}

		if result[0].LoginName == cmd.LoginName {
			return nil, member.ErrLoginNameExists
		}

		return nil, member.ErrMemberExists
	}

	cmd.Phone = strings.ReplaceAll(cmd.Phone, "+", "")
	cmd.Phone = string(cmd.PrefixPhone) + cmd.Phone

	now := time.Now().UTC().Format(TimestampFormat)
	entity := member.Member{
		LoginName:         cmd.LoginName,
		Status:            member.Active,
		Currency:          cmd.Currency,
		FullName:          cmd.FullName,
		Email:             cmd.Email,
		Phone:             cmd.Phone,
		EmailVerifyStatus: member.EmailUnverified,
		PhoneVerifyStatus: member.PhoneUnverified,
		Address:           cmd.Address,
		UpdatedAt:         now,
		CreatedAt:         now,
	}

	password, err := password.HashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	entity.Password = password
	id, err := s.store.createMember(&entity)
	if err != nil {
		return nil, err
	}

	return &member.CreateMemberResult{ID: id}, nil
}

func (s *service) GetMemberByID(ctx context.Context, id string) (*member.Member, error) {
	return s.store.getMemberByID(ctx, id)
}
