package config

type Config struct {
	Postgres PostgresConfig
	Server   Server
	Deposit  DepositConfig
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	cfg.postgresConfig()
	cfg.serverConfig()
	cfg.depositConfig()

	return cfg, nil
}
