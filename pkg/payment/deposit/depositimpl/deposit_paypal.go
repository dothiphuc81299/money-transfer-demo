package depositimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"money-transfer-demo/pkg/identity/member"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/deposit"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/transaction"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/plutov/paypal/v4"
	"gorm.io/gorm"
)

func (s *service) CreateDepositPaypal(ctx context.Context, cmd *deposit.CreateDepositPaypalCommand) (*deposit.CreateDepositPaypalResult, error) {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return nil, deposit.ErrUnauthorized
	}

	if account.AccountType != token.Member {
		return nil, deposit.ErrUnauthorizedUserIsNotMember
	}

	cmd.MemberID = account.ID
	cmd.LoginName = account.LoginName

	_, err := s.memberAccSrv.GetByMemberID(ctx, cmd.MemberID)
	if err != nil {
		return nil, err
	}

	client, err := paypal.NewClient(
		s.cfg.Deposit.ClientIP,
		s.cfg.Deposit.Secret,
		paypal.APIBaseSandBox,
	)
	if err != nil {
		return nil, err
	}

	accessToken, err := client.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	client.SetAccessToken(accessToken.Token)
	amountStr := strconv.FormatFloat(cmd.Amount, 'f', 2, 64)
	order, err := client.CreateOrder(context.Background(),
		paypal.OrderIntentAuthorize,
		[]paypal.PurchaseUnitRequest{
			{
				ReferenceID: uuid.NewString(),
				Amount: &paypal.PurchaseUnitAmount{
					Currency: string(cmd.Currency),
					Value:    amountStr,
				},
			},
		},
		nil,
		&paypal.ApplicationContext{
			ReturnURL: s.cfg.Deposit.ReturnURL + "/paypal/return",
			CancelURL: s.cfg.Deposit.ReturnURL + "/paypal/cancel",
		},
	)

	if err != nil {
		return nil, err
	}

	var approveURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approveURL = link.Href
			break
		}
	}

	if approveURL == "" {
		return nil, deposit.ErrApproveURLNotFound
	}

	cmd.TransactionID = order.ID

	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		entity := deposit.Deposit{
			PaymentMethodCode: string(cmd.PaymentMethodCode),
			TransactionID:     cmd.TransactionID,
			Status:            deposit.Processing,
			MemberID:          cmd.MemberID,
			LoginName:         cmd.LoginName,
			Currency:          string(cmd.Currency),
			GrossAmount:       cmd.Amount,
			CreatedAt:         time.Now().UTC().Format(time.RFC3339),
			CreatedBy:         cmd.LoginName,
			UpdatedAt:         time.Now().UTC().Format(time.RFC3339),
		}

		id, err := s.store.createDeposit(tx, &entity)
		if err != nil {
			return err
		}

		timeline, err := json.Marshal(&deposit.TimelineDetail{
			TransactionID: cmd.TransactionID,
			Amount:        cmd.Amount,
			PaymentMethod: string(cmd.PaymentMethodCode),
		})
		if err != nil {
			return err
		}

		err = s.store.createDepositTimeline(tx, &deposit.DepositTimeline{
			DepositID:         id,
			Message:           deposit.Message(false, deposit.Processing, cmd.LoginName),
			AdditionalContent: timeline,
			CreatedAt:         time.Now().UTC().Format(time.RFC3339),
			CreatedBy:         cmd.LoginName,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &deposit.CreateDepositPaypalResult{
		ApproveURL: approveURL,
	}, nil
}

func (s *service) verifyPaypal(ctx context.Context, cmd *deposit.VerifyPaypalCommand) (*deposit.UpdateDepositStatusCommand, error) {
	var result deposit.UpdateDepositStatusCommand
	client, err := paypal.NewClient(
		s.cfg.Deposit.ClientIP,
		s.cfg.Deposit.Secret,
		paypal.APIBaseSandBox,
	)

	if err != nil {
		return nil, err
	}

	_, err = client.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	authorizationResp, err := client.AuthorizeOrder(context.Background(), cmd.Token, paypal.AuthorizeOrderRequest{})
	if err != nil {
		return nil, err
	}

	if len(authorizationResp.PurchaseUnits) == 0 ||
		authorizationResp.PurchaseUnits[0].Payments == nil ||
		len(authorizationResp.PurchaseUnits[0].Payments.Autthorizations) == 0 {
		return nil, deposit.ErrPaymentNotAunthorized
	}

	authorizationID := authorizationResp.PurchaseUnits[0].Payments.Autthorizations[0].ID
	captureResp, err := client.CaptureAuthorization(context.Background(), authorizationID, &paypal.PaymentCaptureRequest{})
	if err != nil {
		return nil, err
	}

	if captureResp.Status != "COMPLETED" {
		result.Status = deposit.Failed
	} else {
		result.Status = deposit.Successful
	}

	if result.Status == deposit.Successful {
		err = s.getDetailOrder(ctx, client, cmd.Token, &result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *service) getDetailOrder(ctx context.Context, client *paypal.Client, token string, cmd *deposit.UpdateDepositStatusCommand) error {
	order, err := client.GetOrder(ctx, token)
	if err != nil {
		return err
	}

	paymentCapture := order.PurchaseUnits[0].Payments.Captures[0]
	grossAmount := paymentCapture.SellerReceivableBreakdown.GrossAmount.Value

	chargeAmount := paymentCapture.SellerReceivableBreakdown.PaypalFee.Value
	netAmount := paymentCapture.SellerReceivableBreakdown.NetAmount.Value

	cmd.GrossAmount, _ = strconv.ParseFloat(grossAmount, 64)
	cmd.NetAmount, _ = strconv.ParseFloat(netAmount, 64)
	cmd.ChargeAmount, _ = strconv.ParseFloat(chargeAmount, 64)

	shippingLine := order.PurchaseUnits[0].Shipping.Address.AddressLine1
	shippingCity := order.PurchaseUnits[0].Shipping.Address.AdminArea2
	shippingState := order.PurchaseUnits[0].Shipping.Address.AdminArea1
	shippingCountry := order.PurchaseUnits[0].Shipping.Address.CountryCode
	shippingPostalCode := order.PurchaseUnits[0].Shipping.Address.PostalCode
	fullShippingAddress := fmt.Sprintf("%s, %s, %s, %s, %s", shippingLine, shippingCity, shippingState, shippingCountry, shippingPostalCode)

	paypalDetail, err := json.Marshal(&deposit.PayPalDetail{
		PaypalTransactionID: order.PurchaseUnits[0].Payments.Captures[0].ID,
		PayerID:             order.Payer.PayerID,
		PayerEmail:          order.Payer.EmailAddress,
		PayerName:           order.Payer.Name.GivenName + " " + order.Payer.Name.Surname,
		OrderID:             order.ID,
		CountryCode:         order.Payer.Address.CountryCode,
		ShippingName:        order.PurchaseUnits[0].Shipping.Name.FullName,
		ShippingAddress:     fullShippingAddress,
	})

	if err != nil {
		return err
	}

	cmd.DetailStr = string(paypalDetail)
	cmd.TransactionID = order.ID

	return nil
}

func (s *service) VerifyPaypal(ctx context.Context, cmd *deposit.VerifyPaypalCommand) error {
	cmdUpdate, err := s.verifyPaypal(ctx, cmd)
	if err != nil {
		return err
	}

	dp, err := s.store.getDepositByTransactionID(ctx, cmdUpdate.TransactionID)
	if err != nil {
		return err
	}

	if dp == nil {
		return deposit.ErrDepositNotFound
	}

	if dp.Status != deposit.Processing {
		return deposit.ErrDepositNotProcessing
	}

	cmdUpdate.ID = dp.ID
	cmdUpdate.UpdatedBy = deposit.DefaultUser

	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		mem, err := s.memberAccSrv.GetByMemberID(ctx, dp.MemberID)
		if err != nil {
			return err
		}

		cmdUpdate.Note = deposit.DefaultPaypalNote

		err = s.updateDetail(tx, cmdUpdate)
		if err != nil {
			return err
		}

		if string(mem.Currency) != dp.Currency && dp.Currency == string(member.UnitedStatesDollar) {
			cmdUpdate.NetAmount = cmdUpdate.NetAmount * s.cfg.ExchangeVNDRate
		}

		if cmdUpdate.Status == deposit.Successful {
			err = s.memberAccSrv.AdjustMemberAccountBalance(ctx, tx, &memberacc.AdjustMemberAccountBalanceCommand{
				MemberID:        dp.MemberID,
				UpdatedBy:       deposit.DefaultUser,
				AdjustedAmount:  cmdUpdate.NetAmount,
				TransactionID:   cmdUpdate.TransactionID,
				TransactionType: transaction.DepositType,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
