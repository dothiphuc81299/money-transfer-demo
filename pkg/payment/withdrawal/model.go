package withdrawal

import (
	"encoding/json"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/datatypes"
)

var (
	ErrUnauthorized                = errors.New("unauthorized")
	ErrWithdrawalNotFound          = errors.New("withdrawal not found")
	ErrUnauthorizedUserIsNotMember = errors.New("unauthorized: user is not member")
)

type Status int

const (
	Pending Status = iota + 1
	Reviewing
	Transferring
	Successful
	Failed
	Expired
)

type PaymentMethodCode string

const (
	LBT    PaymentMethodCode = "LBT"
	MOMO   PaymentMethodCode = "MOMO"
	PAYPAL PaymentMethodCode = "PAYPAL"
)

type Withdrawal struct {
	ID                int64   `json:"id"`
	PaymentMethodCode string  `json:"payment_method_code"`
	TransactionID     string  `json:"transaction_id"`
	Status            Status  `json:"status"`
	MemberID          int64   `json:"member_id"`
	LoginName         string  `json:"login_name"`
	Currency          string  `json:"currency"`
	GrossAmount       float64 `json:"gross_amount"`
	ChargeAmount      float64 `json:"charge_amount"`
	NetAmount         float64 `json:"net_amount"`
	Detail            string  `json:"detail"`
	BankAccountID     int64   `json:"bank_account_id"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

func (Withdrawal) TableName() string {
	return "withdrawal"
}

type WithdrawalDTO struct {
	ID                int64   `json:"id"`
	PaymentMethodCode string  `json:"payment_method_code"`
	TransactionID     string  `json:"transaction_id"`
	Status            Status  `json:"status"`
	MemberID          int64   `json:"member_id"`
	LoginName         string  `json:"login_name"`
	Currency          string  `json:"currency"`
	GrossAmount       float64 `json:"gross_amount"`
	ChargeAmount      float64 `json:"charge_amount"`
	NetAmount         float64 `json:"net_amount"`
	MemberFullName    string  `json:"member_full_name"`
	MemberBankCode    string  `json:"member_bank_code"`
	MemberAccountNo   string  `json:"member_account_no"`
	MemberAccountName string  `json:"member_account_name"`
	BankAccountID     int64   `json:"bank_account_id"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type WithdrawalTimeline struct {
	ID                int64          `json:"id"`
	WithdrawalID      int64          `json:"withdrawal_id"`
	Message           string         `json:"message"`
	AdditionalContent datatypes.JSON `json:"additional_content"`
	CreatedBy         string         `json:"created_by"`
	CreatedAt         string         `json:"created_at"`
}

type TimelineDetail struct {
	WithdrawalStatus Status  `json:"withdrawal_status,omitempty"`
	Note             string  `json:"note,omitempty"`
	TransactionID    string  `json:"transaction_id,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	PaymentMethod    string  `json:"payment_method,omitempty"`
	Bank             string  `json:"bank,omitempty"`
}

func (WithdrawalTimeline) TableName() string {
	return "withdrawal_timeline"
}

type CreateWithdrawalCommand struct {
	PaymentMethodCode string `json:"payment_method_code"`
	TransactionID     string
	MemberID          int64
	LoginName         string
	Currency          string          `json:"currency"`
	Amount            float64         `json:"amount"`
	Detail            json.RawMessage `json:"detail"`
	BankAccountID     int64           `json:"bank_account_id"`
	DetailStr         string
	CreatedBy         string
	MemberBankCode    string
}

type WithdrawalDetail struct {
	MemberFullName    string `json:"member_full_name"`
	MemberBankCode    string `json:"member_bank_code"`
	MemberAccountNo   string `json:"member_account_no"`
	MemberAccountName string `json:"member_account_name"`
}

type SearchWithdrawalQuery struct {
	Status            Status `form:"status"`
	TransactionID     string `form:"transaction_id"`
	LoginName         string `form:"login_name"`
	PaymentMethodCode string `form:"payment_method_code"`
	Page              int    `form:"page"`
	PerPage           int    `form:"per_page"`
}

type SearchWithdrawalResult struct {
	Withdrawals []*WithdrawalDTO `json:"result"`
	Total       int64            `json:"total"`
	Page        int              `json:"page"`
	PerPage     int              `json:"per_page"`
}

type UpdateWithdrawalStatusCommand struct {
	ID              int64
	Status          Status
	BankAccountID   int64 `json:"bank_account_id"`
	BankAccountCode string
	Note            string  `json:"note"`
	ChargeAmount    float64 `json:"charge_amount"`
	UpdatedBy       string
	TransactionID   string
	Detail          string
}

type CreateWithdrawalTimelineCommand struct {
	WithdrawalID      int64          `json:"withdrawal_id"`
	Message           string         `json:"message"`
	AdditionalContent datatypes.JSON `json:"additional_content"`
	CreatedBy         string         `json:"created_by"`
	CreatedAt         string         `json:"created_at"`
}

func (cmd CreateWithdrawalCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.PaymentMethodCode, validation.Required, validation.In(LBT, MOMO, PAYPAL)),
		validation.Field(&cmd.Amount, validation.Required),
		validation.Field(&cmd.Detail, validation.Required),
	)
}

func (cmd WithdrawalDetail) ValidateWithdrawal(detail json.RawMessage) error {
	err := json.Unmarshal(detail, &cmd)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.MemberBankCode, validation.Required),
		validation.Field(&cmd.MemberAccountNo, validation.Required),
		validation.Field(&cmd.MemberAccountName, validation.Required),
	)
}

func (cmd UpdateWithdrawalStatusCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.ID, validation.Required),
		validation.Field(&cmd.Note, validation.Required),
	)
}

func Message(status Status, name string) string {
	var message string

	switch status {
	case Pending:
		message = "%s created this withdrawal"
	case Successful:
		message = "%s manually marked this withdrawal as successful"
	case Failed:
		message = "%s manually marked this withdrawal as failed"

	case Transferring:
		message = "%s manually marked this withdrawal as Transferring"

	case Expired:
		message = "%s manually retried the withdrawal expired"
	}

	return fmt.Sprintf(message, name)
}
