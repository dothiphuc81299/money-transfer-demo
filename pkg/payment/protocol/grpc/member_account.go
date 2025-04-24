package grpc

import (
	"context"
	"log"
	payment "money-transfer-demo/pkg/apis/payment"
	"money-transfer-demo/pkg/identity/member"
	"time"

	"money-transfer-demo/pkg/payment/memberacc"
)

const (
	GrpcRequestTimeout = 5 * time.Second
)

func (s *Server) CreateMemberAccount(ctx context.Context, req *payment.CreateMemberAccountCommand) (*payment.CreateMemberAccountResponse, error) {
	log.Println("❌ Request canceled by client")
	ctx, cancel := context.WithTimeout(ctx, GrpcRequestTimeout)
	defer cancel()

	// select {
	// case <-ctx.Done():
	// 	log.Println("❌ Request canceled by client")
	// 	return nil, status.Error(codes.Canceled, "request was canceled")
	// default:
	// }

	cmd := &memberacc.CreateMemberAccountCommand{
		MemberID:  req.MemberId,
		Currency:  member.CurrencyType(req.Currency),
		LoginName: req.LoginName,
		Status:    member.Status(req.Status),
		FullName:  req.FullName,
	}

	err := s.dependencies.MemberAccountSvc.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &payment.CreateMemberAccountResponse{}, nil
}
