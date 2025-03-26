package config

import "money-transfer-demo/pkg/util/env"

const (
	defaultPayment = "localhost:50003"
)

type GRPCClientConfig struct {
	Payment string
}

func (cfg *Config) grpcClientConfig() {
	cfg.GrpcClient.Payment = env.GetEnvAsString("GRPC_PAYMENT", defaultPayment)
}
