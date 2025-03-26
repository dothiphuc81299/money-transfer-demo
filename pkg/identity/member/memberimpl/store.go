package memberimpl

import (
	"context"
	"errors"
	"money-transfer-demo/pkg/identity/db"
	"money-transfer-demo/pkg/identity/member"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database db.DBConnector) *store {
	return &store{db: database.GetDB()}
}

func (s *store) GetMemberByEmail(email string) (*member.Member, error) {
	var member member.Member
	if err := s.db.Where("email = ?", email).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (s *store) isMemberTaken(ctx context.Context, id int64, loginName, email, phone string) ([]*member.Member, error) {
	var members []*member.Member

	query := s.db.WithContext(ctx).Where("login_name = ? OR email = ? OR phone = ? OR id = ?", loginName, email, phone, id)
	if err := query.Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (s *store) createMember(member *member.Member) (int64, error) {
	err := s.db.Create(member).Error
	if err != nil {
		return 0, err
	}

	return member.ID, nil
}
