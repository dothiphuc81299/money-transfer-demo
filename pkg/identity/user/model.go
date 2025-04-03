package user

import validation "github.com/go-ozzo/ozzo-validation/v4"

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
	Salt      string `json:"-"`
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

func (cmd CreateUserCommand) Validate() error {
	return validation.ValidateStruct(&cmd,
		validation.Field(cmd.LoginName, validation.Required),
		validation.Field(cmd.Password, validation.Required),
	)
}
