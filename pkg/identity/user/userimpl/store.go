package userimpl

import (
	"context"
	"errors"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/identity/user"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database postgres.DBConnector) *store {
	return &store{db: database.GetDB()}
}

func (s *store) createUser(tx *gorm.DB, user *user.User) error {
	return tx.Create(user).Error
}

func (s *store) getUserByLoginName(ctx context.Context, loginName string) (*user.User, error) {
	var user user.User

	if err := s.db.WithContext(ctx).Where("login_name = ?", loginName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
