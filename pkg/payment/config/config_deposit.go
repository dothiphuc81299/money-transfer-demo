package config

import "money-transfer-demo/pkg/util/env"

const (
	defaultDepositClientIP = "AW2SoYSo5-oYeHkLYG2mMDo_JBkth-bX0W3k3ZNXAm96Pz4c10qKFVzsLn5F2AjpYURYI2Y0Mjp9NUOq"
	defaultDepositSecret   = "EJGdh7O1GF-eVc8N5qgu06yAOzBXCwJ51lPQrXOrTUDk7H1CBLals0UXEMlLhHcaR8UH1oyMDC6zKUyW"
	defaultReturnURL       = "http://localhost:3089/api/mem/deposit"
)

type DepositConfig struct {
	ClientIP  string     `json:"client_ip"`
	Secret    string     `json:"secret"`
	ReturnURL string     `json:"return_url"`
	MoMo      MoMoConfig `json:"momo"`
}

func (cfg *Config) depositConfig() {
	cfg.Deposit.ClientIP = env.GetEnvAsString("DEPOSIT_CLIENT_IP", defaultDepositClientIP)
	cfg.Deposit.Secret = env.GetEnvAsString("DEPOSIT_SECRET", defaultDepositSecret)
	cfg.Deposit.ReturnURL = env.GetEnvAsString("DEPOSIT_RETURN_URL", defaultReturnURL)

	cfg.momoConfig()
}
