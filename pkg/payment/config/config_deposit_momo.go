package config

import "money-transfer-demo/pkg/util/env"

const (
	defaultDepositMoMoSubscriptionKey = "d2434deb21c4462e83e11f1aa818f7f9"
	defaultDepositMoMoAPIUser         = "9da57f94-7384-4c65-9f74-dba5eddf3edb" // X-Reference-Id used for token
	defaultDepositMoMoAPIKey          = "fe84fffb57d6463094d4d55f27d0f8b1"
	defaultDepositMoMoBaseURL         = "https://sandbox.momodeveloper.mtn.com"
)

type MoMoConfig struct {
	SubscriptionKey string
	APIUser         string
	APIKey          string
	BaseURL         string // Eg: https://sandbox.momodeveloper.mtn.com
}

func (cfg *Config) momoConfig() {
	cfg.Deposit.MoMo.SubscriptionKey = env.GetEnvAsString("DEPOSIT_MOMO_SUBSCRIPTION_KEY", defaultDepositMoMoSubscriptionKey)
	cfg.Deposit.MoMo.APIUser = env.GetEnvAsString("DEPOSIT_MOMO_API_USER", defaultDepositMoMoAPIUser)
	cfg.Deposit.MoMo.APIKey = env.GetEnvAsString("DEPOSIT_MOMO_API_KEY", defaultDepositMoMoAPIKey)
	cfg.Deposit.MoMo.BaseURL = env.GetEnvAsString("DEPOSIT_MOMO_BASE_URL", defaultDepositMoMoBaseURL)
}
