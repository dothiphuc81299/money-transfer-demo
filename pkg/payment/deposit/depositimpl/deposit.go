package depositimpl

import (
	"encoding/json"
	"money-transfer-demo/pkg/payment/bankacc"
	"money-transfer-demo/pkg/payment/config"
	"money-transfer-demo/pkg/payment/deposit"
	"money-transfer-demo/pkg/payment/memberacc"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type service struct {
	store        *store
	memberAccSrv memberacc.Service
	bankAccSrv   bankacc.Service
	cfg          *config.Config
}

func NewService(store *store, memberAccSrv memberacc.Service, bankAccSrv bankacc.Service, cfg *config.Config) deposit.Service {
	return &service{store: store, memberAccSrv: memberAccSrv, bankAccSrv: bankAccSrv, cfg: cfg}
}

func (s *service) updateDetail(tx *gorm.DB, cmd *deposit.UpdateDepositStatusCommand) error {
	now := time.Now().UTC().Format(time.RFC3339)
	err := s.store.updateDeposit(tx, &deposit.Deposit{
		ID:           cmd.ID,
		Status:       cmd.Status,
		UpdatedBy:    cmd.UpdatedBy,
		GrossAmount:  cmd.GrossAmount,
		NetAmount:    cmd.NetAmount,
		ChargeAmount: cmd.ChargeAmount,
		Detail:       datatypes.JSON(cmd.DetailStr),
		UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	detailsJSON, err := json.Marshal(deposit.TimelineDetail{
		Remark: cmd.Note,
	})
	if err != nil {
		return err
	}

	timeline := &deposit.DepositTimeline{
		DepositID:         cmd.ID,
		Message:           deposit.Message(false, cmd.Status, cmd.UpdatedBy),
		AdditionalContent: detailsJSON,
		CreatedAt:         now,
		CreatedBy:         cmd.UpdatedBy,
	}

	if len(cmd.Note) != 0 {
		timeline.AdditionalContent = detailsJSON
	}

	err = s.store.createDepositTimeline(tx, timeline)
	if err != nil {
		return err
	}

	return nil
}
