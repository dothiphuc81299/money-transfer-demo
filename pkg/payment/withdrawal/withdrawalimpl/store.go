package withdrawalimpl

import (
	"context"
	"money-transfer-demo/pkg/payment/withdrawal"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *store {
	return &store{db: db}
}

func (s *store) createWithdrawal(tx *gorm.DB, entity *withdrawal.Withdrawal) (int64, error) {
	err := tx.Model(&withdrawal.Withdrawal{}).Create(entity).Error
	if err != nil {
		return 0, err
	}

	return entity.ID, nil
}

func (s *store) createWithdrawalTimeline(tx *gorm.DB, timeline *withdrawal.WithdrawalTimeline) error {
	return tx.Model(&withdrawal.WithdrawalTimeline{}).Create(timeline).Error
}

func (s *store) searchWithdrawal(ctx context.Context, query *withdrawal.SearchWithdrawalQuery) (*withdrawal.SearchWithdrawalResult, error) {
	var (
		result  = &withdrawal.SearchWithdrawalResult{}
		records []*withdrawal.WithdrawalDTO
		total   int64
		db      = s.db.WithContext(ctx).Model(&withdrawal.Withdrawal{})
	)

	if len(query.TransactionID) > 0 {
		db = db.Where("transaction_id = ?", query.TransactionID)
	}

	if len(query.LoginName) > 0 {
		db = db.Where("login_name = ?", query.LoginName)
	}

	if len(query.PaymentMethodCode) > 0 {
		db = db.Where("payment_method_code = ?", query.PaymentMethodCode)
	}

	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
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

		db = db.Select(`
			id,
			transaction_id,
			login_name,
			payment_method_code,
			status,
			amount,
			member_id,
			currency,
			created_at
			detail->>'meber_full_name' as member_full_name,
			detail->>'member_bank_code' as member_bank_code,
			detail->>'member_account_no' as member_account_no,
			detail->>'member_account_name' as member_account_name,
			updated_at.
			bank_account_id
			`).Limit(query.PerPage).
			Offset(offset).
			Order("w.created_at DESC")

		err := db.Scan(&result).Error
		if err != nil {
			return nil, err
		}
	}

	result.Withdrawals = records
	return result, nil
}

func (s *store) getWithdrawal(ctx context.Context, id int64) (*withdrawal.WithdrawalDTO, error) {
	var result withdrawal.WithdrawalDTO
	err := s.db.WithContext(ctx).Table("withdrawal").Select(`
			id,
			transaction_id,
			login_name,
			payment_method_code,
			status,
			amount,
			member_id,
			currency,
			created_at
			detail->>'meber_full_name' as member_full_name,
			detail->>'member_bank_code' as member_bank_code,
			detail->>'member_account_no' as member_account_no,
			detail->>'member_account_name' as member_account_name,
			updated_at.
			bank_account_id
			`).Where("id =?", id).Scan(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	return &result, nil
}
