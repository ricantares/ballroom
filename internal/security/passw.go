package security

import (
	"math/rand"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// Genera un hash per la password fornita
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Verifica la password fornita con lo hash memorizzato.
func VerifyPassword(encrypted string, plain string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(plain))
	return err == nil
}

// Genera una nuova password utente (default 12 caratteri)
func ResetPassword() (string, error) {
	const (
		lettere  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numeri   = "0123456789"
		speciali = "!@#$%^&*()_+-=;':,.<>/?"
	)
	var password []byte

	length, err := strconv.ParseInt(os.Getenv("PASSWLEN"), 10, 16)
	if err != nil {
		length = 12
	}

	charset := lettere + numeri + speciali

	for i := 0; i < int(length); i++ {
		randNum := rand.Intn(len(charset))
		password = append(password, charset[randNum])
	}
	return string(password), nil
}
