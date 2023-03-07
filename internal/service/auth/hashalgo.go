package auth

import "golang.org/x/crypto/bcrypt"

type HashAlgoInterface interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

type HashAlgoBcrypt struct{}

func NewHashAlgoBcrypt() (*HashAlgoBcrypt, error) {
	return &HashAlgoBcrypt{}, nil
}

func (HashAlgoBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (HashAlgoBcrypt) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
