package memberpayacc

import (
	"encoding/json"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/datatypes"
)

var (
	ErrUnauthorized                      = errors.New("unauthorized")
	ErrMemberPaymentAccountAlreadyExists = errors.New("member payment account already exists")
	ErrMemberNotFound                    = errors.New("member not found")
)

type VerifyStatus int

const (
	Unverified VerifyStatus = iota + 1
	Verified
)

type PaymentMethodCode string

const (
	LBT    PaymentMethodCode = "LBT"
	MOMO   PaymentMethodCode = "MOMO"
	PAYPAL PaymentMethodCode = "PAYPAL"
)

type MemberPayAccount struct {
	ID                int64          `json:"id"`
	MemberID          int64          `json:"member_id"`
	PaymentMethodCode string         `json:"payment_method_code"`
	Detail            datatypes.JSON `json:"detail"`
	VerifyStatus      VerifyStatus   `json:"verify_status"`
	CreatedBy         string         `json:"created_by"`
	UpdatedBy         string         `json:"updated_by"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
}

func (*MemberPayAccount) TableName() string {
	return "member_payment_account"
}

type MemberPaymentAccountDTO struct {
	ID                int64  `json:"id"`
	MemberID          int64  `json:"member_id"`
	PaymentMethodCode string `json:"payment_method_code"`
	VerifyStatus      string `json:"verify_status"`
	CreatedBy         string `json:"created_by"`
	UpdatedBy         string `json:"updated_by,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`

	// for local bank
	MemberBankCode    string `json:"member_bank_code,omitempty"`
	MemberAccountNo   string `json:"member_account_no,omitempty"`
	MemberAccountName string `json:"member_account_name,omitempty"`
	MemberFullName    string `json:"member_full_name,omitempty"`

	// for paypal
	PaypalEmail string `json:"paypal_email,omitempty"`
}

type CreateMemberPayAccountCommand struct {
	MemberID          int64
	PaymentMethodCode string          `json:"payment_method_code"`
	Detail            json.RawMessage `json:"detail"`
	CreatedBy         string
	DetailStr         string
	PaypalEmail       string
	MemberBankCode    string
	MemberFullName    string
}

type CreateMemberPaymentAccountResult struct {
	ID int64 `json:"id"`
}

type LBTDetail struct {
	MemberBankCode    string `json:"member_bank_code"`
	MemberAccountNo   string `json:"member_account_no"`
	MemberAccountName string `json:"member_account_name"`
	MemberFullName    string `json:"member_full_name,omitempty"`
}

type PayPalDetail struct {
	PayPalEmail string `json:"paypal_email"`
}

type SearchMemberPayAccountQuery struct {
	MemberID          int64  `form:"member_id"`
	PaymentMethodCode string `form:"payment_method_code"`
	Page              int    `form:"page"`
	PerPage           int    `form:"per_page"`
}

type SearchMemberPayAccountResult struct {
	MemberPayAccounts []*MemberPaymentAccountDTO `json:"result"`
	Total             int64                      `json:"total"`
	Page              int                        `json:"page"`
	PerPage           int                        `json:"per_page"`
}

func (cmd CreateMemberPayAccountCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.PaymentMethodCode, validation.Required, validation.In(string(LBT), string(MOMO), string(PAYPAL))),
		validation.Field(&cmd.Detail, validation.Required),
	)
}

func (lbt *LBTDetail) ValidateLBT(detail json.RawMessage) error {
	err := json.Unmarshal(detail, lbt)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(lbt,
		validation.Field(&lbt.MemberBankCode, validation.Required),
		validation.Field(&lbt.MemberAccountNo, validation.Required),
		validation.Field(&lbt.MemberAccountName, validation.Required),
	)
}

func (paypal *PayPalDetail) ValidatePayPal(detail json.RawMessage) error {
	err := json.Unmarshal(detail, paypal)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(paypal,
		validation.Field(&paypal.PayPalEmail, validation.Required),
	)
}
