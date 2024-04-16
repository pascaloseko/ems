package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SecretKey is the secret key used to sign JWTs.
// It should be kept private.

var (
	SecretKey = []byte("secret")
)

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// check if username is empty
	if username == "" {
		return "", errors.New("'username' cannot be empty")
	}
	// Set token claims
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if token == nil {
		return "", errors.ErrUnsupported
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
