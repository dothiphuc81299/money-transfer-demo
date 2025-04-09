package user

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrUserLoginnameAlreadyExists = errors.New("user login name already exists")
	ErrUserNotFound               = errors.New("user not found")
	ErrUserInactive               = errors.New("user inactive")
)

type Status int

const (
	Active Status = iota + 1
	InActive
)

type User struct {
	ID        int64  `json:"id"`
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
	Status    Status `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *User) TableName() string {
	return "user"
}

type CreateUserCommand struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type LoginUserCommand struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type LoginUserResult struct {
	LoginName   string `json:"login_name,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

func (cmd CreateUserCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.LoginName, validation.Required),
		validation.Field(&cmd.Password, validation.Required),
	)
}

func (cmd LoginUserCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.LoginName, validation.Required),
		validation.Field(&cmd.Password, validation.Required),
	)
}
