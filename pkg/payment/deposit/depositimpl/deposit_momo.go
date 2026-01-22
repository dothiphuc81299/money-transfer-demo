package depositimpl

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"money-transfer-demo/pkg/identity/token"
	"net/http"
	"time"

	"money-transfer-demo/pkg/payment/deposit"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

func (s *service) CreateDepositMoMo(ctx context.Context, cmd *deposit.CreateDepositMoMoCommand) (*deposit.CreateDepositMoMoResult, error) {
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

	accessToken, err := s.getAccessToken()
	if err != nil {
		return nil, err
	}

	referenceID := uuid.NewString()

	err = s.requestToPay(accessToken, RequestToPayPayload{
		Amount:       fmt.Sprintf("%.0f", cmd.Amount),
		Currency:     "EUR",
		ExternalID:   referenceID,
		Payer:        PayerInfo{PartyIDType: "MSISDN", PartyID: cmd.PhoneNumber},
		PayerMessage: "Thanh toan nap tien",
		PayeeNote:    "Nap tien vao vi",
		ReferenceID:  referenceID,
	})
	if err != nil {
		return nil, err
	}

	cmd.TransactionID = referenceID

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

	// return status pending to client (momo sandbox don't perform real transaction)
	return &deposit.CreateDepositMoMoResult{
		Status: "PENDING",
	}, nil
}

func (s *service) getAccessToken() (string, error) {
	url := fmt.Sprintf("%s/collection/token/", s.cfg.Deposit.MoMo.BaseURL)

	authStr := fmt.Sprintf("%s:%s", s.cfg.Deposit.MoMo.APIUser, s.cfg.Deposit.MoMo.APIKey)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(authStr))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(nil))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+authEncoded)
	req.Header.Set("Ocp-Apim-Subscription-Key", s.cfg.Deposit.MoMo.SubscriptionKey)
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get access token: %s", string(bodyBytes))
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (s *service) requestToPay(accessToken string, payload RequestToPayPayload) error {
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", s.cfg.Deposit.MoMo.BaseURL)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Reference-Id", payload.ReferenceID)
	req.Header.Set("X-Target-Environment", "sandbox")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", s.cfg.Deposit.MoMo.SubscriptionKey)
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("MoMo request failed with status %d", resp.StatusCode)
	}

	return nil
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

type RequestToPayPayload struct {
	Amount       string    `json:"amount"`
	Currency     string    `json:"currency"`
	ExternalID   string    `json:"externalId"`
	Payer        PayerInfo `json:"payer"`
	PayerMessage string    `json:"payerMessage"`
	PayeeNote    string    `json:"payeeNote"`
	ReferenceID  string    `json:"referenceId"`
}

type PayerInfo struct {
	PartyIDType string `json:"partyIdType"` // always "MSISDN"
	PartyID     string `json:"partyId"`     // phone number, e.g., "46733123453"
}
