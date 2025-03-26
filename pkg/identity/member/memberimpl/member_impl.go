package memberimpl

import (
	"context"
	"errors"
	"money-transfer-demo/pkg/identity/member"
)

type service struct {
	store *store
}

func NewService(store *store) member.Service {
	return &service{store: store}
}

func (s *service) CreateMember(ctx context.Context, cmd *member.CreateMemberCommand) (*member.Member, error) {
	existingMember, err := s.store.GetMemberByEmail(cmd.Email)
	if err != nil {
		return nil, err
	}
	if existingMember != nil {
		return nil, errors.New("email already Createed")
	}

	entity := &member.Member{
		Name:     cmd.Name,
		Email:    cmd.Email,
		Password: cmd.Password,
	}

	if err := s.store.CreateMember(entity); err != nil {
		return nil, err
	}

	return entity, nil
}
