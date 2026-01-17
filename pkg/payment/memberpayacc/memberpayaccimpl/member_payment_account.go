package memberpayaccimpl

import (
	"context"
	"encoding/json"
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/payment/memberacc"
	"money-transfer-demo/pkg/payment/memberpayacc"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type service struct {
	store        *store
	memberAccSrv memberacc.Service
}

func NewService(store *store, memberAccSrv memberacc.Service) memberpayacc.Service {
	return &service{store: store, memberAccSrv: memberAccSrv}
}

func (s *service) Create(ctx context.Context, cmd *memberpayacc.CreateMemberPayAccountCommand) (*memberpayacc.CreateMemberPaymentAccountResult, error) {
	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return nil, memberpayacc.ErrUnauthorized
	}

	cmd.CreatedBy = account.LoginName
	ma, err := s.memberAccSrv.GetByMemberID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	if ma == nil {
		return nil, memberpayacc.ErrMemberNotFound
	}

	if cmd.PaymentMethodCode == string(memberpayacc.LBT) {
		cmd.MemberFullName = ma.FullName
		detail, err := s.getDetailLBT(cmd)
		if err != nil {
			return nil, err
		}

		cmd.DetailStr = detail
	} else if cmd.PaymentMethodCode == string(memberpayacc.PAYPAL) {
		detail, err := s.getDetailPayPal(cmd)
		if err != nil {
			return nil, err
		}

		cmd.DetailStr = detail
	}

	exist, err := s.store.taken(ctx, account.ID, cmd.PaymentMethodCode, cmd.PaypalEmail, cmd.MemberBankCode)
	if err != nil {
		return nil, err
	}

	if len(exist) > 0 {
		return nil, memberpayacc.ErrMemberPaymentAccountAlreadyExists
	}

	entity := &memberpayacc.MemberPayAccount{
		MemberID:          account.ID,
		PaymentMethodCode: cmd.PaymentMethodCode,
		Detail:            datatypes.JSON(cmd.DetailStr),
		VerifyStatus:      memberpayacc.Unverified,
		CreatedBy:         cmd.CreatedBy,
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:         time.Now().UTC().Format(time.RFC3339),
	}

	var id int64
	err = s.store.db.Transaction(func(tx *gorm.DB) error {
		newID, err := s.store.create(tx, entity)
		if err != nil {
			return err
		}
		id = newID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &memberpayacc.CreateMemberPaymentAccountResult{ID: id}, nil
}

func (s *service) getDetailLBT(cmd *memberpayacc.CreateMemberPayAccountCommand) (string, error) {
	var lbt memberpayacc.LBTDetail
	err := lbt.ValidateLBT(cmd.Detail)
	if err != nil {
		return "", err
	}

	cmd.MemberBankCode = lbt.MemberBankCode
	lbt.MemberFullName = cmd.MemberFullName
	dt, err := json.Marshal(lbt)
	if err != nil {
		return "", err
	}

	return string(dt), nil
}

func (s *service) getDetailPayPal(cmd *memberpayacc.CreateMemberPayAccountCommand) (string, error) {
	var paypal memberpayacc.PayPalDetail
	err := paypal.ValidatePayPal(cmd.Detail)
	if err != nil {
		return "", err
	}

	cmd.PaypalEmail = paypal.PayPalEmail
	dt, err := json.Marshal(paypal)
	if err != nil {
		return "", err
	}

	return string(dt), nil
}

func (s *service) Search(ctx context.Context, query *memberpayacc.SearchMemberPayAccountQuery) (*memberpayacc.SearchMemberPayAccountResult, error) {
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.PerPage <= 0 {
		query.PerPage = 10
	}

	account, ok := ctx.Value("current_account").(*token.AccountData)
	if !ok || account == nil {
		return nil, memberpayacc.ErrUnauthorized
	}

	query.MemberID = account.ID
	result, err := s.store.search(ctx, query)
	if err != nil {
		return nil, err
	}

	result.Page = query.Page
	result.PerPage = query.PerPage
	return result, nil
}
