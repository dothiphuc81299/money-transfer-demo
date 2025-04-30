package memberacc

import (
	"errors"
	"money-transfer-demo/pkg/identity/member"
	"money-transfer-demo/pkg/payment/transaction"
)

var (
	ErrMemberAccountAlreadyExists    = errors.New("member account already exists")
	ErrMemberAccountNotFound         = errors.New("member account not found")
	ErrMemberAccountNotEnoughBalance = errors.New("member account not enough balance")
)

type MemberAccount struct {
	ID                 int64               `json:"id"`
	MemberID           int64               `json:"member_id"`
	LoginName          string              `json:"login_name"`
	FullName           string              `json:"full_name"`
	Status             member.Status       `json:"status"`
	Currency           member.CurrencyType `json:"currency"`
	Balance            float64             `json:"balance"`
	OutstandingBalance float64             `json:"outstanding_balance"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
}

type CreateMemberAccountCommand struct {
	MemberID           int64               `json:"member_id"`
	FullName           string              `json:"full_name"`
	LoginName          string              `json:"login_name"`
	Status             member.Status       `json:"status"`
	Currency           member.CurrencyType `json:"currency"`
	Balance            float64             `json:"balance"`
	OutstandingBalance float64             `json:"outstanding_balance"`
}

func (MemberAccount) TableName() string {
	return "member_account"
}

type AdjustMemberAccountBalanceCommand struct {
	MemberID                  int64
	UpdatedBy                 string
	AdjustedAmount            float64
	AdjustedOutstandingAmount float64
	TransactionType           transaction.Type
}

type UpdateMemberAccountBalanceCommand struct {
	ID                        int64
	AdjustedAmount            float64
	AdjustedOutstandingAmount float64
	UpdatedAt                 string
}

func VerifyBalance(accountBalance float64, outstandingBalance float64, adjustedAmount float64, adjustedOutstandingAmount float64, transactionType transaction.Type) bool {
	var status bool = false
	switch transactionType {
	case transaction.DepositType:

		status = true
	case transaction.WithdrawalType, transaction.TransferType:
		status =
			(accountBalance >= outstandingBalance) &&
				(accountBalance+adjustedAmount >= outstandingBalance+adjustedOutstandingAmount) &&
				(outstandingBalance+adjustedOutstandingAmount >= 0)
	}

	return status
}
