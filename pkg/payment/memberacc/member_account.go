package memberacc

import (
	"context"

	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, cmd *CreateMemberAccountCommand) error
	GetByMemberID(ctx context.Context, memberID int64) (*MemberAccount, error)
	AdjustMemberAccountBalance(ctx context.Context, tx *gorm.DB, cmd *AdjustMemberAccountBalanceCommand) error
}
