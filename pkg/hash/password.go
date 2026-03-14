package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	password_hash_bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(password_hash_bytes), err
}

func ComparePassword(password string, password_hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
}
