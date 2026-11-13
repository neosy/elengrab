package auth

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword hashes a plain text password.
func (a *auth) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(bytes), err
}

// checkPassword checks if the provided password matches the hash.
func (a *auth) checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

// GenerateRandomPassword generates a cryptographically secure random password.
func GenerateRandomPassword() (string, error) {
	const (
		letters  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		digits   = "0123456789"
		specials = "!@#$%^&*"
		all      = letters + digits + specials
		length   = 32
	)

	password := make([]byte, length)

	// Ensure all required character types are present.
	for i, chars := range []string{letters, digits, specials} {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		password[i] = chars[n.Int64()]
	}

	// Fill the remaining characters.
	for i := 3; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}
		password[i] = all[n.Int64()]
	}

	// Shuffle the password characters.
	for i := len(password) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := int(n.Int64())
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}
