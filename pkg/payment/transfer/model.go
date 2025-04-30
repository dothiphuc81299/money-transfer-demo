package transfer

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrMemberAccountNotFound         = errors.New("member account not found")
	ErrMemberAccountNotEnoughBalance = errors.New("member account not enough balance")
	ErrTransferNotFound              = errors.New("transfer not found")
	ErrInvalidCurrency               = errors.New("invalid currency")
	ErrUnauthorized                  = errors.New("unauthorized")
	ErrInvalidStatus                 = errors.New("invalid status")
)

type Status int

const (
	Pending Status = iota + 1
	Failed
	Successful
)

type Transfer struct {
	ID            int64   `json:"id"`
	FromMemberID  int64   `json:"from_member_id"`
	FromLoginName string  `json:"from_login_name"`
	ToMemberID    int64   `json:"to_member_id"`
	ToLoginName   string  `json:"to_login_name"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Status        Status  `json:"status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

func (t Transfer) TableName() string {
	return "transfer"
}

type TransferDTO struct {
	ID               int64               `json:"id"`
	TransactionID    string              `json:"transaction_id"`
	FromMemberID     int64               `json:"from_member_id"`
	FromLoginName    string              `json:"from_login_name"`
	ToMemberID       int64               `json:"to_member_id"`
	ToLoginName      string              `json:"to_login_name"`
	Amount           float64             `json:"amount"`
	Status           Status              `json:"status"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
	TransferTimeline []*TransferTimeline `json:"transfer_timeline,omitempty" gorm:"-"`
}

type TransferTimeline struct {
	ID            int64  `json:"id"`
	TransferID    int64  `json:"transfer_id"`
	TransactionID string `json:"transaction_id"`
	Message       string `json:"message"`
	Note          string `json:"note,omitempty"`
	CreatedBy     string `json:"created_by"`
	CreatedAt     string `json:"created_at"`
}

func (t TransferTimeline) TableName() string {
	return "transfer_timeline"
}

type CreateTransferCommand struct {
	FromMemberID  int64   `json:"from_member_id"`
	ToMemberID    int64   `json:"to_member_id"`
	Amount        float64 `json:"amount"`
	TransactionID string
}
type CreateTransferResult struct {
	ID int64 `json:"id"`
}

func (c CreateTransferCommand) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ToMemberID, validation.Required),
		validation.Field(&c.Amount, validation.Required),
	)
}

type UpdateTransferStatusCommand struct {
	ID     int64  `json:"id"`
	Note   string `json:"note"`
	Status Status `json:"status"`
}

type SearchTransferQuery struct {
	Page          int    `form:"page"`
	PerPage       int    `form:"per_page"`
	TransactionID string `form:"transaction_id"`
	Status        Status `form:"status"`
	FromMemberID  int64  `form:"from_member_id"`
	ToMemberID    int64  `form:"to_member_id"`
}

type SearchTransferResult struct {
	Total   int64          `json:"total"`
	Data    []*TransferDTO `json:"data"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

func (c UpdateTransferStatusCommand) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ID, validation.Required),
		validation.Field(&c.Status, validation.Required, validation.In(Failed, Successful)),
	)
}

func Message(status Status, name string) string {
	var message string
	switch status {
	case Pending:
		message = "%s created this transfer"
	case Successful:
		message = "%s manually marked this transfer as successful"
	case Failed:
		message = "%s manually marked this transfer as failed"
	}

	return fmt.Sprintf(message, name)
}
