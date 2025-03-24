package auth

import (
  "golang.org/x/crypto/bcrypt"
)

const hashCost = 10

func HashPassword(password string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
  return string(hash), err
}

func CheckPasswordHash(password, hash string) error {
  return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
