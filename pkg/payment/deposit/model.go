package deposit

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/datatypes"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrUnauthorized              = errors.New("unauthorized")
	ErrInvalidMemberBankCode     = errors.New("invalid member bank code")
	ErrInvalidMemberAccountNo    = errors.New("invalid member account no")
	ErrInvalidMemberAccountName  = errors.New("invalid member account name")
	ErrInvalidCompanyBankCode    = errors.New("invalid company bank code")
	ErrInvalidCompanyAccountNo   = errors.New("invalid company account no")
	ErrInvalidCompanyAccountName = errors.New("invalid company account name")
	ErrDepositNotFound           = errors.New("deposit not found")
	ErrBankAccountNotFound       = errors.New("bank account not found")
)

type Status int

const (
	Processing Status = iota + 1
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

type Deposit struct {
	ID                int64          `json:"id"`
	TransactionID     string         `json:"transaction_id"`
	MemberID          int64          `json:"member_id"`
	LoginName         string         `json:"login_name"`
	PaymentMethodCode string         `json:"payment_method_code"`
	RefCode           string         `json:"ref_code"`
	Currency          string         `json:"currency"`
	BankAccountID     int64          `json:"bank_account_id"`
	Detail            datatypes.JSON `json:"detail"`
	Amount            float64        `json:"amount"`
	Status            Status         `json:"status"`
	CreatedBy         string         `json:"created_by"`
	UpdatedBy         string         `json:"updated_by"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
}

func (Deposit) TableName() string {
	return "deposit"
}

type CreateDepositCommand struct {
	TransactionID     string            `json:"transaction_id"`
	MemberID          int64             `json:"member_id"`
	LoginName         string            `json:"login_name"`
	PaymentMethodCode PaymentMethodCode `json:"payment_method_code"`
	RefCode           string            `json:"ref_code"`
	Detail            json.RawMessage   `json:"detail"`
	BankAccountID     int64             `json:"bank_account_id"`
	Amount            float64           `json:"amount"`
	DetailStr         string
}

type LBTDetail struct {
	MemberBankCode    string `json:"member_bank_code"`
	MemberAccountNo   string `json:"member_account_no"`
	MemberAccountName string `json:"member_account_name"`
	MemberBankRef     string `json:"member_bank_ref"`
	MemberFullName    string `json:"member_full_name,omitempty"`
}

type UpdateDepositStatusCommand struct {
	ID            int64
	MemberID      int64
	Status        Status `json:"status"`
	Note          string `json:"note"`
	Amount        float64
	TransactionID string
	UpdatedBy     string
}

type DepositTimeline struct {
	ID                int64          `json:"id"`
	DepositID         int64          `json:"deposit_id"`
	Message           string         `json:"message"`
	AdditionalContent datatypes.JSON `json:"additional_content"`
	CreatedBy         string         `json:"created_by"`
	CreatedAt         string         `json:"created_at"`
}

type TimelineDetail struct {
	DepositStatus Status  `json:"deposit_status,omitempty"`
	RefCode       string  `json:"ref_code,omitempty"`
	Remark        string  `json:"remark,omitempty"`
	TransactionID string  `json:"transaction_id,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
	PaymentMethod string  `json:"payment_method,omitempty"`
	BankAccount   string  `json:"bank_account,omitempty"`
}

func (DepositTimeline) TableName() string {
	return "deposit_timeline"
}

func (cmd CreateDepositCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.BankAccountID, validation.Required),
		validation.Field(&cmd.PaymentMethodCode, validation.Required, validation.In(LBT, MOMO, PAYPAL)),
		validation.Field(&cmd.Amount, validation.Required),
		validation.Field(&cmd.Detail, validation.Required),
		validation.Field(&cmd.RefCode, validation.Required),
	)
}

func (lbt *LBTDetail) ValidateLBT(detail json.RawMessage) error {
	err := json.Unmarshal(detail, &lbt)
	if err != nil {
		return err
	}

	if len(lbt.MemberBankCode) <= 0 {
		return ErrInvalidMemberBankCode
	}

	if len(lbt.MemberAccountNo) <= 0 {
		return ErrInvalidMemberAccountNo
	}

	if len(lbt.MemberAccountName) <= 0 {
		return ErrInvalidMemberAccountName
	}

	return nil
}

func (cmd UpdateDepositStatusCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.ID, validation.Required),
		validation.Field(&cmd.Status, validation.Required, validation.In(Successful, Failed)),
	)
}

func Message(isAuto bool, status Status, name string) string {
	var message string

	switch status {
	case Processing:
		message = "%s created this deposit"
	case Successful:
		if isAuto {
			message = "%s marked this deposit as successful"
		} else {
			message = "%s manually marked this deposit as successful"
		}
	case Failed:
		message = "%s manually marked this deposit as failed"
	case Expired:
		message = "%s marked this deposit as expired"
	}

	return fmt.Sprintf(message, name)
}
