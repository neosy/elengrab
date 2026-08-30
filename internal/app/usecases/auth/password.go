package auth

import "golang.org/x/crypto/bcrypt"

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
