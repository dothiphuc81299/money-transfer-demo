package userimpl

import (
	"context"
	"time"

	"money-transfer-demo/pkg/identity/user"
	"money-transfer-demo/pkg/util/password"

	"gorm.io/gorm"
)

type service struct {
	store *store
}

func NewService(store *store) user.Service {
	return &service{store: store}
}

func (s *service) CreateUser(ctx context.Context, cmd *user.CreateUserCommand) error {
	return s.store.db.Transaction(func(tx *gorm.DB) error {
		exist, err := s.store.getUserByLoginName(ctx, cmd.LoginName)
		if err != nil {
			return err
		}

		if exist != nil {
			return user.ErrUserLoginnameAlreadyExists
		}

		entity := &user.User{
			LoginName: cmd.LoginName,
			Status:    user.Active,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		password, err := password.HashPassword(cmd.Password)
		if err != nil {
			return err
		}

		entity.Password = password

		return s.store.createUser(tx, entity)
	})
}
