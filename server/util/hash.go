package util

import "golang.org/x/crypto/bcrypt"

/* HashPassword hashes a password with bcrypt */
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

/* VerifyPassword verifies a password with bcrypt */
func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
