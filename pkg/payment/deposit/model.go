package deposit

import (
	"encoding/json"
	"errors"
	"fmt"
	"money-transfer-demo/pkg/identity/member"

	"gorm.io/datatypes"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrUnauthorized                = errors.New("unauthorized")
	ErrInvalidMemberBankCode       = errors.New("invalid member bank code")
	ErrInvalidMemberAccountNo      = errors.New("invalid member account no")
	ErrInvalidMemberAccountName    = errors.New("invalid member account name")
	ErrInvalidCompanyBankCode      = errors.New("invalid company bank code")
	ErrInvalidCompanyAccountNo     = errors.New("invalid company account no")
	ErrInvalidCompanyAccountName   = errors.New("invalid company account name")
	ErrDepositNotFound             = errors.New("deposit not found")
	ErrBankAccountNotFound         = errors.New("bank account not found")
	ErrInvalidCurrency             = errors.New("invalid currency")
	ErrDepositNotProcessing        = errors.New("deposit not processing")
	ErrApproveURLNotFound          = errors.New("approve url not found")
	ErrUnauthorizedUserIsNotMember = errors.New("unauthorized: user is not member")
	ErrPaymentNotAunthorized       = errors.New("payment not authorized")
)

const (
	DefaultUser       = "system"
	DefaultPaypalNote = "Verify paypal from user"
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
	GrossAmount       float64        `json:"gross_amount"`
	NetAmount         float64        `json:"net_amount"`
	ChargeAmount      float64        `json:"charge_amount"`
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
	Currency          string            `json:"currency"`
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
	GrossAmount   float64
	NetAmount     float64
	ChargeAmount  float64
	TransactionID string
	UpdatedBy     string
	DetailStr     string
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

type CreateDepositPaypalResult struct {
	ApproveURL string `json:"approve_url"`
}

type VerifyPaypalCommand struct {
	Token   string `form:"token"`
	PayerID string `form:"PayerID"`
}

type CreateDepositPaypalCommand struct {
	TransactionID     string
	MemberID          int64               `json:"member_id"`
	LoginName         string              `json:"login_name"`
	PaymentMethodCode PaymentMethodCode   `json:"payment_method_code"`
	Amount            float64             `json:"amount"`
	Currency          member.CurrencyType `json:"currency"`
}

type PayPalDetail struct {
	PaypalTransactionID string `json:"paypal_transaction_id"`
	PayerID             string `json:"payer_id"`
	PayerEmail          string `json:"payer_email"`
	PayerName           string `json:"payer_name"` // given name + surname
	OrderID             string `json:"order_id"`
	CountryCode         string `json:"country_code"`
	ShippingName        string `json:"shipping_name"`
	ShippingAddress     string `json:"shipping_address"`
}

type CreateDepositMoMoCommand struct {
	MemberID          int64   `json:"member_id"`
	LoginName         string  `json:"login_name"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`     //// "EUR" (MoMo sandbox uses EUR)
	PhoneNumber       string  `json:"phone_number"` //MSISDN
	PaymentMethodCode string  `json:"payment_method_code"`
	TransactionID     string
}

type CreateDepositMoMoResult struct {
	ReferenceID string `json:"reference_id"` // UUID used as the MoMo transaction reference
	Status      string `json:"status"`       // Current status of the transaction (e.g., "PENDING")
}

type MoMoRequestResponse struct {
	StatusCode  int
	ReferenceID string
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

func (cmd VerifyPaypalCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Token, validation.Required),
		validation.Field(&cmd.PayerID, validation.Required),
	)
}

func (cmd CreateDepositPaypalCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.PaymentMethodCode, validation.Required, validation.In(PAYPAL)),
		validation.Field(&cmd.Amount, validation.Required),
		validation.Field(&cmd.Currency, validation.Required, validation.In(member.VietnamDong, member.UnitedStatesDollar).Error(ErrInvalidCurrency.Error())),
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
