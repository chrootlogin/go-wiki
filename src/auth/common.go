package auth

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username     string   `json:"-"`
	PasswordHash string   `json:"password-hash"`
	Email        string   `json:"email"`
	IsEnabled    bool     `json:"is-enabled"`
	Permissions  []string `json:"permissions"`
}

// Hash a password with bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check a password hash
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}