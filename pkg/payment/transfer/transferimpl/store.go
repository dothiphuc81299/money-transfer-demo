package transferimpl

import (
	"context"
	"errors"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/transfer"

	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(db postgres.DBConnector) *store {
	return &store{
		db: db.GetDB(),
	}
}

func (s *store) createTransfer(tx *gorm.DB, entity *transfer.Transfer) (int64, error) {
	err := tx.Model(&transfer.Transfer{}).Create(entity).Error
	if err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (s *store) updateTransferStatus(tx *gorm.DB, entity *transfer.Transfer) error {
	result := tx.Model(&transfer.Transfer{}).
		Where("id = ?", entity.ID).
		Updates(map[string]interface{}{
			"status":     entity.Status,
			"updated_at": entity.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *store) getTransferByID(ctx context.Context, id int64) (*transfer.Transfer, error) {
	var transfer transfer.Transfer
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&transfer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transfer, nil
}

func (s *store) createTransferTimeline(tx *gorm.DB, entity *transfer.TransferTimeline) error {
	err := tx.Model(&transfer.TransferTimeline{}).Create(entity).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *store) search(ctx context.Context, query *transfer.SearchTransferQuery) (*transfer.SearchTransferResult, error) {
	var (
		result = &transfer.SearchTransferResult{}
		total  int64
		db     = s.db.WithContext(ctx).Model(&transfer.Transfer{})
	)

	if query.TransactionID != "" {
		db = db.Where("transaction_id = ?", query.TransactionID)
	}

	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	if query.FromMemberID > 0 {
		db = db.Where("from_member_id = ?", query.FromMemberID)
	}

	if query.ToMemberID > 0 {
		db = db.Where("to_member_id = ?", query.ToMemberID)
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

		db = db.Offset(offset).Limit(query.PerPage)

		if err := db.Scan(&result.Data).Error; err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *store) getDetail(ctx context.Context, id int64) (*transfer.TransferDTO, error) {
	var result transfer.TransferDTO
	if err := s.db.WithContext(ctx).Table("transfer").Where("id = ?", id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := s.db.WithContext(ctx).Table("transfer_timeline").Where("transfer_id = ?", id).Find(&result.TransferTimeline).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
