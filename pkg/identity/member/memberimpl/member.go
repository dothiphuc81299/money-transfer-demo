package memberimpl

import (
	"context"
	"fmt"
	paymentapi "money-transfer-demo/pkg/apis/payment"
	"money-transfer-demo/pkg/identity/member"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/util/password"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	TimestampFormat = time.RFC3339
)

type service struct {
	store         *store
	paymentClient paymentapi.PaymentClient
}

func NewService(store *store, paymentClient paymentapi.PaymentClient) member.Service {
	return &service{
		store:         store,
		paymentClient: paymentClient,
	}
}

func (s *service) CreateMember(ctx context.Context, cmd *member.CreateMemberCommand) (*member.CreateMemberResult, error) {
	var (
		now = time.Now().UTC().Format(TimestampFormat)
		id  int64
	)

	result, err := s.store.isMemberTaken(ctx, 0, cmd.LoginName, cmd.Email, cmd.Phone)
	if err != nil {
		return nil, err
	}

	if len(result) != 0 {
		if result[0].Email == cmd.Email {
			return nil, member.ErrEmailExists
		}

		if result[0].Phone == cmd.Phone {
			return nil, member.ErrPhoneNumberExists
		}

		if result[0].LoginName == cmd.LoginName {
			return nil, member.ErrLoginNameExists
		}

		return nil, member.ErrMemberExists
	}

	cmd.Phone = strings.ReplaceAll(cmd.Phone, "+", "")
	cmd.Phone = string(cmd.PrefixPhone) + cmd.Phone

	entity := member.Member{
		LoginName:         cmd.LoginName,
		Status:            member.Active,
		Currency:          cmd.Currency,
		FullName:          cmd.FullName,
		Email:             cmd.Email,
		Phone:             cmd.Phone,
		EmailVerifyStatus: member.EmailUnverified,
		PhoneVerifyStatus: member.PhoneUnverified,
		Address:           cmd.Address,
		UpdatedAt:         now,
		CreatedAt:         now,
	}

	password, err := password.HashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	entity.Password = password

	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		id, err = s.store.createMember(tx, &entity)
		if err != nil {
			return err
		}

		_, err = s.paymentClient.CreateMemberAccount(ctx, &paymentapi.CreateMemberAccountCommand{
			MemberId:  id,
			LoginName: cmd.LoginName,
			Currency:  string(cmd.Currency),
			Status:    int64(member.Active),
			FullName:  cmd.FullName,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &member.CreateMemberResult{ID: id}, nil
}

func (s *service) GetMemberByID(ctx context.Context, id string) (*member.Member, error) {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	fmt.Println("account", account, ok)
	if !ok || account == nil {
		return nil, nil
	}
	return s.store.getMemberByID(ctx, id)
}
