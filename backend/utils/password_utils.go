package utils

import "golang.org/x/crypto/bcrypt"

// Verificar contraseña (hash bcrypt)
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
