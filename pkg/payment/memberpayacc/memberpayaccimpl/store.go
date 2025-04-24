package memberpayaccimpl

import (
	"context"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/memberpayacc"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database postgres.DBConnector) *store {
	return &store{db: database.GetDB()}
}

func (s *store) create(tx *gorm.DB, entity *memberpayacc.MemberPayAccount) (int64, error) {
	err := tx.Model(&memberpayacc.MemberPayAccount{}).Create(entity).Error
	if err != nil {
		return 0, err
	}

	return entity.ID, nil
}

func (s *store) taken(ctx context.Context, memberID int64, paymentMethodCode string, emailPaypal string, bankCode string) ([]*memberpayacc.MemberPayAccount, error) {
	var entities []*memberpayacc.MemberPayAccount

	err := s.db.WithContext(ctx).Find(&entities, "member_id = ? AND payment_method_code = ? AND (detail ->> 'paypal_email' = ? OR detail ->> 'member_bank_code' = ?)", memberID, paymentMethodCode, emailPaypal, bankCode).Error
	return entities, err
}

func (s *store) Get(ctx context.Context, id int64) (*memberpayacc.MemberPaymentAccountDTO, error) {
	var entity memberpayacc.MemberPaymentAccountDTO

	err := s.db.WithContext(ctx).Table("member_payment_account").Select(
		"payment_method_code",
		"member_id",
		"detail ->> 'member_bank_code' as member_bank_code",
		"detail ->> 'member_account_no' as member_account_no",
		"detail ->> 'member_account_name' as member_account_name",
		"detail ->> 'member_full_name' as member_full_name",
		"detail ->> 'paypal_email' as paypal_email",
	).Where("id = ?", id).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity, err
}

func (s *store) search(ctx context.Context, query *memberpayacc.SearchMemberPayAccountQuery) (*memberpayacc.SearchMemberPayAccountResult, error) {

	var (
		result = &memberpayacc.SearchMemberPayAccountResult{}
		total  int64
		db     = s.db.WithContext(ctx).Model(&memberpayacc.MemberPayAccount{})
	)

	if query.MemberID > 0 {
		db = db.Where("member_id = ?", query.MemberID)
	}

	if len(query.PaymentMethodCode) > 0 {
		db = db.Where("payment_method_code = ?", query.PaymentMethodCode)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	result.Total = total
	if total > 0 {
		offset := 0
		if query.Page > 0 && query.PerPage > 0 {
			offset = query.PerPage * (query.Page - 1)
		}

		db = db.Select(
			"payment_method_code",
			"member_id",
			"id",
			"verify_status",
			"updated_at",
			"created_at",
			"created_by",
			"updated_by",
			"detail ->> 'member_bank_code' as member_bank_code",
			"detail ->> 'member_account_no' as member_account_no",
			"detail ->> 'member_account_name' as member_account_name",
			"detail ->> 'member_full_name' as member_full_name",
			"detail ->> 'paypal_email' as paypal_email",
		).Limit(query.PerPage).Offset(offset).
			Order("created_at DESC")

		err := db.Scan(&result.MemberPayAccounts).Error
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *store) UpdateVerifyStatus(tx *gorm.DB, entity *memberpayacc.MemberPayAccount) error {
	return tx.Model(&memberpayacc.MemberPayAccount{}).Where("id = ?", entity.ID).Updates(map[string]interface{}{
		"verify_status": entity.VerifyStatus,
		"updated_at":    entity.UpdatedAt,
		"updated_by":    entity.UpdatedBy,
	}).Error

}
