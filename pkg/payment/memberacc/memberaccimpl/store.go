package memberaccimpl

import (
	"context"
	"errors"
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

func (s *store) getMemberAccountByMemberID(ctx context.Context, memberID int64) (*memberacc.MemberAccount, error) {
	var memberAccount *memberacc.MemberAccount

	query := s.db.WithContext(ctx).Where("member_id = ?", memberID)
	if err := query.First(&memberAccount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return memberAccount, nil
}

func (s *store) updateMemberAccountBalance(tx *gorm.DB, cmd *memberacc.UpdateMemberAccountBalanceCommand) error {
	result := tx.Model(&memberacc.MemberAccount{}).
		Where("id = ?", cmd.ID).
		Updates(map[string]interface{}{
			"balance":             gorm.Expr("balance + ?", cmd.AdjustedAmount),
			"outstanding_balance": gorm.Expr("outstanding_balance + ?", cmd.AdjustedOutstandingAmount),
			"updated_at":          cmd.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	return nil

}
