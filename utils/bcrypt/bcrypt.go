package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateFromPassword hashes a plaintext password
func GenerateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHashAndPassword compares a bcrypt hash with a plaintext password
func CompareHashAndPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
