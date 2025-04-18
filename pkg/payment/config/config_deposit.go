package config

import "money-transfer-demo/pkg/util/env"

const (
	defaultDepositClientIP      = "AW2SoYSo5-oYeHkLYG2mMDo_JBkth-bX0W3k3ZNXAm96Pz4c10qKFVzsLn5F2AjpYURYI2Y0Mjp9NUOq"
	defaultDepositSecret        = "EJGdh7O1GF-eVc8N5qgu06yAOzBXCwJ51lPQrXOrTUDk7H1CBLals0UXEMlLhHcaR8UH1oyMDC6zKUyW"
	defaultReturnURL            = "http://localhost:3089/api/mem/deposit"
	defaultShippingFullName     = "No Promise"
	defaultShippingAddressLine1 = "1 Main St"
	defaultShippingAddressCity  = "San Jose"
	ShippingAddressState        = "CA"
	ShippingCountryCode         = "US"
	ShippingPostalCode          = "95131"
)

type DepositConfig struct {
	ClientIP             string `json:"client_ip"`
	Secret               string `json:"secret"`
	ReturnURL            string `json:"return_url"`
	ShippingFullName     string `json:"shipping_full_name"`
	ShippingAddressLine1 string `json:"shipping_address_line1"`
	ShippingAddressCity  string `json:"shipping_address_city"`
	ShippingAddressState string `json:"shipping_address_state"`
	ShippingCountryCode  string `json:"shipping_country"`
	ShippingPostalCode   string `json:"shipping_postal_code"`
}

func (cfg *Config) depositConfig() {
	cfg.Deposit.ClientIP = env.GetEnvAsString("DEPOSIT_CLIENT_IP", defaultDepositClientIP)
	cfg.Deposit.Secret = env.GetEnvAsString("DEPOSIT_SECRET", defaultDepositSecret)
	cfg.Deposit.ReturnURL = env.GetEnvAsString("DEPOSIT_RETURN_URL", defaultReturnURL)

	cfg.Deposit.ShippingFullName = env.GetEnvAsString("DEPOSIT_SHIPPING_FULL_NAME", defaultShippingFullName)
	cfg.Deposit.ShippingAddressLine1 = env.GetEnvAsString("DEPOSIT_SHIPPING_ADDRESS_LINE1", defaultShippingAddressLine1)
	cfg.Deposit.ShippingAddressCity = env.GetEnvAsString("DEPOSIT_SHIPPING_ADDRESS_CITY", defaultShippingAddressCity)
	cfg.Deposit.ShippingAddressState = env.GetEnvAsString("DEPOSIT_SHIPPING_ADDRESS_STATE", ShippingAddressState)
	cfg.Deposit.ShippingCountryCode = env.GetEnvAsString("DEPOSIT_SHIPPING_COUNTRY_CODE", ShippingCountryCode)
	cfg.Deposit.ShippingPostalCode = env.GetEnvAsString("DEPOSIT_SHIPPING_POSTAL_CODE", ShippingPostalCode)

}
