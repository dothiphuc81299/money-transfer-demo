package withdrawalimpl

import (
	"context"
	"fmt"
	"log"

	"github.com/plutov/paypal/v4"
)

func (s *service) CreateSinglePayout(ctx context.Context) error {
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

	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: "Subject will be displayed on PayPal",
		},
		Items: []paypal.PayoutItem{
			paypal.PayoutItem{
				RecipientType: "EMAIL",
				Receiver:      "sb-af7o140179909@personal.example.com",
				Amount: &paypal.AmountPayout{
					Value:    "15.11",
					Currency: "USD",
				},
				Note:         "Optional note",
				SenderItemID: "Optional Item ID",
			},
		},
	}

	payoutResp, err := client.CreatePayout(ctx, payout)

	batch, err := client.GetPayout(ctx, payoutResp.BatchHeader.PayoutBatchID)
	if err != nil {
		log.Fatal("Failed to fetch payout status:", err)
	}
	fmt.Println("Payout Status:", batch.BatchHeader.BatchStatus)
	return nil
}
