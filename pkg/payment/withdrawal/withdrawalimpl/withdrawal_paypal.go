package withdrawalimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/memberpayacc"
	"money-transfer-demo/pkg/payment/transaction"
	"money-transfer-demo/pkg/payment/withdrawal"
	"strconv"
	"time"

	"github.com/plutov/paypal/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (s *service) Run(ctx context.Context) error {
	fmt.Println("Start running withdrawal paypal")
	defer fmt.Println("Stop running withdrawal paypal")

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.updateWithdrawalPaypal(ctx)
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *service) updateWithdrawalPaypal(ctx context.Context) {
	allWithdrawal, err := s.store.getALllTransferWithdrawalPaypal(ctx)
	if err != nil {
		fmt.Println("get all transfer withdrawal paypal err ", err)
		return
	}

	if len(allWithdrawal) == 0 {
		fmt.Println("no data")
		return
	}

	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		client, err := paypal.NewClient(
			s.cfg.Deposit.ClientIP,
			s.cfg.Deposit.Secret,
			paypal.APIBaseSandBox,
		)
		if err != nil {
			return err
		}

		accessToken, err := client.GetAccessToken(ctx)
		if err != nil {
			return err
		}

		client.SetAccessToken(accessToken.Token)
		for _, w := range allWithdrawal {
			batch, err := client.GetPayout(ctx, w.PayoutBatchID)
			if err != nil {
				log.Fatal("Failed to fetch payout status:", err)
				continue
			}

			if batch.BatchHeader.BatchStatus == "SUCCESS" {
				floatVal, err := strconv.ParseFloat(batch.BatchHeader.Amount.Value, 64)
				if err != nil {
					return err
				}

				floatFee, err := strconv.ParseFloat(batch.BatchHeader.Fees.Value, 64)
				if err != nil {
					return err
				}

				adjustedOutstandingAmount := withdrawal.SafeEstimatePaypalFee(w.GrossAmount) + w.GrossAmount

				var adjustAmount float64
				if w.Currency != string(w.MemberCurrency) {
					adjustedOutstandingAmount = adjustedOutstandingAmount * s.cfg.ExchangeVNDRate
					adjustAmount = (floatVal + floatFee) * s.cfg.ExchangeVNDRate
				} else {
					adjustedOutstandingAmount = adjustedOutstandingAmount
					adjustAmount = floatVal + floatFee
				}

				err = s.store.updatePaypal(tx, &withdrawal.Withdrawal{
					ID:           w.ID,
					Status:       withdrawal.Successful,
					GrossAmount:  floatVal + floatFee,
					NetAmount:    floatVal,
					ChargeAmount: floatFee,
					UpdatedAt:    time.Now().Format(time.RFC3339),
				})
				if err != nil {
					return err
				}

				message := withdrawal.Message(withdrawal.Successful, "system")
				detail, err := json.Marshal(&withdrawal.TimelineDetail{
					WithdrawalStatus: withdrawal.Successful,
					Note:             "auto ",
					TransactionID:    w.TransactionID,
				})
				if err != nil {
					return err
				}

				err = s.store.createWithdrawalTimeline(tx, &withdrawal.WithdrawalTimeline{
					WithdrawalID:      w.ID,
					Message:           message,
					AdditionalContent: datatypes.JSON(detail),
					CreatedAt:         time.Now().UTC().Format(time.RFC3339),
					CreatedBy:         "system",
				})
				if err != nil {
					return err
				}

				err = s.memberPaymentAccStore.UpdateVerifyStatus(tx, &memberpayacc.MemberPayAccount{
					ID:           w.MemberPaymentAccountID,
					VerifyStatus: memberpayacc.Verified,
					UpdatedBy:    "systẹm",
					UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
				})
				if err != nil {
					return err
				}

				err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
					MemberID:                  w.MemberID,
					UpdatedBy:                 "system",
					AdjustedAmount:            -adjustAmount,
					AdjustedOutstandingAmount: -adjustedOutstandingAmount,
					TransactionID:             w.TransactionID,
					TransactionType:           transaction.WithdrawalType,
				})
				if err != nil {
					return err
				}

				return nil
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("update withdrawal err ", err)
		return
	}
}

func (s *service) createSinglePaypal(ctx context.Context, tx *gorm.DB, entity *withdrawal.WithdrawalDTO, cmd *withdrawal.UpdateWithdrawalStatusCommand) error {
	client, err := paypal.NewClient(
		s.cfg.Deposit.ClientIP,
		s.cfg.Deposit.Secret,
		paypal.APIBaseSandBox,
	)
	if err != nil {
		return err
	}

	accessToken, err := client.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	client.SetAccessToken(accessToken.Token)

	amount := fmt.Sprintf("%.2f", entity.GrossAmount)

	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: "Subject will be displayed on PayPal",
		},
		Items: []paypal.PayoutItem{
			paypal.PayoutItem{
				RecipientType: "EMAIL",
				Receiver:      entity.PaypalEmail,
				Amount: &paypal.AmountPayout{
					Value:    amount,
					Currency: entity.Currency,
				},
				Note: cmd.Note,
			},
		},
	}

	payoutResp, err := client.CreatePayout(ctx, payout)
	if err != nil {
		return err
	}

	detail := &withdrawal.WithdrawalDetail{
		PaypalEmail:    entity.PaypalEmail,
		MemberCurrency: entity.MemberCurrency,
		PayoutBatchID:  payoutResp.BatchHeader.PayoutBatchID,
	}

	detailStr, err := json.Marshal(detail)
	if err != nil {
		return err
	}

	err = s.store.updateDetail(tx, &withdrawal.Withdrawal{
		ID:        entity.ID,
		Detail:    string(detailStr),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	return nil
}
