package bankacc

import (
	"errors"
	"money-transfer-demo/pkg/identity/member"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrBankAccountAlreadyExists = errors.New("bank account already exists")
	ErrBankAccountNotFound      = errors.New("bank account not found")
	ErrBankAccountNotActive     = errors.New("bank account not active")
	ErrInvalidCurrency          = errors.New("invalid currency")
	ErrInsufficientBalance      = errors.New("insufficient balance")
)

type BankAccountStatus int

const (
	Active BankAccountStatus = iota + 1
	InActive
)

var (
	bank_codes = []string{
		"VCB001",
		"TCB002",
		"ACB003",
		"BIDV004",
		"VIB005",
		"MB001",
		"VPB006",
		"SHB007",
		"TBB008",
		"SCB009",
	}
)

type BankAccount struct {
	ID                 int64             `json:"id"`
	BankCode           string            `json:"bank_code"`
	AccountNo          string            `json:"account_no"`
	Balance            float64           `json:"balance"`
	OutstandingBalance float64           `json:"outstanding_balance"`
	Status             BankAccountStatus `json:"status"`
	Currency           string            `json:"currency"`
	CreatedAt          string            `json:"created_at"`
	UpdatedAt          string            `json:"updated_at"`
}

func (b BankAccount) TableName() string {
	return "bank_account"
}

type CreateBankAccountCommand struct {
	BankCode           string  `json:"bank_code"`
	AccountNo          string  `json:"account_no"`
	Balance            float64 `json:"balance"`
	Currency           string  `json:"currency"`
	OutstandingBalance float64
}

type AdjustBankAccountBalanceCommand struct {
	BankAccountID int64
	Currency      string
	ChangedAmount float64
	UpdatedAt     string
}

func (cmd CreateBankAccountCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.BankCode, validation.Required, validation.In(toInterfaceSlice(bank_codes)...)),
		validation.Field(&cmd.AccountNo, validation.Required),
		validation.Field(&cmd.Currency, validation.Required, validation.In(string(member.VietnamDong), string(member.UnitedStatesDollar)).Error(ErrInvalidCurrency.Error())),
	)
}

func toInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
