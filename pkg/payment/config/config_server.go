package config

import "money-transfer-demo/pkg/util/env"

const (
	defaultHTTPPort = "5090"
	defaultGRPCPort = "50004"
)

type Server struct {
	Host           string
	HTTPPort       string
	GRPCPort       string
	StaticRootPath string
}

func (cfg *Config) serverConfig() {
	cfg.Server.HTTPPort = env.GetEnvAsString("PAYMENT_HTTP_PORT", defaultHTTPPort)
	cfg.Server.GRPCPort = env.GetEnvAsString("PAYMENT_GRPC_PORT", defaultGRPCPort)

}
