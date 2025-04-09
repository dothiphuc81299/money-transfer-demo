package bankacc

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrBankAccountAlreadyExists = errors.New("bank account already exists")
	ErrBankAccountNotFound      = errors.New("bank account not found")
	ErrBankAccountNotActive     = errors.New("bank account not active")
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
	OutstandingBalance float64
}

type AdjustBankAccountBalanceCommand struct {
	BankAccountID int64
	ChangedAmount float64
	UpdatedAt     string
}

func (cmd CreateBankAccountCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.BankCode, validation.Required, validation.In(toInterfaceSlice(bank_codes)...)),
		validation.Field(&cmd.AccountNo, validation.Required),
	)
}


func toInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
