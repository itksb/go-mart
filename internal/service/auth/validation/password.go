package validation

import "errors"

var ErrValidationPasswordToShort = errors.New("password to short")

func ValidatePassword(p string) (ok bool, err error) {
	if p == "" {
		return false, ErrValidationPasswordToShort
	}
	return true, nil
}
