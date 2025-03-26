package memberaccimpl

import (
	"context"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/memberacc"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database postgres.DBConnector) *store {
	return &store{db: database.GetDB()}
}

func (s *store) createMemberAccount(tx *gorm.DB, memberAccount *memberacc.MemberAccount) error {
	err := tx.Create(memberAccount).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *store) getMemberAccount(ctx context.Context, memberID int64, loginName string) ([]*memberacc.MemberAccount, error) {
	var memberAccount []*memberacc.MemberAccount

	query := s.db.WithContext(ctx).Where("member_id = ? or login_name = ?", memberID, loginName)
	if err := query.Find(&memberAccount).Error; err != nil {
		return nil, err
	}

	return memberAccount, nil
}
