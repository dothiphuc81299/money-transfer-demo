package memberpayacc

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, cmd *CreateMemberPayAccountCommand) (*CreateMemberPaymentAccountResult,error)
	Search(ctx context.Context, query *SearchMemberPayAccountQuery) (*SearchMemberPayAccountResult, error)
}

type Store interface {
	Get(ctx context.Context, id int64) (*MemberPaymentAccountDTO, error)
	UpdateVerifyStatus(tx *gorm.DB, entity *MemberPayAccount) error
}
