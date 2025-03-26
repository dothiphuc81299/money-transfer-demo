package member

import "time"

type Status int

const (
	Active Status = iota + 1
	InActive
	Lock
)

type Member struct {
	ID        uint   `json:"id"`
	Uuid      string `json:"uuid"`
	LoginName string `json:"login_name"`
	Password  string `json:"-"`
	Status    Status `json:"status"`
	Currency  string `json:"currency"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`

	Name      string `json:"name"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateMemberCommand struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *CreateMemberCommand) Validate() error {
	return nil
}
