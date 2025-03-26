package app

import (
	"money-transfer-demo/pkg/identity/config"
	"time"

	paymentapi "money-transfer-demo/pkg/apis/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GrpcRequestTimeout = 5 * time.Second
)

func getPaymentClient(cfg *config.Config) (paymentapi.PaymentClient, error) {
	conn, err := grpc.NewClient(cfg.GrpcClient.Payment, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return paymentapi.NewPaymentClient(conn), nil
}
