package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SK"))

type Claim struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claim := Claim{
		UserID:           userID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expirationTime)},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (int, error) {
	claim := &Claim{}
	token, err := jwt.ParseWithClaims(tokenStr, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claim.UserID, nil
}
