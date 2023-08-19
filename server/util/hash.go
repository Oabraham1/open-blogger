package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

/* HashPassword hashes a password with bcrypt */
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return string(hashedPassword), nil
}

/* VerifyPassword verifies a password with bcrypt */
func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
