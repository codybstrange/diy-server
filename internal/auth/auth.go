package auth

import (
  "golang.org/x/crypto/bcrypt"
  "time"
  "github.com/google/uuid"
  "github.com/golang-jwt/jwt/v5"
  "fmt"
  "errors"
  "strings"
  "net/http"
)

const hashCost = 10
type TokenType string
const (
  TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
  return string(hash), err
}

func CheckPasswordHash(password, hash string) error {
  return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error){
  val := headers.Get("Authorization")
  if val == "" {
    return "", errors.New("No Authorization header found")
  }
  splitToken := strings.Split(val, " ")
  if len(splitToken) != 2 {
    return "", errors.New("Incorrect format for the authorization header value")
  }
  return splitToken[1], nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
  claims := jwt.RegisteredClaims{
    Issuer: string(TokenTypeAccess),
    IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
    ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
    Subject: userID.String(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

