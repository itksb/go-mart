package token

import (
	"time"
)

type CreateTokenFunc func(martClaims *MartClaims, secretReader Secret) (newToken string, err error)
type ParseWithClaimsFunc func(tokenString string, claims *MartClaims, secretReader Secret) error

// Secret defines interface returning secret key
type Secret interface {
	Get() (string, error)
}

// SecretFunc type is an adapter to allow the use of ordinary functions as Secret. If f is a function
// with the appropriate signature, SecretFunc(f) is a Handler that calls f.
type SecretFunc func() (string, error)

// Get calls f()
func (f SecretFunc) Get() (string, error) {
	return f()
}

// Opts - options for the MartToken
type MartClaims struct {
	UserID    string
	ExpiresAt time.Time
}
