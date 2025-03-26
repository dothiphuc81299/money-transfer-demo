package member

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	ErrInvalidCurrency                          = errors.New("invalid currency")
	ErrLoginNameMustNotContainSpecialCharacters = errors.New("login name must not contain special characters")
	ErrInvalidPhoneNumber                       = errors.New("invalid phone number")
	ErrInvalidPrefixPhone                       = errors.New("invalid prefix phone")
	ErrEmailExists                              = errors.New("email already exists")
	ErrPhoneNumberExists                        = errors.New("phone number already exists")
	ErrMemberExists                             = errors.New("member already exists")
	ErrLoginNameExists                          = errors.New("login name already exists")
	ErrMemberNotFound                           = errors.New("member not found")
	ErrMemberInactive                           = errors.New("member inactive")
	ErrInvalidPassword                          = errors.New("invalid password")
)

type Status int

const (
	Active Status = iota + 1
	InActive
	Lock
)

type EmailVerifyStatus int

const (
	EmailUnverified EmailVerifyStatus = iota + 1
	EmailVerified
)

type PhoneVerifyStatus int

const (
	PhoneUnverified PhoneVerifyStatus = iota + 1
	PhoneVerified
)

type CurrencyType string

const (
	VietnamDong        CurrencyType = "VND"
	UnitedStatesDollar CurrencyType = "USD"
)

type Member struct {
	ID                int64             `json:"id"`
	LoginName         string            `json:"login_name"`
	Password          string            `json:"-"`
	Status            Status            `json:"status"`
	Currency          CurrencyType      `json:"currency"`
	FullName          string            `json:"full_name"`
	Email             string            `json:"email"`
	Phone             string            `json:"phone"`
	Address           string            `json:"address"`
	EmailVerifyStatus EmailVerifyStatus `json:"email_verify_status"`
	PhoneVerifyStatus PhoneVerifyStatus `json:"phone_verify_status"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
}

type CreateMemberCommand struct {
	LoginName   string       `json:"login_name"`
	Password    string       `json:"password"`
	Currency    CurrencyType `json:"currency"`
	FullName    string       `json:"full_name"`
	Email       string       `json:"email"`
	Phone       string       `json:"phone"`
	Address     string       `json:"address"`
	PrefixPhone string       `json:"prefix_phone"`
}

type CreateMemberResult struct {
	ID int64 `json:"id"`
}

// GORM will use `member` table
func (Member) TableName() string {
	return "member"
}

type LoginMemberCommand struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type LoginMemberResult struct {
	LoginName   string `json:"login_name,omitempty"`
	Currency    string `json:"currency,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

func (cmd CreateMemberCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.LoginName, validation.Required, validation.Length(6, 20), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$")).Error("must be alphanumeric")),
		validation.Field(&cmd.Password, validation.Required, validation.Length(6, 0)),
		validation.Field(&cmd.Currency, validation.Required, validation.In(VietnamDong, UnitedStatesDollar).Error(ErrInvalidCurrency.Error())),
		validation.Field(&cmd.FullName, validation.Required, validation.Length(3, 0)),
		validation.Field(&cmd.Email, validation.Required, is.Email),
		validation.Field(&cmd.Phone, validation.Required, validation.By(validatePhone)),
		validation.Field(&cmd.PrefixPhone, validation.Required),
	)
}

func validatePhone(value interface{}) error {
	phone, _ := value.(string)
	phone = strings.Replace(phone, "+", "", -1)
	if !regexp.MustCompile(`^\d{9,11}$`).MatchString(phone) {
		return ErrInvalidPhoneNumber
	}
	return nil
}

func (cmd LoginMemberCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.LoginName, validation.Required, validation.Length(6, 20), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$")).Error("must be alphanumeric")),
		validation.Field(&cmd.Password, validation.Required, validation.Length(6, 0)),
	)
}
