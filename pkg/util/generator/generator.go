package generator

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateTransactionID(paymentCode string) string {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	random := rand.Intn(100000)
	return fmt.Sprintf("%s%s%05d", paymentCode, timestamp, random)
}
