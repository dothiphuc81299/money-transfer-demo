package config

import "money-transfer-demo/pkg/util/env"

type Config struct {
	Postgres        PostgresConfig
	Server          Server
	Deposit         DepositConfig
	ExchangeVNDRate float64
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	cfg.postgresConfig()
	cfg.serverConfig()
	cfg.depositConfig()

	cfg.ExchangeVNDRate = env.GetEnvAsFloat64("EXCHANGE_VND_RATE", 25)

	return cfg, nil
}
