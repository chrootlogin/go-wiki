package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/patrickmn/go-cache"
)

var authCache *cache.Cache

func init() {
	authCache = cache.New(30 * time.Minute, 10 * time.Minute)
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