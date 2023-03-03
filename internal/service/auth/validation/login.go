package validation

import "errors"

var ErrValidationLoginToShort = errors.New("login to short")
var ErrValidationLoginToLong = errors.New("login to short")

const LoginMaxLen = 30

func ValidateLogin(p string) (ok bool, err error) {
	if p == "" {
		return false, ErrValidationLoginToShort
	}
	if len(p) > LoginMaxLen {
		return false, ErrValidationLoginToLong
	}
	return true, nil
}
