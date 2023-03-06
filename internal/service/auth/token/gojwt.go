package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type customClaims struct {
	UserID string
	jwt.RegisteredClaims
}

func CreateToken(martClaims *MartClaims, secretReader Secret, timeNow func() time.Time) (token string, err error) {
	signingKey, err := secretReader.Get()
	if err != nil {
		return "", err
	}

	claims := customClaims{
		UserID: martClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-mart",
			ExpiresAt: jwt.NewNumericDate(martClaims.ExpiresAt),
			NotBefore: jwt.NewNumericDate(timeNow()),
			IssuedAt:  jwt.NewNumericDate(timeNow()),
			ID:        "",
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString(signingKey)
}

func ParseWithClaims(tokenString string, martClaims *MartClaims, secretReader Secret) error {
	claims := customClaims{}
	signingKey, err := secretReader.Get()
	if err != nil {
		return err
	}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrTokenInvalid
	}

	martClaims.UserID = claims.UserID
	martClaims.ExpiresAt = claims.ExpiresAt.Time

	return nil
}
