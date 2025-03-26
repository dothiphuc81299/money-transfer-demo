package memberacc

import (
	"errors"
	"money-transfer-demo/pkg/identity/member"
)

var (
	ErrMemberAccountAlreadyExists = errors.New("member account already exists")
)

type MemberAccount struct {
	ID                 int64               `json:"id"`
	MemberID           int64               `json:"member_id"`
	LoginName          string              `json:"login_name"`
	Status             member.Status       `json:"status"`
	Currency           member.CurrencyType `json:"currency"`
	Balance            float64             `json:"balance"`
	OutstandingBalance float64             `json:"outstanding_balance"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
}

type CreateMemberAccountCommand struct {
	MemberID           int64               `json:"member_id"`
	LoginName          string              `json:"login_name"`
	Status             member.Status       `json:"status"`
	Currency           member.CurrencyType `json:"currency"`
	Balance            float64             `json:"balance"`
	OutstandingBalance float64             `json:"outstanding_balance"`
}

func (MemberAccount) TableName() string {
	return "member_account"
}
