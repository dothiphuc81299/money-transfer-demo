package config

import "money-transfer-demo/pkg/util/env"

const (
	defaultPayment = "payment-service:50004"
)

type GRPCClientConfig struct {
	Payment string
}

func (cfg *Config) grpcClientConfig() {
	cfg.GrpcClient.Payment = env.GetEnvAsString("GRPC_PAYMENT", defaultPayment)
}
