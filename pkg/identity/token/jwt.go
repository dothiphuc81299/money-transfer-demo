package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

type AccountType int

const (
	User AccountType = iota + 1
	Member
)

type Claims struct {
	UserID      int64       `json:"user_id"`
	LoginName   string      `json:"login_name"`
	AccountType AccountType `json:"account_type"`
	jwt.RegisteredClaims
}

type AccountData struct {
	LoginName   string      `json:"login_name"`
	ID          int64       `json:"id"`
	AccountType AccountType `json:"account_type"`
}


func GenerateJWT(userID int64, loginName string, accountType AccountType) (string, error) {
	expirationTime := time.Now().UTC().Add(2 * time.Hour)

	claims := &Claims{
		UserID:      userID,
		LoginName:   loginName,
		AccountType: accountType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT1(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func ValidateJWT(tokenStr string, expectedType AccountType) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	if claims.AccountType != expectedType {
		return nil, fmt.Errorf("unauthorized: account type mismatch")
	}

	return claims, nil
}
