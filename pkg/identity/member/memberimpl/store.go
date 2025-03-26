package memberimpl

import (
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

func (s *store) CreateMember(member *member.Member) error {
	return s.db.Create(member).Error
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
