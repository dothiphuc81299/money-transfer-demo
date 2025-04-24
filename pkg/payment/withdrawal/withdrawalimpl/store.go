package withdrawalimpl

import (
	"context"
	"fmt"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/withdrawal"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(database postgres.DBConnector) *store {
	return &store{db: database.GetDB()}
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
		result = &withdrawal.SearchWithdrawalResult{}
		total  int64
		db     = s.db.WithContext(ctx).Model(&withdrawal.Withdrawal{})
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
			member_payment_account_id,
			login_name,
			payment_method_code,
			status,
			gross_amount,
			charge_amount,
			net_amount,
			member_id,
			currency,
			created_at,
			detail->>'member_currency' as member_currency,
			detail->>'member_full_name' as member_full_name,
			detail->>'member_bank_code' as member_bank_code,
			detail->>'member_account_no' as member_account_no,
			detail->>'member_account_name' as member_account_name,
			detail->>'paypal_email' as paypal_email,
			detail->>'payout_batch_id' as payout_batch_id,
			updated_at,
			bank_account_id
			`).Limit(query.PerPage).
			Offset(offset).
			Order("created_at DESC")

		err := db.Scan(&result.Withdrawals).Error
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *store) getWithdrawal(ctx context.Context, id int64) (*withdrawal.WithdrawalDTO, error) {
	var result withdrawal.WithdrawalDTO
	err := s.db.WithContext(ctx).Table("withdrawal").Select(`
			id,
			transaction_id,
			member_payment_account_id,
			login_name,
			payment_method_code,
			status,
			gross_amount,
			net_amount,
			charge_amount,
			member_id,
			currency,
			created_at,
			detail->>'meber_full_name' as member_full_name,
			detail->>'member_bank_code' as member_bank_code,
			detail->>'member_account_no' as member_account_no,
			detail->>'member_account_name' as member_account_name,
			updated_at,
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

/*************  ✨ Windsurf Command ⭐  *************/
/*******  93c68713-78ad-42b3-a054-fa72fabe66d2  *******/
func (s *store) updateWithdrawal(tx *gorm.DB, entity *withdrawal.Withdrawal) error {
	err := tx.Model(&withdrawal.Withdrawal{}).Where("id = ?", entity.ID).Updates(map[string]interface{}{
		"status":          entity.Status,
		"bank_account_id": entity.BankAccountID,
		"charge_amount":   entity.ChargeAmount,
		"net_amount":      entity.NetAmount,
		"detail":          entity.Detail,
		"updated_at":      entity.UpdatedAt,
	}).Error

	if err != nil {
		return fmt.Errorf("failed to update withdrawal: %w", err)
	}

	return nil
}

func (s *store) getALllTransferWithdrawalPaypal(ctx context.Context) ([]*withdrawal.WithdrawalDTO, error) {
	var result []*withdrawal.WithdrawalDTO

	err := s.db.WithContext(ctx).
		Table("withdrawal").
		Select(`
			id,
			member_payment_account_id,
			currency,
			member_id,
			gross_amount,
			transaction_id,
			detail ->> 'payout_batch_id' AS payout_batch_id,
			detail ->> 'member_currency' AS member_currency
		`).
		Where("status = ? AND payment_method_code = ?", withdrawal.Transferring, withdrawal.PAYPAL).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) updateDetail(tx *gorm.DB, entity *withdrawal.Withdrawal) error {
	err := tx.Model(&withdrawal.Withdrawal{}).Where("id = ?", entity.ID).Updates(map[string]interface{}{
		"detail":     entity.Detail,
		"updated_at": entity.UpdatedAt,
	}).Error

	if err != nil {
		return fmt.Errorf("failed to update withdrawal: %w", err)
	}

	return nil
}

func (s *store) updateStatus(tx *gorm.DB, entity *withdrawal.Withdrawal) error {
	err := tx.Model(&withdrawal.Withdrawal{}).Where("id = ?", entity.ID).Updates(map[string]interface{}{
		"status":     entity.Status,
		"updated_at": entity.UpdatedAt,
	}).Error

	if err != nil {
		return fmt.Errorf("failed to update withdrawal: %w", err)
	}
	return nil
}
